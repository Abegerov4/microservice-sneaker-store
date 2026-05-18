package usecase

import (
	"context"

	"sneaker-store/ai-service/internal/model"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, msg *model.ChatMessage) error
	GetHistory(ctx context.Context, sessionID string, limit int) ([]*model.ChatMessage, error)
}

type AICache interface {
	GetAdvice(ctx context.Context, key string) (string, error)
	SetAdvice(ctx context.Context, key string, reply string) error
	GetTrending(ctx context.Context) ([]*model.TrendingSneaker, error)
	SetTrending(ctx context.Context, sneakers []*model.TrendingSneaker) error
}

type EventPublisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}

type ProductCatalog interface {
	ListProducts(ctx context.Context) ([]*model.ProductInfo, error)
}

type GeminiClient interface {
	Chat(ctx context.Context, systemPrompt, history, userMessage string) (string, error)
}
