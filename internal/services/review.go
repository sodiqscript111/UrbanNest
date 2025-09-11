package services

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/store"
	"UrbanNest/pkg/kafka"
	"context"
	"fmt"
)

type ReviewService struct {
	db       *store.PostgresStore
	redis    *store.RedisStore
	producer *kafka.Producer
}

func NewReviewService(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) *ReviewService {
	return &ReviewService{db, redis, producer}
}

func (s *ReviewService) CreateReview(ctx context.Context, review *entities.Review) error {
	// Validate review
	if review.Rating < 1 || review.Rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	// Check if listing exists
	var listing entities.Listing
	if err := s.db.DB.First(&listing, review.ListingID).Error; err != nil {
		return fmt.Errorf("listing not found")
	}

	// Save review
	tx := s.db.DB.Create(review)
	if err := tx.Error; err != nil {
		return err
	}

	// Cache review in Redis (optional)
	if s.redis != nil {
		if err := s.redis.CacheReview(ctx, review); err != nil {
			return err
		}
	}

	// Publish review creation event
	return s.producer.PublishMessage(ctx, "review.created", review)
}

func (s *ReviewService) GetReview(ctx context.Context, id uint) (*entities.Review, error) {
	// Try Redis cache first
	if s.redis != nil {
		review, err := s.redis.GetReview(ctx, fmt.Sprintf("%d", id))
		if err == nil {
			return review, nil
		}
	}

	// Fallback to PostgreSQL
	var review entities.Review
	if err := s.db.DB.First(&review, id).Error; err != nil {
		return nil, err
	}

	// Cache in Redis
	if s.redis != nil {
		if err := s.redis.CacheReview(ctx, &review); err != nil {
			return nil, err
		}
	}

	return &review, nil
}

func (s *ReviewService) GetReviewsByListing(ctx context.Context, listingID uint) ([]entities.Review, error) {
	// Check if listing exists
	var listing entities.Listing
	if err := s.db.DB.First(&listing, listingID).Error; err != nil {
		return nil, fmt.Errorf("listing not found")
	}

	// Try Redis cache for listing reviews
	if s.redis != nil {
		reviews, err := s.redis.GetReviewsByListing(ctx, listingID)
		if err == nil && len(reviews) > 0 {
			return reviews, nil
		}
	}

	// Fallback to PostgreSQL
	var reviews []entities.Review
	if err := s.db.DB.Where("listing_id = ?", listingID).Find(&reviews).Error; err != nil {
		return nil, err
	}

	// Cache in Redis
	if s.redis != nil {
		if err := s.redis.CacheReviewsByListing(ctx, listingID, reviews); err != nil {
			return nil, err
		}
	}

	return reviews, nil
}
