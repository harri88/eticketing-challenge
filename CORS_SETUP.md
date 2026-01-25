# CORS Configuration Guide

## Overview

CORS (Cross-Origin Resource Sharing) is configured on both the .NET and Go backends to allow requests from the React frontend running on different origins.

## Current Configuration

### Go Payment Service (backend-go-payment)
**File**: `cmd/api/main.go`

```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080", "http://127.0.0.1:3000", "http://127.0.0.1:8080"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With", "Idempotency-Key"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           3600,
}))
```

**Settings:**
- `AllowOrigins`: Permitted frontend origins
  - `http://localhost:3000` - React dev server default
  - `http://localhost:8080` - Alternative port
  - `http://127.0.0.1:3000/8080` - Localhost IP variants
- `AllowMethods`: HTTP methods allowed (GET, POST, PUT, DELETE, OPTIONS)
- `AllowHeaders`: Request headers clients can send
  - `Content-Type` - JSON payloads
  - `Authorization` - Auth tokens
  - `X-Requested-With` - XMLHttpRequest identification
  - `Idempotency-Key` - Payment idempotency
- `AllowCredentials`: Allow cookies/credentials
- `MaxAge`: How long browser caches preflight (3600s = 1 hour)

### .NET Ticket Service (backend-dotnet-tickets)
**File**: `src/Api/Program.cs`

```csharp
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowReactApp",
        policy =>
        {
            policy.WithOrigins("http://localhost:3000")
                  .AllowAnyHeader()
                  .AllowAnyMethod();
        });
});

app.UseCors("AllowReactApp");
```

**Settings:**
- `AllowReactApp` policy for `http://localhost:3000`
- `AllowAnyHeader()` - Accept any request headers
- `AllowAnyMethod()` - Accept any HTTP method

## API Endpoints with CORS

### Payment Service (Port 8081)
```
POST /api/v1/payments
- Origin: http://localhost:3000
- Method: POST
- Content-Type: application/json
```

### Ticket Service (Port 7175)
```
GET /api/v1/tickets
POST /api/v1/checkout/orders
- Origin: http://localhost:3000
- Methods: GET, POST
- Content-Type: application/json
```

## CORS Flow Diagram

```
React Frontend (localhost:3000)
    ↓
[Preflight OPTIONS Request]
    ↓
Go Payment API (localhost:8081)
    ↓
[Check CORS Headers]
    ↓
[Return Access-Control-* Headers]
    ↓
[Browser allows actual request]
    ↓
POST /api/v1/payments
    ↓
✓ Payment processed
```

## Common CORS Headers

### Request Headers (Frontend sends)
```
Origin: http://localhost:3000
Content-Type: application/json
X-Requested-With: XMLHttpRequest
Idempotency-Key: 12345-67890
```

### Response Headers (Backend responds)
```
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With, Idempotency-Key
Access-Control-Allow-Credentials: true
Access-Control-Max-Age: 3600
```

## Troubleshooting CORS Errors

### Error: "No 'Access-Control-Allow-Origin' header"
**Cause**: Backend not configured for your frontend origin
**Solution**: Add your origin to `AllowOrigins` list

### Error: "Method not allowed"
**Cause**: CORS policy doesn't allow the HTTP method you're using
**Solution**: Add method to `AllowMethods` list

### Error: "Header not allowed"
**Cause**: Custom header not in allowed list
**Solution**: Add header name to `AllowHeaders` list

### Error: "Preflight request failed"
**Cause**: Browser is sending OPTIONS request, backend not responding correctly
**Solution**: Ensure CORS middleware is registered early in pipeline

## Production Considerations

### For Production Deployment

**Security**: Don't allow all origins with wildcard (`*`)

**Instead, specify exact origins:**
```go
AllowOrigins: []string{
    "https://www.yourdomain.com",
    "https://yourdomain.com",
    "https://api.yourdomain.com",
}
```

**Separate Configuration by Environment:**

Development:
```
AllowOrigins: []string{"http://localhost:3000", "http://localhost:8080"}
```

Production:
```
AllowOrigins: []string{"https://www.yourdomain.com"}
```

**Set Strict Credentials Policy:**
```go
AllowCredentials: true  // Only if needed
MaxAge: 3600           // Balance between performance and security
```

## Environment-based CORS

### Docker Compose
When services run in Docker, use service names as origins:

```
frontend:3000 → http://frontend:3000
payment-api:8081 → http://backend-go-payment:8081
ticket-api:7175 → http://backend-dotnet-tickets:8080
```

Update CORS config:
```go
AllowOrigins: []string{
    "http://localhost:3000",           // Local development
    "http://frontend:3000",             // Docker Compose
    "http://127.0.0.1:3000",           // Localhost variant
}
```

## Testing CORS

### Using curl
```bash
# Preflight request
curl -i -X OPTIONS http://localhost:8081/api/v1/payments \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type"

# Expected response headers:
# Access-Control-Allow-Origin: http://localhost:3000
# Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
```

### Using Browser DevTools
1. Open DevTools (F12)
2. Go to Network tab
3. Make a request to the API
4. Check Response Headers for `Access-Control-Allow-*` headers

## References

- [MDN CORS Documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [Echo Framework CORS](https://echo.labstack.com/docs/middleware/cors)
- [ASP.NET Core CORS](https://docs.microsoft.com/en-us/aspnet/core/security/cors)
