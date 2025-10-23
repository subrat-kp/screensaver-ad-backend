package models

import (
	"time"

	"gorm.io/gorm"
)

// Template represents the template metadata model
type Task struct {
	ID         uint                   `gorm:"primaryKey" json:"id"`
	TemplateID uint                   `gorm:"not null" json:"template_id"`
	Template   Template               `gorm:"foreignKey:TemplateID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"template"`
	AssetID    uint                   `gorm:"not null" json:"asset_id"`
	Asset      Asset                  `gorm:"foreignKey:AssetID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"asset"`
	Metadata   map[string]interface{} `gorm:"type:json" json:"metadata,omitempty"`
	CreatedAt  time.Time              `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time              `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt         `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName overrides the default table name for Task
func (Task) TableName() string {
	return "task_metadata"
}
