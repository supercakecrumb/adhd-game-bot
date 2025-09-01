-- Migration 004: Simplify to single currency per chat and add shop
BEGIN;

-- Drop old currency-related tables
DROP TABLE IF EXISTS user_balances CASCADE;
DROP TABLE IF EXISTS currencies CASCADE;

-- Create chat_configs table
CREATE TABLE chat_configs (
    chat_id BIGINT PRIMARY KEY,
    currency_name VARCHAR(50) NOT NULL DEFAULT 'Points',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add balance column to users table
ALTER TABLE users ADD COLUMN balance NUMERIC(20, 8) NOT NULL DEFAULT 0;

-- Create shop_items table
CREATE TABLE shop_items (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL DEFAULT 0, -- 0 for global items
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(20, 8) NOT NULL,
    category VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    stock INTEGER, -- NULL for unlimited
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(chat_id, code)
);

-- Create purchases table
CREATE TABLE purchases (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    item_id BIGINT NOT NULL,
    item_name VARCHAR(255) NOT NULL, -- Denormalized for history
    item_price NUMERIC(20, 8) NOT NULL, -- Price at time of purchase
    quantity INTEGER NOT NULL DEFAULT 1,
    total_cost NUMERIC(20, 8) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'completed',
    purchased_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (item_id) REFERENCES shop_items(id) ON DELETE RESTRICT
);

-- Create indices
CREATE INDEX idx_shop_items_chat_active ON shop_items(chat_id, is_active);
CREATE INDEX idx_purchases_user ON purchases(user_id, purchased_at DESC);
CREATE INDEX idx_purchases_item ON purchases(item_id);

-- Update timestamp triggers
CREATE TRIGGER update_chat_configs_timestamp
BEFORE UPDATE ON chat_configs
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_shop_items_timestamp
BEFORE UPDATE ON shop_items
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

COMMIT;