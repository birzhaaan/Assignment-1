package utils

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type client struct {
	requests int
	lastSeen time.Time
}

var (
	mu      sync.Mutex         
	clients = make(map[string]*client) 
)

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		mu.Lock()       
		defer mu.Unlock() 

		identifier := c.ClientIP()
		if userID, exists := c.Get("userID"); exists {
			identifier = userID.(string)
		}

		now := time.Now()

		if _, found := clients[identifier]; !found {
			clients[identifier] = &client{requests: 0, lastSeen: now}
		}

		if now.Sub(clients[identifier].lastSeen) > 10*time.Second {
			clients[identifier].requests = 0
			clients[identifier].lastSeen = now
		}

		clients[identifier].requests++

		if clients[identifier].requests > 5 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			return
		}

		c.Next()
	}
}