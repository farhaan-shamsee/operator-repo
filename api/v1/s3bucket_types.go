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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// S3BucketSpec defines the desired state of S3Bucket
type S3BucketSpec struct {
	BucketName string `json:"bucketName"`
	Region     string `json:"region"`
	ACL        string `json:"acl,omitempty"`
	// Versioning indicates whether versioning is enabled for the bucket.
	// Possible values are "Enabled" or "Suspended".
	Versioning string `json:"versioning,omitempty"`
	// StorageClass defines the default storage class for objects in the bucket.
	// Examples include "STANDARD", "REDUCED_REDUNDANCY", "GLACIER", etc.
	StorageClass string `json:"storageClass,omitempty"`
}

// S3BucketStatus defines the observed state of S3Bucket.
type S3BucketStatus struct {
	// BucketARN is the Amazon Resource Name of the S3 bucket
	BucketARN string `json:"bucketARN,omitempty"`
	// Location is the AWS region where the bucket was created
	Location string `json:"location,omitempty"`
	// Created indicates whether the bucket has been successfully created
	Created bool `json:"created,omitempty"`
	// LastSyncTime is the last time the bucket status was synchronized with AWS
	LastSyncTime string `json:"lastSyncTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// S3Bucket is the Schema for the s3buckets API
type S3Bucket struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of S3Bucket
	// +required
	Spec S3BucketSpec `json:"spec"`

	// status defines the observed state of S3Bucket
	// +optional
	Status S3BucketStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// S3BucketList contains a list of S3Bucket
type S3BucketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []S3Bucket `json:"items"`
}

type CreatedBucketInfo struct {
	BucketName string `json:"bucketName"`
	BucketARN  string `json:"bucketARN"`
	Location   string `json:"location"`
	Region     string `json:"region"`
}

func init() {
	SchemeBuilder.Register(&S3Bucket{}, &S3BucketList{})
}
