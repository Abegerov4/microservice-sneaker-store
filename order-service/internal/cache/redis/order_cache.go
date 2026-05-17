package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"sneaker-store/order-service/internal/model"
)

const ttl = 60 * time.Second

type OrderCache struct {
	client *redis.Client
}

func NewOrderCache(client *redis.Client) *OrderCache {
	return &OrderCache{client: client}
}

func (c *OrderCache) GetByID(ctx context.Context, id string) (*model.Order, error) {
	data, err := c.client.Get(ctx, fmt.Sprintf("order:%s", id)).Bytes()
	if err != nil {
		return nil, err
	}
	o := &model.Order{}
	if err := json.Unmarshal(data, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (c *OrderCache) SetByID(ctx context.Context, o *model.Order) error {
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, fmt.Sprintf("order:%s", o.ID), data, ttl).Err()
}

func (c *OrderCache) DeleteByID(ctx context.Context, id string) error {
	return c.client.Del(ctx, fmt.Sprintf("order:%s", id)).Err()
}

func (c *OrderCache) GetByUser(ctx context.Context, userID string) ([]*model.Order, error) {
	data, err := c.client.Get(ctx, fmt.Sprintf("orders:user:%s", userID)).Bytes()
	if err != nil {
		return nil, err
	}
	var orders []*model.Order
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (c *OrderCache) SetByUser(ctx context.Context, userID string, orders []*model.Order) error {
	data, err := json.Marshal(orders)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, fmt.Sprintf("orders:user:%s", userID), data, ttl).Err()
}

func (c *OrderCache) DeleteByUser(ctx context.Context, userID string) error {
	return c.client.Del(ctx, fmt.Sprintf("orders:user:%s", userID)).Err()
}
