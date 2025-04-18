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

const (
	NoError    = "Done Successfully"
	Processing = "Processing"
	NotMatch   = "No matching nodes"
)

const (
	ReasonCRNotAvailable          = "OperatorResourceNotAvailable"
	ReasonDeploymentNotAvailable  = "OperandDeploymentNotAvailable"
	ReasonOperandDeploymentFailed = "OperandDeploymentFailed"
	ReasonSucceeded               = "OperatorSucceeded"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DigiNetSpec defines the desired state of DigiNet
type DigiNetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ClusterName string `json:"clusternamClusterNamee,omitempty"`
}

// DigiNetStatus defines the observed state of DigiNet
type DigiNetStatus struct {
	// Conditions is the list of status condition updates
	Conditions []metav1.Condition `json:"conditions"`

	Applied string `json:"applied,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Applied",type=string,JSONPath=`.status.applied`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// DigiNet is the Schema for the diginets API
type DigiNet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DigiNetSpec   `json:"spec,omitempty"`
	Status DigiNetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DigiNetList contains a list of DigiNet
type DigiNetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DigiNet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DigiNet{}, &DigiNetList{})
}
