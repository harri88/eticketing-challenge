package usecase

import (
	"context"

	"github.com/harri88/eticketing-challenge/internal/domain"
)

type TransactionUsecase struct {
	repo domain.TransactionRepository
}

func NewTransactionUsecase(repo domain.TransactionRepository) *TransactionUsecase {
	return &TransactionUsecase{
		repo: repo,
	}
}

func (u *TransactionUsecase) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	return u.repo.GetAll(ctx)
}

func (u *TransactionUsecase) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *TransactionUsecase) GetByTransactionID(ctx context.Context, transactionID string) (*domain.Transaction, error) {
	return u.repo.GetByTransactionID(ctx, transactionID)
}

func (u *TransactionUsecase) GetByOrderID(ctx context.Context, orderID string) (*domain.Transaction, error) {
	return u.repo.GetByOrderID(ctx, orderID)
}
