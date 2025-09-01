-- Migration 001: Initial schema setup for PostgreSQL
BEGIN;

-- Users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    role VARCHAR(10) NOT NULL CHECK (role IN ('admin', 'member')) DEFAULT 'member',
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    display_name VARCHAR(255) NOT NULL,
    preferences_json JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User balances
CREATE TABLE user_balances (
    user_id BIGINT NOT NULL,
    currency_code VARCHAR(10) NOT NULL,
    amount NUMERIC(20, 8) NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, currency_code),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

-- Currencies
CREATE TABLE currencies (
    code VARCHAR(10) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    decimals INTEGER NOT NULL DEFAULT 2,
    conversion_rates_json JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tasks
CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(10) NOT NULL CHECK (category IN ('daily', 'weekly', 'adhoc')),
    difficulty VARCHAR(10) NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
    schedule_json JSONB NOT NULL DEFAULT '{}',
    base_duration INTEGER NOT NULL, -- seconds
    grace_period INTEGER NOT NULL DEFAULT 0, -- seconds
    cooldown INTEGER NOT NULL DEFAULT 0, -- seconds
    reward_curve_json JSONB NOT NULL,
    partial_credit_json JSONB,
    streak_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(10) NOT NULL CHECK (status IN ('inactive', 'active')) DEFAULT 'active',
    last_completed_at TIMESTAMP WITH TIME ZONE,
    streak_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMIT;