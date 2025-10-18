package repository

import (
	"screensaver-ad-backend/internal/models"

	"gorm.io/gorm"
)

type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) Create(template *models.Template) error {
	return r.db.Create(template).Error
}

func (r *TemplateRepository) List() ([]models.Template, error) {
	var templates []models.Template
	err := r.db.Find(&templates).Error
	return templates, err
}
