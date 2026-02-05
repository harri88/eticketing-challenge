#!/bin/bash
# Package script for backend-go-payment
# Creates a distributable package of the Go application

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$SCRIPT_DIR/../backend-go-payment"
OUTPUT_DIR="$SCRIPT_DIR/../dist/backend-go-payment"
VERSION="${VERSION:-1.0.0}"

echo "=========================================="
echo "Packaging backend-go-payment"
echo "Version: $VERSION"
echo "=========================================="

# Clean previous build
echo "Cleaning previous builds..."
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR/bin"

# Navigate to project directory
cd "$PROJECT_DIR"

echo "Downloading dependencies..."
go mod download

echo "Running tests..."
go test ./... || echo "Tests failed or no tests found"

echo "Building application..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.Version=$VERSION" -o "$OUTPUT_DIR/bin/api" ./cmd/api

# Build for other platforms (optional)
echo "Building for additional platforms..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s -X main.Version=$VERSION" -o "$OUTPUT_DIR/bin/api-darwin-amd64" ./cmd/api
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s -X main.Version=$VERSION" -o "$OUTPUT_DIR/bin/api-windows-amd64.exe" ./cmd/api

# Create version file
echo "$VERSION" > "$OUTPUT_DIR/VERSION"

# Copy configuration and documentation
cp .env.example "$OUTPUT_DIR/" 2>/dev/null || true
cp -r docs "$OUTPUT_DIR/" 2>/dev/null || true
for file in ./*.md; do
  [ -f "$file" ] && cp "$file" "$OUTPUT_DIR/" || true
done
cp Dockerfile "$OUTPUT_DIR/" 2>/dev/null || true

# Create archive
echo "Creating distributable archive..."
cd "$SCRIPT_DIR/../dist"
tar -czf "backend-go-payment-$VERSION.tar.gz" backend-go-payment/

echo "=========================================="
echo "Package created successfully!"
echo "Location: $SCRIPT_DIR/../dist/backend-go-payment-$VERSION.tar.gz"
echo "Binaries:"
echo "  - Linux: backend-go-payment/bin/api"
echo "  - macOS: backend-go-payment/bin/api-darwin-amd64"
echo "  - Windows: backend-go-payment/bin/api-windows-amd64.exe"
echo "=========================================="
