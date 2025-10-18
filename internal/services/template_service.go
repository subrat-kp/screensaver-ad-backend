package services

import (
	"screensaver-ad-backend/internal/models"
	"screensaver-ad-backend/internal/repository"
)

type TemplateService struct {
	repo *repository.TemplateRepository
}

func NewTemplateService(repo *repository.TemplateRepository) *TemplateService {
	return &TemplateService{repo: repo}
}

func (s *TemplateService) CreateTemplate(template *models.Template) error {
	return s.repo.Create(template)
}

func (s *TemplateService) ListTemplates() ([]models.Template, error) {
	return s.repo.List()
}
