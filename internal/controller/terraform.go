package controller

import (
	"context"
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
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
)

// A TerraformService does nothing.
type TerraformService struct{}

// A TerraformConnector is expected to produce a TerraformService when its Connect method
// is called.
type TerraformConnector struct{}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *TerraformConnector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	return &TerraformExternal{service: &TerraformService{}}, nil
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

	// These fmt statements should be removed in the real implementation.
	fmt.Printf("Observing: %+v", cr)

	return managed.ExternalObservation{
		// Return false when the external resource does not exist. This lets
		// the managed resource reconciler know that it needs to call Create to
		// (re)create the resource, or that it has successfully been deleted.
		ResourceExists: true,

		// Return false when the external resource exists, but it not up to date
		// with the desired managed resource state. This lets the managed
		// resource reconciler know that it needs to call Update.
		ResourceUpToDate: true,

		// Return any details that may be required to connect to the external
		// resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *TerraformExternal) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Terraform)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotTerraform)
	}

	fmt.Printf("Creating: %+v", cr)

	return managed.ExternalCreation{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *TerraformExternal) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Terraform)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotTerraform)
	}

	fmt.Printf("Updating: %+v", cr)

	return managed.ExternalUpdate{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *TerraformExternal) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Terraform)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotTerraform)
	}

	fmt.Printf("Deleting: %+v", cr)

	return managed.ExternalDelete{}, nil
}

func (c *TerraformExternal) Disconnect(ctx context.Context) error {
	// Nothing to disconnect for this implementation
	return nil
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
