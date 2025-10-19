package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"screensaver-ad-backend/config"
	_ "screensaver-ad-backend/docs" // Import generated docs
	"screensaver-ad-backend/internal/controllers"
	"screensaver-ad-backend/internal/models"
	"screensaver-ad-backend/internal/repository"
	"screensaver-ad-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// @title Screensaver Ad Backend API
// @version 1.0
// @description API server for managing screensaver advertisements, assets, and templates
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/subrat-kp/screensaver-ad-backend
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http https

func main() {
	// Load .env only in development environment
	env := os.Getenv("GO_ENV")
	env = strings.ToLower(env)
	if env == "development" || env == "dev" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		} else {
			log.Println(".env file loaded for development environment")
		}
	}

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
	if err := db.AutoMigrate(models.Models()...); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")

	// Initialize layers
	assetRepo := repository.NewAssetRepository(db)
	assetService := services.NewAssetService(assetRepo)
	assetController := controllers.NewAssetController(assetService)

	s3Service := services.NewS3Service()

	templateRepo := repository.NewTemplateRepository(db)
	templateService := services.NewTemplateService(templateRepo)
	templateController := controllers.NewTemplateController(templateService, s3Service)

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
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
			assets.GET("/:id/url", assetController.GetAssetURL)
			assets.PUT("/:id", assetController.UpdateAsset)
			assets.PATCH("/:id/status", assetController.UpdateAssetStatus)
			assets.DELETE("/:id", assetController.DeleteAsset)
		}

		templates := api.Group("/templates")
		{
			templates.GET("", templateController.ListTemplates)
			templates.POST("", templateController.UploadTemplate)
		}
	}

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server starting on :8080")
	log.Println("Swagger documentation available at http://localhost:8080/swagger/index.html")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
