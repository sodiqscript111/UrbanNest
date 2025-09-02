package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr, password string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisClient{Client: client}
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}
