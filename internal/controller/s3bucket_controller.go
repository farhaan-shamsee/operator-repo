/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	computev1 "github.com/farhaan-shamsee/operator-repo/api/v1"
)

// S3BucketReconciler reconciles a S3Bucket object
type S3BucketReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=compute.cloud.com,resources=s3buckets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=compute.cloud.com,resources=s3buckets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=compute.cloud.com,resources=s3buckets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the S3Bucket object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *S3BucketReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := logf.FromContext(ctx)
	l.Info("=== S3BUCKET RECONCILE LOOP STARTED ===", "namespace", req.Namespace, "name", req.Name)

	s3bucket := &computev1.S3Bucket{}

	if err := r.Get(ctx, req.NamespacedName, s3bucket); err != nil {
		if errors.IsNotFound(err) {
			l.Info("Bucket deleted. No need to reconcile.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err

	}

	if !s3bucket.DeletionTimestamp.IsZero() {
		l.Info("Has deletion timestamp, bucket is being deleted")

		if s3bucket.Status.Created {
			l.Info("Deleting S3 bucket from AWS", "BucketARN", s3bucket.Status.BucketARN)
			_, err := deleteS3Bucket(ctx, s3bucket)
			if err != nil {
				l.Error(err, "Failed to delete S3 bucket from AWS", "BucketARN", s3bucket.Status.BucketARN)
				return ctrl.Result{}, err
			}
		} else {
			l.Info("Bucket was never created in AWS, skipping deletion")
		}

		// Remove finalizer
		if controllerutil.ContainsFinalizer(s3bucket, "s3bucket.compute.cloud.com") {
			controllerutil.RemoveFinalizer(s3bucket, "s3bucket.compute.cloud.com")
			if err := r.Update(ctx, s3bucket); err != nil {
				l.Error(err, "Failed to remove finalizer from S3 bucket", "BucketARN", s3bucket.Status.BucketARN)
				return ctrl.Result{}, err
			}
			l.Info("Finalizer removed from S3 bucket", "BucketARN", s3bucket.Status.BucketARN)
		}

		// Return after handling deletion
		return ctrl.Result{}, nil
	}

	if s3bucket.Status.BucketARN != "" {
		l.Info("Bucket already exists in AWS, updating sync time", "BucketARN", s3bucket.Status.BucketARN)

		// Update LastSyncTime to track when we last checked the bucket
		s3bucket.Status.LastSyncTime = time.Now().Format(time.RFC3339)
		if err := r.Status().Update(ctx, s3bucket); err != nil {
			l.Error(err, "Failed to update LastSyncTime", "BucketARN", s3bucket.Status.BucketARN)
			// Don't return error, this is not critical
		}

		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(s3bucket, "s3bucket.compute.cloud.com") {
		l.Info("=== ADDING FINALIZER ===")
		controllerutil.AddFinalizer(s3bucket, "s3bucket.compute.cloud.com")
		if err := r.Update(ctx, s3bucket); err != nil {
			l.Error(err, "Failed to add finalizer to S3 bucket", "BucketARN", s3bucket.Status.BucketARN)
			return ctrl.Result{}, err
		}
		l.Info("Finalizer added to S3 bucket", "BucketARN", s3bucket.Status.BucketARN)
		return ctrl.Result{}, nil
	}

	l.Info("Creating new s3 bucket")

	// Create new bucket
	createdBucketInfo, err := createS3Bucket(ctx, s3bucket)
	if err != nil {
		l.Error(err, "Failed to create S3 bucket in AWS")
		return ctrl.Result{}, err
	}

	// Update status with created bucket info
	s3bucket.Status.BucketARN = createdBucketInfo.BucketARN
	s3bucket.Status.Created = true
	s3bucket.Status.Location = createdBucketInfo.Location
	s3bucket.Status.LastSyncTime = time.Now().Format(time.RFC3339)
	err = r.Status().Update(ctx, s3bucket)
	if err != nil {
		l.Error(err, "Failed to update S3 bucket status after creation", "BucketARN", s3bucket.Status.BucketARN)
		return ctrl.Result{}, err
	}

	l.Info("S3 bucket created and status updated successfully", "BucketARN", s3bucket.Status.BucketARN)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *S3BucketReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&computev1.S3Bucket{}).
		Named("s3bucket").
		Complete(r)
}
