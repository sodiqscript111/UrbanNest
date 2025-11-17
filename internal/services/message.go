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

type MessageService struct {
	db       interfaces.Database
	redis    interfaces.CacheStore
	producer *kafka.Producer
}

func NewMessageService(db interfaces.Database, redis *store.RedisStore, producer *kafka.Producer) *MessageService {
	return &MessageService{db, redis, producer}
}

func (s *MessageService) CreateMessage(ctx context.Context, message *entities.Message) error {
	if message.Content == "" {
		return fmt.Errorf("message content cannot be empty")
	}

	bool, err := s.db.UserExistsByID(ctx, message.SenderID)
	if err != nil || !bool {
		return fmt.Errorf("sender not found")
	}

	bool, err = s.db.UserExistsByID(ctx, message.ReceiverID)
	if err != nil || !bool {
		return fmt.Errorf("receiver not found")
	}

	bool, err = s.db.UserExistsByID(ctx, message.SenderID)
	if err != nil || !bool {
		return fmt.Errorf("sender not found")
	}

	bool, err = s.db.UserExistsByID(ctx, message.ReceiverID)
	if err != nil || !bool {
		return fmt.Errorf("receiver not found")
	}

	message.SentAt = time.Now()

	if err := s.db.CreateMessage(ctx, message); err != nil {
		return err
	}

	if s.redis != nil {
		if err := s.redis.CacheMessage(ctx, message); err != nil {
			return err
		}

		if err := s.redis.Delete(ctx, fmt.Sprintf("user:%d:messages", message.SenderID)); err != nil {
			return err
		}
		if err := s.redis.Delete(ctx, fmt.Sprintf("user:%d:messages", message.ReceiverID)); err != nil {
			return err
		}
	}

	return s.producer.PublishMessage(ctx, "message.sent", message)
}

func (s *MessageService) GetMessage(ctx context.Context, id uint) (*entities.Message, error) {

	if s.redis != nil {
		message, err := s.redis.GetMessage(ctx, fmt.Sprintf("%d", id))
		if err == nil {
			return message, nil
		}
	}

	message, err := s.db.GetMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.redis != nil {
		if err := s.redis.CacheMessage(ctx, message); err != nil {
			return nil, err
		}
	}

	return message, nil
}

func (s *MessageService) GetMessagesByUser(ctx context.Context, userID uint) ([]entities.Message, error) {

	bool, err := s.db.UserExistsByID(ctx, userID)
	if err != nil || !bool {
		return nil, fmt.Errorf("user not found")
	}

	if s.redis != nil {
		messages, err := s.redis.GetMessagesByUser(ctx, userID)
		if err == nil && len(messages) > 0 {
			return messages, nil
		}
	}

	messages, err := s.db.GetMessagesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if s.redis != nil {
		if err := s.redis.CacheMessagesByUser(ctx, userID, messages); err != nil {
			return nil, err
		}
	}

	return messages, nil
}
