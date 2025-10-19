package controllers

import (
	"net/http"
	"time"

	"screensaver-ad-backend/internal/models"
	"screensaver-ad-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type TemplateController struct {
	service   *services.TemplateService
	s3Service *services.S3Service
}

func NewTemplateController(service *services.TemplateService, s3Service *services.S3Service) *TemplateController {
	return &TemplateController{service: service, s3Service: s3Service}
}

// UploadTemplate handles uploading a template video and name
// @Summary Upload a new template
// @Description Upload a template video file with a name
// @Tags templates
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Template name"
// @Param file formData file true "Template video file"
// @Success 200 {object} map[string]interface{} "Template uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /templates [post]
func (tc *TemplateController) UploadTemplate(c *gin.Context) {
	name := c.PostForm("name")
	file, err := c.FormFile("file")
	if err != nil || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and file are required"})
		return
	}

	fileObj, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer fileObj.Close()

	s3Key, err := tc.s3Service.UploadFileToS3(fileObj, file, name, "template")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload to S3"})
		return
	}

	template := &models.Template{
		Name:     name,
		S3Key:    s3Key,
		S3Bucket: tc.s3Service.Bucket,
	}
	if err := tc.service.CreateTemplate(template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template uploaded", "template": template})
}

// ListTemplates returns all templates with presigned URLs
// @Summary List all templates
// @Description Get a list of all templates with presigned URLs
// @Tags templates
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{} "List of templates with URLs"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /templates [get]
func (tc *TemplateController) ListTemplates(c *gin.Context) {
	templates, err := tc.service.ListTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list templates"})
		return
	}

	result := []gin.H{}
	for _, t := range templates {
		url, err := tc.s3Service.GetFileURL(t.S3Key, 15*time.Minute)
		if err != nil {
			url = ""
		}
		result = append(result, gin.H{
			"name": t.Name,
			"url":  url,
		})
	}
	c.JSON(http.StatusOK, result)
}
