package services

import (
	"fmt"
	"mime/multipart"
	"time"

	"screensaver-ad-backend/config"
	"screensaver-ad-backend/internal/models"
	"screensaver-ad-backend/internal/repository"
)

// AssetService handles business logic for assets
type AssetService struct {
	repo      *repository.AssetRepository
	s3Service *S3Service
}

// NewAssetService creates a new asset service instance
func NewAssetService(repo *repository.AssetRepository) *AssetService {
	return &AssetService{
		repo:      repo,
		s3Service: NewS3Service(),
	}
}

// CreateAsset creates a new asset (without file upload)
func (s *AssetService) CreateAsset(asset *models.Asset) error {
	// Add business logic validation here if needed
	return s.repo.Create(asset)
}

// CreateAssetWithUpload creates a new asset with file upload to S3
func (s *AssetService) CreateAssetWithUpload(file multipart.File, fileHeader *multipart.FileHeader, name string) (*models.Asset, error) {
	// Validate file
	if fileHeader.Size == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	// Check content type (images and videos only)
	contentType := fileHeader.Header.Get("Content-Type")
	if !isValidContentType(contentType) {
		return nil, fmt.Errorf("invalid file type: only images and videos are allowed")
	}

	// Upload to S3
	s3Key, err := s.s3Service.UploadFileToS3(file, fileHeader, name)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Create asset record with initial status as "uploaded"
	asset := &models.Asset{
		FileName:    name,
		FileSize:    fileHeader.Size,
		ContentType: contentType,
		S3Key:       s3Key,
		S3Bucket:    config.GetS3Bucket(),
		Status:      models.AssetStatusUploaded,
	}

	if err := s.repo.Create(asset); err != nil {
		// Rollback: delete file from S3 if database insert fails
		_ = s.s3Service.DeleteFileFromS3(s3Key)
		return nil, fmt.Errorf("failed to create asset record: %w", err)
	}

	return asset, nil
}

// UpdateAssetStatus updates the status of an asset
func (s *AssetService) UpdateAssetStatus(id uint, status models.AssetStatus) error {
	// Validate status
	if status != models.AssetStatusUploaded && status != models.AssetStatusProcessed {
		return fmt.Errorf("invalid status: must be 'uploaded' or 'processed'")
	}

	// Check if asset exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("asset not found")
	}

	// Update status
	return s.repo.UpdateStatus(id, status)
}

// isValidContentType checks if the content type is valid (image or video)
func isValidContentType(contentType string) bool {
	validTypes := []string{
		"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp",
		"video/mp4", "video/mpeg", "video/quicktime", "video/x-msvideo", "video/webm",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
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

// GetAssetURL generates a presigned URL for accessing an asset
func (s *AssetService) GetAssetURL(id uint, expirationMinutes int) (string, error) {
	// Get asset
	asset, err := s.repo.GetByID(id)
	if err != nil {
		return "", fmt.Errorf("asset not found")
	}

	// Default expiration to 60 minutes if not specified
	if expirationMinutes <= 0 {
		expirationMinutes = 60
	}

	// Generate presigned URL
	url, err := s.s3Service.GetFileURL(asset.S3Key, time.Duration(expirationMinutes)*time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to generate URL: %w", err)
	}

	return url, nil
}
