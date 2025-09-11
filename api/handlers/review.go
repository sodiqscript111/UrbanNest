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

func CreateReview(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var review entities.Review
		if err := c.ShouldBindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service := services.NewReviewService(db, redis, producer)
		if err := service.CreateReview(c.Request.Context(), &review); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, review)
	}
}

func GetReview(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		service := services.NewReviewService(db, redis, producer)
		review, err := service.GetReview(c.Request.Context(), uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
			return
		}
		c.JSON(http.StatusOK, review)
	}
}

func GetReviewsByListing(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		listingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
			return
		}

		service := services.NewReviewService(db, redis, producer)
		reviews, err := service.GetReviewsByListing(c.Request.Context(), uint(listingID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, reviews)
	}
}
