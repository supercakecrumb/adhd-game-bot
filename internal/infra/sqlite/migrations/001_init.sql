-- Migration 001: Initial schema setup
BEGIN TRANSACTION;

-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    role TEXT NOT NULL CHECK (role IN ('admin', 'member')) DEFAULT 'member',
    timezone TEXT NOT NULL DEFAULT 'UTC',
    display_name TEXT NOT NULL,
    preferences_json TEXT NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User balances
CREATE TABLE user_balances (
    user_id INTEGER NOT NULL,
    currency_code TEXT NOT NULL,
    amount TEXT NOT NULL DEFAULT '0',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, currency_code),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

-- Currencies
CREATE TABLE currencies (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    decimals INTEGER NOT NULL DEFAULT 2,
    conversion_rates_json TEXT DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tasks
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    category TEXT NOT NULL CHECK (category IN ('daily', 'weekly', 'adhoc')),
    difficulty TEXT NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
    schedule_json TEXT NOT NULL DEFAULT '{}',
    base_duration INTEGER NOT NULL, -- seconds
    grace_period INTEGER NOT NULL DEFAULT 0, -- seconds
    cooldown INTEGER NOT NULL DEFAULT 0, -- seconds
    reward_curve_json TEXT NOT NULL,
    partial_credit_json TEXT,
    streak_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    status TEXT NOT NULL CHECK (status IN ('inactive', 'active')) DEFAULT 'active',
    last_completed_at TIMESTAMP,
    streak_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;