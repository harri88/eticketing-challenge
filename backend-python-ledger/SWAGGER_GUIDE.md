# Backend Python Ledger Service - Swagger/OpenAPI Guide

## Service Information

- **Service Name**: E-Ticketing Ledger Service
- **Version**: 1.0.0
- **API Base URL**: `http://localhost:8000`
- **Swagger UI**: `http://localhost:8000/docs`
- **ReDoc**: `http://localhost:8000/redoc`
- **OpenAPI Schema**: `http://localhost:8000/openapi.json`

## What is Swagger/OpenAPI?

**OpenAPI** (formerly Swagger) is an open standard for describing REST APIs. FastAPI automatically generates interactive API documentation that allows you to:

✅ View all available endpoints
✅ See detailed parameter and response documentation
✅ Try API calls directly from the browser
✅ View request/response examples
✅ Check error codes and status codes
✅ Test the API without external tools

## Accessing Swagger Documentation

### 1. Interactive Swagger UI
Open in your browser: **http://localhost:8000/docs**

This provides:
- Interactive endpoint explorer
- Try It Out button for each endpoint
- Parameter validation
- Real-time request/response examples
- Schema definitions

### 2. ReDoc (Alternative Documentation)
Open in your browser: **http://localhost:8000/redoc**

ReDoc provides:
- Clean, readable documentation
- Better for documentation sharing
- Search functionality
- Organized by tags

### 3. Raw OpenAPI JSON
Access at: **http://localhost:8000/openapi.json**

Contains the complete OpenAPI specification that can be:
- Used by code generation tools
- Imported into other API tools (Postman, etc.)
- Consumed by API gateways

## API Endpoints

### Endpoint Tags

Endpoints are organized into two main tags:

#### Health
- `GET /health` - Service health check

#### Ledger Entries
- `POST /api/v1/ledger` - Record payment transaction
- `GET /api/v1/ledger` - Retrieve all ledger entries
- `GET /api/v1/ledger/{transaction_id}` - Retrieve transaction entries

## Using the Swagger UI

### 1. Explore Endpoints

The Swagger UI lists all endpoints on the left side:
- Click on any endpoint to expand it
- View the endpoint path, method, and description
- See parameters, request body format, and response models

### 2. Try Out an Endpoint

**Example: Record a Payment**

1. Click on `POST /api/v1/ledger`
2. Click "Try it out" button
3. Fill in the parameters:
   ```json
   {
     "transaction_id": "txn-20260125-001",
     "amount": 200.00
   }
   ```
4. Click "Execute"
5. View the response and response code

### 3. View Response Models

Scroll down in any endpoint section to see:
- **Schemas**: Data model definitions
- **Examples**: Sample request/response data
- **Response Codes**: Possible HTTP status codes

## API Documentation Highlights

### Request Documentation

For `POST /api/v1/ledger`:

```
Request Body Schema: PaymentEvent
- transaction_id (string, required)
  Description: Unique transaction identifier from payment service
  Min length: 1
  Max length: 50
  Example: txn-20260125-001

- amount (number, required)
  Description: Payment amount in AED
  Must be: > 0
  Example: 200.00
```

### Response Documentation

```
Response Code: 201 (Created)
Response Body: RecordPaymentDataResponse
{
  "data": {
    "message": "Ledger entries recorded successfully",
    "transaction_id": "txn-20260125-001"
  }
}
```

### Error Documentation

Error responses include:
- **400**: Invalid request data
- **404**: Transaction not found
- **500**: Server error

## How FastAPI Generated This Documentation

The Swagger documentation is automatically generated from:

1. **Function Docstrings**
```python
@app.post("/api/v1/ledger")
def record_payment(event: PaymentEvent, db: Session = Depends(get_db)):
    """
    Record a successful payment as balanced double-entry ledger entries.
    """
```

2. **Type Hints**
```python
def record_payment(
    event: PaymentEvent,  # ← Type hint defines request model
    db: Session = Depends(get_db)
) -> RecordPaymentDataResponse:  # ← Type hint defines response model
```

3. **Pydantic Models**
```python
class PaymentEvent(BaseModel):
    """Request model for recording a new payment transaction."""
    transaction_id: str = Field(
        ..., 
        description="Unique transaction identifier",
        example="txn-20260125-001"
    )
    amount: float = Field(
        ..., 
        description="Payment amount in AED",
        example=200.00,
        gt=0
    )
```

4. **Response Model Decorators**
```python
@app.post(
    "/api/v1/ledger",
    response_model=RecordPaymentDataResponse,  # ← Defines response schema
    status_code=201,
    tags=["Ledger Entries"],  # ← Organizes endpoints
    summary="Record Payment Transaction",  # ← Short description
    description="...",  # ← Detailed description
    responses={  # ← Document error responses
        201: {...},
        400: {...},
        500: {...}
    }
)
```

## Key Features in Our Swagger Documentation

### 1. Comprehensive Descriptions
Each endpoint includes:
- Summary (short description)
- Detailed description
- Use cases and examples
- Business logic explanation

### 2. Field-Level Documentation
Each parameter includes:
- Description
- Data type and constraints
- Valid ranges and formats
- Example values

### 3. Request/Response Examples
Every endpoint shows:
- Complete request example
- Complete response example
- Alternative response codes

### 4. Double-Entry Accounting Explanation
The service description includes:
- What is double-entry bookkeeping
- Why we use it
- Example of debit/credit entries
- Balance verification logic

### 5. Error Documentation
Complete error response models:
- Error code
- Error message
- When each error occurs

## Testing API Endpoints

### Using Swagger UI "Try It Out"

1. **Record a Payment**
   - Go to `POST /api/v1/ledger`
   - Click "Try it out"
   - Enter transaction ID and amount
   - Click "Execute"
   - See the response

2. **Get All Entries**
   - Go to `GET /api/v1/ledger`
   - Click "Try it out"
   - Click "Execute"
   - View the complete ledger

3. **Get Transaction Entries**
   - Go to `GET /api/v1/ledger/{transaction_id}`
   - Click "Try it out"
   - Enter a transaction ID
   - Click "Execute"
   - View debit and credit entries

## API Documentation Standards

Our implementation follows best practices:

✅ **Clear Descriptions**: Every endpoint explains what it does
✅ **Type Safety**: All requests and responses are validated
✅ **Error Handling**: All possible errors are documented
✅ **Examples**: Real-world examples for every endpoint
✅ **Field Documentation**: Each parameter is explained
✅ **Status Codes**: All possible HTTP responses listed
✅ **Business Logic**: Double-entry accounting principles explained
✅ **CORS**: Clear cross-origin access information

## Integration with Other Services

### Payment Service → Ledger Service

When a payment is processed:
```
POST http://localhost:8000/api/v1/ledger
{
  "transaction_id": "txn-001",
  "amount": 200.00
}
```

The Swagger docs show:
- Required parameters
- Valid value ranges
- Response format
- Possible error codes

### Backoffice → Ledger Service

When admin views ledger:
```
GET http://localhost:8000/api/v1/ledger
```

The Swagger docs show:
- Response structure
- Data fields
- Sort order (newest first)
- Complete examples

## Generating Code from OpenAPI

The OpenAPI schema can be used to generate client libraries:

### Generate JavaScript Client
```bash
# Using OpenAPI Generator
openapi-generator-cli generate \
  -i http://localhost:8000/openapi.json \
  -g javascript \
  -o ./generated-client
```

### Generate Python Client
```bash
openapi-generator-cli generate \
  -i http://localhost:8000/openapi.json \
  -g python \
  -o ./generated-client
```

### Import into Postman
1. Open Postman
2. Click "Import"
3. Select "Link"
4. Paste: `http://localhost:8000/openapi.json`
5. All endpoints auto-imported!

## Customization

The OpenAPI documentation is customized in `app/main.py`:

```python
app = FastAPI(
    title="E-Ticketing Ledger Service",
    description="...",
    version="1.0.0",
    contact={...},
    license_info={...}
)

def custom_openapi():
    # Add server information
    openapi_schema["servers"] = [
        {
            "url": "http://localhost:8000",
            "description": "Local Development Server"
        },
        {
            "url": "https://api.eticketing.local",
            "description": "Production Server"
        }
    ]
```

## Security Notes

Current implementation:
- ✅ All endpoints documented
- ✅ CORS enabled for development
- ✅ Input validation on all fields
- ⚠️ No authentication (add API-Key in production)

For production, add:
```python
from fastapi.security import APIKeyHeader

security_schemes = {
    "API-Key": {
        "type": "apiKey",
        "in": "header",
        "name": "X-API-Key"
    }
}
```

## Troubleshooting

### Swagger UI Not Loading
- Check service is running: `http://localhost:8000/health`
- Verify port 8000 is accessible
- Check browser console for errors

### Missing Endpoint Documentation
- Ensure endpoint has docstring
- Check response_model is defined
- Verify Pydantic models are properly typed

### Request Validation Failing
- Check field types match request model
- Verify required fields are provided
- Review field constraints (min_length, gt, etc.)

## Production Deployment

For production:
1. Add authentication schemes to Swagger
2. Update server URLs to production endpoints
3. Add rate limiting documentation
4. Document API versioning strategy
5. Include deprecation notices
6. Add API key in request examples

## References

- **OpenAPI 3.0 Spec**: https://spec.openapis.org/oas/v3.0.3
- **FastAPI Docs**: https://fastapi.tiangolo.com/
- **Pydantic Docs**: https://docs.pydantic.dev/
- **OpenAPI Generator**: https://openapi-generator.tech/

---

**Status**: ✅ Swagger documentation automatically generated and updated with each code change.

Access the interactive documentation at: **http://localhost:8000/docs**
