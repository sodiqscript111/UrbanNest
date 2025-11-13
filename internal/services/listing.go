package services

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/interfaces"
	"UrbanNest/internal/store"
	"UrbanNest/pkg/kafka"
	"context"
	"fmt"
	"time"
)

type ListingService struct {
	db       interfaces.Database
	redis    *store.RedisStore
	producer *kafka.Producer
}

func NewListingService(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) *ListingService {
	return &ListingService{db, redis, producer}
}

func (s *ListingService) CreateListing(ctx context.Context, listing *entities.Listing) error {
	err := s.db.CreateListing(ctx, listing)
	if err != nil {
		return err
	}

	if s.redis != nil {
		if err := s.redis.CacheListing(ctx, listing); err != nil {
			return err
		}
	}

	return s.producer.PublishMessage(ctx, "listing.created", listing)
}

func (s *ListingService) GetListing(ctx context.Context, id uint) (*entities.Listing, error) {

	if s.redis != nil {
		listing, err := s.redis.GetListing(ctx, fmt.Sprintf("%d", id))
		if err == nil {
			return listing, nil
		}
	}

	dbListing, err := s.db.GetListing(ctx, id)
	if err != nil {
		return nil, err

	}

	if s.redis != nil {
		if err := s.redis.CacheListing(ctx, dbListing); err != nil {
			return nil, err
		}
	}

	return dbListing, nil
}

func (s *ListingService) UpdateListing(ctx context.Context, id uint, listing *entities.Listing) error {

	existing, err := s.db.GetListing(ctx, id)
	if err != nil {
		return err
	}

	existing.Title = listing.Title
	existing.Description = listing.Description
	existing.Location = listing.Location
	existing.Price = listing.Price
	existing.Available = listing.Available

	if err := s.db.UpdateListing(ctx, existing); err != nil {
		return err
	}

	if s.redis != nil {
		if err := s.redis.CacheListing(ctx, existing); err != nil {
			return err
		}
	}

	return s.producer.PublishMessage(ctx, "listing.updated", existing)
}

func (s *ListingService) DeleteListing(ctx context.Context, id uint) error {

	_, err := s.db.GetListing(ctx, id)
	if err != nil {
		return err
	}

	if err := s.db.DeleteListing(ctx, id); err != nil {
		return err
	}

	if s.redis != nil {
		if err := s.redis.Client.Del(ctx, fmt.Sprintf("listing:%d", id)).Err(); err != nil {
			return err
		}
	}

	return s.producer.PublishMessage(ctx, "listing.deleted", map[string]uint{"id": id})
}

func (s *ListingService) CheckAvailability(ctx context.Context, listingID uint, startDate, endDate time.Time) (bool, error) {

	if startDate.After(endDate) || startDate.Before(time.Now()) {
		return false, fmt.Errorf("invalid date range")
	}

	startStr := startDate.Format("2006-01-02 15:04:05")
	endStr := endDate.Format("2006-01-02 15:04:05")

	available, err := s.db.CheckAvailability(ctx, listingID, startStr, endStr)
	if err != nil {
		return false, err
	}

	return available, nil
}
