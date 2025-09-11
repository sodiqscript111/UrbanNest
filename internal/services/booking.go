package services

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/store"
	"UrbanNest/pkg/kafka"
	"context"
	"fmt"
	"time"
)

type BookingService struct {
	db       *store.PostgresStore
	redis    *store.RedisStore
	producer *kafka.Producer
}

func NewBookingService(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) *BookingService {
	return &BookingService{db, redis, producer}
}

func (s *BookingService) CreateBooking(ctx context.Context, booking *entities.Booking) error {
	// Validate booking dates
	if booking.StartDate.After(booking.EndDate) || booking.StartDate.Before(time.Now()) {
		return fmt.Errorf("invalid date range")
	}

	// Check if user and listing exist
	var user entities.User
	if err := s.db.DB.First(&user, booking.UserID).Error; err != nil {
		return fmt.Errorf("user not found")
	}
	var listing entities.Listing
	if err := s.db.DB.First(&listing, booking.ListingID).Error; err != nil {
		return fmt.Errorf("listing not found")
	}

	// Check for availability conflicts
	var conflictingBookings []entities.BookedDates
	if err := s.db.DB.Where("listing_id = ? AND (start_date <= ? AND end_date >= ?)",
		booking.ListingID, booking.EndDate, booking.StartDate).Find(&conflictingBookings).Error; err != nil {
		return err
	}
	if len(conflictingBookings) > 0 {
		return fmt.Errorf("listing is not available for the selected dates")
	}

	// Set default status
	booking.Status = "pending"

	// Create booking
	tx := s.db.DB.Create(booking)
	if err := tx.Error; err != nil {
		return err
	}

	// Cache booking in Redis
	if s.redis != nil {
		if err := s.redis.CacheBooking(ctx, booking); err != nil {
			return err
		}
	}

	// Publish booking creation event
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

	// Fallback to PostgreSQL
	var booking entities.Booking
	if err := s.db.DB.First(&booking, id).Error; err != nil {
		return nil, fmt.Errorf("booking not found")
	}

	// Cache in Redis
	if s.redis != nil {
		if err := s.redis.CacheBooking(ctx, &booking); err != nil {
			return nil, err
		}
	}

	return &booking, nil
}

func (s *BookingService) GetBookingsByUser(ctx context.Context, userID uint) ([]entities.Booking, error) {
	// Check if user exists
	var user entities.User
	if err := s.db.DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Try Redis cache
	if s.redis != nil {
		bookings, err := s.redis.GetBookingsByUser(ctx, userID)
		if err == nil && len(bookings) > 0 {
			return bookings, nil
		}
	}

	// Fallback to PostgreSQL
	var bookings []entities.Booking
	if err := s.db.DB.Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}

	// Cache in Redis
	if s.redis != nil {
		if err := s.redis.CacheBookingsByUser(ctx, userID, bookings); err != nil {
			return nil, err
		}
	}

	return bookings, nil
}

func (s *BookingService) GetBookingsByHost(ctx context.Context, hostID uint) ([]entities.Booking, error) {
	// Check if host exists
	var host entities.User
	if err := s.db.DB.First(&host, hostID).Error; err != nil {
		return nil, fmt.Errorf("host not found")
	}

	// Try Redis cache
	if s.redis != nil {
		bookings, err := s.redis.GetBookingsByHost(ctx, hostID)
		if err == nil && len(bookings) > 0 {
			return bookings, nil
		}
	}

	// Fallback to PostgreSQL
	var bookings []entities.Booking
	if err := s.db.DB.Joins("JOIN listings ON listings.id = bookings.listing_id").
		Where("listings.host_id = ?", hostID).Find(&bookings).Error; err != nil {
		return nil, err
	}

	// Cache in Redis
	if s.redis != nil {
		if err := s.redis.CacheBookingsByHost(ctx, hostID, bookings); err != nil {
			return nil, err
		}
	}

	return bookings, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, id uint) error {
	var booking entities.Booking
	if err := s.db.DB.First(&booking, id).Error; err != nil {
		return fmt.Errorf("booking not found")
	}

	// Check if booking is already canceled
	if booking.Status == "canceled" {
		return fmt.Errorf("booking is already canceled")
	}

	// Update status to canceled
	booking.Status = "canceled"
	if err := s.db.DB.Save(&booking).Error; err != nil {
		return err
	}

	// Remove from BookedDates
	if err := s.db.DB.Where("listing_id = ? AND start_date = ? AND end_date = ?",
		booking.ListingID, booking.StartDate, booking.EndDate).Delete(&entities.BookedDates{}).Error; err != nil {
		return err
	}

	// Invalidate caches
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

	// Publish booking cancellation event
	return s.producer.PublishMessage(ctx, "booking.canceled", booking)
}
