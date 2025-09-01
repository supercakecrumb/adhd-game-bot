-- Migration 003: Refactor currency system for multi-chat support
BEGIN;

-- Drop old currency table if exists
DROP TABLE IF EXISTS currencies CASCADE;

-- Create new currencies table with chat association
CREATE TABLE currencies (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    code VARCHAR(10) NOT NULL,
    name VARCHAR(100) NOT NULL,
    decimals INTEGER NOT NULL DEFAULT 2,
    is_base_currency BOOLEAN NOT NULL DEFAULT FALSE,
    exchange_rates_json JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(chat_id, code)
);

-- Add chat_id to users table
ALTER TABLE users ADD COLUMN chat_id BIGINT;

-- Update user_balances to use currency_id instead of currency_code
ALTER TABLE user_balances DROP CONSTRAINT IF EXISTS user_balances_pkey;
ALTER TABLE user_balances ADD COLUMN currency_id BIGINT;
ALTER TABLE user_balances ADD CONSTRAINT user_balances_pkey PRIMARY KEY (user_id, currency_id);
ALTER TABLE user_balances ADD CONSTRAINT fk_currency FOREIGN KEY (currency_id) REFERENCES currencies(id) ON DELETE RESTRICT;

-- Create index for base currency lookup
CREATE INDEX idx_currencies_base ON currencies(chat_id, is_base_currency) WHERE is_base_currency = TRUE;

-- Ensure only one base currency per chat
CREATE UNIQUE INDEX idx_one_base_currency_per_chat ON currencies(chat_id) WHERE is_base_currency = TRUE;

COMMIT;