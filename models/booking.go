package models

import (
	"UrbanNest/db"
	"fmt"
	"time"
)

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	Id         int           `json:"id"`
	ListingId  int           `json:"listing_id" binding:"required"`
	CustomerId int           `json:"customer_id"`
	StartDate  time.Time     `json:"start_date" binding:"required"`
	EndDate    time.Time     `json:"end_date" binding:"required"`
	Status     BookingStatus `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
}

func CreateBooking(booking Booking) error {
	query := `INSERT INTO bookings (listing_id, customer_id, start_date, end_date, status) 
              VALUES ($1, $2, $3, $4, $5)`

	stmt, err := db.DB.Prepare(query) // Prepare the query
	if err != nil {
		return fmt.Errorf("❌ Failed to prepare query: %v", err)
	}
	defer stmt.Close()

	// Execute the prepared statement with values from the booking struct
	_, err = stmt.Exec(booking.ListingId, booking.CustomerId, booking.StartDate, booking.EndDate, booking.Status)
	if err != nil {
		return fmt.Errorf("❌ Failed to execute query: %v", err)
	}

	return nil // Return nil if no errors occurred
}
