package services

import (
	"screensaver-ad-backend/internal/models"
	"screensaver-ad-backend/internal/repository"

	"gorm.io/gorm"
)

// TaskService handles business logic for tasks
type TaskService struct {
	repo *repository.TaskRepository
}

// NewTaskService creates a new task service instance
func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTaskIfNotExists creates a task if no record exists with same asset and template IDs
func (s *TaskService) CreateTaskIfNotExists(task *models.Task) (bool, error) {
	// Check if task already exists
	_, err := s.repo.FindByAssetAndTemplate(task.AssetID, task.TemplateID)
	if err == nil {
		// Task already exists
		return false, nil
	}
	if err != gorm.ErrRecordNotFound {
		// Database error
		return false, err
	}

	// Task doesn't exist, create it
	if err := s.repo.Create(task); err != nil {
		return false, err
	}
	return true, nil
}