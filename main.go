package main

import (
	"log"
	"os"
	"smart-file-api/config"
	"smart-file-api/models"
	"smart-file-api/routes"
	"smart-file-api/middleware"
	
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	
	_ "smart-file-api/docs"
)

// @title Smart File API
// @version 1.0
// @description A smart file processing API with JWT authentication, caching, and background processing
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@smartfileapi.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

func main() {
	// Initialize logger
	config.InitLogger()
	config.Log.Info("Starting Smart File API...")

	// Create uploads directory if not exists
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		log.Fatal("Failed to create uploads directory:", err)
	}
	config.Log.Info("Uploads directory ready")

	// Connect to database
	config.ConnectDatabase()
	
	// Connect to Redis
	config.ConnectRedis()
	
	// Auto migrate database schema
	config.DB.AutoMigrate(&models.User{}, &models.File{})
	config.Log.Info("Database migration completed")
	
	// Set Gin mode
	gin.SetMode(gin.DebugMode)
	
	// Create router
	router := gin.Default()
	
	// Add logging middleware
	router.Use(middleware.LoggerMiddleware())
	
	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Setup routes
	routes.SetupRoutes(router)
	
	// Start server
	log.Println("🚀 Server running on http://localhost:8080")
	log.Println("📖 Swagger documentation: http://localhost:8080/swagger/index.html")
	config.Log.Info("Server started successfully on port 8080")
	
	if err := router.Run(":8080"); err != nil {
		config.Log.Fatal("Failed to start server:", err)
	}
}
