# ADHD Game Bot - Database Schema Design

## SQLite Configuration

```sql
-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- Use WAL mode for better concurrency
PRAGMA journal_mode = WAL;

-- Optimize for performance
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = -64000; -- 64MB
PRAGMA temp_store = MEMORY;
```

## Migration Files

### 001_init.sql - Base Tables

```sql
-- Schema version tracking
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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

-- User balances (separate table for atomic updates)
CREATE TABLE user_balances (
    user_id INTEGER NOT NULL,
    currency_code TEXT NOT NULL,
    amount TEXT NOT NULL DEFAULT '0', -- Stored as string for decimal precision
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, currency_code),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

-- Currencies table
CREATE TABLE currencies (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    decimals INTEGER NOT NULL DEFAULT 2,
    conversion_rates_json TEXT DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tasks table
CREATE TABLE tasks (
    id TEXT PRIMARY KEY, -- UUID
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
    streak_enabled BOOLEAN NOT NULL DEFAULT 0,
    status TEXT NOT NULL CHECK (status IN ('inactive', 'active')) DEFAULT 'active',
    last_completed_at TIMESTAMP,
    streak_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Task tags (many-to-many)
CREATE TABLE task_tags (
    task_id TEXT NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY (task_id, tag),
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

-- Task prerequisites (many-to-many)
CREATE TABLE task_prerequisites (
    task_id TEXT NOT NULL,
    prerequisite_task_id TEXT NOT NULL,
    PRIMARY KEY (task_id, prerequisite_task_id),
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (prerequisite_task_id) REFERENCES tasks(id) ON DELETE RESTRICT
);

-- Timer instances
CREATE TABLE timers (
    id TEXT PRIMARY KEY, -- UUID
    task_id TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    started_at_wall TIMESTAMP NOT NULL,
    started_at_mono INTEGER NOT NULL, -- nanoseconds since arbitrary point
    initial_duration INTEGER NOT NULL, -- seconds
    tier_deadlines_json TEXT NOT NULL, -- Array of absolute timestamps
    state TEXT NOT NULL CHECK (state IN ('running', 'paused', 'completed', 'expired', 'cancelled')),
    last_tick_at TIMESTAMP,
    snooze_count INTEGER NOT NULL DEFAULT 0,
    total_extended INTEGER NOT NULL DEFAULT 0, -- seconds
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE RESTRICT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

-- Ledger entries
CREATE TABLE ledger (
    id TEXT PRIMARY KEY, -- UUID
    user_id INTEGER NOT NULL,
    currency_code TEXT NOT NULL,
    delta TEXT NOT NULL, -- Signed decimal as string
    reason TEXT NOT NULL CHECK (reason IN ('task_complete', 'redemption', 'manual_adjust')),
    ref_id TEXT NOT NULL, -- Idempotency key
    metadata_json TEXT DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (currency_code) REFERENCES currencies(code) ON DELETE RESTRICT
);

-- Store items
CREATE TABLE store_items (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    cost_json TEXT NOT NULL, -- Map of currency to amount
    available BOOLEAN NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Redemptions
CREATE TABLE redemptions (
    id TEXT PRIMARY KEY, -- UUID
    user_id INTEGER NOT NULL,
    item_code TEXT NOT NULL,
    item_name TEXT NOT NULL, -- Denormalized for history
    cost_json TEXT NOT NULL, -- Denormalized for history
    status TEXT NOT NULL CHECK (status IN ('pending', 'fulfilled', 'cancelled')) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fulfilled_at TIMESTAMP,
    fulfilled_by INTEGER,
    notes TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (fulfilled_by) REFERENCES users(id) ON DELETE RESTRICT
);

-- Audit log
CREATE TABLE audit_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    actor_user_id INTEGER NOT NULL,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    payload_json TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (actor_user_id) REFERENCES users(id) ON DELETE RESTRICT
);

-- Insert migration record
INSERT INTO schema_migrations (version) VALUES (1);
```

### 002_indices.sql - Performance and Uniqueness Indices

```sql
-- Ledger indices
CREATE INDEX idx_ledger_user_currency_date ON ledger(user_id, currency_code, created_at DESC);
CREATE UNIQUE INDEX idx_ledger_idempotency ON ledger(reason, ref_id);

-- Timer indices
CREATE INDEX idx_timers_task_user_state ON timers(task_id, user_id, state);
CREATE INDEX idx_timers_started_at ON timers(started_at_wall);
CREATE INDEX idx_timers_user_active ON timers(user_id, state) WHERE state = 'running';

-- Enforce single active timer per task/user
CREATE UNIQUE INDEX idx_timers_single_active ON timers(task_id, user_id) 
    WHERE state = 'running';

-- Task indices
CREATE INDEX idx_tasks_category_status ON tasks(category, status);
CREATE INDEX idx_tasks_last_completed ON tasks(last_completed_at);

-- Redemption indices
CREATE INDEX idx_redemptions_user_status ON redemptions(user_id, status);
CREATE INDEX idx_redemptions_created ON redemptions(created_at DESC);

-- Audit log indices
CREATE INDEX idx_audit_actor_date ON audit_log(actor_user_id, created_at DESC);
CREATE INDEX idx_audit_entity ON audit_log(entity_type, entity_id);

-- User balance index for quick lookups
CREATE INDEX idx_user_balances_user ON user_balances(user_id);

-- Insert migration record
INSERT INTO schema_migrations (version) VALUES (2);
```

### 003_schedule.sql - Schedule Management Tables

```sql
-- Schedule occurrences (materialized view of upcoming task instances)
CREATE TABLE schedule_occurrences (
    id TEXT PRIMARY KEY, -- UUID
    task_id TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    window_start TIMESTAMP NOT NULL,
    window_end TIMESTAMP NOT NULL,
    due_at TIMESTAMP NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('queued', 'triggered', 'skipped')) DEFAULT 'queued',
    triggered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Ensure no duplicate occurrences
CREATE UNIQUE INDEX idx_schedule_unique_occurrence 
    ON schedule_occurrences(task_id, user_id, due_at);

-- Performance indices
CREATE INDEX idx_schedule_due_status 
    ON schedule_occurrences(due_at, status) 
    WHERE status = 'queued';

CREATE INDEX idx_schedule_user_window 
    ON schedule_occurrences(user_id, window_start, window_end);

-- Schedule run history (for debugging and analytics)
CREATE TABLE schedule_run_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_at TIMESTAMP NOT NULL,
    occurrences_created INTEGER NOT NULL DEFAULT 0,
    occurrences_triggered INTEGER NOT NULL DEFAULT 0,
    errors_json TEXT,
    duration_ms INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert migration record
INSERT INTO schema_migrations (version) VALUES (3);
```

### 004_triggers.sql - Database Triggers for Consistency

```sql
-- Update timestamp triggers
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

CREATE TRIGGER update_timers_timestamp 
    AFTER UPDATE ON timers
    FOR EACH ROW
BEGIN
    UPDATE timers SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_user_balances_timestamp 
    AFTER UPDATE ON user_balances
    FOR EACH ROW
BEGIN
    UPDATE user_balances 
    SET updated_at = CURRENT_TIMESTAMP 
    WHERE user_id = NEW.user_id AND currency_code = NEW.currency_code;
END;

-- Prevent timer state regression
CREATE TRIGGER prevent_timer_state_regression
    BEFORE UPDATE ON timers
    FOR EACH ROW
    WHEN OLD.state IN ('completed', 'expired', 'cancelled') 
        AND NEW.state IN ('running', 'paused')
BEGIN
    SELECT RAISE(ABORT, 'Cannot reactivate finished timer');
END;

-- Auto-create user balances for new currencies
CREATE TRIGGER create_user_balances_for_currency
    AFTER INSERT ON currencies
    FOR EACH ROW
BEGIN
    INSERT INTO user_balances (user_id, currency_code, amount)
    SELECT id, NEW.code, '0' FROM users;
END;

-- Auto-create balance entries for new users
CREATE TRIGGER create_balances_for_user
    AFTER INSERT ON users
    FOR EACH ROW
BEGIN
    INSERT INTO user_balances (user_id, currency_code, amount)
    SELECT NEW.id, code, '0' FROM currencies;
END;

-- Insert migration record
INSERT INTO schema_migrations (version) VALUES (4);
```

### 005_views.sql - Useful Views

```sql
-- Active timers view
CREATE VIEW v_active_timers AS
SELECT 
    t.id,
    t.task_id,
    t.user_id,
    tk.title as task_title,
    u.display_name as user_name,
    t.started_at_wall,
    t.state,
    t.snooze_count,
    CAST((julianday('now') - julianday(t.started_at_wall)) * 86400 AS INTEGER) as elapsed_seconds
FROM timers t
JOIN tasks tk ON t.task_id = tk.id
JOIN users u ON t.user_id = u.id
WHERE t.state = 'running';

-- User balance summary
CREATE VIEW v_user_balance_summary AS
SELECT 
    u.id as user_id,
    u.display_name,
    ub.currency_code,
    c.name as currency_name,
    ub.amount,
    ub.updated_at
FROM users u
JOIN user_balances ub ON u.id = ub.user_id
JOIN currencies c ON ub.currency_code = c.code
ORDER BY u.id, c.code;

-- Recent completions
CREATE VIEW v_recent_completions AS
SELECT 
    t.id as timer_id,
    t.task_id,
    tk.title as task_title,
    t.user_id,
    u.display_name as user_name,
    t.completed_at,
    l.currency_code,
    l.delta as reward_amount
FROM timers t
JOIN tasks tk ON t.task_id = tk.id
JOIN users u ON t.user_id = u.id
LEFT JOIN ledger l ON l.ref_id = t.id AND l.reason = 'task_complete'
WHERE t.state = 'completed'
    AND t.completed_at > datetime('now', '-7 days')
ORDER BY t.completed_at DESC;

-- Insert migration record
INSERT INTO schema_migrations (version) VALUES (5);
```

## Data Types and Conventions

### UUID Storage
- Store as TEXT (36 characters)
- Format: `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`

### Decimal Storage
- Store as TEXT to preserve precision
- Format: Standard decimal notation (e.g., "123.45")
- Always store with full precision

### Timestamp Storage
- Use SQLite's TIMESTAMP type (stored as TEXT)
- Always store in UTC
- Format: ISO 8601 (`YYYY-MM-DD HH:MM:SS`)

### JSON Storage
- Store complex objects as JSON TEXT
- Validate JSON structure in application layer
- Use JSON1 extension for queries when needed

### Duration Storage
- Store as INTEGER (seconds)
- Convert to/from time.Duration in Go

## Example Queries

### Get User's Active Timer
```sql
SELECT * FROM timers 
WHERE user_id = ? 
  AND task_id = ? 
  AND state = 'running'
LIMIT 1;
```

### Calculate User Balance
```sql
SELECT 
    currency_code,
    SUM(CAST(delta AS DECIMAL)) as balance
FROM ledger
WHERE user_id = ?
GROUP BY currency_code;
```

### Get Upcoming Schedule Occurrences
```sql
SELECT * FROM schedule_occurrences
WHERE user_id = ?
  AND status = 'queued'
  AND window_start <= datetime('now', '+7 days')
ORDER BY due_at ASC;
```

### Atomic Balance Update
```sql
-- Start transaction
BEGIN IMMEDIATE;

-- Insert ledger entry (will fail if duplicate ref_id)
INSERT INTO ledger (id, user_id, currency_code, delta, reason, ref_id)
VALUES (?, ?, ?, ?, ?, ?);

-- Update balance
UPDATE user_balances 
SET amount = printf('%.2f', CAST(amount AS DECIMAL) + CAST(? AS DECIMAL))
WHERE user_id = ? AND currency_code = ?;

-- Audit log
INSERT INTO audit_log (actor_user_id, action, entity_type, entity_id, payload_json)
VALUES (?, ?, ?, ?, ?);

COMMIT;
```

## Backup Strategy

```sql
-- Online backup (while database is in use)
PRAGMA wal_checkpoint(TRUNCATE);
.backup backup.db

-- Or use SQLite's backup API from Go
```

## Performance Considerations

1. **WAL Mode**: Enables concurrent reads while writing
2. **Proper Indices**: Cover all common query patterns
3. **Denormalization**: Store item details in redemptions for history
4. **Materialized Views**: schedule_occurrences for fast queries
5. **Connection Pooling**: Reuse connections in Go application

## Security Considerations

1. **Foreign Keys**: Always enforced to maintain referential integrity
2. **Check Constraints**: Validate enums at database level
3. **Restrict Deletes**: Financial data cannot be deleted
4. **Audit Everything**: Complete audit trail for compliance