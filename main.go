package main

import (
	"log"
	"net/http"

	"screensaver-ad-backend/config"
	"screensaver-ad-backend/internal/controllers"
	"screensaver-ad-backend/internal/models"
	"screensaver-ad-backend/internal/repository"
	"screensaver-ad-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	if err := config.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize S3 client (optional, won't fail if not configured)
	if err := config.InitS3(); err != nil {
		log.Printf("Warning: Failed to initialize S3: %v", err)
	}

	// Auto-migrate database models
	db := config.GetDB()
	if err := db.AutoMigrate(&models.Asset{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")

	// Initialize layers
	assetRepo := repository.NewAssetRepository(db)
	assetService := services.NewAssetService(assetRepo)
	assetController := controllers.NewAssetController(assetService)

	// Setup Gin router
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

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "connected",
		})
	})

	// Asset API routes
	api := router.Group("/api")
	{
		assets := api.Group("/assets")
		{
			assets.GET("", assetController.ListAssets)
			assets.POST("", assetController.CreateAsset)
			assets.GET("/:id", assetController.GetAsset)
			assets.PUT("/:id", assetController.UpdateAsset)
			assets.DELETE("/:id", assetController.DeleteAsset)
		}
	}

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
