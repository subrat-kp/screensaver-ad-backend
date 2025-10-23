package repository

import (
	"screensaver-ad-backend/internal/models"

	"gorm.io/gorm"
)

// TaskRepository handles database operations for tasks
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new task repository instance
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create inserts a new task into the database
func (r *TaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

// FindByAssetAndTemplate checks if a task exists with the given asset and template IDs
func (r *TaskRepository) FindByAssetAndTemplate(assetID, templateID uint) (*models.Task, error) {
	var task models.Task
	err := r.db.Where("asset_id = ? AND template_id = ?", assetID, templateID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetByIDWithAsset retrieves a task by ID with its associated asset
func (r *TaskRepository) GetByIDWithAsset(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.Preload("Asset").First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Update updates a task record
func (r *TaskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}