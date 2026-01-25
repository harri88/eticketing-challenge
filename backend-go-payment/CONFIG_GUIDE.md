# Payment Service Configuration Guide

This document explains how to configure the Payment Service using the `.env` file.

## Overview

The Payment Service uses a centralized configuration system that loads settings from a `.env` file and falls back to environment variables. All configuration is managed through the `config` package.

## Setup

### 1. Copy Environment File

```bash
cp .env.example .env
```

### 2. Update Configuration

Edit `.env` with your specific settings:

```bash
vim .env
```

## Configuration Sections

### Database Configuration

Controls PostgreSQL connection and connection pooling:

```env
DB_HOST=localhost              # PostgreSQL server hostname
DB_PORT=5432                   # PostgreSQL server port
DB_USER=user                   # PostgreSQL username
DB_PASSWORD=password           # PostgreSQL password
DB_NAME=payment_db             # Database name
DB_SSLMODE=disable             # SSL mode: disable, require, verify-ca, verify-full
DB_MAX_CONNS=25                # Maximum open connections
DB_MIN_CONNS=5                 # Minimum idle connections
```

**Connection String Format:**
```
postgres://user:password@localhost:5432/payment_db?sslmode=disable
```

### Server Configuration

Controls the API server behavior:

```env
SERVER_HOST=0.0.0.0            # Bind address (0.0.0.0 = all interfaces)
SERVER_PORT=8081               # HTTP server port
SERVER_READ_TIMEOUT=10         # Request read timeout (seconds)
SERVER_WRITE_TIMEOUT=10        # Response write timeout (seconds)
```

### Environment

Controls the operating environment:

```env
ENVIRONMENT=development        # development or production
```

**Effects:**
- `development`: Enables debug output, Swagger UI available at `/swagger/index.html`
- `production`: Minimal logging, stricter error handling

### Ticket Service Configuration

Connects to the E-Ticketing backend:

```env
TICKET_SERVICE_URL=http://localhost:7175    # Ticket service API URL
TICKET_SERVICE_TIMEOUT=30                    # Request timeout (seconds)
```

### Payment Gateway Configuration

Configures payment processing:

```env
PAYMENT_GATEWAY_PROVIDER=stripe  # stripe, paypal, or local

# Stripe
STRIPE_API_KEY=sk_test_...       # Get from https://dashboard.stripe.com/apikeys

# PayPal
PAYPAL_CLIENT_ID=...             # Get from https://developer.paypal.com/dashboard/
PAYPAL_SECRET=...
```

### Logging & Debug

```env
LOG_LEVEL=info                   # debug, info, warn, error
DEBUG=false                      # Enable verbose output
```

## Usage Examples

### Local Development

```env
ENVIRONMENT=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
SERVER_PORT=8081
PAYMENT_GATEWAY_PROVIDER=stripe
STRIPE_API_KEY=sk_test_4eC39HqLyjWDarh...
```

### Docker Compose

```env
ENVIRONMENT=development
DB_HOST=postgres                 # Docker service name
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
SERVER_HOST=0.0.0.0
SERVER_PORT=8081
TICKET_SERVICE_URL=http://backend:7175
PAYMENT_GATEWAY_PROVIDER=stripe
```

### Production

```env
ENVIRONMENT=production
DB_HOST=prod.database.internal
DB_PORT=5432
DB_USER=prod_user
DB_PASSWORD=secure_password_here
DB_SSLMODE=require
DB_MAX_CONNS=50
SERVER_HOST=0.0.0.0
SERVER_PORT=8081
LOG_LEVEL=warn
PAYMENT_GATEWAY_PROVIDER=stripe
STRIPE_API_KEY=sk_live_...
```

## Configuration Loading Priority

1. **`.env` file** - Highest priority (if exists in working directory)
2. **Environment variables** - Used if `.env` not found
3. **Default values** - Used if variable not set

## Common Configuration Scenarios

### Increase Database Connection Pool

For high-traffic environments:

```env
DB_MAX_CONNS=100
DB_MIN_CONNS=20
```

### Enable Production Mode with Strict SSL

```env
ENVIRONMENT=production
DB_SSLMODE=require
```

### Use PayPal Instead of Stripe

```env
PAYMENT_GATEWAY_PROVIDER=paypal
PAYPAL_CLIENT_ID=your_client_id
PAYPAL_SECRET=your_secret
```

### Increase Request Timeout for Slow Networks

```env
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30
TICKET_SERVICE_TIMEOUT=60
```

## Validation

The configuration system performs the following checks:

- ✓ Parses integer values correctly
- ✓ Falls back to defaults for invalid values
- ✓ Logs warnings for malformed configuration
- ✓ Tests database connectivity on startup
- ✓ Validates server binding address

## Security Best Practices

1. **Never commit `.env` to version control**
   ```bash
   echo ".env" >> .gitignore
   ```

2. **Use `.env.example` as a template**
   - Remove sensitive values from `.env.example`
   - Document all required variables

3. **Rotate credentials regularly**
   - Database passwords
   - API keys (Stripe, PayPal)

4. **Use environment-specific values**
   - Different databases for dev/prod
   - Different payment gateway keys
   - Different timeouts based on network stability

5. **Set restrictive database user permissions**
   ```sql
   -- Production user with limited permissions
   CREATE USER payment_service WITH ENCRYPTED PASSWORD 'secure_password';
   GRANT CONNECT ON DATABASE payment_db TO payment_service;
   GRANT USAGE ON SCHEMA public TO payment_service;
   GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA public TO payment_service;
   ```

## Troubleshooting

### "Database connection failed"

Check:
- DB_HOST and DB_PORT are correct
- Database is running
- Credentials (DB_USER, DB_PASSWORD) are valid
- Network connectivity (firewall rules)

### "Timeout connecting to ticket service"

Check:
- TICKET_SERVICE_URL is correct
- Ticket service is running
- TICKET_SERVICE_TIMEOUT is sufficient
- Network connectivity between services

### "Payment gateway error"

Check:
- PAYMENT_GATEWAY_PROVIDER matches your setup
- API keys (STRIPE_API_KEY, PAYPAL_*) are valid
- Keys match the environment (test vs. live)

### "Server won't start"

Check:
- SERVER_PORT is not in use
- SERVER_HOST is valid
- PORT is privileged (< 1024 requires root)

## Development Workflow

1. **Initialize configuration:**
   ```bash
   cp .env.example .env
   ```

2. **Start database:**
   ```bash
   docker-compose up -d postgres
   ```

3. **Run migrations:**
   ```bash
   # Add your migration command here
   ```

4. **Start the service:**
   ```bash
   cd cmd/api && go run main.go
   ```

5. **Verify:**
   ```bash
   curl http://localhost:8081/api/v1/payments
   # Swagger UI: http://localhost:8081/swagger/index.html
   ```

## Next Steps

- Review `config.go` for available configuration options
- Check `main.go` for how configuration is loaded
- Set up `.env` file for your environment
- Run the service and verify logs
