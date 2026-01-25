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
	// Implementation omitted for brevity, standard SELECT...
	return nil, nil
}
