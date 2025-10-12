package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Creative struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	FileURL     string    `json:"file_url"`
	ContentType string    `json:"content_type"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

// In-memory storage for demo purposes
var creatives []Creative

func main() {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Routes
	router.GET("/health", healthCheck)
	router.POST("/api/upload", uploadCreative)
	router.GET("/api/assets", listAssets)

	// Start server
	router.Run(":8080")
}

// Health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"timestamp": time.Now().Unix(),
	})
}

// Upload creative (S3 placeholder)
func uploadCreative(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// TODO: Implement actual S3 upload
	// For now, create a placeholder URL
	placeholderURL := "https://s3.amazonaws.com/placeholder/" + header.Filename

	// Create creative record
	creative := Creative{
		ID:          generateID(),
		Name:        header.Filename,
		FileURL:     placeholderURL,
		ContentType: header.Header.Get("Content-Type"),
		UploadedAt:  time.Now(),
	}

	// Store in memory
	creatives = append(creatives, creative)

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully (placeholder)",
		"creative": creative,
	})
}

// List all assets
func listAssets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"assets": creatives,
		"count":  len(creatives),
	})
}

// Helper function to generate unique IDs
func generateID() string {
	return time.Now().Format("20060102150405")
}
