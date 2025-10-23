package controllers

import (
	"net/http"

	"screensaver-ad-backend/internal/models"
	"screensaver-ad-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	service *services.TaskService
}

// NewTaskController creates a new task controller instance
func NewTaskController(service *services.TaskService) *TaskController {
	return &TaskController{service: service}
}

// CreateTask handles POST /tasks
// @Summary Create a new task
// @Description Create a new task if no record exists with same asset and template IDs
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body object{template_id=uint,asset_id=uint,metadata=object} true "Task object"
// @Success 201 {object} map[string]interface{} "Task created successfully"
// @Success 202 {object} map[string]interface{} "Task already exists"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /tasks [post]
func (c *TaskController) CreateTask(ctx *gin.Context) {
	var request struct {
		TemplateID uint                   `json:"template_id" binding:"required"`
		AssetID    uint                   `json:"asset_id" binding:"required"`
		Metadata   map[string]interface{} `json:"metadata,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := &models.Task{
		TemplateID: request.TemplateID,
		AssetID:    request.AssetID,
		Metadata:   request.Metadata,
	}

	created, err := c.service.CreateTaskIfNotExists(task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if created {
		ctx.JSON(http.StatusCreated, gin.H{"message": "Task created successfully"})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "Task already exists"})
	}
}
