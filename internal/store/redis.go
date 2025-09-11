package store

import (
	"UrbanNest/internal/entities"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore(addr, password string) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisStore{Client: client}
}

func (s *RedisStore) CacheListing(ctx context.Context, listing *entities.Listing) error {
	data, err := json.Marshal(listing)
	if err != nil {
		return err
	}
	return s.Client.Set(ctx, fmt.Sprintf("listing:%d", listing.ID), data, 24*time.Hour).Err()
}

func (s *RedisStore) GetListing(ctx context.Context, key string) (*entities.Listing, error) {
	data, err := s.Client.Get(ctx, fmt.Sprintf("listing:%s", key)).Bytes()
	if err != nil {
		return nil, err
	}
	var listing entities.Listing
	if err := json.Unmarshal(data, &listing); err != nil {
		return nil, err
	}
	return &listing, nil
}

func (s *RedisStore) CacheReview(ctx context.Context, review *entities.Review) error {
	data, err := json.Marshal(review)
	if err != nil {
		return err
	}
	return s.Client.Set(ctx, fmt.Sprintf("review:%d", review.ID), data, 24*time.Hour).Err()
}

func (s *RedisStore) GetReview(ctx context.Context, key string) (*entities.Review, error) {
	data, err := s.Client.Get(ctx, fmt.Sprintf("review:%s", key)).Bytes()
	if err != nil {
		return nil, err
	}
	var review entities.Review
	if err := json.Unmarshal(data, &review); err != nil {
		return nil, err
	}
	return &review, nil
}

func (s *RedisStore) CacheReviewsByListing(ctx context.Context, listingID uint, reviews []entities.Review) error {
	data, err := json.Marshal(reviews)
	if err != nil {
		return err
	}
	return s.Client.Set(ctx, fmt.Sprintf("listing:%d:reviews", listingID), data, 24*time.Hour).Err()
}

func (s *RedisStore) GetReviewsByListing(ctx context.Context, listingID uint) ([]entities.Review, error) {
	data, err := s.Client.Get(ctx, fmt.Sprintf("listing:%d:reviews", listingID)).Bytes()
	if err != nil {
		return nil, err
	}
	var reviews []entities.Review
	if err := json.Unmarshal(data, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}

func (s *RedisStore) CacheMessage(ctx context.Context, message *entities.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return s.Client.Set(ctx, fmt.Sprintf("message:%d", message.ID), data, 24*time.Hour).Err()
}

func (s *RedisStore) GetMessage(ctx context.Context, key string) (*entities.Message, error) {
	data, err := s.Client.Get(ctx, fmt.Sprintf("message:%s", key)).Bytes()
	if err != nil {
		return nil, err
	}
	var message entities.Message
	if err := json.Unmarshal(data, &message); err != nil {
		return nil, err
	}
	return &message, nil
}

func (s *RedisStore) CacheMessagesByUser(ctx context.Context, userID uint, messages []entities.Message) error {
	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	return s.Client.Set(ctx, fmt.Sprintf("user:%d:messages", userID), data, 24*time.Hour).Err()
}

func (s *RedisStore) GetMessagesByUser(ctx context.Context, userID uint) ([]entities.Message, error) {
	data, err := s.Client.Get(ctx, fmt.Sprintf("user:%d:messages", userID)).Bytes()
	if err != nil {
		return nil, err
	}
	var messages []entities.Message
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *RedisStore) CacheBooking(ctx context.Context, booking *entities.Booking) error {
	data, err := json.Marshal(booking)
	if err != nil {
		return err
	}
	return s.Client.Set(ctx, fmt.Sprintf("booking:%d", booking.ID), data, 24*time.Hour).Err()
}

func (s *RedisStore) GetBooking(ctx context.Context, key string) (*entities.Booking, error) {
	data, err := s.Client.Get(ctx, fmt.Sprintf("booking:%s", key)).Bytes()
	if err != nil {
		return nil, err
	}
	var booking entities.Booking
	if err := json.Unmarshal(data, &booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *RedisStore) CacheBookingsByUser(ctx context.Context, userID uint, bookings []entities.Booking) error {
	data, err := json.Marshal(bookings)
	if err != nil {
		return err
	}
	return s.Client.Set(ctx, fmt.Sprintf("user:%d:bookings", userID), data, 24*time.Hour).Err()
}

func (s *RedisStore) GetBookingsByUser(ctx context.Context, userID uint) ([]entities.Booking, error) {
	data, err := s.Client.Get(ctx, fmt.Sprintf("user:%d:bookings", userID)).Bytes()
	if err != nil {
		return nil, err
	}
	var bookings []entities.Booking
	if err := json.Unmarshal(data, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *RedisStore) CacheBookingsByHost(ctx context.Context, hostID uint, bookings []entities.Booking) error {
	data, err := json.Marshal(bookings)
	if err != nil {
		return err
	}
	return s.Client.Set(ctx, fmt.Sprintf("host:%d:bookings", hostID), data, 24*time.Hour).Err()
}

func (s *RedisStore) GetBookingsByHost(ctx context.Context, hostID uint) ([]entities.Booking, error) {
	data, err := s.Client.Get(ctx, fmt.Sprintf("host:%d:bookings", hostID)).Bytes()
	if err != nil {
		return nil, err
	}
	var bookings []entities.Booking
	if err := json.Unmarshal(data, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}
