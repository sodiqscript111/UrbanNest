package store

import (
	"UrbanNest/internal/entities"
	"UrbanNest/pkg/redis"
	"context"
	"encoding/json"
	"time"
)

type RedisStore struct {
	Client *redis.RedisClient
}

func NewRedisStore(addr, password string) *RedisStore {
	return &RedisStore{
		Client: redis.NewRedisClient(addr, password),
	}
}

func (r *RedisStore) CacheListing(ctx context.Context, listing *entities.Listing) error {
	data, err := json.Marshal(listing)
	if err != nil {
		return err
	}
	return r.Client.Set(ctx, "listing:"+string(listing.ID), data, time.Hour)
}

func (r *RedisStore) GetListing(ctx context.Context, id string) (*entities.Listing, error) {
	data, err := r.Client.Get(ctx, "listing:"+id)
	if err != nil {
		return nil, err
	}
	var listing entities.Listing
	if err := json.Unmarshal([]byte(data), &listing); err != nil {
		return nil, err
	}
	return &listing, nil
}
