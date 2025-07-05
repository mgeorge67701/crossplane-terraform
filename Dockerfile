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

# Use distroless as minimal base image to package the provider binary
FROM gcr.io/distroless/static:nonroot

# Copy the binary from build stage
COPY --from=build /workspace/bin/provider /usr/local/bin/provider

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/provider"]

# Add some metadata
LABEL name="provider-crossplane-terraform" \
      vendor="Mohan George" \
      version="v3.1.7" \
      release="1" \
      summary="Crossplane Terraform Provider" \
      description="A Crossplane provider for managing Terraform resources with native Kubernetes APIs"

# Use non-root user
USER 65532:65532
