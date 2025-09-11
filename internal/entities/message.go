package entities

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	gorm.Model
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	ListingID  uint      `json:"listing_id"` // Optional: Tie message to a listing
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sent_at"`
}
