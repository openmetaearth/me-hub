#!/bin/bash
# Prepare Docker build environment
# This script ensures vendor directory is ready for Docker build with private dependencies

set -e

echo "=========================================="
echo "Preparing Docker build environment"
echo "=========================================="

# Check if we're in the correct directory
if [ ! -f "go.mod" ]; then
    echo "Error: go.mod not found. Please run this script from the project root."
    exit 1
fi

# Clean old vendor directory if exists
if [ -d "vendor" ]; then
    echo "Removing old vendor directory..."
    rm -rf vendor
fi

# Create vendor directory with all dependencies
echo "Running go mod vendor..."
go mod vendor

if [ ! -d "vendor" ]; then
    echo "Error: Failed to create vendor directory"
    exit 1
fi

# Check vendor directory size
VENDOR_SIZE=$(du -sh vendor | cut -f1)
echo "Vendor directory created successfully (size: ${VENDOR_SIZE})"

# Count vendor modules
VENDOR_COUNT=$(find vendor -maxdepth 2 -type d | wc -l)
echo "Vendor modules count: ${VENDOR_COUNT}"

# Verify critical dependencies are vendored
echo ""
echo "Verifying critical dependencies..."

CRITICAL_DEPS=(
    "github.com/cosmos/cosmos-sdk"
    "github.com/CosmWasm/wasmvm"
    "github.com/cometbft/cometbft"
)

MISSING_DEPS=0
for dep in "${CRITICAL_DEPS[@]}"; do
    if [ -d "vendor/${dep}" ]; then
        echo "✓ ${dep}"
    else
        echo "✗ ${dep} - MISSING"
        MISSING_DEPS=$((MISSING_DEPS + 1))
    fi
done

if [ $MISSING_DEPS -gt 0 ]; then
    echo ""
    echo "Warning: Some critical dependencies are missing from vendor directory"
    echo "This might cause build issues"
fi

echo ""
echo "=========================================="
echo "Vendor directory ready for Docker build!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "  1. Build Docker image: make docker-private-net"
echo "  2. Or run directly:    docker build -f Dockerfile.private-net -t me-hub/private-net:latest ."
echo ""
