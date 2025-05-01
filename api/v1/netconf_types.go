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
	NetplanErr = "Apply Failed"
)

const (
	ReasonCRNotAvailable          = "OperatorResourceNotAvailable"
	ReasonDeploymentNotAvailable  = "OperandDeploymentNotAvailable"
	ReasonOperandDeploymentFailed = "OperandDeploymentFailed"
	ReasonProcessing              = "OperatorProcessingConfiguration"
	ReasonSucceeded               = "OperatorSucceeded"
)

// NetConfSpec defines the desired state of NetConf
type NetConfSpec struct {
	// NodeSelector is a selector which must be true for the policy to be applied to the node.
	// Selector which must match a node's labels for the policy to be scheduled on that node.
	NodeName string `json:"nodeName,omitempty"`

	// The desired configuration of the policy
	NetConf string `json:"netConf,omitempty"`
}

// NetConfStatus defines the observed state of NetConf
type NetConfStatus struct {
	// Conditions is the list of status condition updates
	Conditions []metav1.Condition `json:"conditions"`

	Applied string `json:"applied,omitempty"`
	State   string `json:"state,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Applied",type=string,JSONPath=`.status.applied`
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// NetConf is the Schema for the netconfs API
type NetConf struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetConfSpec   `json:"spec,omitempty"`
	Status NetConfStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NetConfList contains a list of NetConf
type NetConfList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetConf `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NetConf{}, &NetConfList{})
}
