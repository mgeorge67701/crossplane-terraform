package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/feature"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"gopkg.in/alecthomas/kingpin.v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/mgeorge67701/crossplane-terraform/apis/terraform/v1alpha1"
	terraformcontroller "github.com/mgeorge67701/crossplane-terraform/internal/controller"
)

func main() {
	var (
		app                        = kingpin.New(filepath.Base(os.Args[0]), "Terraform provider for Crossplane").DefaultEnvars()
		debug                      = app.Flag("debug", "Run with debug logging.").Short('d').Bool()
		syncInterval               = app.Flag("sync", "Sync interval controls how often all resources will be double checked for drift.").Short('s').Default("1h").Duration()
		pollInterval               = app.Flag("poll", "Poll interval controls how often an individual resource should be checked for drift.").Default("10m").Duration()
		leaderElection             = app.Flag("leader-election", "Use leader election for the controller manager.").Short('l').Default("false").Bool()
		maxReconcileRate           = app.Flag("max-reconcile-rate", "The global maximum rate per second at which resources may be checked for drift from the desired state.").Default("10").Int()
		enableExternalSecretStores = app.Flag("enable-external-secret-stores", "Enable support for ExternalSecretStores.").Default("false").Bool()
		enableManagementPolicies   = app.Flag("enable-management-policies", "Enable support for ManagementPolicies.").Default("true").Bool()
		essTLSCertsPath            = app.Flag("ess-tls-cert-dir", "Path of ESS TLS certificates.").String()
	)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	zl := zap.New(zap.UseDevMode(*debug))
	log := logging.NewLogrLogger(zl.WithName("provider-terraform"))
	if *debug {
		// The controller-runtime runs with a no-op logger by default. It is
		// *very* verbose even at info level, so we only provide it a real
		// logger when we're running in debug mode.
		ctrl.SetLogger(zl)
	}

	// currently, we configure the jitter to be the 5% of the poll interval
	pollJitter := time.Duration(float64(*pollInterval) * 0.05)
	log.Debug("Starting", "sync-interval", syncInterval.String(),
		"poll-interval", pollInterval.String(), "poll-jitter", pollJitter, "max-reconcile-rate", *maxReconcileRate)

	cfg, err := ctrl.GetConfig()
	if err != nil {
		log.Info("Cannot get config", "error", err)
		os.Exit(1)
	}

	// Get the native types scheme.
	scheme, err := v1alpha1.SchemeBuilder.Build()
	if err != nil {
		log.Info("Cannot build scheme", "error", err)
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:           scheme,
		Metrics:          server.Options{BindAddress: "0"},
		LeaderElection:   *leaderElection,
		LeaderElectionID: "crossplane-terraform-provider-leader-election-helper",
		Cache:            cache.Options{SyncPeriod: syncInterval},
		Logger:           ctrl.Log.WithName("manager"),
	})
	if err != nil {
		log.Info("Cannot create manager", "error", err)
		os.Exit(1)
	}

	// Simple controller options for now
	var o controller.Options
	o.Logger = log
	o.PollInterval = *pollInterval

	if *enableManagementPolicies {
		o.Features = &feature.Flags{}
		o.Features.Enable(feature.EnableBetaManagementPolicies)
	}

	if *enableExternalSecretStores {
		o.ESSOptions = &controller.ESSOptions{}
		if *essTLSCertsPath != "" {
			o.ESSOptions.TLSSecretName = essTLSCertsPath
		}
	}

	// Setup controllers
	if err := terraformcontroller.SetupTerraform(mgr, o); err != nil {
		log.Info("Cannot setup Terraform controller", "error", err)
		os.Exit(1)
	}

	if err := terraformcontroller.SetupWorkspace(mgr, o); err != nil {
		log.Info("Cannot setup Workspace controller", "error", err)
		os.Exit(1)
	}

	if err := terraformcontroller.SetupProviderConfig(mgr, o); err != nil {
		log.Info("Cannot setup ProviderConfig controller", "error", err)
		os.Exit(1)
	}

	log.Info("Starting manager - Crossplane Terraform Provider")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Info("Cannot start manager", "error", err)
		os.Exit(1)
	}
}
