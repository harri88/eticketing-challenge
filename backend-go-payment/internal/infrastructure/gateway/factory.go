package gateway

import (
	"context"
	"errors"
	"time"

	"github.com/harri88/eticketing-challenge/internal/domain"
)

// PaymentGatewayFactory decides which strategy to use
type PaymentGatewayFactory struct {
	cc     domain.PaymentGateway
	qr     domain.PaymentGateway
	crypto domain.PaymentGateway
}

func NewPaymentGatewayFactory() *PaymentGatewayFactory {
	return &PaymentGatewayFactory{
		cc:     &CreditCardStrategy{},
		qr:     &QRScanStrategy{},
		crypto: &CryptoStrategy{},
	}
}

func (f *PaymentGatewayFactory) GetStrategy(method string) (domain.PaymentGateway, error) {
	switch method {
	case "credit_card":
		return f.cc, nil
	case "qr_scan":
		return f.qr, nil
	case "crypto":
		return f.crypto, nil
	default:
		return nil, errors.New("unsupported payment method")
	}
}

// --- Strategies ---

// CreditCardStrategy (Visa/Mastercard simulation)
type CreditCardStrategy struct{}

func (s *CreditCardStrategy) ProcessPayment(ctx context.Context, req domain.PaymentRequest, idempotencyKey string) (string, error) {
	// Simulation: Instant success
	// In real life: Send idempotencyKey to Stripe/Adyen
	return "REF-CC-" + idempotencyKey, nil
}

// QRScanStrategy (UPI/QRIS simulation)
type QRScanStrategy struct{}

func (s *QRScanStrategy) ProcessPayment(ctx context.Context, req domain.PaymentRequest, idempotencyKey string) (string, error) {
	// Simulation: 8 seconds delay as per requirements
	select {
	case <-time.After(8 * time.Second):
		return "REF-QR-" + idempotencyKey, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// CryptoStrategy (Future proofing)
type CryptoStrategy struct{}

func (s *CryptoStrategy) ProcessPayment(ctx context.Context, req domain.PaymentRequest, idempotencyKey string) (string, error) {
	// Simulation: Instant for now
	return "REF-CRYPTO-" + idempotencyKey, nil
}
