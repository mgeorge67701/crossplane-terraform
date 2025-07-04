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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// TerraformParameters are the configurable fields of a Terraform resource.
type TerraformParameters struct {
	// Configuration contains the raw Terraform configuration.
	// +kubebuilder:validation:Required
	Configuration runtime.RawExtension `json:"configuration"`

	// Variables is a map of Terraform variables.
	// +optional
	Variables map[string]string `json:"variables,omitempty"`

	// Backend configuration for storing Terraform state.
	// +optional
	Backend *BackendConfig `json:"backend,omitempty"`

	// Workspace name for this Terraform configuration.
	// +optional
	Workspace string `json:"workspace,omitempty"`

	// Source specifies the location of the Terraform module.
	// +optional
	Source *TerraformSource `json:"source,omitempty"`
}

// BackendConfig represents Terraform backend configuration.
type BackendConfig struct {
	// Type of the backend (e.g., "s3", "gcs", "azurerm").
	// +kubebuilder:validation:Required
	Type string `json:"type"`

	// Configuration for the backend.
	// +optional
	Configuration map[string]string `json:"configuration,omitempty"`
}

// TerraformSource represents the source of a Terraform module.
type TerraformSource struct {
	// Path to the Terraform module (for local modules).
	// +optional
	Path string `json:"path,omitempty"`

	// Git repository URL for remote modules.
	// +optional
	Git *GitSource `json:"git,omitempty"`

	// HTTP URL for remote modules.
	// +optional
	HTTP *HTTPSource `json:"http,omitempty"`
}

// GitSource represents a Git repository source.
type GitSource struct {
	// URL of the Git repository.
	// +kubebuilder:validation:Required
	URL string `json:"url"`

	// Branch, tag, or commit to use.
	// +optional
	Ref string `json:"ref,omitempty"`

	// Subdirectory within the repository.
	// +optional
	Path string `json:"path,omitempty"`
}

// HTTPSource represents an HTTP source.
type HTTPSource struct {
	// URL of the HTTP source.
	// +kubebuilder:validation:Required
	URL string `json:"url"`

	// Checksum for verification.
	// +optional
	Checksum string `json:"checksum,omitempty"`
}

// TerraformObservation are the observable fields of a Terraform resource.
type TerraformObservation struct {
	// Outputs from the Terraform execution.
	// +optional
	Outputs map[string]string `json:"outputs,omitempty"`

	// State of the Terraform execution.
	// +optional
	State string `json:"state,omitempty"`

	// LastApplied timestamp.
	// +optional
	LastApplied *metav1.Time `json:"lastApplied,omitempty"`

	// ApplyJobName is the name of the job that applies the Terraform configuration.
	// +optional
	ApplyJobName string `json:"applyJobName,omitempty"`

	// DestroyJobName is the name of the job that destroys the Terraform resources.
	// +optional
	DestroyJobName string `json:"destroyJobName,omitempty"`
}

// A TerraformSpec defines the desired state of a Terraform resource.
type TerraformSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       TerraformParameters `json:"forProvider"`
}

// A TerraformStatus represents the observed state of a Terraform resource.
type TerraformStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          TerraformObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,terraform}

// A Terraform is a managed resource that represents a Terraform configuration.
type Terraform struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TerraformSpec   `json:"spec"`
	Status TerraformStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TerraformList contains a list of Terraform resources.
type TerraformList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Terraform `json:"items"`
}

// GetCondition of this Terraform.
func (mg *Terraform) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return mg.Status.GetCondition(ct)
}

// GetDeletionPolicy of this Terraform.
func (mg *Terraform) GetDeletionPolicy() xpv1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetManagementPolicies of this Terraform.
func (mg *Terraform) GetManagementPolicies() xpv1.ManagementPolicies {
	return mg.Spec.ManagementPolicies
}

// GetProviderConfigReference of this Terraform.
func (mg *Terraform) GetProviderConfigReference() *xpv1.Reference {
	return mg.Spec.ProviderConfigReference
}

// GetPublishConnectionDetailsTo of this Terraform.
func (mg *Terraform) GetPublishConnectionDetailsTo() *xpv1.PublishConnectionDetailsTo {
	return mg.Spec.PublishConnectionDetailsTo
}

// GetWriteConnectionSecretToReference of this Terraform.
func (mg *Terraform) GetWriteConnectionSecretToReference() *xpv1.SecretReference {
	return mg.Spec.WriteConnectionSecretToReference
}

// SetConditions of this Terraform.
func (mg *Terraform) SetConditions(c ...xpv1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this Terraform.
func (mg *Terraform) SetDeletionPolicy(r xpv1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetManagementPolicies of this Terraform.
func (mg *Terraform) SetManagementPolicies(r xpv1.ManagementPolicies) {
	mg.Spec.ManagementPolicies = r
}

// SetProviderConfigReference of this Terraform.
func (mg *Terraform) SetProviderConfigReference(r *xpv1.Reference) {
	mg.Spec.ProviderConfigReference = r
}

// SetPublishConnectionDetailsTo of this Terraform.
func (mg *Terraform) SetPublishConnectionDetailsTo(r *xpv1.PublishConnectionDetailsTo) {
	mg.Spec.PublishConnectionDetailsTo = r
}

// SetWriteConnectionSecretToReference of this Terraform.
func (mg *Terraform) SetWriteConnectionSecretToReference(r *xpv1.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = r
}

// TerraformGroupKind is the GroupKind for the Terraform resource.
var TerraformGroupKind = schema.GroupKind{
	Group: Group,
	Kind:  "Terraform",
}

// TerraformGroupVersionKind is the GroupVersionKind for the Terraform resource.
var TerraformGroupVersionKind = schema.GroupVersionKind{
	Group:   Group,
	Version: Version,
	Kind:    "Terraform",
}
