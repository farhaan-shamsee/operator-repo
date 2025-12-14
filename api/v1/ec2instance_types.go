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

// Ec2instanceSpec defines the desired state of Ec2instance
type Ec2instanceSpec struct {
	InstanceType      string            `json:"instanceType"`
	AMIId             string            `json:"amiId"`
	Region            string            `json:"region"`
	AvailabilityZone  string            `json:"availabilityZone,omitempty"`
	KeyPair           string            `json:"keyPair,omitempty"`
	SecurityGroups    []string          `json:"securityGroups,omitempty"`
	Subnet            string            `json:"subnet,omitempty"`
	UserData          string            `json:"userData,omitempty"`
	Tags              map[string]string `json:"tags,omitempty"`
	Storage           StorageConfig     `json:"storage,omitempty"`
	AssociatePublicIP bool              `json:"associatePublicIP,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="InstanceType",type="string",JSONPath=".spec.instanceType",description="The EC2 instance type"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state",description="The current state of the EC2 instance"
// +kubebuilder:printcolumn:name="PublicIP",type="string",JSONPath=".status.publicIP",description="The public IP of the EC2 instance"
// +kubebuilder:printcolumn:name="InstanceID",type="string",JSONPath=".status.instanceID",description="The AWS instance ID"
// Ec2Instance is the Schema for the ec2instances API.
type Ec2instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	Spec   Ec2instanceSpec   `json:"spec"`
	Status Ec2instanceStatus `json:"status,omitempty,omitzero"`
}

// Ec2instanceStatus defines the observed state of Ec2instance.
type Ec2instanceStatus struct {
	InstanceID string `json:"instanceID,omitempty"`
	State      string `json:"state,omitempty"`
	PublicIP   string `json:"publicIP,omitempty"`
	PrivateIP  string `json:"privateIP,omitempty"`
	PublicDNS  string `json:"publicDNS,omitempty"`
	PrivateDNS string `json:"privateDNS,omitempty"`
	LaunchTime string `json:"launchTime,omitempty"`
}

// StorageConfig defines the storage configuration for the EC2 instance.
type StorageConfig struct {
	RootVolume        VolumeConfig   `json:"rootVolume,omitempty"`
	AdditionalVolumes []VolumeConfig `json:"additionalVolumes,omitempty"`
}

// VolumeConfig defines the configuration for a volume in the EC2 instance.
type VolumeConfig struct {
	Size       int32  `json:"size"`
	Type       string `json:"type,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`
	Encrypted  bool   `json:"encrypted,omitempty"`
}

type Condition struct {
	Type               string      `json:"type"`
	Status             string      `json:"status"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

// +kubebuilder:object:root=true

// Ec2instanceList contains a list of Ec2instance
type Ec2instanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ec2instance `json:"items"`
}

type CreatedInstanceInfo struct {
	InstanceId string `json:"instanceId"`
	State      string `json:"state"`
	PrivateIP  string `json:"privateIP"`
	PublicIP   string `json:"publicIP"`
	PrivateDNS string `json:"privateDNS"`
	PublicDNS  string `json:"publicDNS"`
}

func init() {
	SchemeBuilder.Register(&Ec2instance{}, &Ec2instanceList{})
}
