package middleware

import (
	"UrbanNest/pkg/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RateLimit(client *redis.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "ratelimit:" + c.ClientIP()
		count, err := client.Client.Incr(c, key).Result()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if count == 1 {
			client.Client.Expire(c, key, time.Minute)
		}
		if count > 10 { // 10 requests per minute
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
