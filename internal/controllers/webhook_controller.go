package controllers

import (
	"net/http"

	"screensaver-ad-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type WebhookController struct {
	taskService  *services.TaskService
	assetService *services.AssetService
}

func NewWebhookController(taskService *services.TaskService, assetService *services.AssetService) *WebhookController {
	return &WebhookController{
		taskService:  taskService,
		assetService: assetService,
	}
}

// HandleWebhook handles POST /webhook
// @Summary Handle webhook events
// @Description Process webhook events with flexible payload structure
// @Tags webhook
// @Accept json
// @Produce json
// @Param event body object{event_type=string,payload=object} true "Webhook event with payload"
// @Success 200 {object} map[string]interface{} "Event processed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /webhook [post]
func (c *WebhookController) HandleWebhook(ctx *gin.Context) {
	var request struct {
		EventType string                 `json:"event_type" binding:"required"`
		Payload   map[string]interface{} `json:"payload" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.EventType == "processed" {
		err := c.taskService.UpdateTaskMetadata(request.Payload)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Event processed successfully"})
}