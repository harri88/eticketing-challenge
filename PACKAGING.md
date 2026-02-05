# E-Ticketing Challenge - Package Distribution Guide

This guide explains how to create distributable packages for all services in the E-Ticketing Challenge system.

## Overview

The packaging system creates distributable archives (`.tar.gz` files) for each service that can be deployed on any compatible environment. Each package includes:

- Compiled binaries or built artifacts
- Configuration templates
- Documentation
- Version information
- Deployment instructions

## Quick Start

### Package All Services

To create packages for all services with a single command:

```bash
cd scripts
VERSION=1.0.0 ./package-all.sh
```

This will create all packages in the `dist/` directory at the root of the project.

### Package Individual Services

You can also package services individually:

```bash
# Backend .NET Tickets Service
cd scripts
VERSION=1.0.0 ./package-dotnet.sh

# Backend Go Payment Service
cd scripts
VERSION=1.0.0 ./package-go.sh

# Backend Python Ledger Service
cd scripts
VERSION=1.0.0 ./package-python.sh

# Frontend React (Customer)
cd scripts
VERSION=1.0.0 ./package-react.sh

# Frontend Backoffice (Admin)
cd scripts
VERSION=1.0.0 ./package-backoffice.sh
```

## Prerequisites

### For Backend .NET Tickets

- .NET 8 SDK or later (for building)
- Approximately 200MB disk space

### For Backend Go Payment

- Go 1.21 or later
- Approximately 50MB disk space

### For Backend Python Ledger

- Python 3.11 or later
- pip package manager
- Approximately 100MB disk space

### For Frontend Applications

- Node.js 16 or later
- npm package manager
- Approximately 500MB disk space per frontend

## Package Structure

Each package follows this general structure:

```
service-name-version/
├── VERSION                 # Version identifier
├── README.md              # Service documentation
├── Dockerfile             # Container definition
├── <service-specific>     # Service code/binaries
└── DEPLOY.md             # Deployment instructions (for frontends)
```

### Backend .NET Tickets Package

```
backend-dotnet-tickets-1.0.0/
├── publish/               # .NET published artifacts
│   ├── backend-dotnet-tickets.dll
│   ├── appsettings.json
│   └── ...
├── sql_scripts/          # Database initialization
├── VERSION
├── README.md
└── Dockerfile
```

### Backend Go Payment Package

```
backend-go-payment-1.0.0/
├── bin/
│   ├── api                        # Linux binary
│   ├── api-darwin-amd64          # macOS binary
│   └── api-windows-amd64.exe     # Windows binary
├── docs/                 # API documentation
├── .env.example         # Configuration template
├── VERSION
└── README.md
```

### Backend Python Ledger Package

```
backend-python-ledger-1.0.0/
├── app/                 # Application code
├── requirements.txt     # Python dependencies
├── setup.sh            # Setup script
├── .env.example
├── VERSION
└── README.md
```

### Frontend Packages

```
frontend-react-1.0.0/
├── static/             # Static assets
├── index.html         # Entry point
├── VERSION
├── DEPLOY.md         # Deployment guide
└── package.json
```

## Deployment

### Deploying Backend .NET Tickets

```bash
# Extract package
tar -xzf backend-dotnet-tickets-1.0.0.tar.gz
cd backend-dotnet-tickets

# Run using .NET
dotnet publish/backend-dotnet-tickets.dll
```

### Deploying Backend Go Payment

```bash
# Extract package
tar -xzf backend-go-payment-1.0.0.tar.gz
cd backend-go-payment

# Run the binary (Linux)
./bin/api
```

### Deploying Backend Python Ledger

```bash
# Extract package
tar -xzf backend-python-ledger-1.0.0.tar.gz
cd backend-python-ledger

# Setup virtual environment and install dependencies
./setup.sh

# Run the application
source venv/bin/activate
python main.py
```

### Deploying Frontend Applications

```bash
# Extract package
tar -xzf frontend-react-1.0.0.tar.gz
cd frontend-react

# Serve using any static web server
# Example with Python:
python3 -m http.server 8000

# Or with Node.js serve:
npx serve -s . -p 3000

# See DEPLOY.md for production deployment guides
```

## Docker Deployment

Each package includes a Dockerfile for containerized deployment:

```bash
# Extract and build
tar -xzf backend-dotnet-tickets-1.0.0.tar.gz
cd backend-dotnet-tickets
docker build -t backend-dotnet-tickets:1.0.0 .
docker run -p 5020:5020 backend-dotnet-tickets:1.0.0
```

## Version Management

Set the VERSION environment variable to tag your packages:

```bash
# Development version
VERSION=1.0.0-dev ./package-all.sh

# Release version
VERSION=1.0.0 ./package-all.sh

# Patch version
VERSION=1.0.1 ./package-all.sh
```

## CI/CD Integration

These packaging scripts can be integrated into your CI/CD pipeline:

```yaml
# Example GitHub Actions workflow
- name: Package Services
  run: |
    cd scripts
    VERSION=${{ github.ref_name }} ./package-all.sh
    
- name: Upload Artifacts
  uses: actions/upload-artifact@v3
  with:
    name: packages
    path: dist/*.tar.gz
```

## Troubleshooting

### Build Failures

If a service fails to package:

1. Check the console output for specific errors
2. Ensure all prerequisites are installed
3. Try building the service manually first
4. Check the service's README for specific requirements

### Missing Dependencies

Backend services:
```bash
# .NET
dotnet restore

# Go
go mod download

# Python
pip install -r requirements.txt
```

Frontend services:
```bash
npm install
```

### Disk Space

Ensure you have sufficient disk space:
- Backend packages: ~100-200MB each
- Frontend packages: ~50-100MB each
- Total for all services: ~1GB

## Best Practices

1. **Always specify a version**: Use semantic versioning (e.g., 1.0.0)
2. **Test before packaging**: Ensure services build and run correctly
3. **Clean builds**: Packages use clean builds (no development artifacts)
4. **Document changes**: Update service READMEs when packaging new versions
5. **Verify archives**: Test extracted packages in a clean environment

## Additional Resources

- [Docker Compose Setup](../docker-compose.yml) - For local development
- [Deployment Guide](../DEPLOY_EC2.md) - For production deployment
- [CORS Setup](../CORS_SETUP.md) - For API configuration
- [SonarQube Setup](../SONARQUBE_SETUP.md) - For code quality

## Support

For issues or questions:
1. Check service-specific README files
2. Review the main project README
3. Check GitHub issues
4. Review CI/CD workflow logs

---

*Last updated: 2026-01-28*
