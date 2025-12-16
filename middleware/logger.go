package middleware

import (
	"time"
	"smart-file-api/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get request info
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		
		// Get user ID if authenticated
		userID, exists := c.Get("user_id")
		if !exists {
			userID = "anonymous"
		}

		// Determine log level based on status code
		logEntry := config.Log.WithFields(logrus.Fields{
			"status_code": statusCode,
			"latency":     latency.Milliseconds(),
			"method":      method,
			"path":        path,
			"client_ip":   clientIP,
			"user_agent":  userAgent,
			"user_id":     userID,
			"timestamp":   time.Now().Format(time.RFC3339),
		})

		// Log based on status code
		if statusCode >= 500 {
			logEntry.Error("Server error")
		} else if statusCode >= 400 {
			logEntry.Warn("Client error")
		} else {
			logEntry.Info("Request processed")
		}

		// Log errors if any
		if len(c.Errors) > 0 {
			logEntry.WithFields(logrus.Fields{
				"errors": c.Errors.String(),
			}).Error("Request had errors")
		}
	}
}
