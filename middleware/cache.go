package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"smart-file-api/config"
	"time"

	"github.com/gin-gonic/gin"
)

func CacheMiddleware(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip caching if Redis not available
		if config.RedisClient == nil {
			c.Next()
			return
		}

		// Skip caching for non-GET requests
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		// Generate cache key based on URL and user
		userID := c.GetUint("user_id")
		cacheKey := generateCacheKey(c.Request.URL.Path, userID)

		// Try to get from cache
		cachedResponse, err := config.GetCache(cacheKey)
		if err == nil && cachedResponse != "" {
			// Cache hit
			c.Header("X-Cache", "HIT")
			c.Data(http.StatusOK, "application/json", []byte(cachedResponse))
			c.Abort()
			return
		}

		// Cache miss - continue to handler
		c.Header("X-Cache", "MISS")

		// Capture response
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           []byte{},
		}
		c.Writer = writer

		c.Next()

		// Cache the response if status is 200
		if c.Writer.Status() == http.StatusOK {
			config.SetCache(cacheKey, string(writer.body), duration)
		}
	}
}

// Custom response writer to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body []byte
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

func generateCacheKey(path string, userID uint) string {
	data := fmt.Sprintf("%s:user:%d", path, userID)
	hash := md5.Sum([]byte(data))
	return "cache:" + hex.EncodeToString(hash[:])
}
