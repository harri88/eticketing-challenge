# Ledger Service Swagger/OpenAPI Documentation

## Quick Links

- **Swagger UI**: http://localhost:8000/docs
- **ReDoc**: http://localhost:8000/redoc
- **OpenAPI JSON**: http://localhost:8000/openapi.json

## Overview

The Ledger Service provides a RESTful API for managing financial ledger entries using double-entry accounting principles. Every payment transaction is recorded as a balanced pair of entries ensuring financial integrity and audit compliance.

## API Endpoints

### 1. Health Check

**Endpoint**: `GET /health`

**Purpose**: Monitor service health and database connectivity

**Response**:
```json
{
  "status": "healthy",
  "service": "ledger-service"
}
```

**Status Codes**:
- `200`: Service is operational

---

### 2. Record Payment Transaction

**Endpoint**: `POST /api/v1/ledger`

**Purpose**: Record a successful payment as balanced double-entry ledger entries

**Request Body**:
```json
{
  "transaction_id": "txn-20260125-001",
  "amount": 200.00
}
```

**Request Parameters**:
- `transaction_id` (string, required): Unique transaction identifier from payment service
  - Min length: 1, Max length: 50
  - Example: `txn-20260125-001`

- `amount` (number, required): Payment amount in AED
  - Must be > 0
  - Example: `200.00`

**Response** (201 Created):
```json
{
  "data": {
    "message": "Ledger entries recorded successfully",
    "transaction_id": "txn-20260125-001"
  }
}
```

**Double-Entry Entries Created**:
- **Debit Entry**: Cash_Asset account (money received)
- **Credit Entry**: Ticket_Revenue account (service provided)

**Status Codes**:
- `201`: Entries recorded successfully
- `400`: Invalid request data
- `500`: Server error during recording

**Example cURL**:
```bash
curl -X POST http://localhost:8000/api/v1/ledger \
  -H "Content-Type: application/json" \
  -d '{
    "transaction_id": "txn-20260125-001",
    "amount": 200.00
  }'
```

---

### 3. Get All Ledger Entries

**Endpoint**: `GET /api/v1/ledger`

**Purpose**: Retrieve complete audit trail of all ledger entries

**Query Parameters**: None

**Response** (200 OK):
```json
{
  "data": [
    {
      "id": 2,
      "transaction_id": "txn-20260125-001",
      "account_name": "Ticket_Revenue",
      "entry_type": "CREDIT",
      "amount": 200.00,
      "created_at": "2026-01-25T10:30:00"
    },
    {
      "id": 1,
      "transaction_id": "txn-20260125-001",
      "account_name": "Cash_Asset",
      "entry_type": "DEBIT",
      "amount": 200.00,
      "created_at": "2026-01-25T10:30:00"
    }
  ]
}
```

**Response Fields**:
- `id` (integer): Database record ID
- `transaction_id` (string): Associated transaction ID
- `account_name` (string): Account name (Cash_Asset | Ticket_Revenue)
- `entry_type` (string): Entry type (DEBIT | CREDIT)
- `amount` (number): Entry amount in AED
- `created_at` (string): ISO 8601 timestamp

**Status Codes**:
- `200`: Successfully retrieved entries
- `500`: Server error during retrieval

**Example cURL**:
```bash
curl http://localhost:8000/api/v1/ledger
```

---

### 4. Get Transaction Ledger Entries

**Endpoint**: `GET /api/v1/ledger/{transaction_id}`

**Purpose**: Retrieve debit and credit entries for a specific transaction

**Path Parameters**:
- `transaction_id` (string, required): Transaction identifier to lookup
  - Example: `txn-20260125-001`

**Response** (200 OK):
```json
{
  "data": [
    {
      "id": 2,
      "transaction_id": "txn-20260125-001",
      "account_name": "Ticket_Revenue",
      "entry_type": "CREDIT",
      "amount": 200.00,
      "created_at": "2026-01-25T10:30:00"
    },
    {
      "id": 1,
      "transaction_id": "txn-20260125-001",
      "account_name": "Cash_Asset",
      "entry_type": "DEBIT",
      "amount": 200.00,
      "created_at": "2026-01-25T10:30:00"
    }
  ]
}
```

**Status Codes**:
- `200`: Successfully retrieved entries
- `404`: Transaction not found
- `500`: Server error during retrieval

**Example cURL**:
```bash
curl http://localhost:8000/api/v1/ledger/txn-20260125-001
```

---

## Data Models

### PaymentEvent
Request model for recording a payment.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| transaction_id | string | Yes | Unique transaction identifier (1-50 chars) |
| amount | number | Yes | Payment amount in AED (> 0) |

### LedgerEntry
Single ledger entry in the accounting ledger.

| Field | Type | Description |
|-------|------|-------------|
| id | integer | Database record ID |
| transaction_id | string | Associated transaction ID |
| account_name | string | Cash_Asset or Ticket_Revenue |
| entry_type | string | DEBIT or CREDIT |
| amount | number | Entry amount in AED |
| created_at | string | ISO 8601 timestamp |

### ErrorResponse
Standard error response.

| Field | Type | Description |
|-------|------|-------------|
| detail | string | Error message |

---

## Double-Entry Accounting Explained

### What is Double-Entry Bookkeeping?

Every financial transaction is recorded twice:
1. **Debit**: Money in or expense out
2. **Credit**: Money out or revenue in

This ensures: **Total Debits = Total Credits** (always balanced)

### Example: Customer Purchases 200 AED Premium Ticket

```
Transaction ID: txn-20260125-001
Amount: 200 AED

Ledger Entry 1 (Debit):
  Account:   Cash_Asset
  Type:      DEBIT
  Amount:    200 AED
  Meaning:   Money received from customer

Ledger Entry 2 (Credit):
  Account:   Ticket_Revenue
  Type:      CREDIT
  Amount:    200 AED
  Meaning:   Service revenue earned
```

**Balance Check**: Debits (200) = Credits (200) ✓ Balanced

### Why Double-Entry?

1. **Prevents Fraud**: Impossible to hide money
2. **Error Detection**: Unbalanced entries indicate problems
3. **Audit Trail**: Complete history of all transactions
4. **Financial Reports**: Accurate balance sheets and income statements
5. **Compliance**: Standard accounting practice (required for most businesses)

---

## Account Types

### Cash_Asset
Represents cash and liquid assets received from payment methods.
- Increases when customer pays (DEBIT)
- Decreases when cash is used (CREDIT)

### Ticket_Revenue
Represents revenue earned from ticket sales.
- Increases when ticket is sold (CREDIT)
- Decreases when refund is issued (DEBIT)

---

## API Usage Examples

### JavaScript/Fetch

```javascript
// Record a payment
async function recordPayment(transactionId, amount) {
  const response = await fetch('http://localhost:8000/api/v1/ledger', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      transaction_id: transactionId,
      amount: amount
    })
  });
  return response.json();
}

// Get all ledger entries
async function getLedgerEntries() {
  const response = await fetch('http://localhost:8000/api/v1/ledger');
  return response.json();
}

// Get entries for specific transaction
async function getTransactionLedger(transactionId) {
  const response = await fetch(`http://localhost:8000/api/v1/ledger/${transactionId}`);
  return response.json();
}
```

### Python/Requests

```python
import requests

# Record payment
response = requests.post(
    'http://localhost:8000/api/v1/ledger',
    json={
        'transaction_id': 'txn-20260125-001',
        'amount': 200.00
    }
)
print(response.json())

# Get all entries
response = requests.get('http://localhost:8000/api/v1/ledger')
print(response.json())

# Get transaction entries
response = requests.get('http://localhost:8000/api/v1/ledger/txn-20260125-001')
print(response.json())
```

### Go/http

```go
package main

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
)

func recordPayment(txnID string, amount float64) {
    payload := map[string]interface{}{
        "transaction_id": txnID,
        "amount": amount,
    }
    
    body, _ := json.Marshal(payload)
    resp, _ := http.Post(
        "http://localhost:8000/api/v1/ledger",
        "application/json",
        bytes.NewBuffer(body),
    )
    
    data, _ := ioutil.ReadAll(resp.Body)
    println(string(data))
}
```

---

## HTTP Status Codes

| Code | Meaning | Usage |
|------|---------|-------|
| 200 | OK | Successful GET request |
| 201 | Created | Successful POST request |
| 400 | Bad Request | Invalid request data |
| 404 | Not Found | Transaction not found |
| 500 | Internal Error | Server error |

---

## CORS Configuration

The service is configured to accept requests from any origin:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: *`

This allows the Ledger Service to be called from any frontend service (React backoffice, etc.)

---

## Error Handling

### Invalid Request
```json
{
  "detail": "Failed to record ledger entries"
}
```

### Transaction Not Found
```json
{
  "detail": "No ledger entries found for transaction txn-unknown-001"
}
```

### Server Error
```json
{
  "detail": "Failed to retrieve ledger entries"
}
```

---

## Performance Considerations

- Entries are ordered by creation date (newest first)
- Database queries are indexed on transaction_id for fast lookups
- Connection pooling is configured for concurrent requests
- Atomic transactions ensure consistency

---

## Database Schema

```sql
CREATE TABLE ledger_entries (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(50) NOT NULL,
    account_name VARCHAR(100) NOT NULL,
    entry_type VARCHAR(10) NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transaction_id ON ledger_entries(transaction_id);
CREATE INDEX idx_created_at ON ledger_entries(created_at);
```

---

## Deployment

### Docker
```bash
docker-compose up -d backend-python-ledger
```

### Manual (Local)
```bash
cd backend-python-ledger
PYTHONPATH=. python3 -m uvicorn app.main:app --host 0.0.0.0 --port 8000
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5434 | PostgreSQL port |
| DB_USER | postgres | Database user |
| DB_PASSWORD | postgres@12345 | Database password |
| DB_NAME | ledger_db | Database name |

---

## Support

For issues or questions:
- Check logs: `docker logs backend-python-ledger`
- API Docs: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc
