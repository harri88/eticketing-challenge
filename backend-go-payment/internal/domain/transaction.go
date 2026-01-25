package domain

import (
	"context"
	"time"
)

// Transaction Entity
type Transaction struct {
	ID            int64     `json:"id" db:"id"`
	TransactionID string    `json:"transaction_id" db:"transaction_id"` // e.g., TXN-CC-17000
	OrderID       string    `json:"order_id" db:"order_id"`
	PaymentMethod string    `json:"payment_method" db:"payment_method"`
	PaymentRef    string    `json:"payment_ref" db:"payment_ref"`
	Amount        float64   `json:"amount" db:"amount"`
	Currency      string    `json:"currency" db:"currency"`
	Status        string    `json:"status" db:"status"` // PENDING, SUCCESS, FAILED
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

const (
	StatusPending = "PENDING"
	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"
)

// PaymentRequest DTO
type PaymentRequest struct {
	OrderID       string `json:"order_id" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"` // credit_card, qr_scan, crypto
	CardNumber    string `json:"card_number,omitempty"`              // for credit_card
	CardExpiry    string `json:"card_expiry,omitempty"`              // for credit_card
	CardCVC       string `json:"card_cvc,omitempty"`                 // for credit_card
	// Removed fields to simplify the request
	// Currency      string `json:"currency"`
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
