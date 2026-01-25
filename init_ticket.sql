-- init.sql
CREATE TABLE IF NOT EXISTS tickets (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    quota INT NOT NULL
);

INSERT INTO tickets (id, name, price, currency, quota) VALUES
('gold_001', 'Gold', 100.00, 'AED', 100),
('prem_001', 'Premium', 200.00, 'AED', 100),
('vip_001', 'VIP', 500.00, 'AED', 100)
ON CONFLICT (id) DO NOTHING;




-- 2. Create Orders Table
CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(50) PRIMARY KEY,
    customer_email VARCHAR(150) NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'created', 'paid', 'cancelled'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Create Order Items Table (for normalization)
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(50) REFERENCES orders(id),
    ticket_id VARCHAR(50) REFERENCES tickets(id),
    quantity INT NOT NULL,
    price_at_purchase DECIMAL(10, 2) NOT NULL
);

-- 1. Add column to track items currently in checkout/reserved
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS held_quota INT DEFAULT 0;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS paid_at TIMESTAMP DEFAULT NULL;