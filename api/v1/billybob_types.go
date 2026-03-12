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

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.numberPods,statuspath=.status.currentPods,selectorpath=.status.selector

// Billybob is the Schema for the billybobs API.
type Billybob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BillybobSpec   `json:"spec,omitempty"`
	Status BillybobStatus `json:"status,omitempty"`
}

// BillybobSpec defines the desired state of Billybob
type BillybobSpec struct {
	NumberPods int `json:"numberPods,omitempty"`
}

// BillybobStatus defines the observed state of Billybob
type BillybobStatus struct {
	PodNames    []string `json:"podNames,omitempty"`
	CurrentPods int      `json:"currentPods,omitempty"`
	Selector    string   `json:"selector,omitempty"`
}

// +kubebuilder:object:root=true

// BillybobList contains a list of Billybob.
type BillybobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Billybob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Billybob{}, &BillybobList{})
}
