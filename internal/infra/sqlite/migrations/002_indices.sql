-- Migration 002: Add indices and constraints
BEGIN TRANSACTION;

-- Create indices for performance
CREATE INDEX idx_user_balances_user ON user_balances(user_id);
CREATE INDEX idx_tasks_category_status ON tasks(category, status);
CREATE INDEX idx_tasks_last_completed ON tasks(last_completed_at);

-- Add triggers for updated_at
CREATE TRIGGER update_users_timestamp 
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_tasks_timestamp 
AFTER UPDATE ON tasks
FOR EACH ROW
BEGIN
    UPDATE tasks SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_user_balances_timestamp 
AFTER UPDATE ON user_balances
FOR EACH ROW
BEGIN
    UPDATE user_balances 
    SET updated_at = CURRENT_TIMESTAMP 
    WHERE user_id = NEW.user_id AND currency_code = NEW.currency_code;
END;

COMMIT;