package models

import (
	"time"

	"gorm.io/gorm"
)

// AssetStatus represents the processing status of an asset
type AssetStatus string

const (
	AssetStatusUploaded  AssetStatus = "uploaded"
	AssetStatusProcessed AssetStatus = "processed"
)

// Asset represents the asset metadata model
type Asset struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	FileName    string         `gorm:"size:255;not null" json:"file_name"`
	FileSize    int64          `gorm:"not null" json:"file_size"`
	ContentType string         `gorm:"size:100;not null" json:"content_type"`
	S3Key       string         `gorm:"size:500;not null;unique" json:"s3_key"`
	OutputS3Key *string        `gorm:"size:500" json:"output_s3_key,omitempty"`
	S3Bucket    string         `gorm:"size:255;not null" json:"s3_bucket"`
	Status      AssetStatus    `gorm:"size:50;not null;default:'uploaded'" json:"status"`
	UploadedAt  time.Time      `gorm:"autoCreateTime" json:"uploaded_at"`
	ProcessedAt *time.Time     `json:"processed_at,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName overrides the default table name
func (Asset) TableName() string {
	return "asset_metadata"
}
