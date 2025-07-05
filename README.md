# Crossplane Terraform Provider

A Crossplane provider for managing Terraform resources natively in Kubernetes.

## Overview

This provider enables you to manage Terraform configurations and workspaces using Kubernetes APIs through Crossplane. It provides native Kubernetes resources that map to Terraform concepts, allowing you to use GitOps workflows and Kubernetes-native tooling to manage your Terraform infrastructure.

## Features

- **Native Kubernetes API**: Manage Terraform resources using standard Kubernetes resources
- **GitOps Ready**: Declarative configuration management with version control
- **Workspace Management**: Full support for Terraform workspaces and state isolation
- **Crossplane Integration**: Leverage Crossplane's composition and dependency management
- **Multi-Backend Support**: Support for various Terraform backends (S3, GCS, Azure, etc.)
- **Secure Credential Management**: Integration with Kubernetes secrets and external secret stores

## Managed Resources

### Terraform

The `Terraform` resource allows you to manage Terraform configurations as Kubernetes resources.

```yaml
apiVersion: terraform.crossplane.io/v1alpha1
kind: Terraform
metadata:
  name: my-infrastructure
spec:
  forProvider:
    configuration: |
      terraform {
        required_providers {
          aws = {
            source  = "hashicorp/aws"
            version = "~> 5.0"
          }
        }
      }
      
      resource "aws_s3_bucket" "example" {
        bucket = "my-example-bucket"
      }
      
      output "bucket_name" {
        value = aws_s3_bucket.example.bucket
      }
    variables:
      region: "us-west-2"
      environment: "production"
    backend:
      type: "s3"
      configuration:
        bucket: "my-terraform-state"
        key: "infrastructure/terraform.tfstate"
        region: "us-west-2"
  providerConfigRef:
    name: default
```

### Workspace

The `Workspace` resource provides advanced workspace management capabilities.

```yaml
apiVersion: terraform.crossplane.io/v1alpha1
kind: Workspace
metadata:
  name: production-workspace
spec:
  forProvider:
    name: production
    variables:
      environment: "production"
      region: "us-west-2"
    autoApply: true
    terraformVersion: "1.12.2"
    workingDirectory: "/terraform/modules/infrastructure"
  providerConfigRef:
    name: default
```

## Installation

### Prerequisites

- Kubernetes cluster with Crossplane installed
- kubectl configured to access your cluster

### Install the Provider

```bash
# Install the provider
kubectl apply -f - <<EOF
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-terraform
spec:
  package: xpkg.upbound.io/mgeorge67701/provider-crossplane-terraform:v3.0.1
EOF
```

## Automatic Publishing Setup

This provider includes a **modernized GitHub Actions CI/CD pipeline** for automatic building and publishing to the Upbound Marketplace. The pipeline has been updated to use the latest actions, resolve all deprecation warnings, and ensure reliable cross-platform builds.

### Pipeline Features

- ✅ **Modern Actions**: Updated to latest GitHub Actions (v4/v5)
- ✅ **Cross-Platform**: Linux and macOS (amd64/arm64) support
- ✅ **Reliable Caching**: Optimized Go module caching
- ✅ **Automatic Publishing**: Seamless Upbound Marketplace integration
- ✅ **Release Assets**: Complete binary distribution with checksums
- ✅ **Zero Deprecation Warnings**: All legacy actions updated

Here's exactly what you need to do:

### 1. Add Upbound Credentials to GitHub Secrets

Go to your GitHub repository → **Settings** → **Secrets and variables** → **Actions**, then add these two secrets:

**Secret 1: Access ID**
- **Name**: `UPBOUND_ACCESS_ID`
- **Value**: Your Upbound Access ID (found in Upbound Console → Profile → Tokens)

**Secret 2: Token**
- **Name**: `UPBOUND_TOKEN`
- **Value**: Your Upbound authentication token (the long JWT token)

### 2. Development Workflow

#### For Daily Development (Safe - No Publishing):
```bash
# Make your changes
git add .
git commit -m "Add new features"
git push origin main
```
**Result**: ✅ Automatic testing and validation, ❌ No publishing to Upbound

#### For Publishing New Versions:

##### Step 1: Create and Push a Version Tag

```bash
# Create a new version tag (increment the version number)
git tag v3.0.1
git push origin v3.0.1
```

##### Step 2: Create a GitHub Release

1. Go to your repository on GitHub: `https://github.com/mgeorge67701/crossplane-terraform`
2. Click **"Releases"** on the right sidebar
3. Click **"Create a new release"**
4. Fill in the release form:
   - **Tag version:** `v3.0.1` (should be pre-selected)
   - **Release title:** `v3.0.1`
   - **Description:** Add your release notes, for example:

     ```text
     ## Changes
     - ✅ New feature: Added workspace management
     - 🐛 Fixed: Binary path issues in CI/CD
     - 📦 Updated: Dependencies to latest versions
     ```

5. Click **"Publish release"**

**Result**: ✅ Automatic testing, building, and publishing to Upbound Marketplace

### 3. What Happens Automatically

When you create a GitHub release, the CI/CD pipeline will:

1. **Test**: Run all tests and validations
2. **Build**: Create cross-platform binaries (Linux, macOS, Windows)
3. **Package**: Generate the Crossplane `.xpkg` package
4. **Publish to GitHub**: Upload release assets with checksums
5. **Publish to Upbound**: Push the package to Upbound Marketplace as:
   - `xpkg.upbound.io/mgeorge67701/provider-crossplane-terraform:v0.2.0`
   - `xpkg.upbound.io/mgeorge67701/provider-crossplane-terraform:latest`
6. **Update Marketplace**: Set repository to public and published

### 4. Version Management

- **Stable releases**: `v1.0.0`, `v1.1.0` → Gets tagged as `latest`
- **Pre-releases**: `v1.0.0-alpha1`, `v1.0.0-beta1` → No `latest` tag
- **Installation**: Users can install specific versions or `latest`

### 5. Monitoring

#### Check Build Status

After creating a release, monitor the progress:

1. **GitHub Actions**: Go to `https://github.com/mgeorge67701/crossplane-terraform/actions`
2. **Look for the latest workflow run** with your release tag (e.g., `v3.0.1`)
3. **Expected jobs to complete**:
   - ✅ **test** (2 jobs) - Go 1.21 and 1.22 testing
   - ✅ **validate-examples** (1 job) - YAML validation
   - ✅ **build** (4 jobs) - Linux AMD64/ARM64, macOS AMD64/ARM64
   - ✅ **package** (1 job) - Create .xpkg package
   - ✅ **release** (1 job) - Upload assets to GitHub release
   - ✅ **publish-to-upbound** (1 job) - Publish to Upbound Marketplace

#### Verify Release Assets

After successful completion:

1. **Check GitHub Release**: `https://github.com/mgeorge67701/crossplane-terraform/releases/latest`
   - Should contain platform-specific tar.gz files
   - Should contain the .xpkg package
   - Should contain SHA256SUMS file

2. **Check Upbound Marketplace**: `https://marketplace.upbound.io/providers/mgeorge67701/provider-crossplane-terraform`
   - Should show your new version
   - Should be marked as "latest" (if not a pre-release)

#### Common Build Issues

##### GitHub Actions Permission Errors

**Error**: `Resource not accessible by integration`
**Solution**: The workflow needs proper permissions. This is already configured in the workflow file with:
```yaml
permissions:
  contents: write
  packages: write
  attestations: write
  id-token: write
```

If you still see this error:
1. Check your repository settings: `Settings → Actions → General`
2. Ensure "Read and write permissions" is selected for Workflow permissions
3. Or manually grant permissions in the workflow file

##### Provider not starting

- Check that Crossplane is installed and running
- Verify the provider package was applied correctly
- Check the provider pod logs

##### Authentication errors

- Verify the credentials secret exists and has the correct keys
- Check that the ProviderConfig references the correct secret
- Ensure the service account has necessary permissions

##### Terraform execution errors

- Check the Terraform configuration syntax
- Verify all required variables are provided
- Check the backend configuration

## Quick Setup Checklist

### One-Time Setup (First Time Only)

- [ ] **Add Upbound credentials to GitHub Secrets**:
  - Go to `https://github.com/mgeorge67701/crossplane-terraform/settings/secrets/actions`
  - Add `UPBOUND_ACCESS_ID` secret
  - Add `UPBOUND_TOKEN` secret
- [ ] **Test the build pipeline**:
  - Push any change to main branch
  - Verify build completes successfully

### For Each New Release

- [ ] **Create version tag**: `git tag v3.0.1 && git push origin v3.0.1`
- [ ] **Create GitHub release**: Go to releases page and create new release
- [ ] **Monitor build**: Check Actions tab for progress
- [ ] **Verify assets**: Check release has all binary files
- [ ] **Confirm publishing**: Check Upbound Marketplace for new version

### Debugging

Enable debug logging:

```bash
kubectl logs -n crossplane-system deployment/provider-terraform -f
```

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.

## Support

- GitHub Issues: [https://github.com/mgeorge67701/crossplane-terraform/issues](https://github.com/mgeorge67701/crossplane-terraform/issues)
- Documentation: [https://github.com/mgeorge67701/crossplane-terraform](https://github.com/mgeorge67701/crossplane-terraform)

## Roadmap

- [ ] Enhanced workspace management
- [ ] Support for Terraform Cloud/Enterprise
- [ ] Advanced state management features
- [ ] Integration with more backends
- [ ] Terraform module composition
- [ ] Advanced GitOps workflows
