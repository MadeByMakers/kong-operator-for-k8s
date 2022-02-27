/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RouteSpec defines the desired state of Route
type RouteSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Id    string   `json:"id,omitempty"`
	Name  string   `json:"name,omitempty"`
	Tags  []string `json:"tags,omitempty"`
	Hosts []string `json:"hosts,omitempty"`
	Paths []string `json:"paths,omitempty"`

	Headers map[string]string `json:"headers,omitempty"`

	Service                    ObjectId `json:"service,omitempty"`
	Path_handling              string   `json:"path_handling,omitempty"`
	Https_redirect_status_code int      `json:"https_redirect_status_code,omitempty"`
	Regex_priority             int      `json:"regex_priority,omitempty"`
	Methods                    []string `json:"methods,omitempty"`
	Strip_path                 bool     `json:"strip_path,omitempty"`
	Preserve_host              bool     `json:"preserve_host,omitempty"`
	Request_buffering          bool     `json:"request_buffering,omitempty"`
	Response_buffering         bool     `json:"response_buffering,omitempty"`

	Protocols    []string `json:"protocols,omitempty"`
	Snis         []string `json:"snis,omitempty"`
	Destinations []string `json:"destinations,omitempty"`
}

// RouteStatus defines the observed state of Route
type RouteStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Code     int        `json:"code,omitempty"`
	Message  string     `json:"message,omitempty"`
	Response HttpStatus `json:"response,omitempty"`

	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Route is the Schema for the routes API
type Route struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RouteSpec   `json:"spec,omitempty"`
	Status RouteStatus `json:"status,omitempty"`
}

func (m *Route) GetConditions() []metav1.Condition {
	return m.Status.Conditions
}

func (m *Route) SetConditions(conditions []metav1.Condition) {
	m.Status.Conditions = conditions
}

//+kubebuilder:object:root=true

// RouteList contains a list of Route
type RouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Route `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Route{}, &RouteList{})
}
