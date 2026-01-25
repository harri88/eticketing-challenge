package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/harri88/eticketing-challenge/internal/domain"
	"github.com/harri88/eticketing-challenge/internal/infrastructure/gateway"
)

type paymentUsecase struct {
	repo           domain.TransactionRepository
	gatewayFactory *gateway.PaymentGatewayFactory
	ticketClient   domain.TicketClient
	ledgerURL      string // Ledger service URL
	httpClient     *http.Client
	timeout        time.Duration
}

func NewPaymentUsecase(repo domain.TransactionRepository, gf *gateway.PaymentGatewayFactory, tc domain.TicketClient) domain.PaymentUsecase {
	return &paymentUsecase{
		repo:           repo,
		gatewayFactory: gf,
		ticketClient:   tc,
		ledgerURL:      "http://localhost:8000", // Default ledger service URL
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		timeout:        30 * time.Second, // Global timeout
	}
}

func (u *paymentUsecase) ProcessPayment(c context.Context, req domain.PaymentRequest, idempotencyKey string) (*domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	// 1. Resolve Strategy
	strategy, err := u.gatewayFactory.GetStrategy(req.PaymentMethod)
	if err != nil {
		return nil, err
	}

	// 2. Create Initial Transaction Record (PENDING)
	// We generate our own Transaction ID here
	txID := fmt.Sprintf("TXN-%s-%s", req.PaymentMethod, strings.Split(uuid.New().String(), "-")[0])

	tx := &domain.Transaction{
		TransactionID: strings.ToUpper(txID),
		OrderID:       req.OrderID,
		PaymentMethod: req.PaymentMethod,
		Amount:        req.Amount,
		Currency:      "AED",
		Status:        domain.StatusPending,
		CreatedAt:     time.Now(),
	}

	if err := u.repo.Create(ctx, tx); err != nil {
		return nil, err
	}

	// 3. Simulate External Payment (with Idempotency Key)
	// We use the idempotency key passed from the client, or generate one if missing for the external call
	if idempotencyKey == "" {
		idempotencyKey = uuid.New().String()
	}

	externalRef, err := strategy.ProcessPayment(ctx, req, idempotencyKey)
	if err != nil {
		// Update DB to Failed
		_ = u.repo.UpdateStatus(ctx, tx.ID, domain.StatusFailed, "")
		return nil, err
	}

	// 4. Update DB to Success
	if err := u.repo.UpdateStatus(ctx, tx.ID, domain.StatusSuccess, externalRef); err != nil {
		// Critical error: Payment succeeded but DB failed.
		// In production: Push to Dead Letter Queue (DLQ) or Retry mechanism
		return nil, err
	}
	tx.Status = domain.StatusSuccess
	tx.PaymentRef = externalRef

	// 5. Callback Ticket Service (Confirmation)
	// We do this asynchronously or synchronously depending on strictness.
	// Synchronous is safer for the "Challenge" context.
	if err := u.ticketClient.ConfirmPayment(ctx, req.OrderID); err != nil {
		// Warning: Ticket service is down, but payment is taken.
		// Log this error heavily. The reconciliation worker (Ledger) will fix this later.
		fmt.Printf("WARNING: Failed to notify ticket service for Order %s: %v\n", req.OrderID, err)
	}

	// 6. Record Payment in Ledger Service (Double-Entry Accounting)
	// This creates balanced ledger entries for financial audit trail
	if err := u.recordPaymentInLedger(ctx, tx.TransactionID, req.Amount); err != nil {
		// Warning: Ledger service is down, but payment is processed
		// Log this error - the payment is still successful in the payment service
		// The ledger can be reconciled later
		fmt.Printf("WARNING: Failed to record payment in ledger for transaction %s: %v\n", tx.TransactionID, err)
	}

	return tx, nil
}

// recordPaymentInLedger sends a payment record to the ledger service for double-entry accounting
// This creates balanced entries: Debit (Cash_Asset) and Credit (Ticket_Revenue)
func (u *paymentUsecase) recordPaymentInLedger(ctx context.Context, transactionID string, amount float64) error {
	// Prepare request payload for ledger service
	ledgerPayload := map[string]interface{}{
		"transaction_id": transactionID,
		"amount":         amount,
	}

	payloadBytes, err := json.Marshal(ledgerPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal ledger payload: %w", err)
	}

	// Create HTTP request
	ledgerURL := fmt.Sprintf("%s/api/v1/ledger", u.ledgerURL)
	req, err := http.NewRequestWithContext(ctx, "POST", ledgerURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create ledger request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Log cURL command for debugging
	curlCmd := fmt.Sprintf(
		`curl -X %s "%s" -H "Content-Type: application/json" -H "Accept: application/json" -d '%s'`,
		req.Method, req.URL.String(), string(payloadBytes),
	)
	fmt.Printf("DEBUG: Ledger Service cURL: %s\n", curlCmd)

	// Send request to ledger service
	resp, err := u.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call ledger service: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read ledger response: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ledger service returned error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// Log success
	fmt.Printf("INFO: Payment recorded in ledger service for transaction %s (amount: %.2f AED)\n", transactionID, amount)
	return nil
}
