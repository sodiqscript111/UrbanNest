package store

import (
	"UrbanNest/internal/entities"
	"UrbanNest/pkg/config"
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type PostgresStore struct {
	DB *gorm.DB
}

func NewPostgresStore(config *config.Config) (*PostgresStore, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&entities.User{}, &entities.Listing{}, &entities.Booking{}, &entities.Review{}, &entities.Message{}, &entities.BookedDates{})
	return &PostgresStore{DB: db}, nil
}

func (p *PostgresStore) CreateListing(ctx context.Context, listing *entities.Listing) error {
	return p.DB.Create(listing).Error
}

func (p *PostgresStore) GetListing(ctx context.Context, id uint) (*entities.Listing, error) {
	var l entities.Listing
	if err := p.DB.First(&l, id).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (p *PostgresStore) UpdateListing(ctx context.Context, listing *entities.Listing) error {
	return p.DB.Save(listing).Error
}

func (p *PostgresStore) DeleteListing(ctx context.Context, id uint) error {
	return p.DB.Delete(&entities.Listing{}, id).Error
}

func (p *PostgresStore) CheckAvailability(ctx context.Context, listingID uint, startDate, endDate string) (bool, error) {
	var conflicts []entities.BookedDates
	if err := p.DB.Where("listing_id = ? AND (start_date <= ? AND end_date >= ?)",
		listingID, endDate, startDate).Find(&conflicts).Error; err != nil {
		return false, err
	}
	return len(conflicts) == 0, nil
}

func (p *PostgresStore) CreateUser(ctx context.Context, user *entities.User) error {
	return p.DB.Create(user).Error
}

func (p *PostgresStore) GetUser(ctx context.Context, id uint) (*entities.User, error) {
	var u entities.User
	if err := p.DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (p *PostgresStore) CreateBooking(ctx context.Context, booking *entities.Booking) error {
	return p.DB.Create(booking).Error
}

func (p *PostgresStore) GetBooking(ctx context.Context, id uint) (*entities.Booking, error) {
	var b entities.Booking
	if err := p.DB.First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}
func (p *PostgresStore) GetBookingsByUser(ctx context.Context, userID uint) ([]entities.Booking, error) {
	var bookings []entities.Booking
	if err := p.DB.Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}
func (p *PostgresStore) GetBookingsByHost(ctx context.Context, hostID uint) ([]entities.Booking, error) {
	var bookings []entities.Booking
	if err := p.DB.Joins("JOIN listings ON listings.id = bookings.listing_id").
		Where("listings.host_id = ?", hostID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (p *PostgresStore) UpdateBooking(ctx context.Context, booking *entities.Booking) error {
	return p.DB.Save(booking).Error
}
func (p *PostgresStore) DeleteBooking(ctx context.Context, bookingID uint) error {
	return p.DB.Delete(&entities.Booking{}, bookingID).Error
}
func (p *PostgresStore) DeleteBookedDate(ctx context.Context, listingID uint, startDate, endDate time.Time) error {
	return p.DB.Where(
		"listing_id = ? AND start_date = ? AND end_date = ?",
		listingID, startDate, endDate,
	).Delete(&entities.BookedDates{}).Error
}

func (p *PostgresStore) FindConflictingBookings(ctx context.Context, listingID uint, start, end time.Time) ([]entities.BookedDates, error) {
	var bookings []entities.BookedDates
	err := p.DB.
		Where("listing_id = ? AND (start_date <= ? AND end_date >= ?)",
			listingID, end, start).
		Find(&bookings).Error
	return bookings, err
}

func (p *PostgresStore) CreateReview(ctx context.Context, review *entities.Review) error {
	return p.DB.Create(review).Error
}

func (p *PostgresStore) GetReview(ctx context.Context, id uint) (*entities.Review, error) {
	var r entities.Review
	if err := p.DB.First(&r, id).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (p *PostgresStore) GetReviewsByListing(ctx context.Context, listingID uint) ([]entities.Review, error) {
	var reviews []entities.Review
	if err := p.DB.Where("listing_id = ?", listingID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}
