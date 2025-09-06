# Phase 1: Complete Database Schema (Fresh Start)

## Overview
Create a complete database schema from scratch that supports the full quest and dungeon system. No migrations needed - we're starting fresh.

## Goals
- Create all necessary tables for users, dungeons, quests, and completions
- Ensure schema supports both web interface and Telegram bot integration
- Add proper indexes and constraints
- Support Telegram authentication and user mapping

## Tasks

### 1.1 Create Complete Schema File
**File**: `database/schema.sql` (new file)

```sql
-- Complete ADHD Game Bot Database Schema
-- PostgreSQL 13+

BEGIN;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Update timestamp function
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Users table (supports both web and Telegram users)
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY, -- Can be Telegram user ID or auto-generated
    telegram_user_id BIGINT UNIQUE, -- Telegram user ID for auth
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE, -- For web users
    balance NUMERIC(20, 8) NOT NULL DEFAULT 0,
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Dungeons table (groups/teams)
CREATE TABLE dungeons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    admin_user_id BIGINT NOT NULL,
    telegram_chat_id BIGINT UNIQUE, -- Link to Telegram group
    invite_code VARCHAR(20) UNIQUE, -- For joining via web
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (admin_user_id) REFERENCES users(id) ON DELETE RESTRICT
);

-- Dungeon members
CREATE TABLE dungeon_members (
    dungeon_id UUID NOT NULL,
    user_id BIGINT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member' CHECK (role IN ('admin', 'member')),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (dungeon_id, user_id),
    FOREIGN KEY (dungeon_id) REFERENCES dungeons(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Quests table
CREATE TABLE quests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dungeon_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(10) NOT NULL CHECK (category IN ('daily', 'weekly', 'adhoc')),
    difficulty VARCHAR(10) NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
    
    -- MVP Scoring Configuration
    mode VARCHAR(20) NOT NULL CHECK (mode IN ('BINARY', 'PARTIAL', 'PER_MINUTE')) DEFAULT 'BINARY',
    points_award NUMERIC(20, 8) NOT NULL,
    rate_points_per_min NUMERIC(20, 8), -- For PER_MINUTE mode
    min_minutes INTEGER, -- Optional floor for PER_MINUTE
    max_minutes INTEGER, -- Optional cap for PER_MINUTE
    daily_points_cap NUMERIC(20, 8), -- Optional anti-abuse limit
    
    -- Behavioral Controls
    cooldown_sec INTEGER NOT NULL DEFAULT 0,
    streak_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Operational State
    status VARCHAR(10) NOT NULL CHECK (status IN ('active', 'paused', 'archived')) DEFAULT 'active',
    last_completed_at TIMESTAMP WITH TIME ZONE,
    streak_count INTEGER NOT NULL DEFAULT 0,
    time_zone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    
    -- Metadata
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (dungeon_id) REFERENCES dungeons(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT
);

-- Quest completions (history and points tracking)
CREATE TABLE quest_completions (
    id BIGSERIAL PRIMARY KEY,
    quest_id UUID NOT NULL,
    user_id BIGINT NOT NULL,
    points_awarded NUMERIC(20, 8) NOT NULL,
    completion_ratio DECIMAL(5,4), -- For PARTIAL mode (0.0000 to 1.0000)
    minutes_spent INTEGER, -- For PER_MINUTE mode
    streak_count INTEGER NOT NULL DEFAULT 0,
    notes TEXT, -- Optional completion notes
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (quest_id) REFERENCES quests(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Shop items (rewards system)
CREATE TABLE shop_items (
    id BIGSERIAL PRIMARY KEY,
    dungeon_id UUID, -- NULL for global items
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(20, 8) NOT NULL,
    category VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    stock INTEGER, -- NULL for unlimited
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(dungeon_id, code),
    FOREIGN KEY (dungeon_id) REFERENCES dungeons(id) ON DELETE CASCADE
);

-- Purchases (shop transaction history)
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

-- Telegram bot sessions (for authentication)
CREATE TABLE telegram_sessions (
    telegram_user_id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Idempotency keys (prevent duplicate operations)
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

-- Indexes for performance
CREATE INDEX idx_users_telegram_user_id ON users(telegram_user_id);
CREATE INDEX idx_users_email ON users(email);

CREATE INDEX idx_dungeons_admin ON dungeons(admin_user_id);
CREATE INDEX idx_dungeons_telegram_chat ON dungeons(telegram_chat_id);
CREATE INDEX idx_dungeons_invite_code ON dungeons(invite_code);

CREATE INDEX idx_dungeon_members_user ON dungeon_members(user_id);
CREATE INDEX idx_dungeon_members_dungeon ON dungeon_members(dungeon_id);

CREATE INDEX idx_quests_dungeon ON quests(dungeon_id);
CREATE INDEX idx_quests_status ON quests(status);
CREATE INDEX idx_quests_category ON quests(category);
CREATE INDEX idx_quests_created_by ON quests(created_by);

CREATE INDEX idx_quest_completions_quest ON quest_completions(quest_id);
CREATE INDEX idx_quest_completions_user ON quest_completions(user_id);
CREATE INDEX idx_quest_completions_completed_at ON quest_completions(completed_at);
CREATE INDEX idx_quest_completions_user_quest ON quest_completions(user_id, quest_id);

CREATE INDEX idx_shop_items_dungeon_active ON shop_items(dungeon_id, is_active);
CREATE INDEX idx_shop_items_code ON shop_items(code);

CREATE INDEX idx_purchases_user ON purchases(user_id, purchased_at DESC);
CREATE INDEX idx_purchases_item ON purchases(item_id);

CREATE INDEX idx_telegram_sessions_token ON telegram_sessions(session_token);
CREATE INDEX idx_telegram_sessions_expires ON telegram_sessions(expires_at);

CREATE INDEX idx_idempotency_keys_expires_at ON idempotency_keys(expires_at);
CREATE INDEX idx_idempotency_keys_user_id ON idempotency_keys(user_id);

-- Timestamp triggers
CREATE TRIGGER update_users_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_dungeons_timestamp
    BEFORE UPDATE ON dungeons
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_quests_timestamp
    BEFORE UPDATE ON quests
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_shop_items_timestamp
    BEFORE UPDATE ON shop_items
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

COMMIT;
```

**Testing**:
- Run schema creation on fresh PostgreSQL database
- Verify all tables created with correct columns
- Test all foreign key constraints
- Verify all indexes exist
- Test timestamp triggers

### 1.2 Create Schema Setup Script
**File**: `scripts/setup_database.sh`

```bash
#!/bin/bash
set -e

# Database setup script for ADHD Game Bot

DB_NAME=${DB_NAME:-adhd_bot}
DB_USER=${DB_USER:-postgres}
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}

echo "Setting up database: $DB_NAME"

# Create database if it doesn't exist
createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME 2>/dev/null || echo "Database already exists"

# Apply schema
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f database/schema.sql

echo "Database setup complete!"
```

**Testing**:
- Test script with fresh PostgreSQL instance
- Verify script is idempotent (can run multiple times)
- Test with different database configurations

### 1.3 Create Sample Data Script
**File**: `database/sample_data.sql`

```sql
-- Sample data for development and testing

BEGIN;

-- Sample users
INSERT INTO users (id, telegram_user_id, username, balance, timezone) VALUES
(1, 123456789, 'Alice', 100.00, 'America/New_York'),
(2, 987654321, 'Bob', 50.00, 'Europe/London'),
(3, 555666777, 'Charlie', 75.00, 'Asia/Tokyo');

-- Sample dungeon
INSERT INTO dungeons (id, title, description, admin_user_id, telegram_chat_id, invite_code) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'ADHD Warriors', 'A supportive group for ADHD management', 1, -1001234567890, 'ADHD2024');

-- Add members to dungeon
INSERT INTO dungeon_members (dungeon_id, user_id, role) VALUES
('550e8400-e29b-41d4-a716-446655440000', 1, 'admin'),
('550e8400-e29b-41d4-a716-446655440000', 2, 'member'),
('550e8400-e29b-41d4-a716-446655440000', 3, 'member');

-- Sample quests
INSERT INTO quests (id, dungeon_id, title, description, category, difficulty, mode, points_award, created_by) VALUES
('660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', 'Morning Routine', 'Complete your morning routine', 'daily', 'easy', 'BINARY', 10.00, 1),
('660e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440000', 'Exercise', '30 minutes of physical activity', 'daily', 'medium', 'PER_MINUTE', 1.00, 1),
('660e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440000', 'Weekly Planning', 'Plan your week ahead', 'weekly', 'medium', 'BINARY', 25.00, 1);

-- Update rate_points_per_min for PER_MINUTE quest
UPDATE quests SET rate_points_per_min = 1.00, min_minutes = 10, max_minutes = 60 
WHERE mode = 'PER_MINUTE';

-- Sample shop items
INSERT INTO shop_items (dungeon_id, code, name, description, price, category) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'coffee', 'Virtual Coffee', 'Treat yourself to a coffee break', 15.00, 'treats'),
('550e8400-e29b-41d4-a716-446655440000', 'movie', 'Movie Night', 'Guilt-free movie watching', 50.00, 'entertainment'),
(NULL, 'badge_bronze', 'Bronze Achievement', 'Bronze level achievement badge', 100.00, 'achievements');

COMMIT;
```

**Testing**:
- Apply sample data to test database
- Verify all foreign key relationships work
- Test queries against sample data

## Testing Strategy

### 1.4 Schema Validation Tests
**File**: `test/database/schema_test.go`

```go
func TestDatabaseSchema(t *testing.T) {
    // Test all tables exist
    // Test all columns have correct types
    // Test all constraints work
    // Test all indexes exist
    // Test foreign key relationships
}

func TestSampleDataIntegration(t *testing.T) {
    // Load sample data
    // Test complex queries
    // Verify data integrity
}
```

### 1.5 Performance Tests
**File**: `test/database/performance_test.go`

```go
func TestQueryPerformance(t *testing.T) {
    // Test common query patterns
    // Verify indexes are used
    // Test with larger datasets
}
```

## Verification Checklist

- [ ] Schema creates all necessary tables
- [ ] All foreign key relationships work correctly
- [ ] All constraints prevent invalid data
- [ ] All indexes improve query performance
- [ ] Timestamp triggers update correctly
- [ ] Sample data loads successfully
- [ ] Schema supports both web and Telegram authentication
- [ ] Performance tests pass with good query times

## Files to Create

1. `database/schema.sql` - Complete database schema
2. `database/sample_data.sql` - Sample data for testing
3. `scripts/setup_database.sh` - Database setup script
4. `test/database/schema_test.go` - Schema validation tests
5. `test/database/performance_test.go` - Performance tests

## Success Criteria

✅ Complete database schema supports all MVP features
✅ Schema handles both web users and Telegram users
✅ All relationships and constraints work correctly
✅ Performance is optimized with proper indexes
✅ Sample data demonstrates full functionality
✅ Tests validate schema integrity

## Next Phase
Once database schema is complete, move to **Phase 2: Telegram Bot Authentication** to implement the auth flow between bot and web interface.