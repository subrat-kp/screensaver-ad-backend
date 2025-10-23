package services

import (
	"fmt"
	"screensaver-ad-backend/internal/models"
	"screensaver-ad-backend/internal/repository"

	"gorm.io/gorm"
)

// TaskService handles business logic for tasks
type TaskService struct {
	repo      *repository.TaskRepository
	assetRepo *repository.AssetRepository
}

// NewTaskService creates a new task service instance
func NewTaskService(repo *repository.TaskRepository, assetRepo *repository.AssetRepository) *TaskService {
	return &TaskService{
		repo:      repo,
		assetRepo: assetRepo,
	}
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

// UpdateTaskMetadata updates task metadata and corresponding asset based on payload
func (s *TaskService) UpdateTaskMetadata(payload map[string]interface{}) error {
	// Extract task_id from payload
	taskIDFloat, ok := payload["task_id"].(float64)
	if !ok {
		return fmt.Errorf("task_id not found or invalid in payload")
	}
	taskID := uint(taskIDFloat)

	// Get task with asset relation
	task, err := s.repo.GetByIDWithAsset(taskID)
	if err != nil {
		return err
	}

	// Update task metadata
	task.Metadata = payload
	if err := s.repo.Update(task); err != nil {
		return err
	}

	// Extract s3_key and update asset if present
	if s3Key, exists := payload["s3_key"].(string); exists {
		return s.assetRepo.UpdateOutputS3Key(task.AssetID, s3Key)
	}

	return nil
}
