#!/bin/bash

# Crossplane Provider Diagnosis Script
# Run this to diagnose the current state of your Crossplane provider

echo "=== Crossplane Provider Diagnosis ==="
echo "Date: $(date)"
echo ""

echo "1. Checking all providers:"
kubectl get providers
echo ""

echo "2. Checking provider revisions:"
kubectl get providerrevisions
echo ""

echo "3. Checking package installations:"
kubectl get packageinstallations
echo ""

echo "4. Checking package revisions:"
kubectl get packagerevisions
echo ""

echo "5. Checking Crossplane system pods:"
kubectl get pods -n crossplane-system
echo ""

echo "6. Checking CRDs related to terraform:"
kubectl get crd | grep terraform
echo ""

echo "7. Checking if there are any terraform resources:"
kubectl get terraforms.terraform.crossplane.io -A 2>/dev/null || echo "No terraform resources found"
echo ""

echo "8. Checking ArgoCD applications (if ArgoCD is installed):"
kubectl get applications -A 2>/dev/null || echo "ArgoCD not found or no applications"
echo ""

echo "9. Checking events in crossplane-system namespace:"
kubectl get events -n crossplane-system --sort-by='.lastTimestamp' | tail -10
echo ""

echo "10. Checking for any provider-related resources:"
kubectl get all -n crossplane-system | grep -i terraform
echo ""

echo "=== Diagnosis Complete ==="
echo "Please share this output for further analysis."
