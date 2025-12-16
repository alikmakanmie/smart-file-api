package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"smart-file-api/config"
	"smart-file-api/models"
	"smart-file-api/utils"
	"time"

	"github.com/gin-gonic/gin"
)

const MaxFileSize = 10 << 20 // 10 MB

// UploadFile godoc
// @Summary Upload file
// @Description Upload a file (max 10MB)
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 201 {object} map[string]interface{} "File uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid file"
// @Security BearerAuth
// @Router /files/upload [post]
func UploadFile(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "File is required")
		return
	}

	// Validate file size
	if file.Size > MaxFileSize {
		utils.ErrorResponse(c, http.StatusBadRequest, "File size exceeds 10MB limit")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join("uploads", newFilename)

	// Save file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Detect file type
	fileType := detectFileType(ext)

	// Save to database
	fileRecord := models.File{
		UserID:       userID,
		FileName:     newFilename,
		OriginalName: file.Filename,
		FilePath:     filePath,
		FileSize:     file.Size,
		FileType:     fileType,
		Status:       "pending",
	}

	if err := config.DB.Create(&fileRecord).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save file record")
		return
	}

	// Invalidate cache for user's file list
	config.DeleteCachePattern("cache:*")

	// Start processing in background (async)
	go processFile(&fileRecord)

	utils.SuccessResponse(c, http.StatusCreated, "File uploaded successfully", gin.H{
		"file": fileRecord,
	})
}

// GetUserFiles godoc
// @Summary Get all user files with pagination and filtering
// @Description Get list of all files with pagination, filtering, sorting, and search
// @Tags Files
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page (max 100)" default(10)
// @Param type query string false "Filter by file type" Enums(image, audio, video, document, other)
// @Param status query string false "Filter by status" Enums(pending, processing, completed, failed)
// @Param sort query string false "Sort by field" Enums(created_at, file_size, file_name, original_name, file_type) default(created_at)
// @Param order query string false "Sort order" Enums(asc, desc) default(desc)
// @Param search query string false "Search by filename"
// @Success 200 {object} map[string]interface{} "Files retrieved successfully"
// @Security BearerAuth
// @Router /files/ [get]
func GetUserFiles(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Generate pagination and filter
	pagination := utils.GeneratePaginationFromRequest(c)
	filter := utils.GenerateFilterFromRequest(c)

	var files []models.File
	query := config.DB.Where("user_id = ?", userID)

	// Apply filters
	if filter.Type != "" {
		query = query.Where("file_type = ?", filter.Type)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("file_name LIKE ? OR original_name LIKE ?", searchTerm, searchTerm)
	}

	// Get total count
	query.Model(&models.File{}).Count(&pagination.TotalRows)
	pagination.CalculateTotalPages()

	// Apply sorting and pagination
	orderClause := filter.SortBy + " " + filter.SortOrder
	if err := query.Order(orderClause).
		Limit(pagination.Limit).
		Offset(pagination.GetOffset()).
		Find(&files).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch files")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Files retrieved successfully", gin.H{
		"files":      files,
		"pagination": pagination,
		"filter":     filter,
	})
}


// GetFileDetail godoc
// @Summary Get file detail
// @Description Get detailed information about a specific file
// @Tags Files
// @Produce json
// @Param id path int true "File ID"
// @Success 200 {object} map[string]interface{} "File retrieved successfully"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Security BearerAuth
// @Router /files/{id} [get]
func GetFileDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID := c.Param("id")

	var file models.File
	if err := config.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "File not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "File retrieved successfully", gin.H{
		"file": file,
	})
}

// DeleteFile godoc
// @Summary Delete file (soft delete)
// @Description Soft delete a file (can be restored)
// @Tags Files
// @Produce json
// @Param id path int true "File ID"
// @Success 200 {object} map[string]interface{} "File deleted successfully"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Security BearerAuth
// @Router /files/{id} [delete]
func DeleteFile(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID := c.Param("id")

	var file models.File
	if err := config.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "File not found")
		return
	}

	// Delete physical file from disk
	if err := os.Remove(file.FilePath); err != nil {
		fmt.Printf("Warning: Failed to delete physical file: %v\n", err)
	}

	// Soft delete from database
	if err := config.DB.Delete(&file).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete file")
		return
	}

	// Invalidate cache
	config.DeleteCachePattern("cache:*")

	utils.SuccessResponse(c, http.StatusOK, "File deleted successfully", nil)
}

func GetDeletedFiles(c *gin.Context) {
	userID := c.GetUint("user_id")

	var files []models.File
	if err := config.DB.Unscoped().Where("user_id = ? AND deleted_at IS NOT NULL", userID).Order("deleted_at DESC").Find(&files).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch deleted files")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Deleted files retrieved successfully", gin.H{
		"files": files,
		"total": len(files),
	})
}

func RestoreFile(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID := c.Param("id")

	var file models.File
	if err := config.DB.Unscoped().Where("id = ? AND user_id = ? AND deleted_at IS NOT NULL", fileID, userID).First(&file).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Deleted file not found")
		return
	}

	if err := config.DB.Model(&file).Update("deleted_at", nil).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to restore file")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "File restored successfully", gin.H{
		"file": file,
	})
}

func HardDeleteFile(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID := c.Param("id")

	var file models.File
	if err := config.DB.Unscoped().Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "File not found")
		return
	}

	if err := os.Remove(file.FilePath); err != nil {
		fmt.Printf("Warning: Failed to delete physical file: %v\n", err)
	}

	if err := config.DB.Unscoped().Delete(&file).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete file")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "File permanently deleted", nil)
}

// Helper function to detect file type
func detectFileType(ext string) string {
	imageExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true}
	audioExts := map[string]bool{".mp3": true, ".wav": true, ".flac": true, ".m4a": true, ".ogg": true}
	videoExts := map[string]bool{".mp4": true, ".avi": true, ".mkv": true, ".mov": true}
	documentExts := map[string]bool{".pdf": true, ".doc": true, ".docx": true, ".txt": true}

	if imageExts[ext] {
		return "image"
	} else if audioExts[ext] {
		return "audio"
	} else if videoExts[ext] {
		return "video"
	} else if documentExts[ext] {
		return "document"
	}
	return "other"
}

// Background processing function
func processFile(file *models.File) {
	config.DB.Model(file).Update("status", "processing")

	time.Sleep(3 * time.Second)

	now := time.Now()
	config.DB.Model(file).Updates(map[string]interface{}{
		"status":       "completed",
		"processed_at": &now,
	})
}

// GetFileStatistics godoc
// @Summary Get file statistics
// @Description Get statistics about user's files (total count, storage used, files by type)
// @Tags Files
// @Produce json
// @Success 200 {object} map[string]interface{} "Statistics retrieved successfully"
// @Security BearerAuth
// @Router /files/statistics [get]
func GetFileStatistics(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Total files count
	var totalFiles int64
	config.DB.Model(&models.File{}).Where("user_id = ?", userID).Count(&totalFiles)

	// Total storage used
	var totalSize int64
	config.DB.Model(&models.File{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(file_size), 0)").
		Scan(&totalSize)

	// Files by type
	type FileTypeCount struct {
		FileType string `json:"file_type"`
		Count    int64  `json:"count"`
	}
	var filesByType []FileTypeCount
	config.DB.Model(&models.File{}).
		Select("file_type, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("file_type").
		Scan(&filesByType)

	// Files by status
	type FileStatusCount struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var filesByStatus []FileStatusCount
	config.DB.Model(&models.File{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&filesByStatus)

	// Recent files (last 7 days)
	var recentFilesCount int64
	config.DB.Model(&models.File{}).
		Where("user_id = ? AND created_at >= datetime('now', '-7 days')", userID).
		Count(&recentFilesCount)

	utils.SuccessResponse(c, http.StatusOK, "Statistics retrieved successfully", gin.H{
		"total_files":        totalFiles,
		"total_storage":      totalSize,
		"total_storage_mb":   float64(totalSize) / (1024 * 1024),
		"files_by_type":      filesByType,
		"files_by_status":    filesByStatus,
		"recent_files_7d":    recentFilesCount,
	})
}
