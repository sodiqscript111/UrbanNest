package interfaces

import (
	"UrbanNest/internal/entities"
	"context"
	"time"
)

type Database interface {
	//Users operations
	CreateUser(ctx context.Context, user *entities.User) error
	GetUser(ctx context.Context, id uint) (*entities.User, error)

	//Listings operations
	CreateListing(ctx context.Context, listing *entities.Listing) error
	GetListing(ctx context.Context, id uint) (*entities.Listing, error)
	UpdateListing(ctx context.Context, listing *entities.Listing) error
	DeleteListing(ctx context.Context, id uint) error
	CheckAvailability(ctx context.Context, listingID uint, startDate, endDate string) (bool, error)

	//Bookings operations
	CreateBooking(ctx context.Context, booking *entities.Booking) error
	GetBooking(ctx context.Context, id uint) (*entities.Booking, error)
	GetBookingsByUser(ctx context.Context, userID uint) ([]entities.Booking, error)
	GetBookingsByHost(ctx context.Context, hostID uint) ([]entities.Booking, error)
	UpdateBooking(ctx context.Context, booking *entities.Booking) error
	DeleteBooking(ctx context.Context, bookingID uint) error
	FindConflictingBookings(ctx context.Context, listingID uint, startDate, endDate time.Time) ([]entities.BookedDates, error)
}
