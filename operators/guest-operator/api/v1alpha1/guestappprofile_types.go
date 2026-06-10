package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GuestAppProfileSpec struct {
	Replicas       int32  `json:"replicas"`
	Enabled        bool   `json:"enabled"`
	DeploymentName string `json:"deploymentName,omitempty"`
}

type GuestAppProfileStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

type GuestAppProfile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GuestAppProfileSpec   `json:"spec,omitempty"`
	Status GuestAppProfileStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type GuestAppProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GuestAppProfile `json:"items"`
}
