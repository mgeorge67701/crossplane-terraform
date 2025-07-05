package controller

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mgeorge67701/crossplane-terraform/apis/terraform/v1alpha1"
)

const (
	errNotTerraform = "managed resource is not a Terraform custom resource"
	errTrackPCUsage = "cannot track ProviderConfig usage"
	errGetPC        = "cannot get ProviderConfig"
	errGetCreds     = "cannot get credentials"
	errNewClient    = "cannot create new Service"
	errInitTF       = "cannot initialize Terraform"
	errPlanTF       = "cannot plan Terraform"
	errApplyTF      = "cannot apply Terraform"
	errDestroyTF    = "cannot destroy Terraform"
	errWriteConfig  = "cannot write Terraform configuration"
)

// A TerraformService manages Terraform configurations.
type TerraformService struct {
	workDir string
}

// A TerraformConnector is expected to produce a TerraformService when its Connect method
// is called.
type TerraformConnector struct{}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *TerraformConnector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	// Create a working directory for this Terraform configuration with secure permissions
	workDir := fmt.Sprintf("/tmp/terraform-%s", mg.GetName())
	if err := os.MkdirAll(workDir, 0700); err != nil { // Changed from 0755 to 0700
		return nil, errors.Wrap(err, "cannot create working directory")
	}

	return &TerraformExternal{
		service: &TerraformService{workDir: workDir},
	}, nil
}

// An TerraformExternal observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type TerraformExternal struct {
	service *TerraformService
}

func (c *TerraformExternal) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Terraform)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotTerraform)
	}

	// Write the Terraform configuration to a file with secure permissions
	configPath := filepath.Join(c.service.workDir, "main.tf")
	if err := os.WriteFile(configPath, cr.Spec.ForProvider.Configuration.Raw, 0600); err != nil { // Changed from 0644 to 0600
		return managed.ExternalObservation{}, errors.Wrap(err, errWriteConfig)
	}

	// Write backend configuration if specified
	if err := c.writeBackendConfig(cr); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "cannot write backend configuration")
	}

	// Write variables configuration if specified
	if err := c.writeVariablesConfig(cr); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "cannot write variables configuration")
	}

	// Initialize Terraform executor
	tf, err := tfexec.NewTerraform(c.service.workDir, "terraform")
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errNewClient)
	}

	// Initialize Terraform
	if err := tf.Init(ctx); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errInitTF)
	}

	// Check if the configuration has been applied
	state, err := tf.Show(ctx)
	if err != nil {
		// If show fails, likely no state exists yet
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	// If we have state, the resource exists
	resourceExists := state != nil

	return managed.ExternalObservation{
		ResourceExists:    resourceExists,
		ResourceUpToDate:  resourceExists, // For now, assume it's up to date if it exists
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *TerraformExternal) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Terraform)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotTerraform)
	}

	// Write the Terraform configuration to a file with secure permissions
	configPath := filepath.Join(c.service.workDir, "main.tf")
	if err := os.WriteFile(configPath, cr.Spec.ForProvider.Configuration.Raw, 0600); err != nil { // Changed from 0644 to 0600
		return managed.ExternalCreation{}, errors.Wrap(err, errWriteConfig)
	}

	// Write backend configuration if specified
	if err := c.writeBackendConfig(cr); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "cannot write backend configuration")
	}

	// Write variables configuration if specified
	if err := c.writeVariablesConfig(cr); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "cannot write variables configuration")
	}

	// Initialize Terraform executor
	tf, err := tfexec.NewTerraform(c.service.workDir, "terraform")
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errNewClient)
	}

	// Initialize Terraform
	if err := tf.Init(ctx); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errInitTF)
	}

	// Plan the changes
	hasChanges, err := tf.Plan(ctx)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errPlanTF)
	}

	// Apply the configuration if there are changes
	if hasChanges {
		if err := tf.Apply(ctx); err != nil {
			return managed.ExternalCreation{}, errors.Wrap(err, errApplyTF)
		}
	}

	return managed.ExternalCreation{
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *TerraformExternal) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Terraform)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotTerraform)
	}

	// Write the Terraform configuration to a file with secure permissions
	configPath := filepath.Join(c.service.workDir, "main.tf")
	if err := os.WriteFile(configPath, cr.Spec.ForProvider.Configuration.Raw, 0600); err != nil { // Changed from 0644 to 0600
		return managed.ExternalUpdate{}, errors.Wrap(err, errWriteConfig)
	}

	// Write backend configuration if specified
	if err := c.writeBackendConfig(cr); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "cannot write backend configuration")
	}

	// Write variables configuration if specified
	if err := c.writeVariablesConfig(cr); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "cannot write variables configuration")
	}

	// Initialize Terraform executor
	tf, err := tfexec.NewTerraform(c.service.workDir, "terraform")
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errNewClient)
	}

	// Initialize Terraform
	if err := tf.Init(ctx); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errInitTF)
	}

	// Plan the changes
	hasChanges, err := tf.Plan(ctx)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errPlanTF)
	}

	// Apply the configuration if there are changes
	if hasChanges {
		if err := tf.Apply(ctx); err != nil {
			return managed.ExternalUpdate{}, errors.Wrap(err, errApplyTF)
		}
	}

	return managed.ExternalUpdate{
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *TerraformExternal) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Terraform)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotTerraform)
	}

	// Write the Terraform configuration to a file with secure permissions
	configPath := filepath.Join(c.service.workDir, "main.tf")
	if err := os.WriteFile(configPath, cr.Spec.ForProvider.Configuration.Raw, 0600); err != nil { // Changed from 0644 to 0600
		return managed.ExternalDelete{}, errors.Wrap(err, errWriteConfig)
	}

	// Write backend configuration if specified
	if err := c.writeBackendConfig(cr); err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, "cannot write backend configuration")
	}

	// Write variables configuration if specified
	if err := c.writeVariablesConfig(cr); err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, "cannot write variables configuration")
	}

	// Initialize Terraform executor
	tf, err := tfexec.NewTerraform(c.service.workDir, "terraform")
	if err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, errNewClient)
	}

	// Initialize Terraform
	if err := tf.Init(ctx); err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, errInitTF)
	}

	// Destroy the configuration
	if err := tf.Destroy(ctx); err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, errDestroyTF)
	}

	// Clean up the working directory
	if err := os.RemoveAll(c.service.workDir); err != nil {
		// Log the error but don't fail the deletion
		fmt.Printf("Warning: failed to clean up working directory %s: %v\n", c.service.workDir, err)
	}

	return managed.ExternalDelete{}, nil
}

func (c *TerraformExternal) Disconnect(ctx context.Context) error {
	// Nothing to disconnect for this implementation
	return nil
}

// writeBackendConfig writes the Terraform backend configuration to a file
func (c *TerraformExternal) writeBackendConfig(cr *v1alpha1.Terraform) error {
	if cr.Spec.ForProvider.Backend == nil {
		return nil // No backend configuration specified
	}

	backendPath := filepath.Join(c.service.workDir, "backend.tf")

	// Build backend configuration
	var backendConfig strings.Builder
	backendConfig.WriteString(fmt.Sprintf("terraform {\n  backend \"%s\" {\n", cr.Spec.ForProvider.Backend.Type))

	// Add backend configuration parameters
	for key, value := range cr.Spec.ForProvider.Backend.Configuration {
		backendConfig.WriteString(fmt.Sprintf("    %s = \"%s\"\n", key, value))
	}

	backendConfig.WriteString("  }\n}\n")

	// Write backend configuration with secure permissions
	return os.WriteFile(backendPath, []byte(backendConfig.String()), 0600)
}

// writeVariablesConfig writes Terraform variables to a tfvars file
func (c *TerraformExternal) writeVariablesConfig(cr *v1alpha1.Terraform) error {
	if len(cr.Spec.ForProvider.Variables) == 0 {
		return nil // No variables specified
	}

	varsPath := filepath.Join(c.service.workDir, "terraform.tfvars")

	var varsConfig strings.Builder
	for key, value := range cr.Spec.ForProvider.Variables {
		varsConfig.WriteString(fmt.Sprintf("%s = \"%s\"\n", key, value))
	}

	// Write variables with secure permissions
	return os.WriteFile(varsPath, []byte(varsConfig.String()), 0600)
}

// SetupTerraform adds a controller that reconciles Terraform managed resources.
func SetupTerraform(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.TerraformGroupKind.Kind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.TerraformGroupVersionKind),
		managed.WithExternalConnecter(&TerraformConnector{}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1alpha1.Terraform{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}
