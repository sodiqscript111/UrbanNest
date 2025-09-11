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

func StartBookingConsumer(brokers []string, db *store.PostgresStore, resendAPIKey string) {
	consumer := NewConsumer(brokers, "booking.created,booking.canceled", "booking-group")
	ctx := context.Background()

	consumer.Consume(ctx, func(msg kafka.Message) {
		if msg.Topic == "booking.created" {
			var booking entities.Booking
			if err := json.Unmarshal(msg.Value, &booking); err != nil {
				log.Printf("Error unmarshaling booking: %v", err)
				return
			}

			// Add to BookedDates
			bookedDates := entities.BookedDates{
				ListingID: booking.ListingID,
				StartDate: booking.StartDate,
				EndDate:   booking.EndDate,
			}
			if err := db.DB.Create(&bookedDates).Error; err != nil {
				log.Printf("Error saving booked dates: %v", err)
				return
			}

			// Send confirmation email
			var user entities.User
			if err := db.DB.Where("id = ?", booking.UserID).First(&user).Error; err != nil {
				log.Printf("Error fetching user: %v", err)
				return
			}
			emailClient := email.NewResendClient(resendAPIKey)
			emailParams := email.EmailParams{
				To:      user.Email,
				Subject: "Booking Confirmation",
				Body:    fmt.Sprintf("Your booking for listing %d from %s to %s is confirmed.", booking.ListingID, booking.StartDate, booking.EndDate),
			}
			if err := emailClient.SendEmail(ctx, emailParams); err != nil {
				log.Printf("Error sending email: %v", err)
				return
			}

			log.Printf("Processed booking %d for listing %d by user %d", booking.ID, booking.ListingID, booking.UserID)
		} else if msg.Topic == "booking.canceled" {
			var booking entities.Booking
			if err := json.Unmarshal(msg.Value, &booking); err != nil {
				log.Printf("Error unmarshaling canceled booking: %v", err)
				return
			}

			// Send cancellation email to user
			var user entities.User
			if err := db.DB.Where("id = ?", booking.UserID).First(&user).Error; err != nil {
				log.Printf("Error fetching user: %v", err)
				return
			}
			emailClient := email.NewResendClient(resendAPIKey)
			emailParams := email.EmailParams{
				To:      user.Email,
				Subject: "Booking Canceled",
				Body:    fmt.Sprintf("Your booking for listing %d from %s to %s has been canceled.", booking.ListingID, booking.StartDate, booking.EndDate),
			}
			if err := emailClient.SendEmail(ctx, emailParams); err != nil {
				log.Printf("Error sending cancellation email: %v", err)
				return
			}

			// Notify host
			var listing entities.Listing
			if err := db.DB.Where("id = ?", booking.ListingID).First(&listing).Error; err != nil {
				log.Printf("Error fetching listing: %v", err)
				return
			}
			var host entities.User
			if err := db.DB.Where("id = ?", listing.HostID).First(&host).Error; err != nil {
				log.Printf("Error fetching host: %v", err)
				return
			}
			hostEmailParams := email.EmailParams{
				To:      host.Email,
				Subject: "Booking Canceled",
				Body:    fmt.Sprintf("The booking for your listing %d from %s to %s has been canceled.", booking.ListingID, booking.StartDate, booking.EndDate),
			}
			if err := emailClient.SendEmail(ctx, hostEmailParams); err != nil {
				log.Printf("Error sending host notification: %v", err)
				return
			}

			log.Printf("Processed cancellation for booking %d", booking.ID)
		}
	})
}
