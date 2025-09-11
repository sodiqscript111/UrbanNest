package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

func RateLimit(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		key := "ratelimit:" + c.ClientIP()

		// Get current request count
		count, err := client.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limit error"})
			c.Abort()
			return
		}

		// Set limit (e.g., 100 requests per minute)
		limit := 100
		window := time.Minute

		if count >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		// Increment count and set expiry if first request
		pipe := client.Pipeline()
		pipe.Incr(ctx, key)
		if count == 0 {
			pipe.Expire(ctx, key, window)
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limit error"})
			c.Abort()
			return
		}

		c.Next()
	}
}
