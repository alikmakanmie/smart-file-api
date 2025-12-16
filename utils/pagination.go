package utils

import (
	"strconv"
	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
}

func GeneratePaginationFromRequest(c *gin.Context) *Pagination {
	// Default values
	page := 1
	limit := 10

	// Get page from query
	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	// Get limit from query (max 100)
	if l := c.Query("limit"); l != "" {
		if limitNum, err := strconv.Atoi(l); err == nil && limitNum > 0 {
			limit = limitNum
			if limit > 100 {
				limit = 100
			}
		}
	}

	return &Pagination{
		Page:  page,
		Limit: limit,
	}
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

func (p *Pagination) CalculateTotalPages() {
	if p.TotalRows > 0 {
		p.TotalPages = int((p.TotalRows + int64(p.Limit) - 1) / int64(p.Limit))
	}
}

type FileFilter struct {
	Type      string `json:"type"`      // image, audio, video, document
	Status    string `json:"status"`    // pending, processing, completed
	SortBy    string `json:"sort_by"`   // created_at, file_size, file_name
	SortOrder string `json:"sort_order"` // asc, desc
	Search    string `json:"search"`     // search by filename
}

func GenerateFilterFromRequest(c *gin.Context) *FileFilter {
	filter := &FileFilter{
		Type:      c.Query("type"),
		Status:    c.Query("status"),
		SortBy:    c.Query("sort"),
		SortOrder: c.Query("order"),
		Search:    c.Query("search"),
	}

	// Default sorting
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	// Validate sort order
	if filter.SortOrder != "asc" && filter.SortOrder != "desc" {
		filter.SortOrder = "desc"
	}

	// Validate sort by
	validSortFields := map[string]bool{
		"created_at":    true,
		"file_size":     true,
		"file_name":     true,
		"original_name": true,
		"file_type":     true,
	}
	if !validSortFields[filter.SortBy] {
		filter.SortBy = "created_at"
	}

	return filter
}
