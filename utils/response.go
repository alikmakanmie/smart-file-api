package utils

import (
	"smart-file-api/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	// Log error
	config.Log.WithFields(logrus.Fields{
		"status_code": statusCode,
		"message":     message,
		"path":        c.Request.URL.Path,
		"method":      c.Request.Method,
		"user_id":     c.GetUint("user_id"),
	}).Error("Error response sent")

	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
	})
}
