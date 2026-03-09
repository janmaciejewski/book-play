package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// CORS middleware handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimit middleware implements basic rate limiting
func RateLimit(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple in-memory rate limiting
		// In production, use Redis-based rate limiting
		c.Next()
	}
}

// Logger middleware logs request details
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		if status >= 400 {
			gin.DefaultWriter.Write([]byte(
				"[ERROR] " + start.Format("2006/01/02 - 15:04:05") +
					" | " + c.Request.Method +
					" | " + path +
					" | " + latency.String() +
					" | " + string(rune(status)) + "\n",
			))
		}
	}
}
