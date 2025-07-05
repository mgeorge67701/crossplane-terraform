#!/bin/bash

# Test Docker image build locally
echo "Building Docker image locally..."

# Build the image
docker build -t provider-crossplane-terraform:test .

# Test the image
echo "Testing the image..."
docker run --rm provider-crossplane-terraform:test --help

echo "Docker image test complete!"
