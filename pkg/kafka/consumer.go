package kafka

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/store"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type Consumer struct {
	Reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}
}

func (c *Consumer) Consume(ctx context.Context, handler func(message kafka.Message)) {
	for {
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}
		handler(msg)
	}
}

func (c *Consumer) Close() error {
	return c.Reader.Close()
}

func StartBookingConsumer(brokers []string, db *store.PostgresStore) {
	consumer := NewConsumer(brokers, "booking.created", "booking-group")
	ctx := context.Background()

	consumer.Consume(ctx, func(msg kafka.Message) {
		var booking entities.Booking
		if err := json.Unmarshal(msg.Value, &booking); err != nil {
			log.Printf("Error unmarshaling booking: %v", err)
			return
		}
		// Example: Update listing availability
		log.Printf("Processed booking for listing %d", booking.ListingID)
	})
}
