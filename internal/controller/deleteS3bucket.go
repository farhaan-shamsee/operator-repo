package controller

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	computev1 "github.com/farhaan-shamsee/operator-repo/api/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

func deleteS3Bucket(ctx context.Context, s3Bucket *computev1.S3Bucket) (bool, error) {
	l := logf.FromContext(ctx)

	l.Info("Deleting S3 bucket", "bucketARN", s3Bucket.Status.BucketARN)

	// Get AWS config and create S3 client
	cfg, err := getAWSConfig(s3Bucket.Spec.Region)
	if err != nil {
		l.Error(err, "Failed to get AWS config")
		return false, err
	}
	s3Client := s3.NewFromConfig(cfg)

	_, err = s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(s3Bucket.Spec.BucketName),
	})
	if err != nil {
		l.Error(err, "Failed to delete S3 bucket")
		return false, err
	}

	l.Info("S3 bucket deletion initiated", "bucketARN", s3Bucket.Status.BucketARN)

	// Wait for the bucket to be deleted
	waiter := s3.NewBucketNotExistsWaiter(s3Client)
	maxWaitTime := 5 * time.Minute
	waitParams := &s3.HeadBucketInput{
		Bucket: aws.String(s3Bucket.Spec.BucketName),
	}

	l.Info("Waiting for the S3 bucket to be deleted",
		"bucketARN", s3Bucket.Status.BucketARN,
		"maxWaitTime", maxWaitTime)

	err = waiter.Wait(ctx, waitParams, maxWaitTime)
	if err != nil {
		l.Error(err, "Failed to wait for S3 bucket deletion",
			"bucketARN", s3Bucket.Status.BucketARN)
		return false, err
	}

	l.Info("S3 bucket successfully deleted", "bucketARN", s3Bucket.Status.BucketARN)
	return true, nil
}
