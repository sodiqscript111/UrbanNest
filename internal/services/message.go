package services

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/store"
	"UrbanNest/pkg/kafka"
	"context"
	"fmt"
	"time"
)

type MessageService struct {
	db       *store.PostgresStore
	redis    *store.RedisStore
	producer *kafka.Producer
}

func NewMessageService(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) *MessageService {
	return &MessageService{db, redis, producer}
}

func (s *MessageService) CreateMessage(ctx context.Context, message *entities.Message) error {
	// Validate message
	if message.Content == "" {
		return fmt.Errorf("message content cannot be empty")
	}

	// Check if sender and receiver exist
	var sender, receiver entities.User
	if err := s.db.DB.First(&sender, message.SenderID).Error; err != nil {
		return fmt.Errorf("sender not found")
	}
	if err := s.db.DB.First(&receiver, message.ReceiverID).Error; err != nil {
		return fmt.Errorf("receiver not found")
	}

	// Check if listing exists (if provided)
	if message.ListingID != 0 {
		var listing entities.Listing
		if err := s.db.DB.First(&listing, message.ListingID).Error; err != nil {
			return fmt.Errorf("listing not found")
		}
	}

	// Set sent timestamp
	message.SentAt = time.Now()

	// Save message
	tx := s.db.DB.Create(message)
	if err := tx.Error; err != nil {
		return err
	}

	// Cache message in Redis
	if s.redis != nil {
		if err := s.redis.CacheMessage(ctx, message); err != nil {
			return err
		}
		// Invalidate user messages cache
		if err := s.redis.Client.Del(ctx, fmt.Sprintf("user:%d:messages", message.SenderID)).Err(); err != nil {
			return err
		}
		if err := s.redis.Client.Del(ctx, fmt.Sprintf("user:%d:messages", message.ReceiverID)).Err(); err != nil {
			return err
		}
	}

	// Publish message sent event
	return s.producer.PublishMessage(ctx, "message.sent", message)
}

func (s *MessageService) GetMessage(ctx context.Context, id uint) (*entities.Message, error) {
	// Try Redis cache first
	if s.redis != nil {
		message, err := s.redis.GetMessage(ctx, fmt.Sprintf("%d", id))
		if err == nil {
			return message, nil
		}
	}

	// Fallback to PostgreSQL
	var message entities.Message
	if err := s.db.DB.First(&message, id).Error; err != nil {
		return nil, err
	}

	// Cache in Redis
	if s.redis != nil {
		if err := s.redis.CacheMessage(ctx, &message); err != nil {
			return nil, err
		}
	}

	return &message, nil
}

func (s *MessageService) GetMessagesByUser(ctx context.Context, userID uint) ([]entities.Message, error) {
	// Check if user exists
	var user entities.User
	if err := s.db.DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Try Redis cache
	if s.redis != nil {
		messages, err := s.redis.GetMessagesByUser(ctx, userID)
		if err == nil && len(messages) > 0 {
			return messages, nil
		}
	}

	// Fallback to PostgreSQL
	var messages []entities.Message
	if err := s.db.DB.Where("sender_id = ? OR receiver_id = ?", userID, userID).Find(&messages).Error; err != nil {
		return nil, err
	}

	// Cache in Redis
	if s.redis != nil {
		if err := s.redis.CacheMessagesByUser(ctx, userID, messages); err != nil {
			return nil, err
		}
	}

	return messages, nil
}
