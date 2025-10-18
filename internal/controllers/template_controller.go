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

	s3Key, err := tc.s3Service.UploadFileToS3(fileObj, file, name)
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
