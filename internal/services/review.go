package services

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/interfaces"
	"UrbanNest/pkg/kafka"
	"context"
	"fmt"
)

type ReviewService struct {
	db       interfaces.Database
	redis    interfaces.CacheStore
	producer *kafka.Producer
}

func NewReviewService(db interfaces.Database, redis interfaces.CacheStore, producer *kafka.Producer) *ReviewService {
	return &ReviewService{db, redis, producer}
}

func (s *ReviewService) CreateReview(ctx context.Context, review *entities.Review) error {

	if review.Rating < 1 || review.Rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	if err := s.db.CreateReview(ctx, review); err != nil {
		return err
	}

	if s.redis != nil {
		if err := s.redis.CacheReview(ctx, review); err != nil {
			return err
		}
	}

	return s.producer.PublishMessage(ctx, "review.created", review)
}

func (s *ReviewService) GetReview(ctx context.Context, id uint) (*entities.Review, error) {

	if s.redis != nil {
		review, err := s.redis.GetReview(ctx, fmt.Sprintf("%d", id))
		if err == nil {
			return review, nil
		}
	}

	review, err := s.db.GetReview(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.redis != nil {
		if err := s.redis.CacheReview(ctx, review); err != nil {
			return nil, err
		}
	}

	return review, nil
}

func (s *ReviewService) GetReviewsByListing(ctx context.Context, listingID uint) ([]entities.Review, error) {

	if s.redis != nil {
		reviews, err := s.redis.GetReviewsByListing(ctx, listingID)
		if err == nil && len(reviews) > 0 {
			return reviews, nil
		}
	}

	reviews, err := s.db.GetReviewsByListing(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if s.redis != nil {
		if err := s.redis.CacheReviewsByListing(ctx, listingID, reviews); err != nil {
			return nil, err
		}
	}

	return reviews, nil
}
