package handlers

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/interfaces"
	"UrbanNest/internal/services"
	"UrbanNest/pkg/kafka"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUser(db interfaces.Database, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user entities.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service := services.NewUserService(db, producer)
		if err := service.CreateUser(c.Request.Context(), &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func GetUser(db interfaces.Database, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		service := services.NewUserService(db, producer)
		user, err := service.GetUser(c.Request.Context(), parseID(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
