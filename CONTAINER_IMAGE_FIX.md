# Crossplane Terraform Provider - Container Image Fix

## Problem Summary

The provider pod was failing with `CreateContainerError` and the message `"no command specified"`. This was happening because:

1. **Missing Dockerfile**: The provider package didn't have a proper Dockerfile to build the container image
2. **Wrong entrypoint**: The container image didn't have the correct entrypoint defined
3. **Version mismatch**: The crossplane.yaml was using the wrong image version

## Root Cause

The error `"no command specified"` occurs when:
- The container image doesn't have a proper `ENTRYPOINT` or `CMD` defined
- The binary path is incorrect or the binary doesn't exist in the container
- The container image is built incorrectly

## Fix Applied

### 1. Added Dockerfile

Created `/Dockerfile` with:
- Multi-stage build using golang:1.21-alpine
- Proper binary compilation with CGO_ENABLED=0
- Distroless base image for security
- Correct entrypoint: `/usr/local/bin/provider`
- Non-root user for security

### 2. Updated CI/CD Pipeline

Modified `.github/workflows/ci.yml` to:
- Build Docker image with proper entrypoint
- Push to Upbound registry before building the package
- Use correct authentication for Upbound registry

### 3. Fixed Package Configuration

Updated `package/crossplane.yaml` to:
- Use the correct image version (latest)
- Reference the properly built container image

### 4. Updated Documentation

Modified `README.md` to:
- Reference the correct version (v3.0.1)
- Include troubleshooting information

## Testing the Fix

### Step 1: Wait for CI/CD Pipeline
After pushing tag `v3.0.2`, wait for the GitHub Actions workflow to complete:
- Check: https://github.com/mgeorge67701/crossplane-terraform/actions
- Look for the workflow run with tag `v3.0.2`
- All jobs should complete successfully

### Step 2: Delete and Reinstall Provider
```bash
# Delete the existing failing provider
kubectl delete provider mgeorge67701-provider-crossplane-terraform

# Wait for cleanup
kubectl get providers --watch

# Reinstall with the new version
kubectl apply -f - <<EOF
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-terraform
spec:
  package: xpkg.upbound.io/mgeorge67701/provider-crossplane-terraform:v3.0.2
EOF
```

### Step 3: Monitor Provider Health
```bash
# Check provider status
kubectl get providers

# Check pod status
kubectl get pods -n crossplane-system | grep terraform

# Check pod logs
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-terraform

# Describe provider for detailed status
kubectl describe provider provider-terraform
```

## Expected Results After Fix

1. **Container starts successfully**: No more `CreateContainerError`
2. **Provider becomes healthy**: Status shows `Installed: True` and `Healthy: True`
3. **Pod runs without errors**: Pod status shows `Running` with `1/1` ready
4. **Binary executes correctly**: Logs show provider starting up correctly

## Verification Commands

```bash
# Check if the provider is healthy
kubectl get provider provider-terraform -o jsonpath='{.status.conditions[?(@.type=="Healthy")].status}'

# Check if the provider is installed
kubectl get provider provider-terraform -o jsonpath='{.status.conditions[?(@.type=="Installed")].status}'

# Check pod logs for startup messages
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-terraform --tail=20
```

## Container Image Structure

The fixed Docker image has:
- **Base**: `gcr.io/distroless/static:nonroot` (secure, minimal)
- **Binary**: `/usr/local/bin/provider` (compiled Go binary)
- **Entrypoint**: `["/usr/local/bin/provider"]`
- **User**: `65532:65532` (non-root)
- **Architecture**: Multi-arch (amd64, arm64)

## Next Steps

1. **Monitor the CI/CD pipeline**: Ensure v3.0.2 builds successfully
2. **Test the new provider**: Install and verify it works
3. **Create GitHub release**: Create a proper release for v3.0.2
4. **Update documentation**: Include the fix in the troubleshooting guide

This fix addresses the fundamental issue with the container image configuration and should resolve the provider deployment problems.
