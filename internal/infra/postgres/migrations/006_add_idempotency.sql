-- Migration 006: Add idempotency support
BEGIN;

CREATE TABLE idempotency_keys (
    key VARCHAR(255) PRIMARY KEY,
    operation VARCHAR(50) NOT NULL,
    user_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'completed', 'failed')),
    result TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Index for cleanup of expired keys
CREATE INDEX idx_idempotency_keys_expires_at ON idempotency_keys(expires_at);

-- Index for user lookups
CREATE INDEX idx_idempotency_keys_user_id ON idempotency_keys(user_id);

COMMIT;