package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AdminAppProfileSpec struct {
	Replicas       int32  `json:"replicas"`
	Enabled        bool   `json:"enabled"`
	DeploymentName string `json:"deploymentName,omitempty"`
}

type AdminAppProfileStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

type AdminAppProfile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AdminAppProfileSpec   `json:"spec,omitempty"`
	Status AdminAppProfileStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type AdminAppProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AdminAppProfile `json:"items"`
}
