package models

import (
	"UrbanNest/db"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/net/context"
	"net/url"

	"github.com/lib/pq" // For array support in batched queries
)

type Media struct {
	ID        int    `json:"id"`
	ListingID int    `json:"listing_id"`
	MediaURL  string `json:"media_url"`
	MediaType string `json:"media_type,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type Listing struct {
	ID          int     `json:"id"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Location    string  `json:"location" binding:"required"`
	ListerID    int     `json:"lister_id"`
	CreatedAt   string  `json:"created_at,omitempty"`
	IsAvailable bool    `json:"is_available"`
	Media       []Media `json:"media"` // New field for multiple media
}

func (listing *Listing) Create() error {
	query := `INSERT INTO listings (title, description, price, location, lister_id, is_available)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING id`
	err := db.DB.QueryRow(query, listing.Title, listing.Description, listing.Price, listing.Location, listing.ListerID, listing.IsAvailable).
		Scan(&listing.ID)
	if err != nil {
		return fmt.Errorf("❌ Failed to insert listing: %v", err)
	}
	return nil
}

func AddMedia(listingID int, mediaURL string, mediaType string) (Media, error) {
	// Decode URL-encoded mediaURL for S3
	decodedURL, err := url.QueryUnescape(mediaURL)
	if err != nil {
		return Media{}, fmt.Errorf("❌ Failed to decode media URL %s: %v", mediaURL, err)
	}
	fmt.Printf("Validating S3 file: %s\n", decodedURL)

	// Validate file exists in S3
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return Media{}, fmt.Errorf("❌ Failed to load AWS config: %v", err)
	}
	client := s3.NewFromConfig(cfg)

	_, err = client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String("urbannestbucket"), // Replace with your bucket name
		Key:    aws.String(decodedURL),
	})
	if err != nil {
		return Media{}, fmt.Errorf("❌ File %s does not exist in S3: %v", decodedURL, err)
	}
	fmt.Printf("S3 file %s validated successfully\n", decodedURL)

	var media Media
	query := `
        INSERT INTO listing_media (listing_id, media_url, media_type)
        VALUES ($1, $2, $3)
        RETURNING id, listing_id, media_url, media_type, created_at`
	err = db.DB.QueryRow(query, listingID, mediaURL, mediaType).
		Scan(&media.ID, &media.ListingID, &media.MediaURL, &media.MediaType, &media.CreatedAt)
	if err != nil {
		return media, fmt.Errorf("❌ Failed to add media to database: %v", err)
	}
	fmt.Printf("Inserted media into listing_media: %s\n", mediaURL)

	return media, nil
}

func GetAllListing() ([]Listing, error) {
	query := `SELECT id, title, description, price, location, lister_id, is_available FROM listings`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to fetch listings: %v", err)
	}
	defer rows.Close()

	var listings []Listing
	listingIDs := []int64{}
	for rows.Next() {
		var listing Listing
		err := rows.Scan(
			&listing.ID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.Location,
			&listing.ListerID,
			&listing.IsAvailable,
		)
		if err != nil {
			return nil, fmt.Errorf("❌ Failed to scan listing: %v", err)
		}
		listing.Media = []Media{}
		listings = append(listings, listing)
		listingIDs = append(listingIDs, int64(listing.ID))
	}

	if len(listingIDs) > 0 {
		mediaQuery := `
            SELECT id, listing_id, media_url, media_type, created_at
            FROM listing_media
            WHERE listing_id = ANY($1)`
		mediaRows, err := db.DB.Query(mediaQuery, pq.Array(listingIDs))
		if err != nil {
			return nil, fmt.Errorf("❌ Failed to fetch media: %v", err)
		}
		defer mediaRows.Close()

		mediaMap := make(map[int64][]Media)
		for mediaRows.Next() {
			var media Media
			if err := mediaRows.Scan(&media.ID, &media.ListingID, &media.MediaURL, &media.MediaType, &media.CreatedAt); err != nil {
				return nil, fmt.Errorf("❌ Failed to scan media: %v", err)
			}
			mediaMap[int64(media.ListingID)] = append(mediaMap[int64(media.ListingID)], media)
		}

		for i := range listings {
			listings[i].Media = mediaMap[int64(listings[i].ID)]
		}
	}

	return listings, nil
}

func GetById(id int64) (*Listing, error) {
	var listing Listing
	listing.Media = []Media{}

	query := `
        SELECT 
            l.id, l.title, l.description, l.price, l.location, l.lister_id, l.is_available,
            m.id AS media_id, m.listing_id AS media_listing_id, m.media_url, m.media_type, m.created_at AS media_created_at
        FROM listings l
        LEFT JOIN listing_media m ON l.id = m.listing_id
        WHERE l.id = $1`

	rows, err := db.DB.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to fetch listing: %v", err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true
		var mediaID sql.NullInt64
		var mediaListingID sql.NullInt64
		var mediaURL sql.NullString
		var mediaType sql.NullString
		var mediaCreatedAt sql.NullString

		err := rows.Scan(
			&listing.ID, &listing.Title, &listing.Description, &listing.Price,
			&listing.Location, &listing.ListerID, &listing.IsAvailable,
			&mediaID, &mediaListingID, &mediaURL, &mediaType, &mediaCreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("❌ Failed to scan row: %v", err)
		}

		if mediaID.Valid {
			listing.Media = append(listing.Media, Media{
				ID:        int(mediaID.Int64),
				ListingID: int(mediaListingID.Int64),
				MediaURL:  mediaURL.String,
				MediaType: mediaType.String,
				CreatedAt: mediaCreatedAt.String,
			})
		}
	}

	if !found {
		return nil, fmt.Errorf("listing with ID %d not found", id)
	}

	return &listing, nil
}

func EditListing(id int64, listing Listing) error {
	query := `
        UPDATE listings
        SET title = $1, description = $2, price = $3, location = $4, is_available = $5
        WHERE id = $6`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("❌ Failed to prepare update query: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(listing.Title, listing.Description, listing.Price, listing.Location, listing.IsAvailable, id)
	if err != nil {
		return fmt.Errorf("❌ Failed to update listing: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("❌ Failed to check update result: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("⚠️ No listing found with ID %d", id)
	}

	// Update media: Delete existing media and insert new ones
	if len(listing.Media) > 0 {
		// Delete existing media for this listing
		_, err := db.DB.Exec("DELETE FROM listing_media WHERE listing_id = $1", id)
		if err != nil {
			return fmt.Errorf("❌ Failed to delete existing media: %v", err)
		}

		// Insert new media
		for _, media := range listing.Media {
			_, err := AddMedia(int(id), media.MediaURL, media.MediaType)
			if err != nil {
				return fmt.Errorf("❌ Failed to add media: %v", err)
			}
		}
	}

	return nil
}

func MyListingById(listerID int64) ([]Listing, error) {
	query := `SELECT id, title, description, price, location, lister_id, created_at, is_available FROM listings WHERE lister_id = $1`
	rows, err := db.DB.Query(query, listerID)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to fetch listings: %v", err)
	}
	defer rows.Close()

	var listings []Listing
	listingIDs := []int64{}
	for rows.Next() {
		var listing Listing
		err := rows.Scan(
			&listing.ID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.Location,
			&listing.ListerID,
			&listing.CreatedAt,
			&listing.IsAvailable,
		)
		if err != nil {
			return nil, fmt.Errorf("❌ Failed to scan listing: %v", err)
		}
		listing.Media = []Media{}
		listings = append(listings, listing)
		listingIDs = append(listingIDs, int64(listing.ID))
	}

	if len(listingIDs) > 0 {
		mediaQuery := `
            SELECT id, listing_id, media_url, media_type, created_at
            FROM listing_media
            WHERE listing_id = ANY($1)`
		mediaRows, err := db.DB.Query(mediaQuery, pq.Array(listingIDs))
		if err != nil {
			return nil, fmt.Errorf("❌ Failed to fetch media: %v", err)
		}
		defer mediaRows.Close()

		mediaMap := make(map[int64][]Media)
		for mediaRows.Next() {
			var media Media
			if err := mediaRows.Scan(&media.ID, &media.ListingID, &media.MediaURL, &media.MediaType, &media.CreatedAt); err != nil {
				return nil, fmt.Errorf("❌ Failed to scan media: %v", err)
			}
			mediaMap[int64(media.ListingID)] = append(mediaMap[int64(media.ListingID)], media)
		}

		for i := range listings {
			listings[i].Media = mediaMap[int64(listings[i].ID)]
		}
	}

	return listings, nil
}

func DeleteById(id int64) error {
	query := `DELETE FROM listings WHERE id = $1`
	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("❌ Failed to delete listing: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("❌ Failed to check deletion result: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("⚠️ No listing found with ID %d", id)
	}

	// Note: Media is automatically deleted due to ON DELETE CASCADE in listing_media
	return nil
}
