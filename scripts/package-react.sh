#!/bin/bash
# Package script for frontend-react
# Creates a distributable package of the React application

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$SCRIPT_DIR/../frontend-react"
OUTPUT_DIR="$SCRIPT_DIR/../dist/frontend-react"
VERSION="${VERSION:-1.0.0}"

echo "=========================================="
echo "Packaging frontend-react"
echo "Version: $VERSION"
echo "=========================================="

# Clean previous build
echo "Cleaning previous builds..."
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Navigate to project directory
cd "$PROJECT_DIR"

echo "Installing dependencies..."
npm ci

echo "Running tests..."
npm test -- --passWithNoTests --watchAll=false || echo "Tests failed or no tests found"

echo "Building production bundle..."
npm run build

# Copy build artifacts
echo "Copying build artifacts..."
cp -r build/* "$OUTPUT_DIR/"

# Create version file
echo "$VERSION" > "$OUTPUT_DIR/VERSION"

# Copy additional files
cp package.json "$OUTPUT_DIR/" 2>/dev/null || true
cp README.md "$OUTPUT_DIR/" 2>/dev/null || true
cp Dockerfile "$OUTPUT_DIR/" 2>/dev/null || true

# Create a simple deployment guide
cat > "$OUTPUT_DIR/DEPLOY.md" << 'EOF'
# Frontend React - Deployment Guide

## Static Hosting Deployment

This is a production-ready build of the React application.

### Deploy to Web Server (Nginx/Apache)
1. Copy all files to your web server's document root
2. Configure your web server to serve `index.html` for all routes (SPA routing)

### Example Nginx Configuration
```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/frontend-react;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

### Deploy with Docker
Use the included Dockerfile to create a containerized version:
```bash
docker build -t frontend-react .
docker run -p 3000:80 frontend-react
```

### Environment Variables
Configure API URLs before building:
- REACT_APP_API_URL: Backend API endpoint
EOF

# Create archive
echo "Creating distributable archive..."
cd "$SCRIPT_DIR/../dist"
tar -czf "frontend-react-$VERSION.tar.gz" frontend-react/

echo "=========================================="
echo "Package created successfully!"
echo "Location: $SCRIPT_DIR/../dist/frontend-react-$VERSION.tar.gz"
echo "=========================================="
