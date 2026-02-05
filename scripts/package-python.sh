#!/bin/bash
# Package script for backend-python-ledger
# Creates a distributable package of the Python application

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$SCRIPT_DIR/../backend-python-ledger"
OUTPUT_DIR="$SCRIPT_DIR/../dist/backend-python-ledger"
VERSION="${VERSION:-1.0.0}"

echo "=========================================="
echo "Packaging backend-python-ledger"
echo "Version: $VERSION"
echo "=========================================="

# Clean previous build
echo "Cleaning previous builds..."
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Navigate to project directory
cd "$PROJECT_DIR"

echo "Copying application files..."
for file in ./*.py; do
  [ -f "$file" ] && cp "$file" "$OUTPUT_DIR/" || true
done
cp -r app "$OUTPUT_DIR/" 2>/dev/null || true
cp -r config "$OUTPUT_DIR/" 2>/dev/null || true
cp -r models "$OUTPUT_DIR/" 2>/dev/null || true
cp -r routers "$OUTPUT_DIR/" 2>/dev/null || true
cp -r services "$OUTPUT_DIR/" 2>/dev/null || true
cp -r database "$OUTPUT_DIR/" 2>/dev/null || true

# Copy configuration and documentation
cp requirements.txt "$OUTPUT_DIR/" 2>/dev/null || true
cp .env.example "$OUTPUT_DIR/" 2>/dev/null || true
for file in ./*.md; do
  [ -f "$file" ] && cp "$file" "$OUTPUT_DIR/" || true
done
cp Dockerfile "$OUTPUT_DIR/" 2>/dev/null || true

# Create version file
echo "$VERSION" > "$OUTPUT_DIR/VERSION"

# Create a simple setup script
cat > "$OUTPUT_DIR/setup.sh" << 'EOF'
#!/bin/bash
# Setup script for backend-python-ledger

echo "Setting up Python virtual environment..."
python3 -m venv venv
source venv/bin/activate

echo "Installing dependencies..."
pip install --upgrade pip
pip install -r requirements.txt

echo "Setup complete!"
echo "To run the application:"
echo "  source venv/bin/activate"
echo "  python main.py"
EOF

chmod +x "$OUTPUT_DIR/setup.sh"

# Create archive
echo "Creating distributable archive..."
cd "$SCRIPT_DIR/../dist"
tar -czf "backend-python-ledger-$VERSION.tar.gz" backend-python-ledger/

echo "=========================================="
echo "Package created successfully!"
echo "Location: $SCRIPT_DIR/../dist/backend-python-ledger-$VERSION.tar.gz"
echo "=========================================="
