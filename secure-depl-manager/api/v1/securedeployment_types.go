/*
Copyright 2023.

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
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SecureDeploymentSpec defines the desired state of SecureDeployment
type SecureDeploymentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of SecureDeployment. Edit securedeployment_types.go to remove/update
	AuthUrl               string `json:"authUrl,omitempty"`
	appsv1.DeploymentSpec `json:"spec,omitempty"`
}

// SecureDeploymentStatus defines the observed state of SecureDeployment
type SecureDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DeploymentStatus appsv1.DeploymentStatus `json:"deploymentStatus,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SecureDeployment is the Schema for the securedeployments API
type SecureDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecureDeploymentSpec   `json:"spec,omitempty"`
	Status SecureDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SecureDeploymentList contains a list of SecureDeployment
type SecureDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecureDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecureDeployment{}, &SecureDeploymentList{})
}
