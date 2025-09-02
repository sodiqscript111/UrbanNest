package entities

import "time"

type Listing struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	HostID      uint      `json:"host_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
