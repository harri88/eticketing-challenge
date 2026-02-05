#!/bin/bash
# Package script for backend-dotnet-tickets
# Creates a distributable package of the .NET application

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$SCRIPT_DIR/../backend-dotnet-tickets"
OUTPUT_DIR="$SCRIPT_DIR/../dist/backend-dotnet-tickets"
VERSION="${VERSION:-1.0.0}"

echo "=========================================="
echo "Packaging backend-dotnet-tickets"
echo "Version: $VERSION"
echo "=========================================="

# Clean previous build
echo "Cleaning previous builds..."
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Navigate to project directory
cd "$PROJECT_DIR"

echo "Restoring dependencies..."
dotnet restore

echo "Building in Release mode..."
dotnet build --configuration Release --no-restore

echo "Publishing application..."
dotnet publish --configuration Release --output "$OUTPUT_DIR/publish" --no-build

# Create version file
echo "$VERSION" > "$OUTPUT_DIR/VERSION"

# Copy additional files
cp -r sql_scripts "$OUTPUT_DIR/" 2>/dev/null || true
cp README.md "$OUTPUT_DIR/" 2>/dev/null || true
cp Dockerfile "$OUTPUT_DIR/" 2>/dev/null || true

# Create archive
echo "Creating distributable archive..."
cd "$SCRIPT_DIR/../dist"
tar -czf "backend-dotnet-tickets-$VERSION.tar.gz" backend-dotnet-tickets/

echo "=========================================="
echo "Package created successfully!"
echo "Location: $SCRIPT_DIR/../dist/backend-dotnet-tickets-$VERSION.tar.gz"
echo "=========================================="
