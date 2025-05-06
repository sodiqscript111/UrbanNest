package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	connStr := "postgres://postgres:password@localhost:5432/urbanest?sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("❌ Failed to connect to database: %v", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("❌ Failed to ping the database: %v", err)
	}

	// Create tables
	err = createTables()
	if err != nil {
		return err
	}

	fmt.Println("✅ Database initialized and tables ready.")
	return nil
}

func createTables() error {
	createCustomers := `
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);`
	if _, err := DB.Exec(createCustomers); err != nil {
		return fmt.Errorf("❌ Failed to create customers table: %v", err)
	}

	createListers := `
CREATE TABLE IF NOT EXISTS listers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    nin_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);`
	if _, err := DB.Exec(createListers); err != nil {
		return fmt.Errorf("❌ Failed to create listers table: %v", err)
	}

	createListings := `
CREATE TABLE IF NOT EXISTS listings (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    location VARCHAR(255) NOT NULL,
    lister_id INTEGER NOT NULL REFERENCES listers(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);`
	if _, err := DB.Exec(createListings); err != nil {
		return fmt.Errorf("❌ Failed to create listings table: %v", err)
	}

	createBookings := `
CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    listing_id INTEGER NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    customer_id INTEGER NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending','confirmed','cancelled')),
    created_at TIMESTAMP DEFAULT NOW()
);`
	if _, err := DB.Exec(createBookings); err != nil {
		return fmt.Errorf("❌ Failed to create bookings table: %v", err)
	}

	createReviews := `
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    reviewer_id INTEGER NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);`
	if _, err := DB.Exec(createReviews); err != nil {
		return fmt.Errorf("❌ Failed to create reviews table: %v", err)
	}

	createListingMedia := `
CREATE TABLE IF NOT EXISTS listing_media (
    id SERIAL PRIMARY KEY,
    listing_id INTEGER NOT NULL,
    media_url VARCHAR(255) NOT NULL,
    media_type VARCHAR(50), -- Optional: e.g., 'image', 'video'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (listing_id) REFERENCES listings(id) ON DELETE CASCADE
)`
	if _, err := DB.Exec(createListingMedia); err != nil {
		return fmt.Errorf("<UNK> Failed to create listing_media table: %v", err)
	}
	return nil
}
