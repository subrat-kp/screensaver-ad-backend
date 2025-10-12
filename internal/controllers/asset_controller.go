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

// CreateAsset handles POST /assets
func (c *AssetController) CreateAsset(ctx *gin.Context) {
	var asset models.Asset
	if err := ctx.ShouldBindJSON(&asset); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateAsset(&asset); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create asset"})
		return
	}

	ctx.JSON(http.StatusCreated, asset)
}

// GetAsset handles GET /assets/:id
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
