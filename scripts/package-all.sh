#!/bin/bash
# Master package script for E-Ticketing Challenge
# Packages all services into distributable archives

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VERSION="${VERSION:-1.0.0}"

echo "=========================================="
echo "E-Ticketing Challenge - Package All Services"
echo "Version: $VERSION"
echo "=========================================="

# Create dist directory
mkdir -p "$SCRIPT_DIR/../dist"

# Array to track packaging results
declare -a RESULTS

# Function to run packaging script
package_service() {
    local service_name=$1
    local script_name=$2
    
    echo ""
    echo ">>> Packaging $service_name..."
    echo ""
    
    if bash "$SCRIPT_DIR/$script_name"; then
        RESULTS+=("✓ $service_name: SUCCESS")
    else
        RESULTS+=("✗ $service_name: FAILED")
    fi
}

# Package backend services
package_service "backend-dotnet-tickets" "package-dotnet.sh"
package_service "backend-go-payment" "package-go.sh"
package_service "backend-python-ledger" "package-python.sh"

# Package frontend services
package_service "frontend-react" "package-react.sh"
package_service "frontend-backoffice" "package-backoffice.sh"

# Print summary
echo ""
echo "=========================================="
echo "PACKAGING SUMMARY"
echo "=========================================="
for result in "${RESULTS[@]}"; do
    echo "$result"
done
echo "=========================================="

# List all created packages
echo ""
echo "Created packages in dist/:"
ls -lh "$SCRIPT_DIR/../dist"/*.tar.gz 2>/dev/null || echo "No archives created"

echo ""
echo "Package creation complete!"
echo "Version: $VERSION"
echo "Location: $SCRIPT_DIR/../dist/"
