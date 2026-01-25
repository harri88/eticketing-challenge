# Ledger Service Integration

## Overview

The Go payment service now integrates with the Python ledger service to record double-entry accounting entries after successful payment processing.

## Integration Architecture

```
Payment Flow:
1. Resolve Payment Strategy
2. Create PENDING Transaction (DB)
3. Process External Payment (Gateway)
4. Update Transaction to SUCCESS (DB)
5. Confirm with Ticket Service (HTTP)
6. Record in Ledger Service (HTTP) ← NEW
```

## Configuration

### Ledger Service URL
- **Default**: `http://localhost:8000`
- **Field**: `ledgerURL` in `paymentUsecase` struct
- **Initialization**: Set in `NewPaymentUsecase()` function

```go
ledgerURL: "http://localhost:8000", // Default ledger service URL
```

### HTTP Client
- **Timeout**: 10 seconds per request
- **Request Method**: POST
- **Content-Type**: application/json
- **Accept**: application/json

## Request/Response Format

### Request Payload
```json
{
  "transaction_id": "TXN-CC-0123456789",
  "amount": 200.00
}
```

### Expected Response
- **Status Code**: 201 (Created) or 200 (OK)
- **Response Body**: 
```json
{
  "data": {
    "message": "Ledger entries recorded successfully",
    "transaction_id": "TXN-CC-0123456789"
  }
}
```

## Error Handling

**Non-blocking**: If the ledger service fails, the payment is still marked as successful. Warnings are logged to stderr.

```go
if err := u.recordPaymentInLedger(ctx, tx.TransactionID, req.Amount); err != nil {
  fmt.Printf("WARNING: Failed to record payment in ledger for transaction %s: %v\n", tx.TransactionID, err)
}
```

### Error Scenarios
- Ledger service unavailable → WARNING logged, payment succeeds
- Network timeout → Treated as service unavailable
- Malformed response → Error logged but payment still succeeds
- Invalid status code → Error logged but payment still succeeds

## Ledger Service Endpoint

### POST /api/v1/ledger

Creates double-entry accounting entries:
- **Debit Entry**: Cash_Asset (incoming funds)
- **Credit Entry**: Ticket_Revenue (revenue recognition)

### Database Table
```sql
CREATE TABLE ledger_entries (
  id SERIAL PRIMARY KEY,
  transaction_id VARCHAR(50) NOT NULL,
  account_name VARCHAR(100) NOT NULL,
  entry_type VARCHAR(10) NOT NULL,  -- 'DEBIT' or 'CREDIT'
  amount NUMERIC(15,2) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Code Structure

### File: `internal/usecase/payment_usecase.go`

**Struct Fields** (lines 18-26):
```go
type paymentUsecase struct {
  repo           domain.TransactionRepository
  gatewayFactory *gateway.PaymentGatewayFactory
  ticketClient   domain.TicketClient
  ledgerURL      string        // NEW
  httpClient     *http.Client  // NEW
  timeout        time.Duration
}
```

**ProcessPayment Method** (lines 38-105):
- Step 1-5: Existing payment flow
- **Step 6 (NEW)**: Calls `recordPaymentInLedger()`

**recordPaymentInLedger Method** (lines 107-155):
- Marshals JSON payload
- Creates HTTP POST request with context
- Sets required headers
- Sends to ledger service
- Validates response status (201 or 200)
- Logs success with amount details
- Returns error or nil

## Testing Integration

### Prerequisites
1. Python ledger service running on `localhost:8000`
```bash
cd backend-python-ledger
python main.py  # or uvicorn app.main:app --host 0.0.0.0 --port 8000
```

2. Go payment service running
```bash
cd backend-go-payment
go run cmd/api/main.go
```

3. PostgreSQL running on both ports (5432/5433/5434)

### Test Payment Flow
1. Submit payment request to payment service
2. Check payment returns SUCCESS status
3. Check ledger service logs for POST request received
4. Query ledger database for new entries
```bash
curl -X POST http://localhost:8000/api/v1/ledger \
  -H "Content-Type: application/json" \
  -d '{"transaction_id":"TXN-CC-test123","amount":100.00}'
```

4. Verify two entries created (Debit + Credit)

## Future Enhancements

### Configuration from Environment
Move hardcoded URL to environment variable:
```bash
export LEDGER_SERVICE_URL="http://localhost:8000"
```

### Retry Logic
Implement exponential backoff or circuit breaker for ledger service failures.

### Idempotency
Ledger service should support idempotent requests via transaction_id to handle retries safely.

### Structured Logging
Use structured logging instead of fmt.Printf for better observability.

### Dead Letter Queue
Queue failed ledger operations for asynchronous retry.

## Troubleshooting

### Issue: "WARNING: Failed to record payment in ledger"
- **Check**: Is ledger service running on port 8000?
- **Check**: Is PostgreSQL accessible from payment service?
- **Check**: Are ledger service logs showing the POST request?
- **Check**: Is the JSON payload valid (transaction_id and amount)?

### Issue: Payments succeed but ledger entries missing
- Payment service may have network issues reaching ledger
- Check network connectivity between services
- Check ledger service logs for 400/500 errors
- Verify ledger database for orphaned transactions

### Issue: Timeouts on payment requests
- Ledger service may be slow (processing double entries)
- Increase HTTP client timeout from 10s if needed
- Check ledger database query performance

## Compliance & Audit

Double-entry accounting ensures:
- ✅ All transactions are balanced (Assets = Revenue)
- ✅ Audit trail for all payments
- ✅ Financial reconciliation possible
- ✅ Compliance with accounting standards (IFRS)
