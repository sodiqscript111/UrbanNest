package services

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/store"
	"UrbanNest/pkg/kafka"
	"context"
	"fmt"
	"time"
)

type ListingService struct {
	db       *store.PostgresStore
	redis    *store.RedisStore
	producer *kafka.Producer
}

func NewListingService(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) *ListingService {
	return &ListingService{db, redis, producer}
}

func (s *ListingService) CreateListing(ctx context.Context, listing *entities.Listing) error {
	tx := s.db.DB.Create(listing)
	if err := tx.Error; err != nil {
		return err
	}

	// Cache listing in Redis
	if s.redis != nil {
		if err := s.redis.CacheListing(ctx, listing); err != nil {
			return err
		}
	}

	// Publish listing creation event
	return s.producer.PublishMessage(ctx, "listing.created", listing)
}

func (s *ListingService) GetListing(ctx context.Context, id uint) (*entities.Listing, error) {
	// Try Redis cache first
	if s.redis != nil {
		listing, err := s.redis.GetListing(ctx, fmt.Sprintf("%d", id))
		if err == nil {
			return listing, nil
		}
	}

	// Fallback to PostgreSQL
	var dbListing entities.Listing
	if err := s.db.DB.First(&dbListing, id).Error; err != nil {
		return nil, err
	}

	// Cache in Redis
	if s.redis != nil {
		if err := s.redis.CacheListing(ctx, &dbListing); err != nil {
			return nil, err
		}
	}

	return &dbListing, nil
}

func (s *ListingService) UpdateListing(ctx context.Context, id uint, listing *entities.Listing) error {
	var existing entities.Listing
	if err := s.db.DB.First(&existing, id).Error; err != nil {
		return err
	}

	// Update fields
	existing.Title = listing.Title
	existing.Description = listing.Description
	existing.Location = listing.Location
	existing.Price = listing.Price
	existing.Available = listing.Available

	if err := s.db.DB.Save(&existing).Error; err != nil {
		return err
	}

	// Update Redis cache
	if s.redis != nil {
		if err := s.redis.CacheListing(ctx, &existing); err != nil {
			return err
		}
	}

	// Publish listing update event
	return s.producer.PublishMessage(ctx, "listing.updated", existing)
}

func (s *ListingService) DeleteListing(ctx context.Context, id uint) error {
	var listing entities.Listing
	if err := s.db.DB.First(&listing, id).Error; err != nil {
		return err
	}

	if err := s.db.DB.Delete(&listing).Error; err != nil {
		return err
	}

	// Remove from Redis cache
	if s.redis != nil {
		if err := s.redis.Client.Del(ctx, fmt.Sprintf("listing:%d", id)).Err(); err != nil {
			return err
		}
	}

	// Publish listing deletion event
	return s.producer.PublishMessage(ctx, "listing.deleted", map[string]uint{"id": id})
}

func (s *ListingService) CheckAvailability(ctx context.Context, listingID uint, startDate, endDate time.Time) (bool, error) {
	if startDate.After(endDate) || startDate.Before(time.Now()) {
		return false, fmt.Errorf("invalid date range")
	}

	var conflictingBookings []entities.BookedDates
	if err := s.db.DB.Where("listing_id = ? AND (start_date <= ? AND end_date >= ?)",
		listingID, endDate, startDate).Find(&conflictingBookings).Error; err != nil {
		return false, err
	}

	return len(conflictingBookings) == 0, nil
}
