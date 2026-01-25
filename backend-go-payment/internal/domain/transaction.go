package domain

import (
	"context"
	"time"
)

// Transaction Entity
// swagger:model Transaction
type Transaction struct {
	// Database ID
	// example: 1
	ID int64 `json:"id" db:"id"`

	// Unique transaction identifier
	// example: TXN-CC-17000
	TransactionID string `json:"transaction_id" db:"transaction_id"`

	// Associated order ID
	// example: order-12345
	OrderID string `json:"order_id" db:"order_id"`

	// Payment method used
	// enum: credit_card,qr_scan,crypto
	// example: credit_card
	PaymentMethod string `json:"payment_method" db:"payment_method"`

	// External payment gateway reference
	// example: ch_1234567890
	PaymentRef string `json:"payment_ref" db:"payment_ref"`

	// Transaction amount
	// example: 250.50
	Amount float64 `json:"amount" db:"amount"`

	// Currency code
	// example: AED
	Currency string `json:"currency" db:"currency"`

	// Transaction status
	// enum: PENDING,SUCCESS,FAILED
	// example: SUCCESS
	Status string `json:"status" db:"status"`

	// Timestamp when transaction was created
	// example: 2026-01-25T10:30:00Z
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

const (
	StatusPending = "PENDING"
	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"
)

// PaymentRequest DTO
// swagger:model PaymentRequest
type PaymentRequest struct {
	// Order ID - unique identifier for the order
	// required: true
	// example: order-12345
	OrderID string `json:"order_id" validate:"required"`

	// Payment method - credit_card, qr_scan, or crypto
	// required: true
	// enum: credit_card,qr_scan,crypto
	// example: credit_card
	PaymentMethod string `json:"payment_method" validate:"required"`

	// Credit card number (for credit_card method)
	// example: 4242424242424242
	CardNumber string `json:"card_number,omitempty"`

	// Credit card expiry date MM/YY (for credit_card method)
	// example: 12/25
	CardExpiry string `json:"card_expiry,omitempty"`

	// Credit card CVC (for credit_card method)
	// example: 123
	CardCVC string `json:"card_cvc,omitempty"`

	// Payment amount
	// required: true
	// minimum: 0.01
	// example: 250.50
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

// Repository Interface (Port)
type TransactionRepository interface {
	Create(ctx context.Context, tx *Transaction) error
	UpdateStatus(ctx context.Context, id int64, status string, ref string) error
	GetByOrderID(ctx context.Context, orderID string) (*Transaction, error)
}

// Gateway Interface (Strategy Pattern Port)
type PaymentGateway interface {
	// ProcessPayment simulates the external call. Returns extRefID, error
	ProcessPayment(ctx context.Context, req PaymentRequest, idempotencyKey string) (string, error)
}

// TicketClient Interface (Port for downstream communication)
type TicketClient interface {
	ConfirmPayment(ctx context.Context, orderID string) error
}

// Usecase Interface
type PaymentUsecase interface {
	ProcessPayment(ctx context.Context, req PaymentRequest, idempotencyKey string) (*Transaction, error)
}
