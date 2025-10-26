package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"sneaker-store/product-service/internal/model"
)

const ttl = 60 * time.Second

type ProductCache struct {
	client *redis.Client
}

func NewProductCache(client *redis.Client) *ProductCache {
	return &ProductCache{client: client}
}

func (c *ProductCache) GetByID(ctx context.Context, id string) (*model.Product, error) {
	key := fmt.Sprintf("product:%s", id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	p := &model.Product{}
	if err := json.Unmarshal(data, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (c *ProductCache) SetByID(ctx context.Context, p *model.Product) error {
	key := fmt.Sprintf("product:%s", p.ID)
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *ProductCache) DeleteByID(ctx context.Context, id string) error {
	return c.client.Del(ctx, fmt.Sprintf("product:%s", id)).Err()
}

func (c *ProductCache) GetList(ctx context.Context) ([]*model.Product, error) {
	data, err := c.client.Get(ctx, "products:list").Bytes()
	if err != nil {
		return nil, err
	}
	var products []*model.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (c *ProductCache) SetList(ctx context.Context, products []*model.Product) error {
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, "products:list", data, ttl).Err()
}

func (c *ProductCache) DeleteList(ctx context.Context) error {
	return c.client.Del(ctx, "products:list").Err()
}
