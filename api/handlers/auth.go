package handlers

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/services"
	"UrbanNest/internal/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(db *store.PostgresStore, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user entities.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service := services.NewAuthService(db, jwtSecret)
		token, err := service.Register(c.Request.Context(), &user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"token": token})
	}
}

func Login(db *store.PostgresStore, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service := services.NewAuthService(db, jwtSecret)
		token, err := service.Login(c.Request.Context(), creds.Email, creds.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
