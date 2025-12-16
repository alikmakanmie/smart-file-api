package models

import (
	"time"
	"gorm.io/gorm"
)

type File struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `json:"user_id"`
	FileName     string         `json:"file_name"`
	OriginalName string         `json:"original_name"`
	FilePath     string         `json:"file_path"`
	FileSize     int64          `json:"file_size"`
	FileType     string         `json:"file_type"`
	ProcessedAt  *time.Time     `json:"processed_at"`
	Status       string         `json:"status"` // pending, processing, completed, failed
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
