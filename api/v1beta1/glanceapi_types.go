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

package v1beta1

import (
	"fmt"

	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	"github.com/openstack-k8s-operators/lib-common/modules/common/endpoint"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// DeploymentHash hash used to detect changes
	DeploymentHash = "deployment"
)


// GlanceAPISpec defines the desired state of GlanceAPI
type GlanceAPISpec struct {

	// Input parameter coming from glance template
	GlanceAPITemplate `json:",inline"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=glance
	// ServiceUser - optional username used for this service to register in glance
	ServiceUser string `json:"serviceUser"`

	// +kubebuilder:validation:Optional
	// DatabaseHostname - Glance Database Hostname
	DatabaseHostname string `json:"databaseHostname"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=glance
	// DatabaseUser - optional username used for glance DB, defaults to glance
	// TODO: -> implement needs work in mariadb-operator, right now only glance
	DatabaseUser string `json:"databaseUser"`

	// +kubebuilder:validation:Optional
	// Secret containing OpenStack password information for glance AdminPassword
	Secret string `json:"secret"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={database: GlanceDatabasePassword, service: GlancePassword}
	// PasswordSelectors - Selectors to identify the DB and ServiceUser password from the Secret
	PasswordSelectors PasswordSelector `json:"passwordSelectors"`

	// +kubebuilder:validation:Optional
	// ExtraMounts containing conf files and credentials
	ExtraMounts []GlanceExtraVolMounts `json:"extraMounts,omitempty"`
}

// GlanceAPIDebug defines the observed state of GlanceAPIDebug
type GlanceAPIDebug struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	// Service enable debug
	Service bool `json:"service"`
}

// GlanceAPIStatus defines the observed state of GlanceAPI
type GlanceAPIStatus struct {
	// ReadyCount of glance API instances
	ReadyCount int32 `json:"readyCount,omitempty"`

	// Map of hashes to track e.g. job status
	Hash map[string]string `json:"hash,omitempty"`

	// API endpoint
	APIEndpoints map[string]string `json:"apiEndpoint,omitempty"`

	// Conditions
	Conditions condition.Conditions `json:"conditions,omitempty" optional:"true"`

	// NetworkAttachments status of the deployment pods
	NetworkAttachments map[string][]string `json:"networkAttachments,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="NetworkAttachments",type="string",JSONPath=".status.networkAttachments",description="NetworkAttachments"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[0].status",description="Status"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[0].message",description="Message"

// GlanceAPI is the Schema for the glanceapis API
type GlanceAPI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec        GlanceAPISpec          `json:"spec,omitempty"`
	Status      GlanceAPIStatus        `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GlanceAPIList contains a list of GlanceAPI
type GlanceAPIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlanceAPI `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GlanceAPI{}, &GlanceAPIList{})
}

// GetEndpoint - returns OpenStack endpoint url for type
func (instance GlanceAPI) GetEndpoint(endpointType endpoint.Endpoint) (string, error) {
	if url, found := instance.Status.APIEndpoints[string(endpointType)]; found {
		return url, nil
	}
	return "", fmt.Errorf("%s endpoint not found", string(endpointType))
}

// IsReady - returns true if service is ready to server requests
func (instance GlanceAPI) IsReady() bool {
	// Ready when:
	// there is at least a single pod to serve the glance service
	return instance.Status.ReadyCount >= 1
}
