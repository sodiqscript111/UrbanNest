package store

import (
	"UrbanNest/internal/entities"
	"UrbanNest/pkg/config"
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
