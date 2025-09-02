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
	consumerType := flag.String("consumer", "", "Consumer type: email, booking, message")
	flag.Parse()

	db, err := store.NewPostgresStore(config)
	if err != nil {
		log.Fatal(err)
	}
	redisStore := store.NewRedisStore(config.RedisAddr, config.RedisPassword)

	if *mode == "server" {
		bookingProducer := kafka.NewProducer(strings.Split(config.KafkaBrokers, ","), "booking.created")
		defer bookingProducer.Close()

		r := gin.Default()
		r.Use(middleware.RateLimit(redisStore.Client))

		r.POST("/bookings", handlers.CreateBooking(db, bookingProducer))

		r.Run(":" + config.Port)
	} else if *mode == "worker" {
		switch *consumerType {
		case "email":
			kafka.StartEmailConsumer(strings.Split(config.KafkaBrokers, ","), config.ResendAPIKey)
		case "booking":
			log.Println("Starting booking consumer")
			kafka.StartBookingConsumer(strings.Split(config.KafkaBrokers, ","), db)
		case "message":
			log.Println("Starting message consumer")
		default:
			log.Fatal("Invalid consumer type")
		}
	}
}
