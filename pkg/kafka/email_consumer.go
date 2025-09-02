package kafka

import (
	"UrbanNest/pkg/email"
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

func StartEmailConsumer(brokers []string, apiKey string) {
	consumer := NewConsumer(brokers, "notification.email", "email-group")
	ctx := context.Background()

	consumer.Consume(ctx, func(msg kafka.Message) {
		err := email.SendEmail(apiKey, "recipient@example.com", "Notification", string(msg.Value))
		if err != nil {
			log.Printf("Error sending email: %v", err)
		}
	})
}
