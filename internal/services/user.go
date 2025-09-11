package services

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/store"
	"UrbanNest/pkg/kafka"
	"context"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db       *store.PostgresStore
	producer *kafka.Producer
}

func NewUserService(db *store.PostgresStore, producer *kafka.Producer) *UserService {
	return &UserService{db, producer}
}

func (s *UserService) CreateUser(ctx context.Context, user *entities.User) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Save to database
	tx := s.db.DB.Create(user)
	if err := tx.Error; err != nil {
		return err
	}

	// Publish user creation event
	return s.producer.PublishMessage(ctx, "user.created", user)
}

func (s *UserService) GetUser(ctx context.Context, id uint) (*entities.User, error) {
	var user entities.User
	if err := s.db.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
