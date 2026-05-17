package usecase

import (
	"context"
	"time"

	"sneaker-store/order-service/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, o *model.Order) error
	GetByID(ctx context.Context, id string) (*model.Order, error)
	List(ctx context.Context) ([]*model.Order, error)
	UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error
	GetByUserID(ctx context.Context, userID string) ([]*model.Order, error)
	GetByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error)
	GetStats(ctx context.Context) (*model.OrderStats, error)
	GetTotalRevenue(ctx context.Context) (float64, error)
	GetByDateRange(ctx context.Context, from, to time.Time) ([]*model.Order, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
}

type OrderCache interface {
	GetByID(ctx context.Context, id string) (*model.Order, error)
	SetByID(ctx context.Context, o *model.Order) error
	DeleteByID(ctx context.Context, id string) error
	GetByUser(ctx context.Context, userID string) ([]*model.Order, error)
	SetByUser(ctx context.Context, userID string, orders []*model.Order) error
	DeleteByUser(ctx context.Context, userID string) error
}

type EventPublisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}

type ProductClient interface {
	GetProduct(ctx context.Context, id string) (*ProductInfo, error)
}

type ProductInfo struct {
	ID    string
	Name  string
	Price float64
	Stock int
}
