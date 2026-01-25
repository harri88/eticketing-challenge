package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/harri88/eticketing-challenge/internal/domain"
	"github.com/harri88/eticketing-challenge/internal/infrastructure/gateway"
)

type paymentUsecase struct {
	repo           domain.TransactionRepository
	gatewayFactory *gateway.PaymentGatewayFactory
	ticketClient   domain.TicketClient
	timeout        time.Duration
}

func NewPaymentUsecase(repo domain.TransactionRepository, gf *gateway.PaymentGatewayFactory, tc domain.TicketClient) domain.PaymentUsecase {
	return &paymentUsecase{
		repo:           repo,
		gatewayFactory: gf,
		ticketClient:   tc,
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
	txID := fmt.Sprintf("TXN-%s-%s", req.PaymentMethod, uuid.New().String())

	tx := &domain.Transaction{
		TransactionID: txID,
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

	return tx, nil
}
