# ADHD Game Bot - Core Domain System

A gamified task management backend system designed to help individuals with ADHD manage their daily routines through positive reinforcement and structured task management.

## Table of Contents
- [Overview](#overview)
- [Core Capabilities](#core-capabilities)
- [Architecture](#architecture)
- [Domain Model](#domain-model)
- [Services](#services)
- [Database Schema](#database-schema)
- [Testing Strategy](#testing-strategy)
- [Getting Started](#getting-started)
- [Development](#development)

## Overview

This is a backend domain system that implements gamification mechanics for task and routine management. The system is designed with Clean Architecture principles, making it frontend-agnostic and ready for integration with various delivery mechanisms (web, mobile, chat bots, etc.).

### Key Design Principles
- **Domain-Driven Design**: Rich domain model with business logic encapsulated in entities
- **Clean Architecture**: Clear separation between domain, application, and infrastructure layers
- **Test-Driven Development**: Comprehensive test coverage including unit, integration, and property-based tests
- **Transactional Integrity**: All financial operations are atomic with proper rollback support
- **Idempotency**: Safe retry mechanisms for all critical operations

## Core Capabilities

### 1. User Management
- User creation and management
- Balance tracking with decimal precision
- Timezone-aware user profiles
- Chat-based user grouping (for multi-tenant scenarios)

### 2. Task System
The system supports three types of tasks:

#### Daily Tasks
- Repeat every day at specified times
- Configurable completion windows
- Automatic streak tracking

#### Weekly Tasks
- Occur on specific days of the week
- Flexible scheduling with timezone support
- Prerequisite task support

#### Adhoc Tasks
- One-time or irregular tasks
- Custom reward configurations
- Tag-based categorization

### 3. Timer Mechanics
The timer system implements "stretchy" timers with multiple reward tiers:

```go
type RewardTier struct {
    Duration time.Duration
    Reward   valueobject.Decimal
    Name     string // Bronze, Silver, Gold, Platinum
}
```

Timer features:
- Multiple duration options with corresponding rewards
- Pause/resume functionality
- Snooze policies (demotion or time decrement)
- Crash recovery with state persistence

### 4. Currency System

#### Decimal Precision
All monetary calculations use [`valueobject.Decimal`](internal/domain/valueobject/decimal.go):
```go
// Example usage
balance := valueobject.NewDecimal("100.50")
cost := valueobject.NewDecimal("25.25")
newBalance := balance.Sub(cost) // 75.25
```

#### Currency Configuration
- Each chat/group can have custom currency names
- Default: "Points"
- Examples: "Focus Coins", "Productivity Tokens"

### 5. Shop System

#### Shop Items
Items can be:
- **Global**: Available to all users (ChatID = 0)
- **Chat-specific**: Available only to users in a specific chat
- **Limited stock**: Finite quantity available
- **Unlimited**: Always available

#### Purchase Flow
1. User browses available items
2. System validates:
   - User balance
   - Item availability
   - Stock levels
3. Atomic transaction:
   - Deduct user balance
   - Update item stock
   - Create purchase record
   - Generate audit trail

### 6. Scheduling Engine

#### Timezone Support
- All times stored in UTC
- User-specific timezone conversion
- Proper DST handling
- IANA timezone database support

#### Recurrence Patterns
```go
type Schedule struct {
    Type      ScheduleType // daily, weekly, custom
    Time      time.Time    // Base time in user's timezone
    Weekdays  []time.Weekday
    Interval  *time.Duration
}
```

### 7. Idempotency System
Prevents duplicate operations:
```go
type IdempotencyKey struct {
    Key       string
    UserID    int64
    Operation string
    Status    string // pending, completed, failed
    Result    []byte // Cached result
    ExpiresAt time.Time
}
```

## Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────┐
│              Delivery Layer                      │
│         (Future: Web, Mobile, Bots)              │
├─────────────────────────────────────────────────┤
│             Application Layer                    │
│           (Use Cases/Services)                   │
├─────────────────────────────────────────────────┤
│              Domain Layer                        │
│       (Entities, Value Objects, Rules)           │
├─────────────────────────────────────────────────┤
│           Infrastructure Layer                   │
│    (Repositories, External Services)             │
└─────────────────────────────────────────────────┘
```

### Project Structure
```
adhd-game-bot/
├── internal/
│   ├── domain/           # Core business logic
│   │   ├── entity/       # Domain entities
│   │   │   ├── user.go
│   │   │   ├── task.go
│   │   │   ├── shop_item.go
│   │   │   └── idempotency.go
│   │   └── valueobject/  # Value objects
│   │       └── decimal.go
│   ├── usecase/          # Application services
│   │   ├── shop_service.go
│   │   └── task_service.go
│   ├── ports/            # Interface definitions
│   │   ├── repositories.go
│   │   └── transaction.go
│   └── infra/            # Implementations
│       ├── postgres/     # PostgreSQL repositories
│       ├── sqlite/       # SQLite repositories
│       └── inmemory/     # In-memory implementations
└── test/                 # Comprehensive test suite
```

## Domain Model

### Core Entities

#### User
```go
type User struct {
    ID        int64
    ChatID    int64
    Username  string
    Balance   valueobject.Decimal
    Timezone  string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### Task
```go
type Task struct {
    ID            string
    UserID        int64
    Name          string
    Description   string
    Category      TaskCategory
    Schedule      *Schedule
    Prerequisites []string
    Tags          []string
    RewardTiers   []RewardTier
    StreakCount   int
    Timezone      string
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

#### ShopItem
```go
type ShopItem struct {
    ID          int64
    ChatID      int64
    Code        string
    Name        string
    Description string
    Price       valueobject.Decimal
    Stock       *int
    IsActive    bool
    Category    string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

#### Purchase
```go
type Purchase struct {
    ID         int64
    UserID     int64
    ItemID     int64
    ItemName   string
    ItemPrice  valueobject.Decimal
    Quantity   int
    TotalCost  valueobject.Decimal
    Status     string
    CreatedAt  time.Time
}
```

## Services

### ShopService
Handles all shop-related operations:
```go
// Create a new shop item
CreateShopItem(ctx context.Context, item *entity.ShopItem) error

// Get available items for a chat
GetShopItems(ctx context.Context, chatID int64) ([]*entity.ShopItem, error)

// Purchase an item
PurchaseItem(ctx context.Context, userID int64, itemCode string, quantity int) (*entity.Purchase, error)

// Get user's purchase history
GetUserPurchases(ctx context.Context, userID int64) ([]*entity.Purchase, error)

// Configure currency name for a chat
SetCurrencyName(ctx context.Context, chatID int64, currencyName string) error
```

### TaskService
Manages task lifecycle:
```go
// Create a new task
CreateTask(ctx context.Context, task *entity.Task) error

// Get user's tasks
GetUserTasks(ctx context.Context, userID int64) ([]*entity.Task, error)

// Start a task timer
StartTaskTimer(ctx context.Context, taskID string, duration time.Duration) (*entity.Timer, error)

// Complete a task
CompleteTask(ctx context.Context, taskID string) (*entity.CompletionResult, error)

// Update task streak
UpdateStreak(ctx context.Context, taskID string) error
```

## Database Schema

### PostgreSQL Schema

#### Core Tables
```sql
-- Users table
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    username TEXT,
    balance DECIMAL(20,8) NOT NULL DEFAULT 0,
    timezone TEXT NOT NULL DEFAULT 'UTC',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tasks table
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,
    description TEXT,
    category TEXT NOT NULL,
    schedule JSONB,
    prerequisites TEXT[],
    tags TEXT[],
    reward_tiers JSONB,
    streak_count INT DEFAULT 0,
    timezone TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Shop items table
CREATE TABLE shop_items (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    price DECIMAL(20,8) NOT NULL,
    stock INT,
    is_active BOOLEAN DEFAULT true,
    category TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(chat_id, code)
);

-- Purchases table
CREATE TABLE purchases (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    item_id BIGINT NOT NULL REFERENCES shop_items(id),
    item_name TEXT NOT NULL,
    item_price DECIMAL(20,8) NOT NULL,
    quantity INT NOT NULL,
    total_cost DECIMAL(20,8) NOT NULL,
    status TEXT NOT NULL DEFAULT 'completed',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Chat configurations
CREATE TABLE chat_configs (
    chat_id BIGINT PRIMARY KEY,
    currency_name TEXT NOT NULL DEFAULT 'Points',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Idempotency keys
CREATE TABLE idempotency_keys (
    key TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    operation TEXT NOT NULL,
    status TEXT NOT NULL,
    result BYTEA,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);
```

### Indexes
```sql
CREATE INDEX idx_users_chat_id ON users(chat_id);
CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_shop_items_chat_active ON shop_items(chat_id, is_active);
CREATE INDEX idx_purchases_user_id ON purchases(user_id);
CREATE INDEX idx_idempotency_expires ON idempotency_keys(expires_at);
```

## Testing Strategy

### Unit Tests
Test individual components in isolation:
- Domain entities and value objects
- Business rule validation
- Service layer logic

Example:
```go
func TestDecimalArithmetic(t *testing.T) {
    a := valueobject.NewDecimal("10.50")
    b := valueobject.NewDecimal("5.25")
    result := a.Add(b)
    assert.Equal(t, "15.75", result.String())
}
```

### Integration Tests
Test repository implementations:
- Database operations
- Transaction handling
- Query performance

### Property-Based Tests
Test invariants with random data:
```go
func TestBalanceNeverNegative(t *testing.T) {
    quick.Check(func(initial, deduction float64) bool {
        if initial < 0 || deduction < 0 {
            return true // Skip negative inputs
        }
        // Test that balance operations maintain non-negative invariant
        return true
    }, nil)
}
```

### Acceptance Tests
End-to-end scenarios:
- Complete user workflows
- Multi-step operations
- Error recovery scenarios

## Getting Started

### Prerequisites
- Go 1.20 or higher
- PostgreSQL 14+ or SQLite 3.35+

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/adhd-game-bot.git
cd adhd-game-bot
```

2. **Install dependencies**
```bash
go mod download
```

3. **Run database migrations**
```bash
go run cmd/migrate/main.go up
```

4. **Run tests**
```bash
go test ./...
```

## Development

### Running Tests
```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# With race detection
go test -race ./...

# Specific package
go test ./internal/usecase/...
```

### Code Organization
- **Domain Layer**: Pure business logic, no external dependencies
- **Application Layer**: Orchestrates domain objects, implements use cases
- **Infrastructure Layer**: External concerns (database, APIs, etc.)
- **Ports**: Interfaces that define contracts between layers

### Adding New Features
1. Start with domain entities and business rules
2. Define repository interfaces in ports
3. Implement use cases in the application layer
4. Create infrastructure implementations
5. Write comprehensive tests at each layer

## Future Integration Points

This backend system is designed to be integrated with various frontends:
- Web applications
- Mobile apps
- Chat bots (Telegram, Discord, Slack)
- CLI tools
- REST/GraphQL APIs

The clean architecture ensures that adding new delivery mechanisms requires no changes to the core business logic.

## License
MIT License - see LICENSE file for details
