# Crossplane Provider Troubleshooting Guide

## ArgoCD Issues: Package Installation Problems

### Current Issues Observed:
- `current package revision is unhealthy`
- `current package revision health is unknown`
- `cannot apply package revision: existing object is not controlled by UID`

### Step-by-Step Troubleshooting

## 1. Check Provider Status

```bash
# Check all providers
kubectl get providers

# Check specific provider
kubectl get provider provider-terraform -o yaml

# Check provider revision
kubectl get providerrevision
```

## 2. Check Package Installation

```bash
# Check package installations
kubectl get packageinstallation

# Check package revisions
kubectl get packagerevision

# Get detailed status
kubectl describe provider provider-terraform
```

## 3. Check Provider Pod Status

```bash
# Check provider pods
kubectl get pods -n crossplane-system

# Check provider logs
kubectl logs -n crossplane-system deployment/provider-terraform -f

# Check for crash loops
kubectl describe pod -n crossplane-system -l pkg.crossplane.io/provider=provider-terraform
```

## 4. Common Fixes

### Fix 1: Clean Install (Recommended)

```bash
# Remove existing provider
kubectl delete provider provider-terraform

# Wait for cleanup
kubectl get provider

# Reinstall with correct version
kubectl apply -f - <<EOF
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-terraform
spec:
  package: xpkg.upbound.io/mgeorge67701/provider-crossplane-terraform:v3.0.0
EOF
```

### Fix 2: Force Update Package

```bash
# Update with revision policy
kubectl patch provider provider-terraform --type merge -p '{"spec":{"revisionActivationPolicy":"Automatic","revisionHistoryLimit":1}}'
```

### Fix 3: Check Resource Conflicts

```bash
# Check for conflicting CRDs
kubectl get crd | grep terraform

# Check for conflicting resources
kubectl get terraforms.terraform.crossplane.io -A
kubectl get workspaces.terraform.crossplane.io -A
```

### Fix 4: Verify Package Registry

```bash
# Test package availability
docker pull xpkg.upbound.io/mgeorge67701/provider-crossplane-terraform:v3.0.0

# Or check with Upbound CLI
up xpkg dep xpkg.upbound.io/mgeorge67701/provider-crossplane-terraform:v3.0.0
```

## 5. Debug ArgoCD-Specific Issues

### Check ArgoCD Application

```bash
# Check ArgoCD application status
kubectl get application -n argocd

# Check application details
kubectl describe application <your-app-name> -n argocd
```

### ArgoCD Sync Issues

```bash
# Force sync
argocd app sync <your-app-name>

# Prune resources
argocd app sync <your-app-name> --prune

# Force refresh
argocd app refresh <your-app-name>
```

## 6. Check Crossplane Core

```bash
# Check Crossplane core status
kubectl get pods -n crossplane-system

# Check core logs
kubectl logs -n crossplane-system deployment/crossplane -f

# Check core version
kubectl get deployment -n crossplane-system crossplane -o yaml | grep image:
```

## 7. Validation Commands

After fixing, validate with:

```bash
# Check provider is healthy
kubectl get provider provider-terraform -o jsonpath='{.status.conditions[?(@.type=="Healthy")].status}'

# Should return: True

# Check installed condition
kubectl get provider provider-terraform -o jsonpath='{.status.conditions[?(@.type=="Installed")].status}'

# Should return: True

# Check provider config
kubectl get providerconfig

# Test with a simple terraform resource
kubectl apply -f examples/crossplane-provider/terraform-s3-bucket.yaml
```

## 8. Expected Healthy State

When working correctly, you should see:

```bash
$ kubectl get provider
NAME                AGE     INSTALLED   HEALTHY   PACKAGE
provider-terraform  1m      True        True      xpkg.upbound.io/mgeorge67701/provider-crossplane-terraform:v3.0.0

$ kubectl get pods -n crossplane-system
NAME                                         READY   STATUS    RESTARTS   AGE
crossplane-*                                 1/1     Running   0          5m
provider-terraform-*                         1/1     Running   0          1m
```

## 9. Recovery Steps

If all else fails:

```bash
# Nuclear option - clean slate
kubectl delete provider provider-terraform
kubectl delete crd terraforms.terraform.crossplane.io
kubectl delete crd workspaces.terraform.crossplane.io  
kubectl delete crd providerconfigs.terraform.crossplane.io
kubectl delete crd providerconfigusages.terraform.crossplane.io

# Wait for cleanup
sleep 30

# Reinstall
kubectl apply -f - <<EOF
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-terraform
spec:
  package: xpkg.upbound.io/mgeorge67701/provider-crossplane-terraform:v3.0.0
EOF
```

## 10. Prevention

To prevent future issues:

1. **Use specific versions** in your manifests (not `latest`)
2. **Test in staging** before production
3. **Monitor provider health** with alerts
4. **Regular cleanup** of old package revisions
5. **Use proper ArgoCD sync policies**

Example ArgoCD Application with proper sync policy:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: crossplane-terraform-provider
spec:
  project: default
  source:
    repoURL: https://github.com/mgeorge67701/crossplane-terraform
    targetRevision: HEAD
    path: examples/crossplane-provider
  destination:
    server: https://kubernetes.default.svc
    namespace: crossplane-system
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
    - RespectIgnoreDifferences=true
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
```
