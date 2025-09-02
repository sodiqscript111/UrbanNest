package services

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/store"
	"UrbanNest/pkg/kafka"
	"context"
)

type BookingService struct {
	db       *store.PostgresStore
	producer *kafka.Producer
}

func NewBookingService(db *store.PostgresStore, producer *kafka.Producer) *BookingService {
	return &BookingService{db, producer}
}

func (s *BookingService) CreateBooking(ctx context.Context, booking *entities.Booking) error {
	tx := s.db.DB.Create(booking)
	if err := tx.Error; err != nil {
		return err
	}
	return s.producer.PublishMessage(ctx, "booking", booking)
}
