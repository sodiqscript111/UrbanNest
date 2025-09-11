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

func CreateMessage(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var message entities.Message
		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service := services.NewMessageService(db, redis, producer)
		if err := service.CreateMessage(c.Request.Context(), &message); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, message)
	}
}

func GetMessage(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		service := services.NewMessageService(db, redis, producer)
		message, err := service.GetMessage(c.Request.Context(), uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
			return
		}
		c.JSON(http.StatusOK, message)
	}
}

func GetMessagesByUser(db *store.PostgresStore, redis *store.RedisStore, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		service := services.NewMessageService(db, redis, producer)
		messages, err := service.GetMessagesByUser(c.Request.Context(), uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, messages)
	}
}
