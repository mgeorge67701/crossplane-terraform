# Crossplane Terraform Provider

![GitHub Release](https://img.shields.io/github/v/release/mgeorge67701/crossplane-terraform)
![GitHub Actions](https://img.shields.io/github/actions/workflow/status/mgeorge67701/crossplane-terraform/ci.yml)
![Go Version](https://img.shields.io/badge/go-1.24-blue)
![License](https://img.shields.io/badge/license-Apache%202.0-blue)

A modern Crossplane provider for managing Terraform resources natively in Kubernetes with automated CI/CD and marketplace publishing.

## 🚀 Quick Start

### Prerequisites

- Kubernetes cluster (1.23+)
- Crossplane installed ([Installation Guide](https://docs.crossplane.io/latest/software/install/))
- kubectl configured

### Install the Provider

```bash
# Install the provider
kubectl apply -f - <<EOF
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-terraform
spec:
  package: xpkg.upbound.io/mgeorge67701/provider-crossplane-terraform:v3.1.7
EOF
```

### Configure Provider Credentials

```bash
# Create a secret with your cloud credentials
kubectl create secret generic terraform-creds \
  --from-literal=credentials='{"aws_access_key_id": "YOUR_KEY", "aws_secret_access_key": "YOUR_SECRET"}'

# Create a ProviderConfig
kubectl apply -f - <<EOF
apiVersion: terraform.crossplane.io/v1alpha1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: default
      name: terraform-creds
      key: credentials
EOF
```

### Create Your First Terraform Resource

```bash
kubectl apply -f - <<EOF
apiVersion: terraform.crossplane.io/v1alpha1
kind: Terraform
metadata:
  name: my-s3-bucket
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
      
      output "bucket_arn" {
        value = aws_s3_bucket.example.arn
      }
    variables:
      bucket_name: "my-crossplane-bucket-${random_id.suffix.hex}"
    backend:
      type: "s3"
      configuration:
        bucket: "my-terraform-state-bucket"
        key: "crossplane/terraform.tfstate"
        region: "us-west-2"
  providerConfigRef:
    name: default
EOF
```

## 📋 What This Provider Does

- **🔧 Native Kubernetes API**: Manage Terraform resources using standard Kubernetes resources
- **🔄 GitOps Ready**: Declarative configuration management with version control
- **🏗️ Workspace Management**: Full support for Terraform workspaces and state isolation
- **🔗 Crossplane Integration**: Leverage Crossplane's composition and dependency management
- **☁️ Multi-Backend Support**: Support for various Terraform backends (S3, GCS, Azure, etc.)
- **🔐 Secure Credential Management**: Integration with Kubernetes secrets and external secret stores
- **🚀 Automated CI/CD**: Modern GitHub Actions pipeline with automatic publishing

## 📚 Available Resources

### Terraform Resource

Manage complete Terraform configurations as Kubernetes resources:

```yaml
apiVersion: terraform.crossplane.io/v1alpha1
kind: Terraform
metadata:
  name: my-infrastructure
spec:
  forProvider:
    configuration: |
      # Your Terraform configuration here
    variables:
      key: value
    backend:
      type: "s3"
      configuration:
        bucket: "my-state-bucket"
        key: "terraform.tfstate"
        region: "us-west-2"
  providerConfigRef:
    name: default
```

### Workspace Resource

Advanced workspace management with environment isolation:

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
    autoApply: true
    terraformVersion: "1.12.2"
  providerConfigRef:
    name: default
```

## 🔧 Configuration Examples

### AWS S3 Bucket with VPC

```yaml
apiVersion: terraform.crossplane.io/v1alpha1
kind: Terraform
metadata:
  name: aws-infrastructure
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
      
      # Create VPC
      resource "aws_vpc" "main" {
        cidr_block           = var.vpc_cidr
        enable_dns_hostnames = true
        enable_dns_support   = true
        
        tags = {
          Name        = var.vpc_name
          Environment = var.environment
        }
      }
      
      # Create S3 bucket
      resource "aws_s3_bucket" "data" {
        bucket = var.bucket_name
        
        tags = {
          Environment = var.environment
        }
      }
      
      # Outputs
      output "vpc_id" {
        value = aws_vpc.main.id
      }
      
      output "bucket_arn" {
        value = aws_s3_bucket.data.arn
      }
    variables:
      vpc_cidr: "10.0.0.0/16"
      vpc_name: "crossplane-vpc"
      bucket_name: "my-crossplane-data-bucket"
      environment: "production"
    backend:
      type: "s3"
      configuration:
        bucket: "my-terraform-state"
        key: "crossplane/aws-infrastructure.tfstate"
        region: "us-west-2"
  providerConfigRef:
    name: default
```

### Multi-Cloud Setup

```yaml
apiVersion: terraform.crossplane.io/v1alpha1
kind: Terraform
metadata:
  name: multi-cloud-setup
spec:
  forProvider:
    configuration: |
      terraform {
        required_providers {
          aws = {
            source  = "hashicorp/aws"
            version = "~> 5.0"
          }
          azurerm = {
            source  = "hashicorp/azurerm"
            version = "~> 3.0"
          }
        }
      }
      
      # AWS S3 bucket
      resource "aws_s3_bucket" "aws_storage" {
        bucket = var.aws_bucket_name
        
        tags = {
          Provider = "AWS"
        }
      }
      
      # Azure Storage Account
      resource "azurerm_storage_account" "azure_storage" {
        name                     = var.azure_storage_name
        resource_group_name      = var.azure_resource_group
        location                 = var.azure_location
        account_tier             = "Standard"
        account_replication_type = "LRS"
        
        tags = {
          Provider = "Azure"
        }
      }
    variables:
      aws_bucket_name: "my-aws-bucket"
      azure_storage_name: "myazurestorage"
      azure_resource_group: "my-rg"
      azure_location: "East US"
  providerConfigRef:
    name: default
```

## 🔐 Security & Credentials

### Using Kubernetes Secrets

```bash
# Create secret for AWS credentials
kubectl create secret generic aws-creds \
  --from-literal=credentials='{"aws_access_key_id": "AKIAIOSFODNN7EXAMPLE", "aws_secret_access_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"}'

# Create secret for Azure credentials
kubectl create secret generic azure-creds \
  --from-literal=credentials='{"subscription_id": "00000000-0000-0000-0000-000000000000", "client_id": "00000000-0000-0000-0000-000000000000", "client_secret": "your-secret", "tenant_id": "00000000-0000-0000-0000-000000000000"}'
```

### ProviderConfig Examples

```yaml
# AWS ProviderConfig
apiVersion: terraform.crossplane.io/v1alpha1
kind: ProviderConfig
metadata:
  name: aws-config
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: default
      name: aws-creds
      key: credentials
---
# Azure ProviderConfig
apiVersion: terraform.crossplane.io/v1alpha1
kind: ProviderConfig
metadata:
  name: azure-config
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: default
      name: azure-creds
      key: credentials
```

## 🚀 For Developers: Automated CI/CD Pipeline

This provider includes a **fully modernized GitHub Actions CI/CD pipeline** that automatically builds, tests, and publishes releases to the Upbound Marketplace. No manual builds or deployments needed!

### 🚀 Quick Start for Publishing

**One-command release:**

```bash
# Create and publish a new release
git tag v3.2.0 && git push origin v3.2.0
gh release create v3.2.0 --generate-notes
```

**Result:** Automatically builds, tests, and publishes to Upbound Marketplace

**Prerequisites:**

- GitHub CLI installed (`gh`)
- Upbound secrets configured in your GitHub repository
- Version updated in `package/crossplane.yaml`

### Pipeline Features

- ✅ **Latest GitHub Actions**: Updated to actions/checkout@v4, setup-go@v5, docker/build-push-action@v5
- ✅ **Go 1.24**: Latest stable Go version with built-in dependency caching
- ✅ **Multi-Architecture**: Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64)
- ✅ **Automated Testing**: Full test suite with linting (golangci-lint v1.60.1)
- ✅ **Docker Multi-Arch**: Automatic container image publishing to GitHub Container Registry
- ✅ **Marketplace Publishing**: Seamless Upbound Marketplace integration with xpkg format
- ✅ **Release Assets**: Complete binary distribution with SHA256 checksums
- ✅ **Security Scanning**: Automated vulnerability scanning and dependency updates
- ✅ **Zero Deprecation**: All legacy actions and configurations updated

### Setup Instructions for Developers

#### 1. Configure GitHub Secrets

Add these secrets to your GitHub repository (`Settings → Secrets and variables → Actions`):

**Required Secrets:**

- `UPBOUND_ACCESS_ID`: Your Upbound Access ID (e.g., `your-org` or `your-username`)
- `UPBOUND_TOKEN`: Your Upbound authentication token

**Optional Secrets (for enhanced features):**

- `REGISTRY_USERNAME`: Docker registry username (defaults to GitHub username)
- `REGISTRY_PASSWORD`: Docker registry password (defaults to GitHub token)

**Getting Upbound Credentials:**

1. Go to [Upbound Console](https://console.upbound.io/)
2. Navigate to Profile → Tokens
3. Create a new token with "Repository" scope
4. Copy the Access ID and Token to your GitHub secrets

#### 2. Development Workflow

**For regular development (push to main):**

```bash
# Make your changes
git add .
git commit -m "feat: add new functionality"
git push origin main
```

**Result:** Runs tests, linting, and basic validation only

**For publishing new versions (complete automation):**

```bash
# Step 1: Update version in package/crossplane.yaml if needed
# Step 2: Create and push version tag
git tag v3.2.0
git push origin v3.2.0

# Step 3: Create GitHub Release (triggers full pipeline)
gh release create v3.2.0 --title "Release v3.2.0" --notes "
## What's New
- Updated to Terraform 1.12.2
- Enhanced security and multi-arch support
- Improved backend configuration handling

## Installation
\`\`\`bash
kubectl apply -f - <<EOF
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-terraform
spec:
  package: xpkg.upbound.io/mgeorge67701/provider-crossplane-terraform:v3.2.0
EOF
\`\`\`
"
```

**Alternative: Using GitHub Web Interface:**

1. Go to GitHub → Releases → "Create a new release"
2. Choose or create tag `v3.2.0`
3. Add release title and notes
4. Click "Publish release"

**Result:** Full build, test, package, and publish to Upbound Marketplace

#### 3. What Happens Automatically

**On Push to Main Branch:**

1. **🧪 Test**: Runs `go test ./...` with race detection
2. **🔍 Lint**: Runs golangci-lint with comprehensive rules
3. **🏗️ Build**: Validates the build process

**On GitHub Release:**

1. **🧪 Test**: Full test suite with Go 1.24
2. **🔍 Lint**: Comprehensive linting and code quality checks
3. **📦 Build**: Creates binaries for all platforms (Linux, macOS, Windows)
4. **🐳 Docker**: Builds multi-architecture container images (amd64, arm64)
5. **📋 Package**: Generates Crossplane `.xpkg` package with proper metadata
6. **📤 Release**: Uploads binaries and checksums to GitHub release
7. **🌐 Publish**: Pushes package to Upbound Marketplace automatically
8. **🏷️ Tag**: Updates container registry with proper tags

### Complete Automation Features

**Multi-Architecture Support:**

- Linux: amd64, arm64
- macOS: amd64, arm64 (Apple Silicon)
- Windows: amd64
- Container images: linux/amd64, linux/arm64

**Security & Quality:**

- Automated dependency updates via Dependabot
- Vulnerability scanning with GitHub Security
- Code quality checks with golangci-lint
- Secure secret handling in CI/CD

**Package Management:**

- Automatic versioning from Git tags
- Proper semver handling
- Crossplane package validation
- Upbound Marketplace compliance

### Monitoring Your Releases

**GitHub Actions Dashboard:**

- **Status**: Check `Actions` tab for real-time build status
- **Logs**: View detailed logs for each step
- **Artifacts**: Download build artifacts if needed

**Container Registry:**

- **Images**: View at `ghcr.io/mgeorge67701/provider-crossplane-terraform`
- **Tags**: Check available versions and architectures
- **Security**: View vulnerability scan results

**Upbound Marketplace:**

- **Packages**: Monitor at `marketplace.upbound.io`
- **Downloads**: Track package adoption
- **Versions**: Verify published versions

### CI/CD Pipeline Workflow

```mermaid
graph TD
    A[Push to Main] --> B[Run Tests & Lint]
    B --> C[Build Validation]
    
    D[Create GitHub Release] --> E[Full Test Suite]
    E --> F[Multi-Arch Build]
    F --> G[Docker Build & Push]
    G --> H[Generate .xpkg Package]
    H --> I[Upload Release Assets]
    I --> J[Publish to Upbound]
    J --> K[Update Registry Tags]
```

### Advanced Configuration

**Customizing the Pipeline:**

- Edit `.github/workflows/ci.yml` for custom build steps
- Update `package/crossplane.yaml` for package metadata
- Modify `Dockerfile` for container customization

**Local Development:**

```bash
# Test locally before pushing
make test
make lint
make build

# Build Docker image locally
make docker-build

# Generate package locally
make package
```

### Verifying CI/CD Pipeline Success

After creating a GitHub release, you can verify that the automated pipeline worked correctly:

**1. Check GitHub Actions Status:**

- Go to your GitHub repository → Actions tab
- Look for the workflow run triggered by your release
- Verify all steps completed successfully (green checkmarks)

**2. Verify Container Image:**

```bash
# Check if the container image was published
docker pull ghcr.io/mgeorge67701/provider-crossplane-terraform:v3.2.0

# Verify multi-architecture support
docker buildx imagetools inspect ghcr.io/mgeorge67701/provider-crossplane-terraform:v3.2.0
```

**3. Verify Upbound Marketplace:**

- Visit [Upbound Marketplace](https://marketplace.upbound.io/)
- Search for "provider-crossplane-terraform"
- Verify your new version is available

**4. Test the Published Package:**

```bash
# Install the provider using the new version
kubectl apply -f - <<EOF
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-terraform
spec:
  package: xpkg.upbound.io/mgeorge67701/provider-crossplane-terraform:v3.2.0
EOF

# Check provider status
kubectl get providers
```

**5. Verify Release Assets:**

- Go to GitHub → Releases → Your release
- Verify binaries and checksums are attached
- Download and verify checksums if needed

### CI/CD Pipeline Troubleshooting

**Common Issues:**

- **Missing Secrets**: Ensure `UPBOUND_ACCESS_ID` and `UPBOUND_TOKEN` are set
- **Permission Errors**: Verify Upbound token has proper permissions
- **Build Failures**: Check Go version compatibility and dependencies
- **Docker Push Failures**: Verify GitHub Container Registry permissions

**Debug Steps:**

```bash
# Check workflow logs in GitHub Actions
# Look for specific error messages in failed steps
# Verify all required secrets are configured
# Check package/crossplane.yaml for valid metadata
```

## 🔍 Troubleshooting

### Common Issues

#### Provider Not Starting

```bash
# Check provider status
kubectl get providers

# Check provider pods
kubectl get pods -n crossplane-system -l pkg.crossplane.io/provider=provider-terraform

# Check logs
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-terraform
```

#### Authentication Errors

```bash
# Verify secrets exist
kubectl get secrets

# Check secret contents (base64 encoded)
kubectl get secret terraform-creds -o yaml

# Verify ProviderConfig
kubectl get providerconfigs
kubectl describe providerconfig default
```

#### Terraform Execution Errors

```bash
# Check Terraform resource status
kubectl get terraform
kubectl describe terraform my-terraform-resource

# Check detailed logs
kubectl logs -n crossplane-system deployment/provider-terraform -f
```

### Debug Mode

Enable debug logging:

```bash
kubectl patch deployment provider-terraform -n crossplane-system -p '{"spec":{"template":{"spec":{"containers":[{"name":"package-runtime","args":["--debug"]}]}}}}'
```

## 📖 Examples

Check the `examples/` directory for more complete examples:

- [`examples/crossplane-provider/terraform-s3-bucket.yaml`](examples/crossplane-provider/terraform-s3-bucket.yaml) - Simple S3 bucket
- [`examples/crossplane-provider/terraform-vpc.yaml`](examples/crossplane-provider/terraform-vpc.yaml) - VPC with subnets
- [`examples/crossplane-provider/providerconfig.yaml`](examples/crossplane-provider/providerconfig.yaml) - Provider configuration

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **GitHub Issues**: [Report bugs and request features](https://github.com/mgeorge67701/crossplane-terraform/issues)
- **Documentation**: [Project documentation](https://github.com/mgeorge67701/crossplane-terraform)
- **Crossplane Community**: [Join the Crossplane Slack](https://slack.crossplane.io/)

## 🗺️ Roadmap

- [ ] Enhanced workspace management with drift detection
- [ ] Support for Terraform Cloud/Enterprise
- [ ] Advanced state management and backup features
- [ ] Integration with more cloud providers
- [ ] Terraform module composition and reuse
- [ ] Advanced GitOps workflows and policy management
- [ ] Performance optimizations and scaling improvements

---

Built with ❤️ for the Crossplane community
