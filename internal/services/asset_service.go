package services

import (
	"screensaver-ad-backend/internal/models"
	"screensaver-ad-backend/internal/repository"
)

// AssetService handles business logic for assets
type AssetService struct {
	repo *repository.AssetRepository
}

// NewAssetService creates a new asset service instance
func NewAssetService(repo *repository.AssetRepository) *AssetService {
	return &AssetService{repo: repo}
}

// CreateAsset creates a new asset
func (s *AssetService) CreateAsset(asset *models.Asset) error {
	// Add business logic validation here if needed
	return s.repo.Create(asset)
}

// GetAssetByID retrieves an asset by its ID
func (s *AssetService) GetAssetByID(id uint) (*models.Asset, error) {
	return s.repo.GetByID(id)
}

// GetAllAssets retrieves all assets with pagination
func (s *AssetService) GetAllAssets(limit, offset int) ([]models.Asset, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}
	if limit > 100 {
		limit = 100 // max limit
	}
	return s.repo.GetAll(limit, offset)
}

// UpdateAsset updates an existing asset
func (s *AssetService) UpdateAsset(asset *models.Asset) error {
	// Add business logic validation here if needed
	return s.repo.Update(asset)
}

// DeleteAsset deletes an asset by its ID
func (s *AssetService) DeleteAsset(id uint) error {
	return s.repo.Delete(id)
}

// GetAssetCount returns the total number of assets
func (s *AssetService) GetAssetCount() (int64, error) {
	return s.repo.Count()
}
