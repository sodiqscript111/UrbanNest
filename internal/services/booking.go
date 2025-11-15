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

type BookingService struct {
	db       interfaces.Database
	redis    *store.RedisStore
	producer *kafka.Producer
}

func NewBookingService(db interfaces.Database, redis *store.RedisStore, producer *kafka.Producer) *BookingService {
	return &BookingService{db, redis, producer}
}
func (s *BookingService) CreateBooking(ctx context.Context, booking *entities.Booking) error {
	// Validate
	if booking.StartDate.After(booking.EndDate) || booking.StartDate.Before(time.Now()) {
		return fmt.Errorf("invalid date range")
	}

	// Check user
	if _, err := s.db.GetUser(ctx, booking.UserID); err != nil {
		return fmt.Errorf("user not found")
	}

	// Check listing
	if _, err := s.db.GetListing(ctx, booking.ListingID); err != nil {
		return fmt.Errorf("listing not found")
	}

	// Check conflicts
	conflicts, err := s.db.FindConflictingBookings(ctx, booking.ListingID, booking.StartDate, booking.EndDate)
	if err != nil {
		return err
	}
	if len(conflicts) > 0 {
		return fmt.Errorf("listing is not available for selected dates")
	}

	// Default
	booking.Status = "pending"

	// Save
	if err := s.db.CreateBooking(ctx, booking); err != nil {
		return err
	}

	// Redis
	if s.redis != nil {
		if err := s.redis.CacheBooking(ctx, booking); err != nil {
			return err
		}
	}

	// Kafka
	return s.producer.PublishMessage(ctx, "booking.created", booking)
}

func (s *BookingService) GetBooking(ctx context.Context, id uint) (*entities.Booking, error) {
	// Try Redis cache first
	if s.redis != nil {
		booking, err := s.redis.GetBooking(ctx, fmt.Sprintf("%d", id))
		if err == nil {
			return booking, nil
		}
	}

	// Fetch from DB
	booking, err := s.db.GetBooking(ctx, id)
	if err != nil {
		return nil, err
	}
	return booking, nil

	// Cache in Redis
	if s.redis != nil {
		if err := s.redis.CacheBooking(ctx, booking); err != nil {
			return nil, err
		}
	}

	return booking, nil
}
func (s *BookingService) GetBookingsByUser(ctx context.Context, userID uint) ([]entities.Booking, error) {
	if _, err := s.db.GetUser(ctx, userID); err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Try Redis
	if s.redis != nil {
		if bookings, err := s.redis.GetBookingsByUser(ctx, userID); err == nil && len(bookings) > 0 {
			return bookings, nil
		}
	}

	// DB
	bookings, err := s.db.GetBookingsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache
	if s.redis != nil {
		_ = s.redis.CacheBookingsByUser(ctx, userID, bookings)
	}

	return bookings, nil
}

func (s *BookingService) GetBookingsByHost(ctx context.Context, hostID uint) ([]entities.Booking, error) {

	if _, err := s.db.GetUser(ctx, hostID); err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if s.redis != nil {
		bookings, err := s.redis.GetBookingsByHost(ctx, hostID)
		if err == nil && len(bookings) > 0 {
			return bookings, nil
		}
	}

	bookings, err := s.db.GetBookingsByHost(ctx, hostID)
	if err != nil {
		return nil, err
	}

	if s.redis != nil {
		if err := s.redis.CacheBookingsByHost(ctx, hostID, bookings); err != nil {
			return nil, err
		}
	}

	return bookings, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, id uint) error {
	booking, err := s.db.GetBooking(ctx, id)
	if err != nil {
		return err
	}

	if booking.Status == "canceled" {
		return fmt.Errorf("booking is already canceled")
	}

	booking.Status = "canceled"
	if err := s.db.UpdateBooking(ctx, booking); err != nil {
		return err
	}
	if err := s.db.DB.Where("listing_id = ? AND start_date = ? AND end_date = ?",
		booking.ListingID, booking.StartDate, booking.EndDate).Delete(&entities.BookedDates{}).Error; err != nil {
		return err
	}

	if s.redis != nil {
		if err := s.redis.Client.Del(ctx, fmt.Sprintf("booking:%d", booking.ID)).Err(); err != nil {
			return err
		}
		if err := s.redis.Client.Del(ctx, fmt.Sprintf("user:%d:bookings", booking.UserID)).Err(); err != nil {
			return err
		}
		var listing entities.Listing
		if err := s.db.DB.Where("id = ?", booking.ListingID).First(&listing).Error; err == nil {
			if err := s.redis.Client.Del(ctx, fmt.Sprintf("host:%d:bookings", listing.HostID)).Err(); err != nil {
				return err
			}
		}
	}

	return s.producer.PublishMessage(ctx, "booking.canceled", booking)
}
