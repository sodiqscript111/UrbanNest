package main

import (
	"UrbanNest/api/handlers"
	"UrbanNest/api/middleware"
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
		r.POST("/register", handlers.Register(db, config.JWTSecret))
		r.POST("/login", handlers.Login(db, config.JWTSecret))

		// Protected routes
		protected := r.Group("/", middleware.Auth(config.JWTSecret))
		{
			// User routes
			protected.POST("/users", handlers.CreateUser(db, nil))
			protected.GET("/users/:id", handlers.GetUser(db, nil))

			// Listing routes
			protected.POST("/listings", handlers.CreateListing(db, redisStore, listingProducer))
			protected.GET("/listings/:id", handlers.GetListing(db, redisStore, listingProducer))
			protected.PUT("/listings/:id", handlers.UpdateListing(db, redisStore, listingProducer))
			protected.DELETE("/listings/:id", handlers.DeleteListing(db, redisStore, listingProducer))
			protected.GET("/listings/:id/availability", handlers.CheckAvailability(db, redisStore, listingProducer))

			// Review routes
			protected.POST("/reviews", handlers.CreateReview(db, redisStore, reviewProducer))
			protected.GET("/reviews/:id", handlers.GetReview(db, redisStore, reviewProducer))
			protected.GET("/listings/:id/reviews", handlers.GetReviewsByListing(db, redisStore, reviewProducer))

			// Message routes
			protected.POST("/messages", handlers.CreateMessage(db, redisStore, messageProducer))
			protected.GET("/messages/:id", handlers.GetMessage(db, redisStore, messageProducer))
			protected.GET("/users/:id/messages", handlers.GetMessagesByUser(db, redisStore, messageProducer))

			// Booking routes
			protected.POST("/bookings", handlers.CreateBooking(db, bookingProducer))
			protected.GET("/bookings/:id", handlers.GetBooking(db, redisStore, bookingProducer))
			protected.GET("/users/:id/bookings", handlers.GetBookingsByUser(db, redisStore, bookingProducer))
			protected.GET("/hosts/:id/bookings", handlers.GetBookingsByHost(db, redisStore, bookingProducer))
			protected.DELETE("/bookings/:id", handlers.CancelBooking(db, redisStore, bookingProducer))
		}

		r.Run(":" + config.Port)
	} else if *mode == "worker" {
		switch *consumerType {
		case "email":
			kafka.StartEmailConsumer(strings.Split(config.KafkaBrokers, ","), config.ResendAPIKey)
		case "booking":
			log.Println("Starting booking consumer")
			kafka.StartBookingConsumer(strings.Split(config.KafkaBrokers, ","), db, config.ResendAPIKey)
		case "listing":
			log.Println("Starting listing consumer")
			kafka.StartListingConsumer(strings.Split(config.KafkaBrokers, ","), db)
		case "message":
			log.Println("Starting message consumer")
			kafka.StartMessageConsumer(strings.Split(config.KafkaBrokers, ","), db, config.ResendAPIKey)
		default:
			log.Fatal("Invalid consumer type")
		}
	}
}
