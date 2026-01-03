package controller

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// getAWSConfig returns a general AWS config that can be used for any AWS service
func getAWSConfig(region string) (aws.Config, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Validate credentials are not empty
	if accessKey == "" || secretKey == "" {
		return aws.Config{}, fmt.Errorf("AWS credentials not found: ensure AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables are set")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return cfg, nil
}
