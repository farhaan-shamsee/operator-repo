package controller

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	computev1 "github.com/farhaan-shamsee/operator-repo/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func createS3Bucket(ctx context.Context, s3Bucket *computev1.S3Bucket) (createdBucketInfo *computev1.CreatedBucketInfo, err error) {
	l := log.FromContext(ctx) // Use context-aware logger instead of global logger

	l.Info("=== STARTING S3 BUCKET CREATION ===",
		"bucketName", s3Bucket.Spec.BucketName,
		"region", s3Bucket.Spec.Region,
		"storageClass", s3Bucket.Spec.StorageClass)

	// Get AWS config and create S3 client
	cfg, err := getAWSConfig(s3Bucket.Spec.Region)
	if err != nil {
		l.Error(err, "Failed to get AWS config")
		return nil, fmt.Errorf("failed to get AWS config: %w", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	// Prepare the CreateBucket input
	createBucketInput := &s3.CreateBucketInput{
		Bucket: aws.String(s3Bucket.Spec.BucketName),
	}

	// Add ACL if specified
	if s3Bucket.Spec.ACL != "" {
		createBucketInput.ACL = s3types.BucketCannedACL(s3Bucket.Spec.ACL)
	}

	// For regions other than us-east-1, we need to specify LocationConstraint
	if s3Bucket.Spec.Region != "us-east-1" {
		createBucketInput.CreateBucketConfiguration = &s3types.CreateBucketConfiguration{
			LocationConstraint: s3types.BucketLocationConstraint(s3Bucket.Spec.Region),
		}
	}

	l.Info("Creating S3 bucket with configuration",
		"bucketName", s3Bucket.Spec.BucketName,
		"region", s3Bucket.Spec.Region,
		"acl", s3Bucket.Spec.ACL)

	// Create the S3 bucket
	createOutput, err := s3Client.CreateBucket(ctx, createBucketInput)
	if err != nil {
		l.Error(err, "Failed to create S3 bucket")
		return nil, fmt.Errorf("failed to create S3 bucket: %w", err)
	}

	l.Info("=== S3 BUCKET CREATED SUCCESSFULLY ===",
		"bucketName", s3Bucket.Spec.BucketName,
		"location", aws.ToString(createOutput.Location))

	// Construct bucket ARN (format: arn:aws:s3:::bucket-name)
	// Note: createOutput.BucketArn is only populated for directory buckets
	bucketARN := fmt.Sprintf("arn:aws:s3:::%s", s3Bucket.Spec.BucketName)

	return &computev1.CreatedBucketInfo{
		BucketName: s3Bucket.Spec.BucketName,
		BucketARN:  bucketARN,
		Location:   aws.ToString(createOutput.Location),
		Region:     s3Bucket.Spec.Region,
	}, nil
}
