package config

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
