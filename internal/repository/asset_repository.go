package repository

import (
	"screensaver-ad-backend/internal/models"

	"gorm.io/gorm"
)

// AssetRepository handles database operations for assets
type AssetRepository struct {
	db *gorm.DB
}

// NewAssetRepository creates a new asset repository instance
func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

// Create inserts a new asset into the database
func (r *AssetRepository) Create(asset *models.Asset) error {
	return r.db.Create(asset).Error
}

// GetByID retrieves an asset by its ID
func (r *AssetRepository) GetByID(id uint) (*models.Asset, error) {
	var asset models.Asset
	err := r.db.First(&asset, id).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// GetAll retrieves all assets with pagination
func (r *AssetRepository) GetAll(limit, offset int) ([]models.Asset, error) {
	var assets []models.Asset
	err := r.db.Limit(limit).Offset(offset).Find(&assets).Error
	return assets, err
}

// Update updates an existing asset
func (r *AssetRepository) Update(asset *models.Asset) error {
	return r.db.Save(asset).Error
}

// Delete soft deletes an asset by its ID
func (r *AssetRepository) Delete(id uint) error {
	return r.db.Delete(&models.Asset{}, id).Error
}

// Count returns the total number of assets
func (r *AssetRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Asset{}).Count(&count).Error
	return count, err
}
