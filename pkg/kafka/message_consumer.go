package kafka

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/store"
	"UrbanNest/pkg/email"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

func StartMessageConsumer(brokers []string, db *store.PostgresStore, resendAPIKey string) {
	consumer := NewConsumer(brokers, "message.sent", "message-group")
	ctx := context.Background()

	consumer.Consume(ctx, func(msg kafka.Message) {
		var message entities.Message
		if err := json.Unmarshal(msg.Value, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}

		// Fetch receiver's email from User table
		var receiver entities.User
		if err := db.DB.Where("id = ?", message.ReceiverID).First(&receiver).Error; err != nil {
			log.Printf("Error fetching receiver: %v", err)
			return
		}

		// Send email notification to receiver
		emailClient := email.NewResendClient(resendAPIKey)
		emailParams := email.EmailParams{
			To:      receiver.Email,
			Subject: "New Message Received",
			Body:    fmt.Sprintf("You received a message from user %d: %s", message.SenderID, message.Content),
		}
		if err := emailClient.SendEmail(ctx, emailParams); err != nil {
			log.Printf("Error sending email notification: %v", err)
			return
		}

		log.Printf("Processed message %d from user %d to user %d", message.ID, message.SenderID, message.ReceiverID)
	})
}
