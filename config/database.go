package config

import (
	"gorm.io/gorm"
	"github.com/glebarez/sqlite"
	"log"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("smart-file-api.db"), &gorm.Config{})
	
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	DB = database
	log.Println("✅ Database connected successfully")
}
