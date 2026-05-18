package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"sneaker-store/ai-service/internal/model"
)

type AIUseCase struct {
	chatRepo  ChatRepository
	cache     AICache
	publisher EventPublisher
	catalog   ProductCatalog
	gemini    GeminiClient
}

func NewAIUseCase(
	chatRepo ChatRepository,
	cache AICache,
	publisher EventPublisher,
	catalog ProductCatalog,
	gemini GeminiClient,
) *AIUseCase {
	return &AIUseCase{
		chatRepo:  chatRepo,
		cache:     cache,
		publisher: publisher,
		catalog:   catalog,
		gemini:    gemini,
	}
}

const sneakerSystemPrompt = `You are VAULT AI, an expert sneaker advisor for SNKR VAULT — a premium sneaker marketplace.
You have deep knowledge of sneaker culture, brands (Nike, Adidas, Jordan, New Balance, Asics, Salomon, etc.),
colorways, sizing, resale value, and trends. Help users find the perfect sneakers.
Keep responses concise, enthusiastic, and knowledgeable. Use sneaker culture terminology naturally.
When recommending products, be specific about models, colorways, and why they match the user's needs.`

func (uc *AIUseCase) AskSneakerAdvice(ctx context.Context, sessionID, userID, message string) (string, string, error) {
	if sessionID == "" {
		sessionID = uuid.NewString()
	}

	cacheKey := cacheHash(sessionID, message)
	if cached, err := uc.cache.GetAdvice(ctx, cacheKey); err == nil {
		return cached, sessionID, nil
	}

	history, err := uc.chatRepo.GetHistory(ctx, sessionID, 10)
	if err != nil {
		log.Printf("get history: %v", err)
	}

	historyText := buildHistoryText(history)

	reply, err := uc.gemini.Chat(ctx, sneakerSystemPrompt, historyText, message)
	if err != nil {
		return "", sessionID, fmt.Errorf("gemini chat: %w", err)
	}

	userMsg := &model.ChatMessage{
		ID:        uuid.NewString(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      "user",
		Content:   message,
		CreatedAt: time.Now().UTC(),
	}
	assistantMsg := &model.ChatMessage{
		ID:        uuid.NewString(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      "assistant",
		Content:   reply,
		CreatedAt: time.Now().UTC(),
	}

	if err := uc.chatRepo.SaveMessage(ctx, userMsg); err != nil {
		log.Printf("save user message: %v", err)
	}
	if err := uc.chatRepo.SaveMessage(ctx, assistantMsg); err != nil {
		log.Printf("save assistant message: %v", err)
	}

	if err := uc.cache.SetAdvice(ctx, cacheKey, reply); err != nil {
		log.Printf("cache set advice: %v", err)
	}

	_ = uc.publisher.Publish(ctx, "ai.advice.requested", map[string]string{
		"session_id": sessionID,
		"user_id":    userID,
	})

	return reply, sessionID, nil
}

func (uc *AIUseCase) RecommendSneakers(ctx context.Context, userID, preferences string, budget float64, size string) ([]*model.SneakerRecommendation, error) {
	products, err := uc.catalog.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	catalogText := buildCatalogText(products, budget, size)

	prompt := fmt.Sprintf(
		"User wants: %s. Budget: $%.0f. Size: %s.\n\nAvailable products:\n%s\n\n"+
			"Recommend up to 5 best matches. For each, respond with JSON array: "+
			`[{"product_id":"...","name":"...","brand":"...","price":0,"image_url":"...","reason":"...","match_score":0.95}]`+
			" Only output valid JSON, nothing else.",
		preferences, budget, size, catalogText,
	)

	raw, err := uc.gemini.Chat(ctx, sneakerSystemPrompt, "", prompt)
	if err != nil {
		return nil, fmt.Errorf("gemini recommend: %w", err)
	}

	recs, err := parseRecommendations(raw, products)
	if err != nil {
		log.Printf("parse recommendations: %v — falling back to scored list", err)
		recs = scoreProducts(products, preferences, budget, size)
	}

	_ = uc.publisher.Publish(ctx, "ai.recommend.requested", map[string]string{
		"user_id":     userID,
		"preferences": preferences,
	})

	return recs, nil
}

func (uc *AIUseCase) SearchSneakersByStyle(ctx context.Context, style, size string, maxPrice float64) ([]*model.SneakerRecommendation, string, error) {
	products, err := uc.catalog.ListProducts(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("list products: %w", err)
	}

	catalogText := buildCatalogText(products, maxPrice, size)

	prompt := fmt.Sprintf(
		"Style description: %s. Size filter: %s. Max price: $%.0f.\n\nAvailable products:\n%s\n\n"+
			"Find the best style matches. Respond with JSON: "+
			`{"results":[{"product_id":"...","name":"...","brand":"...","price":0,"image_url":"...","reason":"...","match_score":0.9}],"summary":"..."}` +
			" Only output valid JSON.",
		style, size, maxPrice, catalogText,
	)

	raw, err := uc.gemini.Chat(ctx, sneakerSystemPrompt, "", prompt)
	if err != nil {
		return nil, "", fmt.Errorf("gemini style search: %w", err)
	}

	results, summary := parseStyleSearch(raw, products)
	return results, summary, nil
}

func (uc *AIUseCase) GetTrendingSneakers(ctx context.Context, limit int32) ([]*model.TrendingSneaker, error) {
	if cached, err := uc.cache.GetTrending(ctx); err == nil && len(cached) > 0 {
		if int32(len(cached)) > limit && limit > 0 {
			return cached[:limit], nil
		}
		return cached, nil
	}

	products, err := uc.catalog.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	if len(products) == 0 {
		return []*model.TrendingSneaker{}, nil
	}

	catalogText := buildCatalogText(products, 0, "")
	prompt := fmt.Sprintf(
		"Based on these sneakers, identify the top trending ones based on brand hype, model popularity, and current market trends:\n%s\n\n"+
			"Respond with JSON array (top 10): "+
			`[{"product_id":"...","name":"...","brand":"...","price":0,"image_url":"...","trend_reason":"...","trend_score":95}]`+
			" Only output valid JSON.",
		catalogText,
	)

	raw, err := uc.gemini.Chat(ctx, sneakerSystemPrompt, "", prompt)
	if err != nil {
		log.Printf("gemini trending: %v — using fallback", err)
		trending := fallbackTrending(products, limit)
		_ = uc.cache.SetTrending(ctx, trending)
		return trending, nil
	}

	trending := parseTrending(raw, products)
	if int32(len(trending)) > limit && limit > 0 {
		trending = trending[:limit]
	}

	_ = uc.cache.SetTrending(ctx, trending)
	return trending, nil
}

func buildHistoryText(msgs []*model.ChatMessage) string {
	if len(msgs) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, m := range msgs {
		sb.WriteString(m.Role)
		sb.WriteString(": ")
		sb.WriteString(m.Content)
		sb.WriteString("\n")
	}
	return sb.String()
}

func buildCatalogText(products []*model.ProductInfo, maxPrice float64, size string) string {
	var sb strings.Builder
	for _, p := range products {
		if maxPrice > 0 && p.Price > maxPrice {
			continue
		}
		if size != "" {
			found := false
			for _, s := range p.Sizes {
				if s == size {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		fmt.Fprintf(&sb, "ID:%s | %s %s | $%.0f | Stock:%d | Sizes:%s\n",
			p.ID, p.Brand, p.Name, p.Price, p.Stock, strings.Join(p.Sizes, ","))
	}
	return sb.String()
}

func cacheHash(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
	}
	return fmt.Sprintf("%x", h.Sum(nil))[:16]
}

func parseRecommendations(raw string, products []*model.ProductInfo) ([]*model.SneakerRecommendation, error) {
	raw = strings.TrimSpace(raw)
	start := strings.Index(raw, "[")
	end := strings.LastIndex(raw, "]")
	if start < 0 || end < 0 || end <= start {
		return nil, fmt.Errorf("no JSON array found")
	}
	raw = raw[start : end+1]

	var recs []*model.SneakerRecommendation
	if err := json.Unmarshal([]byte(raw), &recs); err != nil {
		return nil, err
	}

	productMap := make(map[string]*model.ProductInfo, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}
	for _, r := range recs {
		if p, ok := productMap[r.ProductID]; ok {
			if r.ImageURL == "" {
				r.ImageURL = p.ImageURL
			}
			if r.Price == 0 {
				r.Price = p.Price
			}
		}
	}
	return recs, nil
}

func parseStyleSearch(raw string, products []*model.ProductInfo) ([]*model.SneakerRecommendation, string) {
	raw = strings.TrimSpace(raw)
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start < 0 || end < 0 {
		recs := scoreProducts(products, raw, 0, "")
		return recs, "Here are the best style matches from our catalog."
	}

	var out struct {
		Results []*model.SneakerRecommendation `json:"results"`
		Summary string                          `json:"summary"`
	}
	if err := json.Unmarshal([]byte(raw[start:end+1]), &out); err != nil {
		recs := scoreProducts(products, raw, 0, "")
		return recs, "Here are the best style matches from our catalog."
	}

	productMap := make(map[string]*model.ProductInfo, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}
	for _, r := range out.Results {
		if p, ok := productMap[r.ProductID]; ok {
			if r.ImageURL == "" {
				r.ImageURL = p.ImageURL
			}
		}
	}

	return out.Results, out.Summary
}

func parseTrending(raw string, products []*model.ProductInfo) []*model.TrendingSneaker {
	raw = strings.TrimSpace(raw)
	start := strings.Index(raw, "[")
	end := strings.LastIndex(raw, "]")
	if start < 0 || end < 0 {
		return fallbackTrending(products, 10)
	}

	var trending []*model.TrendingSneaker
	if err := json.Unmarshal([]byte(raw[start:end+1]), &trending); err != nil {
		return fallbackTrending(products, 10)
	}

	productMap := make(map[string]*model.ProductInfo, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}
	for _, t := range trending {
		if p, ok := productMap[t.ProductID]; ok {
			if t.ImageURL == "" {
				t.ImageURL = p.ImageURL
			}
			if t.Price == 0 {
				t.Price = p.Price
			}
		}
	}
	return trending
}

func scoreProducts(products []*model.ProductInfo, preferences string, budget float64, size string) []*model.SneakerRecommendation {
	type scored struct {
		p     *model.ProductInfo
		score float64
	}

	prefLower := strings.ToLower(preferences)
	var items []scored

	for _, p := range products {
		if budget > 0 && p.Price > budget {
			continue
		}
		if size != "" {
			found := false
			for _, s := range p.Sizes {
				if s == size {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		score := 0.5
		nameLower := strings.ToLower(p.Name + " " + p.Brand)
		words := strings.Fields(prefLower)
		for _, w := range words {
			if strings.Contains(nameLower, w) {
				score += 0.1
			}
		}
		items = append(items, scored{p, score})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].score > items[j].score })

	recs := make([]*model.SneakerRecommendation, 0, 5)
	for i, it := range items {
		if i >= 5 {
			break
		}
		recs = append(recs, &model.SneakerRecommendation{
			ProductID:  it.p.ID,
			Name:       it.p.Name,
			Brand:      it.p.Brand,
			Price:      it.p.Price,
			ImageURL:   it.p.ImageURL,
			Reason:     "Matches your style preferences",
			MatchScore: it.score,
		})
	}
	return recs
}

func fallbackTrending(products []*model.ProductInfo, limit int32) []*model.TrendingSneaker {
	hypeOrder := []string{"Jordan", "Nike", "Adidas", "New Balance", "Asics", "Salomon", "Puma", "Reebok"}
	rankMap := make(map[string]int, len(hypeOrder))
	for i, b := range hypeOrder {
		rankMap[strings.ToLower(b)] = len(hypeOrder) - i
	}

	type scored struct {
		p    *model.ProductInfo
		rank int
	}
	var items []scored
	for _, p := range products {
		rank := rankMap[strings.ToLower(p.Brand)]
		items = append(items, scored{p, rank})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].rank > items[j].rank })

	result := make([]*model.TrendingSneaker, 0, limit)
	for i, it := range items {
		if limit > 0 && int32(i) >= limit {
			break
		}
		result = append(result, &model.TrendingSneaker{
			ProductID:   it.p.ID,
			Name:        it.p.Name,
			Brand:       it.p.Brand,
			Price:       it.p.Price,
			ImageURL:    it.p.ImageURL,
			TrendReason: fmt.Sprintf("%s consistently leads sneaker culture — high demand, strong resale value.", it.p.Brand),
			TrendScore:  int32(50 + it.rank*5),
		})
	}
	return result
}
