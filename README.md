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
    terraformVersion: "1.9.3"
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
  package: xpkg.upbound.io/mgeorge67701/provider-crossplane-terraform:v0.1.0
EOF
```

## Automatic Publishing Setup

This provider includes GitHub Actions for automatic building and publishing to the Upbound Marketplace. Here's exactly what you need to do:

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
```bash
# 1. Create and push a version tag
git tag v0.2.0
git push origin v0.2.0

# 2. Create a GitHub release (this triggers automatic publishing)
gh release create v0.2.0 --title "v0.2.0" --notes "New features and bug fixes"
```
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

- **GitHub Actions**: Check the Actions tab for build status
- **Upbound Console**: Verify packages appear in your repositories
- **Marketplace**: Confirm public visibility in Upbound Marketplace

### Quick Setup Checklist

- [ ] Add `UPBOUND_ACCESS_ID` secret to GitHub
- [ ] Add `UPBOUND_TOKEN` secret to GitHub  
- [ ] Push changes to trigger first build test
- [ ] Create a release tag to test full publishing pipeline
- [ ] Verify provider appears in Upbound Marketplace

### Configure Provider

Create a ProviderConfig to configure the provider:

```yaml
apiVersion: terraform.crossplane.io/v1alpha1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: terraform-creds
      key: credentials
  terraformVersion: "1.9.3"
  backend:
    type: "s3"
    configuration:
      bucket: "my-terraform-state"
      region: "us-west-2"
```

Create the credentials secret:

```bash
kubectl create secret generic terraform-creds -n crossplane-system \
  --from-literal=credentials='{"aws_access_key_id": "your-key", "aws_secret_access_key": "your-secret"}'
```

## Configuration

### Backend Configuration

The provider supports various Terraform backends:

#### S3 Backend

```yaml
backend:
  type: "s3"
  configuration:
    bucket: "my-terraform-state"
    key: "path/to/terraform.tfstate"
    region: "us-west-2"
    encrypt: "true"
```

#### GCS Backend

```yaml
backend:
  type: "gcs"
  configuration:
    bucket: "my-terraform-state"
    prefix: "path/to/terraform.tfstate"
```

#### Azure Backend

```yaml
backend:
  type: "azurerm"
  configuration:
    resource_group_name: "my-rg"
    storage_account_name: "mystorageaccount"
    container_name: "terraform-state"
    key: "terraform.tfstate"
```

### Variable Management

Variables can be provided in multiple ways:

1. **Direct variables** in the resource spec
2. **Environment variables** for sensitive data
3. **ConfigMaps and Secrets** for shared configuration

```yaml
spec:
  forProvider:
    variables:
      region: "us-west-2"
      instance_type: "t3.micro"
    environment:
      TF_VAR_database_password: "secret-password"
```

### Source Configuration

The provider supports multiple source types:

#### Inline Configuration

```yaml
spec:
  forProvider:
    configuration: |
      resource "aws_s3_bucket" "example" {
        bucket = "my-bucket"
      }
```

#### Git Repository

```yaml
spec:
  forProvider:
    source:
      git:
        url: "https://github.com/myorg/terraform-modules.git"
        ref: "v1.0.0"
        path: "modules/s3-bucket"
```

#### HTTP Source

```yaml
spec:
  forProvider:
    source:
      http:
        url: "https://example.com/terraform-module.tar.gz"
        checksum: "sha256:abc123..."
```

## Examples

### Basic S3 Bucket

```yaml
apiVersion: terraform.crossplane.io/v1alpha1
kind: Terraform
metadata:
  name: s3-bucket-example
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
        bucket = var.bucket_name
      }
      
      resource "aws_s3_bucket_versioning" "example" {
        bucket = aws_s3_bucket.example.id
        versioning_configuration {
          status = "Enabled"
        }
      }
    variables:
      bucket_name: "my-crossplane-bucket"
  providerConfigRef:
    name: default
```

### Multi-Resource Infrastructure

```yaml
apiVersion: terraform.crossplane.io/v1alpha1
kind: Terraform
metadata:
  name: vpc-infrastructure
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
      
      resource "aws_vpc" "main" {
        cidr_block           = "10.0.0.0/16"
        enable_dns_hostnames = true
        enable_dns_support   = true
        
        tags = {
          Name = "main-vpc"
        }
      }
      
      resource "aws_subnet" "public" {
        count             = 2
        vpc_id            = aws_vpc.main.id
        cidr_block        = "10.0.${count.index + 1}.0/24"
        availability_zone = data.aws_availability_zones.available.names[count.index]
        
        map_public_ip_on_launch = true
        
        tags = {
          Name = "public-subnet-${count.index + 1}"
        }
      }
      
      data "aws_availability_zones" "available" {
        state = "available"
      }
    backend:
      type: "s3"
      configuration:
        bucket: "my-terraform-state"
        key: "vpc/terraform.tfstate"
        region: "us-west-2"
  providerConfigRef:
    name: default
```

## Development

### Building the Provider

```bash
# Build the provider binary
make build

# Run tests
make test

# Build Docker image
make docker-build

# Generate CRDs
make generate-crds
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## Troubleshooting

### Common Issues

**Provider not starting**
- Check that Crossplane is installed and running
- Verify the provider package was applied correctly
- Check the provider pod logs

**Authentication errors**
- Verify the credentials secret exists and has the correct keys
- Check that the ProviderConfig references the correct secret
- Ensure the service account has necessary permissions

**Terraform execution errors**
- Check the Terraform configuration syntax
- Verify all required variables are provided
- Check the backend configuration

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
