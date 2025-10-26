package usecase

import (
	"context"

	"sneaker-store/product-service/internal/model"
)

type ProductRepository interface {
	Create(ctx context.Context, p *model.Product) error
	GetByID(ctx context.Context, id string) (*model.Product, error)
	List(ctx context.Context) ([]*model.Product, error)
	Update(ctx context.Context, p *model.Product) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, brand string, minPrice, maxPrice float64, size string) ([]*model.Product, error)
	UpdateStock(ctx context.Context, id string, delta int) (int, error)
	GetByBrand(ctx context.Context, brand string) ([]*model.Product, error)
	GetLowStock(ctx context.Context, threshold int) ([]*model.Product, error)
	GetBrands(ctx context.Context) ([]string, error)
	GetStats(ctx context.Context) (*model.ProductStats, error)
	BulkDelete(ctx context.Context, ids []string) (int, error)
}

type ProductCache interface {
	GetByID(ctx context.Context, id string) (*model.Product, error)
	SetByID(ctx context.Context, p *model.Product) error
	DeleteByID(ctx context.Context, id string) error
	GetList(ctx context.Context) ([]*model.Product, error)
	SetList(ctx context.Context, products []*model.Product) error
	DeleteList(ctx context.Context) error
}

type EventPublisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}
