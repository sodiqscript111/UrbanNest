package handlers

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/services"
	"UrbanNest/internal/store"
	"UrbanNest/pkg/kafka"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func CreateListing(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var listing entities.Listing
		if err := c.ShouldBindJSON(&listing); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service := services.NewListingService(db, redis, producer)
		if err := service.CreateListing(c.Request.Context(), &listing); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, listing)
	}
}

func GetListing(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		service := services.NewListingService(db, redis, producer)
		listing, err := service.GetListing(c.Request.Context(), uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
			return
		}
		c.JSON(http.StatusOK, listing)
	}
}

func UpdateListing(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var listing entities.Listing
		if err := c.ShouldBindJSON(&listing); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service := services.NewListingService(db, redis, producer)
		if err := service.UpdateListing(c.Request.Context(), uint(id), &listing); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, listing)
	}
}

func DeleteListing(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		service := services.NewListingService(db, redis, producer)
		if err := service.DeleteListing(c.Request.Context(), uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func CheckAvailability(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
			return
		}
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
			return
		}

		service := services.NewListingService(db, redis, producer)
		available, err := service.CheckAvailability(c.Request.Context(), uint(id), startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"available": available})
	}
}

// Helper to parse ID (used in GetListing, UpdateListing, DeleteListing)
func parseID(idStr string) uint {
	id, _ := strconv.ParseUint(idStr, 10, 32)
	return uint(id)
}
