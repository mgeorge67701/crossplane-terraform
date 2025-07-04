/*
Copyright 2024 The Crossplane Authors.

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

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// ProviderConfigSpec defines the desired state of ProviderConfig
type ProviderConfigSpec struct {
	// Credentials required to authenticate to the Terraform provider.
	Credentials ProviderCredentials `json:"credentials"`

	// TerraformVersion specifies the version of Terraform to use.
	// +optional
	TerraformVersion string `json:"terraformVersion,omitempty"`

	// WorkingDirectory is the default working directory for Terraform execution.
	// +optional
	WorkingDirectory string `json:"workingDirectory,omitempty"`

	// Environment variables to set for all Terraform executions.
	// +optional
	Environment map[string]string `json:"environment,omitempty"`

	// Backend configuration for storing Terraform state.
	// +optional
	Backend *BackendConfig `json:"backend,omitempty"`

	// Parallelism limits the number of concurrent operations as Terraform
	// walks the graph. Defaults to 10.
	// +optional
	Parallelism int `json:"parallelism,omitempty"`

	// Refresh determines whether or not the providers should refresh state
	// before applying changes. Defaults to true.
	// +optional
	Refresh *bool `json:"refresh,omitempty"`

	// Endpoints defines custom endpoints for Terraform providers.
	// +optional
	Endpoints map[string]string `json:"endpoints,omitempty"`
}

// ProviderCredentials defines the credentials for the Terraform provider
type ProviderCredentials struct {
	// Source of the provider credentials.
	// +kubebuilder:validation:Enum=None;Secret;InjectedIdentity;Environment;Filesystem
	Source xpv1.CredentialsSource `json:"source"`

	// CommonCredentialSelectors provides common selectors for extracting
	// credentials.
	xpv1.CommonCredentialSelectors `json:",inline"`
}

// ProviderConfigStatus defines the observed state of ProviderConfig
type ProviderConfigStatus struct {
	xpv1.ProviderConfigStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="SECRET-NAME",type="string",JSONPath=".spec.credentials.secretRef.name",priority=1
// +kubebuilder:resource:scope=Cluster,categories={crossplane,providerconfig,terraform}

// A ProviderConfig configures the Terraform provider.
type ProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigSpec   `json:"spec"`
	Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderConfigList contains a list of ProviderConfig
type ProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfig `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="SECRET-NAME",type="string",JSONPath=".spec.credentials.secretRef.name",priority=1
// +kubebuilder:resource:scope=Cluster,categories={crossplane,providerconfig,terraform}

// A ProviderConfigUsage indicates that a resource is using a ProviderConfig.
type ProviderConfigUsage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigUsageSpec   `json:"spec"`
	Status ProviderConfigUsageStatus `json:"status,omitempty"`
}

// ProviderConfigUsageSpec defines the desired state of ProviderConfigUsage
type ProviderConfigUsageSpec struct {
	// ProviderConfigRef is a reference to a ProviderConfig.
	ProviderConfigRef xpv1.Reference `json:"providerConfigRef"`

	// ResourceRef is a reference to a resource using the ProviderConfig.
	ResourceRef xpv1.TypedReference `json:"resourceRef"`
}

// ProviderConfigUsageStatus defines the observed state of ProviderConfigUsage
type ProviderConfigUsageStatus struct {
	xpv1.ResourceStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// ProviderConfigUsageList contains a list of ProviderConfigUsage
type ProviderConfigUsageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfigUsage `json:"items"`
}
