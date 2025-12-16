package routes

import (
	"smart-file-api/controllers"
	"smart-file-api/middleware"
	"time"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Public health check
	router.GET("/health", controllers.HealthCheck)

	// Public routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			protected.GET("/profile", func(c *gin.Context) {
				userID := c.GetUint("user_id")
				email := c.GetString("user_email")
				
				c.JSON(200, gin.H{
					"status": "success",
					"data": gin.H{
						"user_id": userID,
						"email":   email,
					},
				})
			})

			// Monitoring endpoints (NEW)
			protected.GET("/metrics", controllers.GetMetrics)
			protected.GET("/logs", controllers.GetLogs)

			// File routes with caching
			files := protected.Group("/files")
			{
				// Statistics endpoint
				files.GET("/statistics", controllers.GetFileStatistics)
				
				// Cached endpoints with pagination & filtering (5 minutes cache)
				files.GET("/", middleware.CacheMiddleware(5*time.Minute), controllers.GetUserFiles)
				files.GET("/deleted", middleware.CacheMiddleware(5*time.Minute), controllers.GetDeletedFiles)
				files.GET("/:id", middleware.CacheMiddleware(5*time.Minute), controllers.GetFileDetail)
				
				// Non-cached endpoints
				files.POST("/upload", controllers.UploadFile)
				files.POST("/:id/restore", controllers.RestoreFile)
				files.DELETE("/:id", controllers.DeleteFile)
				files.DELETE("/:id/permanent", controllers.HardDeleteFile)
			}
		}
	}
}
