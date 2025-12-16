package controllers

import (
	"net/http"
	"os"
	"runtime"
	"smart-file-api/config"
	"smart-file-api/utils"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// HealthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags Monitoring
// @Produce json
// @Success 200 {object} map[string]interface{} "API is healthy"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
		"uptime":  time.Since(startTime).String(),
	})
}

// GetMetrics godoc
// @Summary Get system metrics
// @Description Get application metrics (uptime, memory usage, goroutines, etc.)
// @Tags Monitoring
// @Produce json
// @Success 200 {object} map[string]interface{} "Metrics retrieved successfully"
// @Security BearerAuth
// @Router /metrics [get]
func GetMetrics(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Database stats
	var totalUsers int64
	var totalFiles int64
	config.DB.Model(&struct{ ID uint }{}).Table("users").Count(&totalUsers)
	config.DB.Model(&struct{ ID uint }{}).Table("files").Count(&totalFiles)

	// Redis stats
	redisConnected := config.RedisClient != nil
	var cacheKeys int64 = 0
	if redisConnected {
		keys, _ := config.RedisClient.Keys(config.Ctx, "cache:*").Result()
		cacheKeys = int64(len(keys))
	}

	// File system stats (uploads folder)
	uploadsSize := getDirSize("uploads")

	utils.SuccessResponse(c, http.StatusOK, "Metrics retrieved successfully", gin.H{
		"system": gin.H{
			"uptime":           time.Since(startTime).String(),
			"goroutines":       runtime.NumGoroutine(),
			"memory_alloc_mb":  float64(m.Alloc) / 1024 / 1024,
			"memory_total_mb":  float64(m.TotalAlloc) / 1024 / 1024,
			"memory_sys_mb":    float64(m.Sys) / 1024 / 1024,
			"gc_runs":          m.NumGC,
		},
		"database": gin.H{
			"total_users": totalUsers,
			"total_files": totalFiles,
		},
		"cache": gin.H{
			"redis_connected": redisConnected,
			"cache_keys":      cacheKeys,
		},
		"storage": gin.H{
			"uploads_size_mb": float64(uploadsSize) / 1024 / 1024,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// Helper function to get directory size
func getDirSize(path string) int64 {
	var size int64
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err == nil {
				size += info.Size()
			}
		}
	}
	return size
}

// GetLogs godoc
// @Summary Get recent logs
// @Description Get recent application logs (last 100 lines)
// @Tags Monitoring
// @Produce json
// @Param level query string false "Filter by log level" Enums(info, warning, error)
// @Success 200 {object} map[string]interface{} "Logs retrieved successfully"
// @Security BearerAuth
// @Router /logs [get]
func GetLogs(c *gin.Context) {
	// Read last 100 lines from log file
	file, err := os.Open("app.log")
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to read logs")
		return
	}
	defer file.Close()

	// Get file info
	stat, _ := file.Stat()
	fileSize := stat.Size()

	// Read last 10KB of logs (approximately last 100 lines)
	bufferSize := int64(10240)
	if fileSize < bufferSize {
		bufferSize = fileSize
	}

	buffer := make([]byte, bufferSize)
	_, err = file.ReadAt(buffer, fileSize-bufferSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to read logs")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logs retrieved successfully", gin.H{
		"logs":      string(buffer),
		"file_size": fileSize,
	})
}
