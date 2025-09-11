package handlers

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/services"
	"UrbanNest/internal/store"
	"UrbanNest/pkg/kafka"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateBooking(db *store.PostgresStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var booking entities.Booking
		if err := c.ShouldBindJSON(&booking); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service := services.NewBookingService(db, nil, producer)
		if err := service.CreateBooking(c.Request.Context(), &booking); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, booking)
	}
}

func GetBooking(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		service := services.NewBookingService(db, redis, producer)
		booking, err := service.GetBooking(c.Request.Context(), uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, booking)
	}
}

func GetBookingsByUser(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		service := services.NewBookingService(db, redis, producer)
		bookings, err := service.GetBookingsByUser(c.Request.Context(), uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, bookings)
	}
}

func GetBookingsByHost(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		hostID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
			return
		}

		service := services.NewBookingService(db, redis, producer)
		bookings, err := service.GetBookingsByHost(c.Request.Context(), uint(hostID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, bookings)
	}
}

func CancelBooking(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		service := services.NewBookingService(db, redis, producer)
		if err := service.CancelBooking(c.Request.Context(), uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Booking canceled"})
	}
}
