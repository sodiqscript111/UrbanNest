package interfaces

import (
	"UrbanNest/internal/entities"
	"context"
)

type Database interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUser(ctx context.Context, id uint) (*entities.User, error)
	CreateListing(ctx context.Context, listing *entities.Listing) error
	GetListing(ctx context.Context, id uint) (*entities.Listing, error)
	UpdateListing(ctx context.Context, listing *entities.Listing) error
	DeleteListing(ctx context.Context, id uint) error
	CheckAvailability(ctx context.Context, listingID uint, startDate, endDate string) (bool, error)
}
