package services

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/interfaces"
	"UrbanNest/pkg/kafka"
	"context"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db       interfaces.Database
	producer *kafka.Producer
}

func NewUserService(db interfaces.Database, producer *kafka.Producer) *UserService {
	return &UserService{db, producer}
}

func (s *UserService) CreateUser(ctx context.Context, user *entities.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	if err := s.db.CreateUser(ctx, user); err != nil {
		return err
	}

	return s.producer.PublishMessage(ctx, "user.created", user)
}

func (s *UserService) GetUser(ctx context.Context, id uint) (*entities.User, error) {
	return s.db.GetUser(ctx, id)
}
