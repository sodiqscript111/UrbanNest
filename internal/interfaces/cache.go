package interfaces

import (
	"UrbanNest/internal/entities"
	"context"
)

type CacheStore interface {
	CacheListing(ctx context.Context, listing *entities.Listing) error
	GetListing(ctx context.Context, key string) (*entities.Listing, error)

	CacheReview(ctx context.Context, review *entities.Review) error
	GetReview(ctx context.Context, key string) (*entities.Review, error)

	CacheReviewsByListing(ctx context.Context, listingID uint, reviews []entities.Review) error
	GetReviewsByListing(ctx context.Context, listingID uint) ([]entities.Review, error)

	CacheMessage(ctx context.Context, message *entities.Message) error
	GetMessage(ctx context.Context, key string) (*entities.Message, error)

	CacheMessagesByUser(ctx context.Context, userID uint, messages []entities.Message) error
	GetMessagesByUser(ctx context.Context, userID uint) ([]entities.Message, error)

	CacheBooking(ctx context.Context, booking *entities.Booking) error
	GetBooking(ctx context.Context, key string) (*entities.Booking, error)

	CacheBookingsByUser(ctx context.Context, userID uint, bookings []entities.Booking) error
	GetBookingsByUser(ctx context.Context, userID uint) ([]entities.Booking, error)

	CacheBookingsByHost(ctx context.Context, hostID uint, bookings []entities.Booking) error
	GetBookingsByHost(ctx context.Context, hostID uint) ([]entities.Booking, error)
	Delete(ctx context.Context, key string) error
}
