package entities

import "gorm.io/gorm"

type Listing struct {
	gorm.Model
	HostID      uint    `json:"host_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	Price       float64 `json:"price"`
	Available   bool    `json:"available"`
}
