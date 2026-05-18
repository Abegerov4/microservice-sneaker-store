package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"sneaker-store/ai-service/internal/model"
)

const (
	adviceTTL  = 5 * time.Minute
	trendingTTL = 30 * time.Minute
)

type AICache struct {
	client *redis.Client
}

func NewAICache(client *redis.Client) *AICache {
	return &AICache{client: client}
}

func (c *AICache) GetAdvice(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, fmt.Sprintf("ai:advice:%s", key)).Result()
}

func (c *AICache) SetAdvice(ctx context.Context, key string, reply string) error {
	return c.client.Set(ctx, fmt.Sprintf("ai:advice:%s", key), reply, adviceTTL).Err()
}

func (c *AICache) GetTrending(ctx context.Context) ([]*model.TrendingSneaker, error) {
	data, err := c.client.Get(ctx, "ai:trending").Bytes()
	if err != nil {
		return nil, err
	}
	var sneakers []*model.TrendingSneaker
	if err := json.Unmarshal(data, &sneakers); err != nil {
		return nil, err
	}
	return sneakers, nil
}

func (c *AICache) SetTrending(ctx context.Context, sneakers []*model.TrendingSneaker) error {
	data, err := json.Marshal(sneakers)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, "ai:trending", data, trendingTTL).Err()
}
