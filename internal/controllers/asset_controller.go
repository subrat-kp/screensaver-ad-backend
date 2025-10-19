package controllers

import (
	"net/http"
	"strconv"

	"screensaver-ad-backend/internal/models"
	"screensaver-ad-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// AssetController handles HTTP requests for assets
type AssetController struct {
	service *services.AssetService
}

// NewAssetController creates a new asset controller instance
func NewAssetController(service *services.AssetService) *AssetController {
	return &AssetController{service: service}
}

// CreateAsset handles POST /assets with file upload
// @Summary Create a new asset
// @Description Upload a new asset file with metadata
// @Tags assets
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Asset file"
// @Param name formData string false "Asset name (defaults to filename)"
// @Success 201 {object} map[string]interface{} "Asset created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /assets [post]
func (c *AssetController) CreateAsset(ctx *gin.Context) {
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil { // 32 MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	// Get file from form
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	// Get name from form (optional)
	name := ctx.Request.FormValue("name")
	if name == "" {
		// Use original filename if name not provided
		name = fileHeader.Filename
	}

	// Create asset with file upload
	asset, err := c.service.CreateAssetWithUpload(file, fileHeader, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Asset created successfully",
		"asset":   asset,
	})
}

// GetAsset handles GET /assets/:id
// @Summary Get asset by ID
// @Description Retrieve a specific asset by its ID
// @Tags assets
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {object} models.Asset
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 404 {object} map[string]interface{} "Asset not found"
// @Router /assets/{id} [get]
func (c *AssetController) GetAsset(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	asset, err := c.service.GetAssetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}

	ctx.JSON(http.StatusOK, asset)
}

// ListAssets handles GET /assets
// @Summary List all assets
// @Description Get a paginated list of all assets
// @Tags assets
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} map[string]interface{} "List of assets with pagination info"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /assets [get]
func (c *AssetController) ListAssets(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	assets, err := c.service.GetAllAssets(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assets"})
		return
	}

	count, _ := c.service.GetAssetCount()
	ctx.JSON(http.StatusOK, gin.H{
		"assets": assets,
		"total":  count,
		"limit":  limit,
		"offset": offset,
	})
}

// UpdateAsset handles PUT /assets/:id
// @Summary Update an asset
// @Description Update an existing asset's information
// @Tags assets
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Param asset body models.Asset true "Asset object"
// @Success 200 {object} models.Asset
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /assets/{id} [put]
func (c *AssetController) UpdateAsset(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var asset models.Asset
	if err := ctx.ShouldBindJSON(&asset); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset.ID = uint(id)
	if err := c.service.UpdateAsset(&asset); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update asset"})
		return
	}

	ctx.JSON(http.StatusOK, asset)
}

// DeleteAsset handles DELETE /assets/:id
// @Summary Delete an asset
// @Description Delete an asset by its ID
// @Tags assets
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {object} map[string]interface{} "Asset deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /assets/{id} [delete]
func (c *AssetController) DeleteAsset(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.DeleteAsset(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete asset"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Asset deleted successfully"})
}

// UpdateAssetStatus handles PATCH /assets/:id/status
// @Summary Update asset status and the output url
// @Description Update the status of an asset (uploaded, processed, upload_failed, process_failed) and the output url
// @Tags assets
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Param status body object{status=string,output_s3_key=string} true "Status update request"
// @Success 200 {object} map[string]interface{} "Asset status updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /assets/{id}/status [patch]
func (c *AssetController) UpdateAssetStatus(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var request struct {
		Status      models.AssetStatus `json:"status" binding:"required"`
		OutputS3Key *string            `json:"output_s3_key,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	if err := c.service.UpdateAssetStatus(uint(id), request.Status, request.OutputS3Key); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch updated asset
	asset, err := c.service.GetAssetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated asset"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset status updated successfully",
		"asset":   asset,
	})
}

// GetAssetURL handles GET /assets/:id/url
// @Summary Get asset URLs
// @Description Generate presigned URLs for both input and output asset files
// @Tags assets
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Param expiration query int false "URL expiration time in minutes" default(60)
// @Success 200 {object} map[string]interface{} "Presigned URLs for input and output files"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /assets/{id}/url [get]
func (c *AssetController) GetAssetURL(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Get expiration from query parameter (default 60 minutes)
	expiration, _ := strconv.Atoi(ctx.DefaultQuery("expiration", "60"))

	// Generate presigned URLs for both input and output files
	urls, err := c.service.GetAssetURLs(uint(id), expiration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"input_url":  urls.InputURL,
		"output_url": urls.OutputURL,
		"expires_in": expiration,
	})
}
