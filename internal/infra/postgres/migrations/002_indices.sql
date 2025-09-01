-- Migration 002: Add indices and constraints for PostgreSQL
BEGIN;

-- Create indices for performance
CREATE INDEX idx_user_balances_user ON user_balances(user_id);
CREATE INDEX idx_tasks_category_status ON tasks(category, status);
CREATE INDEX idx_tasks_last_completed ON tasks(last_completed_at);

-- Update timestamp function
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add triggers for updated_at
CREATE TRIGGER update_users_timestamp
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_tasks_timestamp
BEFORE UPDATE ON tasks
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_user_balances_timestamp
BEFORE UPDATE ON user_balances
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

COMMIT;