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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	computev1 "github.com/farhaan-shamsee/operator-repo/api/v1"
)

// Ec2instanceReconciler reconciles a Ec2instance object
type Ec2instanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=compute.cloud.com,resources=ec2instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=compute.cloud.com,resources=ec2instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=compute.cloud.com,resources=ec2instances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Ec2instance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *Ec2instanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := logf.FromContext(ctx)

	l.Info("=== RECONCILE LOOP STARTED ===", "namespace", req.Namespace, "name", req.Name)

	ec2instance := &computev1.Ec2instance{}
	if err := r.Get(ctx, req.NamespacedName, ec2instance); err != nil { //this is fetching from cluster and putting into ec2instance variable 
		if errors.IsNotFound(err) {
			l.Info("Instance Deleted. No need to reconcile.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !ec2instance.DeletionTimestamp.IsZero() {
		l.Info("Has Deletion timestamp, Instance is being deleted")

		// Only attempt to delete from AWS if an instance was actually created
		if ec2instance.Status.InstanceID != "" {
			l.Info("Deleting EC2 instance from AWS", "instanceID", ec2instance.Status.InstanceID)
			_, err := deleteEc2Instance(ctx, ec2instance)
			if err != nil {
				l.Error(err, "Failed to delete EC2 instance from AWS")
				// Still continue to remove finalizer - don't block deletion on AWS errors
				// This prevents orphaned Kubernetes resources if AWS is unavailable
			}
		} else {
			l.Info("No instance ID found in status, skipping AWS deletion (instance was never created)")
		}

		// Remove finalizer to allow Kubernetes to delete the resource
		if controllerutil.ContainsFinalizer(ec2instance, "ec2instance.compute.cloud.com") {
			controllerutil.RemoveFinalizer(ec2instance, "ec2instance.compute.cloud.com")
			if err := r.Update(ctx, ec2instance); err != nil {
				l.Error(err, "Failed to remove finalizer")
				return ctrl.Result{}, err
			}
			l.Info("Finalizer removed successfully")
		}

		// Deletion complete
		return ctrl.Result{}, nil
	}

	if ec2instance.Status.InstanceID != "" {
		l.Info("Requested object already exist in K8s. Not creating a new instance", "instance", ec2instance.Status.InstanceID)
		return ctrl.Result{}, nil
	}

	// Add finalizer if not already present
	if !controllerutil.ContainsFinalizer(ec2instance, "ec2instance.compute.cloud.com") {
		l.Info("=== ADDING FINALIZER ===")
		controllerutil.AddFinalizer(ec2instance, "ec2instance.compute.cloud.com")
		if err := r.Update(ctx, ec2instance); err != nil {
			l.Error(err, "Failed to add finalizer")
			return ctrl.Result{}, err
		}
		l.Info("=== FINALIZER ADDED - This update will trigger a NEW reconcile loop ===")
		// Return here - let the new reconcile loop handle instance creation
		// This prevents race conditions and ensures clean state
		return ctrl.Result{}, nil
	}

	l.Info("Creating new instance")

	// Create a new instance
	l.Info("=== CONTINUING WITH EC2 INSTANCE CREATION IN THE CURRENT RECONCILE ===")

	createdInstanceInfo, err := createEc2Instance(ctx, ec2instance)
	if err != nil {
		l.Error(err, "Failed to create EC2 instance")
		// Kubernetes will retry with backoff
		return ctrl.Result{}, err
	}

	l.Info("=== UPDATING EC2INSTANCE STATUS - This will trigger another reconcile ===",
		"instanceId", createdInstanceInfo.InstanceId,
		"state", createdInstanceInfo.State)

	ec2instance.Status.InstanceID = createdInstanceInfo.InstanceId
	ec2instance.Status.State = createdInstanceInfo.State
	ec2instance.Status.PrivateIP = createdInstanceInfo.PrivateIP
	ec2instance.Status.PublicIP = createdInstanceInfo.PublicIP
	ec2instance.Status.PrivateDNS = createdInstanceInfo.PrivateDNS
	ec2instance.Status.PublicDNS = createdInstanceInfo.PublicDNS

	// The Reconcile function must return a ctrl.Result and an error.
	// Returning ctrl.Result{} with nil error means the reconciliation was successful
	// and no requeue is requested. If an error is returned, the controller will
	// automatically requeue the request for another attempt.
	// Sends Requeue ( bool ) and RequeueAfter ( time.Duration ).
	//return ctrl.Result{}, nil
	err = r.Status().Update(ctx, ec2instance)
	if err != nil {
		l.Error(err, "Failed to update the status")
		return ctrl.Result{}, err
	}

	l.Info("=== STATUS UPDATED - Reconcile loop will be triggered again ===")

	// Successfully created and updated status - done
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
// SetupWithManager registers the Ec2InstanceReconciler with the controller manager.
// It configures the controller to watch for changes to Ec2Instance resources.
// The controller will be named "ec2instance" for logging and metrics purposes.
// The Complete(r) call finalizes the setup, associating the reconciler logic with this controller.
func (r *Ec2instanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&computev1.Ec2instance{}).
		Named("ec2instance").
		Complete(r)
}
