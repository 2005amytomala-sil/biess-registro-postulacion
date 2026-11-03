package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or get request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		// Log request
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		log.Printf("[%s] %s %s - RequestID: %s", method, path, c.Request.RemoteAddr, requestID)

		// Process request
		c.Next()

		// Log response
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		log.Printf("[%s] %s - Status: %d - Latency: %v - RequestID: %s - IP: %s",
			method, path, statusCode, latency, requestID, clientIP)
	}
}
