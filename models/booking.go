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

func CreateBooking(booking Booking) (Booking, error) {

	queryAvailability := `
		SELECT COUNT(*)
		FROM bookings
		WHERE listing_id = $1
		  AND status = 'confirmed'
		  AND (start_date, end_date) OVERLAPS ($2::DATE, $3::DATE)
	`

	var count int
	err := db.DB.QueryRow(queryAvailability, booking.ListingId, booking.StartDate, booking.EndDate).Scan(&count)
	if err != nil {
		return Booking{}, fmt.Errorf("Failed to check availability: %v", err)
	}

	if count > 0 {
		return Booking{}, fmt.Errorf("Listing is already booked for the selected dates")
	}

	queryInsert := `
		INSERT INTO bookings (listing_id, customer_id, start_date, end_date, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err = db.DB.QueryRow(queryInsert,
		booking.ListingId,
		booking.CustomerId,
		booking.StartDate,
		booking.EndDate,
		booking.Status,
	).Scan(&booking.Id, &booking.CreatedAt)
	if err != nil {
		return Booking{}, fmt.Errorf("Failed to insert booking: %v", err)
	}

	return booking, nil
}

func GetBookingsByCustomer(id int) ([]Booking, error) {
	query := `SELECT id, listing_id, customer_id, start_date, end_date, status, created_at FROM bookings WHERE customer_id = $1`
	rows, err := db.DB.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get bookings: %v", err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var booking Booking
		err := rows.Scan(
			&booking.Id,
			&booking.ListingId,
			&booking.CustomerId,
			&booking.StartDate,
			&booking.EndDate,
			&booking.Status,
			&booking.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan booking row: %v", err)
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Rows iteration error: %v", err)
	}

	return bookings, nil
}
