package models

import (
	"time"

	"gorm.io/gorm"
)

// Template represents the template metadata model
type Template struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null;unique" json:"name"`
	S3Key     string         `gorm:"size:500;not null;unique" json:"s3_key"`
	S3Bucket  string         `gorm:"size:255;not null" json:"s3_bucket"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName overrides the default table name for Template
func (Template) TableName() string {
	return "template_metadata"
}
