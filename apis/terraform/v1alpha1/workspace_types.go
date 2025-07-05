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
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// WorkspaceParameters are the configurable fields of a Workspace resource.
type WorkspaceParameters struct {
	// Name of the Terraform workspace.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Variables is a map of Terraform variables for this workspace.
	// +optional
	Variables map[string]string `json:"variables,omitempty"`

	// Environment variables for this workspace.
	// +optional
	Environment map[string]string `json:"environment,omitempty"`

	// AutoApply determines if changes should be automatically applied.
	// +optional
	AutoApply bool `json:"autoApply,omitempty"`

	// AutoDestroy determines if the workspace should be automatically destroyed.
	// +optional
	AutoDestroy bool `json:"autoDestroy,omitempty"`

	// TerraformVersion specifies the Terraform version to use.
	// +optional
	TerraformVersion string `json:"terraformVersion,omitempty"`

	// WorkingDirectory is the working directory for Terraform execution.
	// +optional
	WorkingDirectory string `json:"workingDirectory,omitempty"`
}

// WorkspaceObservation are the observable fields of a Workspace resource.
type WorkspaceObservation struct {
	// ID of the workspace.
	// +optional
	ID string `json:"id,omitempty"`

	// CurrentRunID is the ID of the current run.
	// +optional
	CurrentRunID string `json:"currentRunId,omitempty"`

	// Status of the workspace.
	// +optional
	Status string `json:"status,omitempty"`

	// LastRun contains information about the last run.
	// +optional
	LastRun *WorkspaceRun `json:"lastRun,omitempty"`

	// ResourceCount is the number of resources managed by this workspace.
	// +optional
	ResourceCount int `json:"resourceCount,omitempty"`

	// CreatedAt timestamp.
	// +optional
	CreatedAt *metav1.Time `json:"createdAt,omitempty"`

	// UpdatedAt timestamp.
	// +optional
	UpdatedAt *metav1.Time `json:"updatedAt,omitempty"`
}

// WorkspaceRun contains information about a workspace run.
type WorkspaceRun struct {
	// ID of the run.
	// +optional
	ID string `json:"id,omitempty"`

	// Status of the run.
	// +optional
	Status string `json:"status,omitempty"`

	// CreatedAt timestamp.
	// +optional
	CreatedAt *metav1.Time `json:"createdAt,omitempty"`

	// Message associated with the run.
	// +optional
	Message string `json:"message,omitempty"`

	// IsDestroy indicates if this is a destroy run.
	// +optional
	IsDestroy bool `json:"isDestroy,omitempty"`
}

// A WorkspaceSpec defines the desired state of a Workspace resource.
type WorkspaceSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       WorkspaceParameters `json:"forProvider"`
}

// A WorkspaceStatus represents the observed state of a Workspace resource.
type WorkspaceStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          WorkspaceObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,terraform}

// A Workspace is a managed resource that represents a Terraform workspace.
type Workspace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkspaceSpec   `json:"spec"`
	Status WorkspaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WorkspaceList contains a list of Workspace resources.
type WorkspaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Workspace `json:"items"`
}

// GetCondition of this Workspace.
func (mg *Workspace) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return mg.Status.ResourceStatus.GetCondition(ct)
}

// GetDeletionPolicy of this Workspace.
func (mg *Workspace) GetDeletionPolicy() xpv1.DeletionPolicy {
	return mg.Spec.ResourceSpec.DeletionPolicy
}

// GetManagementPolicies of this Workspace.
func (mg *Workspace) GetManagementPolicies() xpv1.ManagementPolicies {
	return mg.Spec.ResourceSpec.ManagementPolicies
}

// GetProviderConfigReference of this Workspace.
func (mg *Workspace) GetProviderConfigReference() *xpv1.Reference {
	return mg.Spec.ResourceSpec.ProviderConfigReference
}

// GetPublishConnectionDetailsTo of this Workspace.
func (mg *Workspace) GetPublishConnectionDetailsTo() *xpv1.PublishConnectionDetailsTo {
	return mg.Spec.ResourceSpec.PublishConnectionDetailsTo
}

// GetWriteConnectionSecretToReference of this Workspace.
func (mg *Workspace) GetWriteConnectionSecretToReference() *xpv1.SecretReference {
	return mg.Spec.ResourceSpec.WriteConnectionSecretToReference
}

// SetConditions of this Workspace.
func (mg *Workspace) SetConditions(c ...xpv1.Condition) {
	mg.Status.ResourceStatus.SetConditions(c...)
}

// SetDeletionPolicy of this Workspace.
func (mg *Workspace) SetDeletionPolicy(r xpv1.DeletionPolicy) {
	mg.Spec.ResourceSpec.DeletionPolicy = r
}

// SetManagementPolicies of this Workspace.
func (mg *Workspace) SetManagementPolicies(r xpv1.ManagementPolicies) {
	mg.Spec.ResourceSpec.ManagementPolicies = r
}

// SetProviderConfigReference of this Workspace.
func (mg *Workspace) SetProviderConfigReference(r *xpv1.Reference) {
	mg.Spec.ResourceSpec.ProviderConfigReference = r
}

// SetPublishConnectionDetailsTo of this Workspace.
func (mg *Workspace) SetPublishConnectionDetailsTo(r *xpv1.PublishConnectionDetailsTo) {
	mg.Spec.ResourceSpec.PublishConnectionDetailsTo = r
}

// SetWriteConnectionSecretToReference of this Workspace.
func (mg *Workspace) SetWriteConnectionSecretToReference(r *xpv1.SecretReference) {
	mg.Spec.ResourceSpec.WriteConnectionSecretToReference = r
}

// WorkspaceGroupKind is the GroupKind for the Workspace resource.
var WorkspaceGroupKind = schema.GroupKind{
	Group: Group,
	Kind:  "Workspace",
}

// WorkspaceGroupVersionKind is the GroupVersionKind for the Workspace resource.
var WorkspaceGroupVersionKind = schema.GroupVersionKind{
	Group:   Group,
	Version: Version,
	Kind:    "Workspace",
}
