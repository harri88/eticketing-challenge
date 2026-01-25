-- Database: ledger_db

CREATE TABLE IF NOT EXISTS ledger_entries (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(100) NOT NULL, -- Links back to Payment Service
    account_name VARCHAR(100) NOT NULL,   -- e.g., 'Cash_Account' or 'Ticket_Revenue'
    entry_type VARCHAR(10) NOT NULL,      -- DEBIT or CREDIT
    amount DECIMAL(19, 4) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Crucial index for financial auditing
CREATE INDEX idx_ledger_transaction_id ON ledger_entries(transaction_id);