package usecase_test

import (
	"context"
	"errors"
	"testing"

	"sneaker-store/ai-service/internal/model"
	"sneaker-store/ai-service/internal/usecase"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type mockChatRepo struct {
	msgs []*model.ChatMessage
}

func (m *mockChatRepo) SaveMessage(_ context.Context, msg *model.ChatMessage) error {
	m.msgs = append(m.msgs, msg)
	return nil
}

func (m *mockChatRepo) GetHistory(_ context.Context, _ string, _ int) ([]*model.ChatMessage, error) {
	return m.msgs, nil
}

type mockCache struct{}

func (m *mockCache) GetAdvice(_ context.Context, _ string) (string, error) {
	return "", errors.New("miss")
}
func (m *mockCache) SetAdvice(_ context.Context, _ string, _ string) error { return nil }
func (m *mockCache) GetTrending(_ context.Context) ([]*model.TrendingSneaker, error) {
	return nil, errors.New("miss")
}
func (m *mockCache) SetTrending(_ context.Context, _ []*model.TrendingSneaker) error { return nil }

type mockPublisher struct{}

func (m *mockPublisher) Publish(_ context.Context, _ string, _ interface{}) error { return nil }

type mockCatalog struct {
	products []*model.ProductInfo
}

func (m *mockCatalog) ListProducts(_ context.Context) ([]*model.ProductInfo, error) {
	return m.products, nil
}

type mockGemini struct {
	reply string
	err   error
}

func (m *mockGemini) Chat(_ context.Context, _, _, _ string) (string, error) {
	return m.reply, m.err
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func TestAskSneakerAdvice_Success(t *testing.T) {
	repo := &mockChatRepo{}
	cache := &mockCache{}
	pub := &mockPublisher{}
	catalog := &mockCatalog{}
	gemini := &mockGemini{reply: "Great choice! The Air Max 90 is iconic."}

	uc := usecase.NewAIUseCase(repo, cache, pub, catalog, gemini)

	reply, sessionID, err := uc.AskSneakerAdvice(context.Background(), "", "", "What Nike shoes should I get?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != gemini.reply {
		t.Errorf("got reply %q, want %q", reply, gemini.reply)
	}
	if sessionID == "" {
		t.Error("expected a session ID to be generated")
	}
	if len(repo.msgs) != 2 {
		t.Errorf("expected 2 saved messages (user+assistant), got %d", len(repo.msgs))
	}
}

func TestAskSneakerAdvice_GeminiError(t *testing.T) {
	uc := usecase.NewAIUseCase(
		&mockChatRepo{},
		&mockCache{},
		&mockPublisher{},
		&mockCatalog{},
		&mockGemini{err: errors.New("api error")},
	)

	_, _, err := uc.AskSneakerAdvice(context.Background(), "", "", "test")
	if err == nil {
		t.Error("expected error from gemini, got nil")
	}
}

func TestGetTrendingSneakers_Fallback(t *testing.T) {
	products := []*model.ProductInfo{
		{ID: "1", Name: "Air Max 90", Brand: "Nike", Price: 130},
		{ID: "2", Name: "Yeezy 350", Brand: "Adidas", Price: 220},
		{ID: "3", Name: "Jordan 1", Brand: "Jordan", Price: 180},
	}
	uc := usecase.NewAIUseCase(
		&mockChatRepo{},
		&mockCache{},
		&mockPublisher{},
		&mockCatalog{products: products},
		&mockGemini{err: errors.New("api error")},
	)

	trending, err := uc.GetTrendingSneakers(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trending) == 0 {
		t.Error("expected trending sneakers via fallback, got none")
	}
}

func TestRecommendSneakers_ScoreFallback(t *testing.T) {
	products := []*model.ProductInfo{
		{ID: "1", Name: "Air Max 90", Brand: "Nike", Price: 130, Sizes: []string{"42", "43"}},
		{ID: "2", Name: "Superstar", Brand: "Adidas", Price: 90, Sizes: []string{"42"}},
	}
	uc := usecase.NewAIUseCase(
		&mockChatRepo{},
		&mockCache{},
		&mockPublisher{},
		&mockCatalog{products: products},
		&mockGemini{reply: "not json at all"},
	)

	recs, err := uc.RecommendSneakers(context.Background(), "", "Nike running", 200, "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) == 0 {
		t.Error("expected fallback recommendations, got none")
	}
}
