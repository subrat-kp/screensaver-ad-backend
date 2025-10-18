package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"screensaver-ad-backend/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

// S3Service handles S3 operations
type S3Service struct {
	Client *s3.S3
	Bucket string
}

// NewS3Service creates a new S3 service instance
func NewS3Service() *S3Service {
	return &S3Service{
		Client: config.GetS3Client(),
		Bucket: config.GetS3Bucket(),
	}
}

// UploadFileToS3 uploads a file to S3 and returns the S3 key
func (s *S3Service) UploadFileToS3(file multipart.File, fileHeader *multipart.FileHeader, customName string) (string, error) {
	if s.Client == nil {
		return "", fmt.Errorf("S3 client is not initialized")
	}

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	var fileName string
	if customName != "" {
		// Sanitize custom name
		sanitized := strings.ReplaceAll(customName, " ", "_")
		sanitized = strings.ToLower(sanitized)
		fileName = fmt.Sprintf("%s_%s%s", sanitized, uuid.New().String()[:8], ext)
	} else {
		fileName = fmt.Sprintf("%s_%s%s", time.Now().Format("20060102_150405"), uuid.New().String()[:8], ext)
	}

	// S3 key with input folder
	s3Key := fmt.Sprintf("input/%s", fileName)

	// Prepare upload input
	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	}

	// Upload to S3
	_, err = s.Client.PutObject(uploadInput)
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return s3Key, nil
}

// DeleteFileFromS3 deletes a file from S3
func (s *S3Service) DeleteFileFromS3(s3Key string) error {
	if s.Client == nil {
		return fmt.Errorf("S3 client is not initialized")
	}

	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s3Key),
	}

	_, err := s.Client.DeleteObject(deleteInput)
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

// GetFileURL generates a presigned URL for accessing the file
func (s *S3Service) GetFileURL(s3Key string, expiration time.Duration) (string, error) {
	if s.Client == nil {
		return "", fmt.Errorf("S3 client is not initialized")
	}

	req, _ := s.Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s3Key),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url, nil
}
