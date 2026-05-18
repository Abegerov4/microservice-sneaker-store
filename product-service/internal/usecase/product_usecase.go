package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"sneaker-store/product-service/internal/model"
)

type ProductUseCase struct {
	repo      ProductRepository
	cache     ProductCache
	publisher EventPublisher
}

func NewProductUseCase(repo ProductRepository, cache ProductCache, pub EventPublisher) *ProductUseCase {
	return &ProductUseCase{repo: repo, cache: cache, publisher: pub}
}

func (uc *ProductUseCase) Create(ctx context.Context, p *model.Product) (*model.Product, error) {
	if p.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if p.Price < 0 {
		return nil, fmt.Errorf("price must be non-negative")
	}

	p.ID = uuid.NewString()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	if err := uc.cache.DeleteList(ctx); err != nil {
		log.Printf("cache invalidation failed: %v", err)
	}

	go func() {
		if err := uc.publisher.Publish(context.Background(), "products.created", p); err != nil {
			log.Printf("publish products.created failed: %v", err)
		}
	}()

	return p, nil
}

func (uc *ProductUseCase) GetByID(ctx context.Context, id string) (*model.Product, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	if p, err := uc.cache.GetByID(ctx, id); err == nil && p != nil {
		return p, nil
	}

	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := uc.cache.SetByID(ctx, p); err != nil {
		log.Printf("cache set failed: %v", err)
	}

	return p, nil
}

func (uc *ProductUseCase) List(ctx context.Context) ([]*model.Product, error) {
	if products, err := uc.cache.GetList(ctx); err == nil && products != nil {
		return products, nil
	}

	products, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	if err := uc.cache.SetList(ctx, products); err != nil {
		log.Printf("cache set list failed: %v", err)
	}

	return products, nil
}

func (uc *ProductUseCase) Update(ctx context.Context, p *model.Product) (*model.Product, error) {
	if p.ID == "" {
		return nil, fmt.Errorf("id is required")
	}

	existing, err := uc.repo.GetByID(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	if p.Name != "" {
		existing.Name = p.Name
	}
	if p.Brand != "" {
		existing.Brand = p.Brand
	}
	if p.Description != "" {
		existing.Description = p.Description
	}
	if p.Price > 0 {
		existing.Price = p.Price
	}
	if len(p.Sizes) > 0 {
		existing.Sizes = p.Sizes
	}
	if p.Stock >= 0 {
		existing.Stock = p.Stock
	}
	if p.ImageURL != "" {
		existing.ImageURL = p.ImageURL
	}
	existing.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	if err := uc.cache.DeleteByID(ctx, existing.ID); err != nil {
		log.Printf("cache delete failed: %v", err)
	}
	if err := uc.cache.DeleteList(ctx); err != nil {
		log.Printf("cache delete list failed: %v", err)
	}

	go func() {
		if err := uc.publisher.Publish(context.Background(), "products.updated", existing); err != nil {
			log.Printf("publish products.updated failed: %v", err)
		}
	}()

	return existing, nil
}

func (uc *ProductUseCase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}

	if err := uc.cache.DeleteByID(ctx, id); err != nil {
		log.Printf("cache delete failed: %v", err)
	}
	if err := uc.cache.DeleteList(ctx); err != nil {
		log.Printf("cache delete list failed: %v", err)
	}

	return nil
}

func (uc *ProductUseCase) Search(ctx context.Context, brand string, minPrice, maxPrice float64, size string) ([]*model.Product, error) {
	return uc.repo.Search(ctx, brand, minPrice, maxPrice, size)
}

func (uc *ProductUseCase) UpdateStock(ctx context.Context, id string, delta int) (int, error) {
	if id == "" {
		return 0, fmt.Errorf("id is required")
	}
	newStock, err := uc.repo.UpdateStock(ctx, id, delta)
	if err != nil {
		return 0, err
	}
	if err := uc.cache.DeleteByID(ctx, id); err != nil {
		log.Printf("cache delete failed: %v", err)
	}
	if err := uc.cache.DeleteList(ctx); err != nil {
		log.Printf("cache delete list failed: %v", err)
	}
	return newStock, nil
}

func (uc *ProductUseCase) GetByBrand(ctx context.Context, brand string) ([]*model.Product, error) {
	if brand == "" {
		return nil, fmt.Errorf("brand is required")
	}
	return uc.repo.GetByBrand(ctx, brand)
}

func (uc *ProductUseCase) GetLowStock(ctx context.Context, threshold int) ([]*model.Product, error) {
	return uc.repo.GetLowStock(ctx, threshold)
}

func (uc *ProductUseCase) GetBrands(ctx context.Context) ([]string, error) {
	return uc.repo.GetBrands(ctx)
}

func (uc *ProductUseCase) GetStats(ctx context.Context) (*model.ProductStats, error) {
	return uc.repo.GetStats(ctx)
}

func (uc *ProductUseCase) BulkDelete(ctx context.Context, ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, fmt.Errorf("at least one id is required")
	}
	count, err := uc.repo.BulkDelete(ctx, ids)
	if err != nil {
		return 0, err
	}
	for _, id := range ids {
		if err := uc.cache.DeleteByID(ctx, id); err != nil {
			log.Printf("cache delete failed for %s: %v", id, err)
		}
	}
	if err := uc.cache.DeleteList(ctx); err != nil {
		log.Printf("cache delete list failed: %v", err)
	}
	return count, nil
}
