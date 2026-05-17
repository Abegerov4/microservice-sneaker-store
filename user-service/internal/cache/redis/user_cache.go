package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"sneaker-store/user-service/internal/model"
)

const ttl = 60 * time.Second

type UserCache struct {
	client *redis.Client
}

func NewUserCache(client *redis.Client) *UserCache {
	return &UserCache{client: client}
}

func (c *UserCache) GetByID(ctx context.Context, id string) (*model.User, error) {
	data, err := c.client.Get(ctx, fmt.Sprintf("user:%s", id)).Bytes()
	if err != nil {
		return nil, err
	}
	u := &model.User{}
	if err := json.Unmarshal(data, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (c *UserCache) SetByID(ctx context.Context, u *model.User) error {
	safe := *u
	safe.PasswordHash = ""
	data, err := json.Marshal(safe)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, fmt.Sprintf("user:%s", u.ID), data, ttl).Err()
}

func (c *UserCache) DeleteByID(ctx context.Context, id string) error {
	return c.client.Del(ctx, fmt.Sprintf("user:%s", id)).Err()
}
