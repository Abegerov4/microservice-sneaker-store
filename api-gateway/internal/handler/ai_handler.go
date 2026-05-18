package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	aipb "sneaker-store/ai-service/proto"
)

type AIHandler struct {
	client aipb.AIServiceClient
}

func NewAIHandler(client aipb.AIServiceClient) *AIHandler {
	return &AIHandler{client: client}
}

func (h *AIHandler) Chat(c *gin.Context) {
	var body struct {
		SessionID string `json:"session_id"`
		UserID    string `json:"user_id"`
		Message   string `json:"message"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message is required"})
		return
	}

	resp, err := h.client.AskSneakerAdvice(c.Request.Context(), &aipb.AskSneakerAdviceRequest{
		SessionId: body.SessionID,
		UserId:    body.UserID,
		Message:   body.Message,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"reply":      resp.Reply,
		"session_id": resp.SessionId,
	})
}

func (h *AIHandler) Recommend(c *gin.Context) {
	var body struct {
		UserID      string  `json:"user_id"`
		Preferences string  `json:"preferences"`
		Budget      float64 `json:"budget"`
		Size        string  `json:"size"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Preferences == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "preferences are required"})
		return
	}

	resp, err := h.client.RecommendSneakers(c.Request.Context(), &aipb.RecommendSneakersRequest{
		UserId:      body.UserID,
		Preferences: body.Preferences,
		Budget:      body.Budget,
		Size:        body.Size,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"recommendations": resp.Recommendations})
}

func (h *AIHandler) SearchByStyle(c *gin.Context) {
	var body struct {
		StyleDescription string  `json:"style_description"`
		Size             string  `json:"size"`
		MaxPrice         float64 `json:"max_price"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.SearchSneakersByStyle(c.Request.Context(), &aipb.SearchSneakersByStyleRequest{
		StyleDescription: body.StyleDescription,
		Size:             body.Size,
		MaxPrice:         body.MaxPrice,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"results":    resp.Results,
		"ai_summary": resp.AiSummary,
	})
}

func (h *AIHandler) Trending(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	if limit <= 0 {
		limit = 10
	}

	resp, err := h.client.GetTrendingSneakers(c.Request.Context(), &aipb.GetTrendingSneakersRequest{
		Limit: int32(limit),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sneakers": resp.Sneakers})
}
