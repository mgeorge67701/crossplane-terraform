FROM --platform=${BUILDPLATFORM} golang:1.24-alpine AS build

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set the working directory
WORKDIR /workspace

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the provider binary
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o bin/provider main.go

# Use alpine as base image to include Terraform binary
FROM alpine:3.21

# Install runtime dependencies
RUN apk --no-cache add ca-certificates bash git curl unzip

# Build arguments for multi-arch support
ARG TARGETOS
ARG TARGETARCH

# Set Terraform version - using latest stable version
ENV TERRAFORM_VERSION=1.12.2
ENV TF_IN_AUTOMATION=1
ENV TF_PLUGIN_CACHE_DIR=/tf/plugin-cache

# Download and install Terraform binary
RUN curl -s -L https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_${TARGETOS}_${TARGETARCH}.zip -o terraform.zip \
    && unzip -d /usr/local/bin terraform.zip \
    && rm terraform.zip \
    && chmod +x /usr/local/bin/terraform \
    && mkdir -p ${TF_PLUGIN_CACHE_DIR} \
    && chown -R 2000 /tf

# Copy the binary from build stage
COPY --from=build /workspace/bin/provider /usr/local/bin/provider

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/provider"]

# Add metadata
LABEL name="provider-crossplane-terraform" \
      vendor="Mohan George" \
      version="v3.1.8" \
      release="1" \
      summary="Crossplane Terraform Provider with Terraform 1.12.2" \
      description="A Crossplane provider for managing Terraform resources with native Kubernetes APIs and latest Terraform support"

# Use non-root user for security
USER 65532:65532
