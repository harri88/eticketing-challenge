-- Database: payment_db

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(100) NOT NULL UNIQUE, -- e.g., TXN-CC-17000000
    order_id VARCHAR(100) NOT NULL,              -- Reference to Ticket Service Order
    payment_method VARCHAR(50) NOT NULL,         -- credit_card, qr_scan
    payment_ref VARCHAR(100),                    -- External reference/gateway ID
    amount DECIMAL(19, 4) NOT NULL,
    currency VARCHAR(3) DEFAULT 'AED',
    status VARCHAR(20) NOT NULL,                 -- SUCCESS, FAILED
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast lookup by Order ID when syncing status
CREATE INDEX idx_payment_order_id ON transactions(order_id);