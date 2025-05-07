package main

import (
	"UrbanNest/db"
	"UrbanNest/middleware"
	"UrbanNest/models"
	"UrbanNest/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := db.InitDB()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	err = godotenv.Load()
	if err != nil {
		fmt.Printf("Failed to load .env: %v\n", err)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/listings", getAllListing)
	router.GET("/listings/:id", getListingById)
	router.POST("/signup", addCustomer)
	router.POST("/signuplister", addLister)
	router.POST("/login", customerLogin)
	router.POST("/loginlister", listerLogin)
	router.POST("/booking", createBooking)
	router.GET("/getmybookings/:id", getBookingsByCustomer)
	router.GET("/api/upload-url", GeneratePresignedURL)
	router.POST("/verify-email", verifyEmail)

	authenticated := router.Group("/")
	authenticated.Use(middleware.Authorize)
	authenticated.POST("/listings", addListing)
	authenticated.GET("/mylisting", mylisting)
	authenticated.PUT("/listings/:id", editListing)
	authenticated.DELETE("/listings/:id", deleteWithId)

	fmt.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

func verifyEmail(c *gin.Context) {
	body, _ := c.GetRawData()
	fmt.Printf("Received /verify-email request: Body=%s, Headers=%v\n", string(body), c.Request.Header)

	var rawPayload map[string]interface{}
	if err := json.Unmarshal(body, &rawPayload); err != nil {
		fmt.Printf("Failed to parse raw JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	email, ok := rawPayload["email"].(string)
	if !ok || email == "" {
		fmt.Printf("Missing or invalid 'email' field: %v\n", rawPayload)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: missing or invalid 'email' field"})
		return
	}

	fmt.Printf("Parsed email: %s\n", email)

	emailRegex := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	if !emailRegex.MatchString(email) {
		fmt.Printf("Invalid email format: %s\n", email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"api_key": os.Getenv("ZERBOUNCE_API_KEY"),
			"email":   email,
		}).
		Get("https://api.zerobounce.net/v2/validate")
	if err != nil {
		fmt.Printf("Failed to verify email with ZeroBounce: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	fmt.Printf("ZeroBounce response: StatusCode=%d, Body=%s\n", resp.StatusCode(), string(resp.Body()))

	if resp.StatusCode() != http.StatusOK {
		fmt.Printf("ZeroBounce returned non-200 status: %d\n", resp.StatusCode())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Email verification failed with status %d", resp.StatusCode())})
		return
	}

	var zbResponse struct {
		Address       string `json:"address"`
		Status        string `json:"status"`
		SubStatus     string `json:"sub_status"`
		FreeEmail     bool   `json:"free_email"`
		Disposable    bool   `json:"disposable"`
		CatchAll      bool   `json:"catch_all"`
		MXFound       string `json:"mx_found"`
		SMTPValid     string `json:"smtp_valid"`
		Domain        string `json:"domain"`
		DomainAgeDays string `json:"domain_age_days"`
		FirstName     string `json:"firstname"`
		LastName      string `json:"lastname"`
		Gender        string `json:"gender"`
		Country       string `json:"country"`
		Region        string `json:"region"`
		City          string `json:"city"`
		ZipCode       string `json:"zipcode"`
		ProcessedAt   string `json:"processed_at"`
		ErrorMessage  string `json:"error_message"`
	}
	if err := json.Unmarshal(resp.Body(), &zbResponse); err != nil {
		fmt.Printf("Failed to parse ZeroBounce response: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse email verification response"})
		return
	}

	if zbResponse.ErrorMessage != "" {
		fmt.Printf("ZeroBounce API error: %s\n", zbResponse.ErrorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Email verification failed: %s", zbResponse.ErrorMessage)})
		return
	}

	mxFound := zbResponse.MXFound == "true"
	smtpValid := zbResponse.SMTPValid == "true"

	if zbResponse.Status != "valid" || !mxFound || !smtpValid {
		fmt.Printf("Email verification failed: %s (Status: %s, SubStatus: %s, MXFound: %s, SMTPValid: %s)\n",
			email, zbResponse.Status, zbResponse.SubStatus, zbResponse.MXFound, zbResponse.SMTPValid)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email does not exist or is invalid", "status": zbResponse.Status})
		return
	}
	if zbResponse.Disposable {
		fmt.Printf("Disposable email detected: %s\n", email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Disposable emails are not allowed", "status": "disposable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "valid"})
}

func addCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBind(&customer); err != nil {
		fmt.Printf("Failed to bind customer JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailRegex := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	if !emailRegex.MatchString(customer.Email) {
		fmt.Printf("Invalid email format: %s\n", customer.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"api_key": os.Getenv("ZERBOUNCE_API_KEY"),
			"email":   customer.Email,
		}).
		Get("https://api.zerobounce.net/v2/validate")
	if err != nil {
		fmt.Printf("Failed to verify email with ZeroBounce: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	fmt.Printf("ZeroBounce response: StatusCode=%d, Body=%s\n", resp.StatusCode(), string(resp.Body()))

	if resp.StatusCode() != http.StatusOK {
		fmt.Printf("ZeroBounce returned non-200 status: %d\n", resp.StatusCode())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Email verification failed with status %d", resp.StatusCode())})
		return
	}

	var zbResponse struct {
		Address       string `json:"address"`
		Status        string `json:"status"`
		SubStatus     string `json:"sub_status"`
		FreeEmail     bool   `json:"free_email"`
		Disposable    bool   `json:"disposable"`
		CatchAll      bool   `json:"catch_all"`
		MXFound       string `json:"mx_found"`
		SMTPValid     string `json:"smtp_valid"`
		Domain        string `json:"domain"`
		DomainAgeDays string `json:"domain_age_days"`
		FirstName     string `json:"firstname"`
		LastName      string `json:"lastname"`
		Gender        string `json:"gender"`
		Country       string `json:"country"`
		Region        string `json:"region"`
		City          string `json:"city"`
		ZipCode       string `json:"zipcode"`
		ProcessedAt   string `json:"processed_at"`
		ErrorMessage  string `json:"error_message"`
	}
	if err := json.Unmarshal(resp.Body(), &zbResponse); err != nil {
		fmt.Printf("Failed to parse ZeroBounce response: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse email verification response"})
		return
	}

	if zbResponse.ErrorMessage != "" {
		fmt.Printf("ZeroBounce API error: %s\n", zbResponse.ErrorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Email verification failed: %s", zbResponse.ErrorMessage)})
		return
	}

	mxFound := zbResponse.MXFound == "true"
	smtpValid := zbResponse.SMTPValid == "true"

	if zbResponse.Status != "valid" || !mxFound || !smtpValid {
		fmt.Printf("Email verification failed: %s (Status: %s, SubStatus: %s, MXFound: %s, SMTPValid: %s)\n",
			customer.Email, zbResponse.Status, zbResponse.SubStatus, zbResponse.MXFound, zbResponse.SMTPValid)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email does not exist or is invalid"})
		return
	}
	if zbResponse.Disposable {
		fmt.Printf("Disposable email detected: %s\n", customer.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Disposable emails are not allowed"})
		return
	}

	err = models.AddCustomer(customer)
	if err != nil {
		fmt.Printf("Failed to add customer: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Customer created successfully", "customer": customer})
}

func addListing(c *gin.Context) {
	var listing models.Listing
	if err := c.ShouldBindJSON(&listing); err != nil {
		fmt.Printf("Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("Received payload: %+v\n", listing)

	listing.ListerID = 1
	listing.IsAvailable = true

	if err := listing.Create(); err != nil {
		fmt.Printf("Failed to create listing: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, media := range listing.Media {
		addedMedia, err := models.AddMedia(listing.ID, media.MediaURL, media.MediaType)
		if err != nil {
			fmt.Printf("Failed to add media: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to add media: %v", err)})
			return
		}
		listing.Media = append(listing.Media[:0], addedMedia)
	}

	c.JSON(http.StatusCreated, listing)
}

func getAllListing(c *gin.Context) {
	fmt.Println("Listings route was hit")
	listings, err := models.GetAllListing()
	if err != nil {
		fmt.Printf("Failed to fetch listings: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to fetch listings: %v", err)})
		return
	}
	if len(listings) == 0 {
		fmt.Println("No listings found")
		c.JSON(http.StatusOK, gin.H{"listings": []models.Listing{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"listings": listings})
}

func getListingById(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid ID format: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	listing, err := models.GetById(id)
	if err != nil {
		fmt.Printf("Failed to fetch listing ID %d: %v\n", id, err)
		if err.Error() == fmt.Sprintf("listing with ID %d not found", id) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"listing": listing})
}

func deleteWithId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid ID format: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = models.DeleteById(id)
	if err != nil {
		fmt.Printf("Failed to delete listing ID %d: %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listing deleted successfully"})
}

func editListing(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid ID format: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var listing models.Listing
	if err := c.ShouldBindJSON(&listing); err != nil {
		fmt.Printf("Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	err = models.EditListing(id, listing)
	if err != nil {
		fmt.Printf("Failed to update listing ID %d: %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listing updated successfully"})
}

func customerLogin(c *gin.Context) {
	var customer models.Customer
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		fmt.Printf("Failed to bind login JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = customer.ValidateCustomer()
	if err != nil {
		fmt.Printf("Customer validation failed: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(customer.Email, customer.Id)
	if err != nil {
		fmt.Printf("Failed to generate token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer logged in", "token": token})
}

func mylisting(c *gin.Context) {
	listerID := c.GetInt("userId")
	listings, err := models.MyListingById(int64(listerID))
	if err != nil {
		fmt.Printf("Failed to fetch listings for lister %d: %v\n", listerID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"listings": listings})
}

func addLister(c *gin.Context) {
	var lister models.Lister
	err := c.ShouldBind(&lister)
	if err != nil {
		fmt.Printf("Failed to bind lister JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = models.AddLister(lister)
	if err != nil {
		fmt.Printf("Failed to add lister: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Lister created successfully"})
}

func listerLogin(c *gin.Context) {
	var lister models.Lister
	err := c.ShouldBindJSON(&lister)
	if err != nil {
		fmt.Printf("Failed to bind lister login JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = lister.ValidateLister()
	if err != nil {
		fmt.Printf("Lister validation failed: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(lister.Email, lister.Id)
	if err != nil {
		fmt.Printf("Failed to generate token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lister logged in", "token": token})
}

type BookingRequest struct {
	ListingID        int    `json:"listing_id" binding:"required"`
	CustomerID       int    `json:"customer_id" binding:"required"`
	StartDate        string `json:"start_date" binding:"required"`
	EndDate          string `json:"end_date" binding:"required"`
	Status           string `json:"status" binding:"required"`
	PaymentReference string `json:"payment_reference" binding:"required"`
}

func createBooking(c *gin.Context) {
	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Failed to bind booking JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	fmt.Printf("Booking request: %+v\n", req)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY")).
		Get("https://api.paystack.co/transaction/verify/" + req.PaymentReference)
	if err != nil || resp.StatusCode() != http.StatusOK {
		fmt.Printf("Failed to verify payment: %v, Status: %d, Response: %s\n", err, resp.StatusCode(), string(resp.Body()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to verify payment"})
		return
	}

	var paystackResp struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Status   string `json:"status"`
			Amount   int    `json:"amount"`
			Currency string `json:"currency"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &paystackResp); err != nil {
		fmt.Printf("Failed to parse Paystack response: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse payment response"})
		return
	}
	fmt.Printf("Paystack response: %+v\n", paystackResp)

	if !paystackResp.Status || paystackResp.Data.Status != "success" {
		fmt.Printf("Payment not successful: Status=%v, Message=%s\n", paystackResp.Data.Status, paystackResp.Message)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment not successful"})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		fmt.Printf("Invalid start date format: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		fmt.Printf("Invalid end date format: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	listing, err := models.GetById(int64(req.ListingID))
	if err != nil {
		fmt.Printf("Failed to fetch listing ID %d: %v\n", req.ListingID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch listing"})
		return
	}
	nights := int(endDate.Sub(startDate).Hours() / 24)
	expectedAmount := int(float64(nights) * listing.Price * 100)
	fmt.Printf("Amount check: Expected=%d, Paystack=%d, Nights=%d, Price=%.2f\n", expectedAmount, paystackResp.Data.Amount, nights, listing.Price)
	if paystackResp.Data.Amount != expectedAmount {
		fmt.Printf("Payment amount mismatch: Expected=%d, Got=%d\n", expectedAmount, paystackResp.Data.Amount)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment amount mismatch"})
		return
	}

	booking := models.Booking{
		ListingId:  req.ListingID,
		CustomerId: req.CustomerID,
		StartDate:  startDate,
		EndDate:    endDate,
		Status:     models.BookingStatus(req.Status),
	}

	createdBooking, err := models.CreateBooking(booking)
	if err != nil {
		fmt.Printf("Failed to create booking: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking created successfully",
		"booking": createdBooking,
	})
}

func getBookingsByCustomer(c *gin.Context) {
	customerIDStr := c.Param("id")
	id, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid customer ID format: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	bookings, err := models.GetBookingsByCustomer(int(id))
	if err != nil {
		fmt.Printf("Failed to fetch bookings for customer %d: %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

func GeneratePresignedURL(c *gin.Context) {
	fileName := c.Query("filename")
	if fileName == "" {
		fmt.Println("Filename is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load AWS config"})
		return
	}

	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)

	input := &s3.PutObjectInput{
		Bucket: aws.String("urbannestbucket"),
		Key:    aws.String(fileName),
	}

	presignedURL, err := presignClient.PresignPutObject(context.TODO(), input, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		fmt.Printf("Failed to generate presigned URL: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate presigned URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": presignedURL.URL})
}
