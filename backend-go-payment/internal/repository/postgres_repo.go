package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/harri88/eticketing-challenge/config"
	"github.com/harri88/eticketing-challenge/internal/domain"
	_ "github.com/lib/pq"
)

type pgRepository struct {
	db *sql.DB
}

// NewPostgresConnection creates and returns a new PostgreSQL database connection
func NewPostgresConnection(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.GetDatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.Database.MaxConns)
	db.SetMaxIdleConns(cfg.Database.MinConns)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func NewPostgresRepository(db *sql.DB) domain.TransactionRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (transaction_id, order_id, payment_method, amount, currency, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		tx.TransactionID, tx.OrderID, tx.PaymentMethod, tx.Amount, tx.Currency, tx.Status, tx.CreatedAt,
	).Scan(&tx.ID)
}

func (r *pgRepository) UpdateStatus(ctx context.Context, id int64, status string, ref string) error {
	query := `UPDATE transactions SET status = $1, payment_ref = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, ref, id)
	return err
}

func (r *pgRepository) GetByOrderID(ctx context.Context, orderID string) (*domain.Transaction, error) {
	query := `
		SELECT id, transaction_id, order_id, payment_method, payment_ref, amount, currency, status, created_at
		FROM transactions
		WHERE order_id = $1
	`
	var tx domain.Transaction
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&tx.ID, &tx.TransactionID, &tx.OrderID, &tx.PaymentMethod, &tx.PaymentRef,
		&tx.Amount, &tx.Currency, &tx.Status, &tx.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}

func (r *pgRepository) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	query := `
		SELECT id, transaction_id, order_id, payment_method, payment_ref, amount, currency, status, created_at
		FROM transactions
		WHERE id = $1
	`
	var tx domain.Transaction
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tx.ID, &tx.TransactionID, &tx.OrderID, &tx.PaymentMethod, &tx.PaymentRef,
		&tx.Amount, &tx.Currency, &tx.Status, &tx.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}

func (r *pgRepository) GetByTransactionID(ctx context.Context, transactionID string) (*domain.Transaction, error) {
	query := `
		SELECT id, transaction_id, order_id, payment_method, payment_ref, amount, currency, status, created_at
		FROM transactions
		WHERE transaction_id = $1
	`
	var tx domain.Transaction
	err := r.db.QueryRowContext(ctx, query, transactionID).Scan(
		&tx.ID, &tx.TransactionID, &tx.OrderID, &tx.PaymentMethod, &tx.PaymentRef,
		&tx.Amount, &tx.Currency, &tx.Status, &tx.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}

func (r *pgRepository) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	query := `
		SELECT id, transaction_id, order_id, payment_method, payment_ref, amount, currency, status, created_at
		FROM transactions
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var tx domain.Transaction
		err := rows.Scan(
			&tx.ID, &tx.TransactionID, &tx.OrderID, &tx.PaymentMethod, &tx.PaymentRef,
			&tx.Amount, &tx.Currency, &tx.Status, &tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, rows.Err()
}
