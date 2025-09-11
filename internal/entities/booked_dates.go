package entities

import (
	"gorm.io/gorm"
	"time"
)

type BookedDates struct {
	gorm.Model
	ListingID uint      `json:"listing_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
