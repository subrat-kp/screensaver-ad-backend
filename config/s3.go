package config

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	S3Client *s3.S3
	S3Bucket string
)

// InitS3 initializes the S3 client
func InitS3() error {
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	S3Bucket = os.Getenv("AWS_S3_BUCKET")

	if region == "" || accessKey == "" || secretKey == "" || S3Bucket == "" {
		return nil // Return nil to allow app to start without S3 configured
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return err
	}

	S3Client = s3.New(sess)
	return nil
}

// GetS3Client returns the S3 client instance
func GetS3Client() *s3.S3 {
	return S3Client
}

// GetS3Bucket returns the configured S3 bucket name
func GetS3Bucket() string {
	return S3Bucket
}

// PutObject uploads an object to S3
func PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64, contentType string) error {
	uploader := s3manager.NewUploaderWithClient(S3Client)
	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	return err
}

// GeneratePresignedURL generates a presigned GET URL for an S3 object
func GeneratePresignedURL(bucket, key string, expire time.Duration) (string, error) {
	req, _ := S3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(expire)
	return urlStr, err
}
