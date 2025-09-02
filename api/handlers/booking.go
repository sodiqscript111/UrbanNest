package handlers

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/services"
	"UrbanNest/internal/store"
	"UrbanNest/pkg/kafka"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateBooking(db *store.PostgresStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var booking entities.Booking
		if err := c.ShouldBindJSON(&booking); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service := services.NewBookingService(db, producer)
		if err := service.CreateBooking(c.Request.Context(), &booking); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, booking)
	}
}
