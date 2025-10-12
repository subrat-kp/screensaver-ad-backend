package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Creative struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	FileURL     string    `json:"file_url" db:"file_url"`
	ContentType string    `json:"content_type" db:"content_type"`
	UploadedAt  time.Time `json:"uploaded_at" db:"uploaded_at"`
	Status      string    `json:"status" db:"status"`
}

type AssetStatus struct {
	Status  string `json:"status"`
	FileURL string `json:"file_url,omitempty"`
	Message string `json:"message,omitempty"`
}

// Global database connection
var db *sql.DB

// In-memory storage for backward compatibility (optional)
var creatives []Creative

func main() {
	// Initialize database connection
	initDB()
	defer db.Close()

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
	router.GET("/api/assets/:id/status", getAssetStatus)

	log.Println("Server starting on :8080")
	router.Run(":8080")
}

// Initialize database connection
func initDB() {
	// Get database connection parameters from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "screensaver_ad")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v. Running in-memory mode.", err)
		return
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Warning: Database ping failed: %v. Running in-memory mode.", err)
		db = nil
		return
	}

	log.Println("Database connected successfully")

	// Create table if not exists
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS asset_metadata (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(500) NOT NULL,
			file_url TEXT NOT NULL,
			content_type VARCHAR(100),
			uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			status VARCHAR(50) NOT NULL DEFAULT 'processing'
		);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Printf("Warning: Failed to create table: %v", err)
	}
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Health check endpoint
func healthCheck(c *gin.Context) {
	dbStatus := "disconnected"
	if db != nil {
		if err := db.Ping(); err == nil {
			dbStatus = "connected"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"database":  dbStatus,
	})
}

// Upload creative with database support
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
		Status:      "processing",
	}

	// Store in database if available, otherwise in memory
	if db != nil {
		query := `INSERT INTO asset_metadata (id, name, file_url, content_type, uploaded_at, status) 
				  VALUES ($1, $2, $3, $4, $5, $6)`
		_, err := db.Exec(query, creative.ID, creative.Name, creative.FileURL,
			creative.ContentType, creative.UploadedAt, creative.Status)
		if err != nil {
			log.Printf("Failed to insert into database: %v. Falling back to in-memory.", err)
			creatives = append(creatives, creative)
		}
	} else {
		// Fallback to in-memory storage
		creatives = append(creatives, creative)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"creative": creative,
	})
}

// List all assets
func listAssets(c *gin.Context) {
	var assets []Creative

	if db != nil {
		rows, err := db.Query(`SELECT id, name, file_url, content_type, uploaded_at, status 
							   FROM asset_metadata ORDER BY uploaded_at DESC`)
		if err != nil {
			log.Printf("Failed to query database: %v", err)
			assets = creatives
		} else {
			defer rows.Close()
			for rows.Next() {
				var asset Creative
				err := rows.Scan(&asset.ID, &asset.Name, &asset.FileURL,
					&asset.ContentType, &asset.UploadedAt, &asset.Status)
				if err != nil {
					log.Printf("Failed to scan row: %v", err)
					continue
				}
				assets = append(assets, asset)
			}
		}
	} else {
		assets = creatives
	}

	c.JSON(http.StatusOK, gin.H{
		"assets": assets,
		"count":  len(assets),
	})
}

// Get asset status endpoint
func getAssetStatus(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Asset ID is required"})
		return
	}

	var creative Creative
	found := false

	// Try to get from database first
	if db != nil {
		query := `SELECT id, name, file_url, content_type, uploaded_at, status 
				  FROM asset_metadata WHERE id = $1`
		err := db.QueryRow(query, id).Scan(&creative.ID, &creative.Name, &creative.FileURL,
			&creative.ContentType, &creative.UploadedAt, &creative.Status)
		if err == nil {
			found = true
		} else if err != sql.ErrNoRows {
			log.Printf("Database query error: %v", err)
		}
	}

	// Fallback to in-memory storage
	if !found {
		for _, asset := range creatives {
			if asset.ID == id {
				creative = asset
				found = true
				break
			}
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}

	// If status is processing, check S3 for processed output
	if creative.Status == "processing" {
		if checkS3ProcessedOutput(creative.ID) {
			// Update status to processed
			creative.Status = "processed"
			if db != nil {
				_, err := db.Exec(`UPDATE asset_metadata SET status = $1 WHERE id = $2`,
					"processed", id)
				if err != nil {
					log.Printf("Failed to update status: %v", err)
				}
			} else {
				// Update in-memory storage
				for i := range creatives {
					if creatives[i].ID == id {
						creatives[i].Status = "processed"
						break
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, AssetStatus{
		Status:  creative.Status,
		FileURL: creative.FileURL,
		Message: fmt.Sprintf("Asset %s is %s", creative.Name, creative.Status),
	})
}

// Check if file exists in S3 processed output folder
func checkS3ProcessedOutput(assetID string) bool {
	// Get S3 configuration from environment
	bucketName := getEnv("S3_BUCKET", "")
	processedFolder := getEnv("S3_PROCESSED_FOLDER", "processed/")

	if bucketName == "" {
		log.Println("S3_BUCKET not configured, assuming not processed")
		return false
	}

	// Create S3 session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(getEnv("AWS_REGION", "us-east-1")),
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		return false
	}

	svc := s3.New(sess)

	// Check if processed file exists
	processedKey := processedFolder + assetID
	_, err = svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(processedKey),
	})

	return err == nil
}

// Helper function to generate unique IDs
func generateID() string {
	return time.Now().Format("20060102150405")
}
