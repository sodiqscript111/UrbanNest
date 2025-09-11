package kafka

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/store"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

func StartListingConsumer(brokers []string, db *store.PostgresStore) {
	// Consumer for listing.updated
	updatedConsumer := NewConsumer(brokers, "listing.updated", "listing-group")
	ctx := context.Background()

	go updatedConsumer.Consume(ctx, func(msg kafka.Message) {
		var listing entities.Listing
		if err := json.Unmarshal(msg.Value, &listing); err != nil {
			log.Printf("Error unmarshaling listing: %v", err)
			return
		}
		log.Printf("Processed updated listing %d: %s", listing.ID, listing.Title)
		// Add logic (e.g., notify users, update search index)
	})

	// Consumer for listing.deleted
	deletedConsumer := NewConsumer(brokers, "listing.deleted", "listing-group")
	deletedConsumer.Consume(ctx, func(msg kafka.Message) {
		var data map[string]uint
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			log.Printf("Error unmarshaling deletion: %v", err)
			return
		}
		id := data["id"]
		log.Printf("Processed deleted listing %d", id)
		// Add logic (e.g., remove from search index)
	})
}
