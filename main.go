package main

import (
	"UrbanNest/api/handlers"
	"UrbanNest/api/middleware"
	"UrbanNest/internal/interfaces"
	"UrbanNest/internal/store"
	"UrbanNest/pkg/config"
	"UrbanNest/pkg/kafka"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func main() {
	config := config.LoadConfig()
	mode := flag.String("mode", "server", "Run mode: server or worker")
	consumerType := flag.String("consumer", "", "Consumer type: email, booking, message, listing, review")
	flag.Parse()

	db, err := store.NewPostgresStore(config)
	if err != nil {
		log.Fatal(err)
	}
	var database interfaces.Database = db
	redisStore := store.NewRedisStore(config.RedisAddr, config.RedisPassword)

	if *mode == "server" {
		bookingProducer := kafka.NewProducer(strings.Split(config.KafkaBrokers, ","), "booking.created,booking.canceled")
		listingProducer := kafka.NewProducer(strings.Split(config.KafkaBrokers, ","), "listing.created")
		reviewProducer := kafka.NewProducer(strings.Split(config.KafkaBrokers, ","), "review.created")
		messageProducer := kafka.NewProducer(strings.Split(config.KafkaBrokers, ","), "message.sent")
		defer bookingProducer.Close()
		defer listingProducer.Close()
		defer reviewProducer.Close()
		defer messageProducer.Close()

		r := gin.Default()
		r.Use(middleware.RateLimit(redisStore.Client))

		// Auth routes (public)
		r.POST("/register", handlers.Register(database, config.JWTSecret))
		r.POST("/login", handlers.Login(database, config.JWTSecret))

		// Protected routes
		protected := r.Group("/", middleware.Auth(config.JWTSecret))
		{
			// User routes
			protected.POST("/users", handlers.CreateUser(database, nil))
			protected.GET("/users/:id", handlers.GetUser(database, nil))

			// Listing routes
			protected.POST("/listings", handlers.CreateListing(database, redisStore, listingProducer))
			protected.GET("/listings/:id", handlers.GetListing(database, redisStore, listingProducer))
			protected.PUT("/listings/:id", handlers.UpdateListing(database, redisStore, listingProducer))
			protected.DELETE("/listings/:id", handlers.DeleteListing(database, redisStore, listingProducer))
			protected.GET("/listings/:id/availability", handlers.CheckAvailability(database, redisStore, listingProducer))

			// Review routes
			protected.POST("/reviews", handlers.CreateReview(database, redisStore, reviewProducer))
			protected.GET("/reviews/:id", handlers.GetReview(database, redisStore, reviewProducer))
			protected.GET("/listings/:id/reviews", handlers.GetReviewsByListing(database, redisStore, reviewProducer))

			// Message routes
			protected.POST("/messages", handlers.CreateMessage(database, redisStore, messageProducer))
			protected.GET("/messages/:id", handlers.GetMessage(database, redisStore, messageProducer))
			protected.GET("/users/:id/messages", handlers.GetMessagesByUser(database, redisStore, messageProducer))

			// Booking routes
			protected.POST("/bookings", handlers.CreateBooking(database, bookingProducer))
			protected.GET("/bookings/:id", handlers.GetBooking(database, redisStore, bookingProducer))
			protected.GET("/users/:id/bookings", handlers.GetBookingsByUser(database, redisStore, bookingProducer))
			protected.GET("/hosts/:id/bookings", handlers.GetBookingsByHost(database, redisStore, bookingProducer))
			protected.DELETE("/bookings/:id", handlers.CancelBooking(database, redisStore, bookingProducer))
		}

		r.Run(":" + config.Port)
	} else if *mode == "worker" {
		switch *consumerType {
		case "email":
			kafka.StartEmailConsumer(strings.Split(config.KafkaBrokers, ","), config.ResendAPIKey)
		case "booking":
			log.Println("Starting booking consumer")
			kafka.StartBookingConsumer(strings.Split(config.KafkaBrokers, ","), database, config.ResendAPIKey)
		case "listing":
			log.Println("Starting listing consumer")
			kafka.StartListingConsumer(strings.Split(config.KafkaBrokers, ","), database)
		case "message":
			log.Println("Starting message consumer")
			kafka.StartMessageConsumer(strings.Split(config.KafkaBrokers, ","), database, config.ResendAPIKey)
		default:
			log.Fatal("Invalid consumer type")
		}
	}
}
