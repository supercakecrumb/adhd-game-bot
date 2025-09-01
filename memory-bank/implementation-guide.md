# ADHD Game Bot - Implementation Guide

## Overview

This guide provides step-by-step instructions for implementing the core system following TDD principles. Each step includes specific files to create and tests to write.

## Phase 1A: Foundation Setup

### Step 1: Initialize Go Module

```bash
go mod init github.com/supercakecrumb/adhd-game-bot
```

### Step 2: Create Directory Structure

```bash
mkdir -p internal/{domain/{entity,valueobject,errors},usecase,ports,infra/{sqlite/migrations,clock,uuid},scheduler,config}
mkdir -p test/{unit,integration,acceptance}
mkdir -p cmd/migrate
```

### Step 3: Add Core Dependencies

```go
// go.mod additions
require (
    github.com/google/uuid v1.3.0
    github.com/mattn/go-sqlite3 v1.14.17
    github.com/shopspring/decimal v1.3.1
    github.com/stretchr/testify v1.8.4
)
```

### Step 4: Create Base Test Helpers

```go
// test/helpers/clock.go
type MockClock struct {
    NowFunc func() time.Time
    MonoFunc func() int64
}

// test/helpers/uuid.go
type MockUUIDGen struct {
    NextFunc func() string
}
```

## Phase 1B: Domain Implementation

### Step 1: Define Domain Errors

```go
// internal/domain/errors/errors.go
var (
    ErrTimerAlreadyActive = errors.New("timer already active for this task")
    ErrInsufficientBalance = errors.New("insufficient balance")
    ErrInvalidTimeWindow = errors.New("invalid time window")
    ErrTaskNotFound = errors.New("task not found")
    // ... more domain errors
)
```

### Step 2: Create Value Objects (TDD)

#### Test First:
```go
// internal/domain/valueobject/decimal_test.go
func TestDecimal_Add(t *testing.T) {
    a := NewDecimal("10.50")
    b := NewDecimal("5.25")
    result := a.Add(b)
    assert.Equal(t, "15.75", result.String())
}
```

#### Then Implement:
```go
// internal/domain/valueobject/decimal.go
type Decimal struct {
    value decimal.Decimal
}
```

### Step 3: Create Entities (TDD)

#### Test First:
```go
// internal/domain/entity/timer_test.go
func TestTimer_CalculateRewardTier(t *testing.T) {
    // Test tier boundary calculations
}
```

#### Then Implement:
```go
// internal/domain/entity/timer.go
type Timer struct {
    ID string
    TaskID string
    // ... fields
}
```

## Phase 1C: Port Interfaces

### Step 1: Define Repository Interfaces

```go
// internal/ports/repositories.go
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id int64) (*entity.User, error)
    UpdateBalance(ctx context.Context, userID int64, currency string, delta decimal.Decimal) error
}

type TaskRepository interface {
    Create(ctx context.Context, task *entity.Task) error
    FindByID(ctx context.Context, id string) (*entity.Task, error)
    FindActiveByCategory(ctx context.Context, category string) ([]*entity.Task, error)
}

// ... other repositories
```

### Step 2: Define Infrastructure Interfaces

```go
// internal/ports/clock.go
type Clock interface {
    Now() time.Time
    NowMonotonic() int64 // nanoseconds
}

// internal/ports/uuid.go
type UUIDGenerator interface {
    New() string
}

// internal/ports/transaction.go
type TransactionManager interface {
    WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
```

## Phase 1D: Use Case Services

### Step 1: Timer Service (TDD)

#### Test First:
```go
// internal/usecase/timer_service_test.go
func TestTimerService_StartTimer(t *testing.T) {
    // Setup mocks
    taskRepo := &MockTaskRepo{}
    timerRepo := &MockTimerRepo{}
    clock := &MockClock{}
    uuidGen := &MockUUIDGen{}
    
    service := NewTimerService(taskRepo, timerRepo, clock, uuidGen)
    
    // Test single active timer enforcement
    // Test tier deadline calculation
    // Test state transitions
}
```

#### Then Implement:
```go
// internal/usecase/timer_service.go
type TimerService struct {
    taskRepo  ports.TaskRepository
    timerRepo ports.TimerRepository
    clock     ports.Clock
    uuidGen   ports.UUIDGenerator
}

func (s *TimerService) StartTimer(ctx context.Context, taskID string, userID int64) (*entity.Timer, error) {
    // Implementation
}
```

## Phase 1E: SQLite Implementation

### Step 1: Migration System

```go
// cmd/migrate/main.go
func main() {
    db, err := sql.Open("sqlite3", "file:adhd_bot.db?_foreign_keys=on&_journal_mode=WAL")
    // Run migrations from internal/infra/sqlite/migrations/
}
```

### Step 2: Repository Implementations (TDD)

#### Integration Test First:
```go
// internal/infra/sqlite/user_repo_test.go
func TestUserRepo_CreateAndFind(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewUserRepository(db)
    
    user := &entity.User{
        DisplayName: "Test User",
        TimeZone: "UTC",
    }
    
    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    
    found, err := repo.FindByID(context.Background(), user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.DisplayName, found.DisplayName)
}
```

#### Then Implement:
```go
// internal/infra/sqlite/user_repo.go
type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
    // Implementation with proper SQL
}
```

## Phase 1F: Scheduler Implementation

### Step 1: Recurrence Calculator (TDD)

#### Test First:
```go
// internal/scheduler/recurrence_test.go
func TestDailyRecurrence_NextOccurrence(t *testing.T) {
    // Test with different timezones
    // Test DST transitions
    // Test window boundaries
}
```

#### Then Implement:
```go
// internal/scheduler/recurrence.go
type RecurrenceCalculator interface {
    NextOccurrence(from time.Time, tz *time.Location) (*Occurrence, error)
}
```

## Testing Strategy by Phase

### Unit Tests (Domain Logic)
1. Decimal arithmetic precision
2. Reward tier calculations
3. Timer state transitions
4. Schedule window validation
5. Streak logic

### Integration Tests (Repository Layer)
1. CRUD operations
2. Transaction atomicity
3. Constraint enforcement
4. Index usage
5. Migration application

### Acceptance Tests (Full Workflows)
1. Complete timer lifecycle
2. Multi-step redemption
3. Concurrent operations
4. Crash recovery
5. Schedule materialization

## Common Patterns

### Repository Pattern
```go
func (r *TaskRepository) FindByID(ctx context.Context, id string) (*entity.Task, error) {
    query := `SELECT id, title, description, ... FROM tasks WHERE id = ?`
    
    var task entity.Task
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &task.ID,
        &task.Title,
        &task.Description,
        // ...
    )
    
    if err == sql.ErrNoRows {
        return nil, domain.ErrTaskNotFound
    }
    
    return &task, err
}
```

### Transaction Pattern
```go
func (s *RewardService) AwardCompletion(ctx context.Context, timerID string) error {
    return s.txManager.WithTx(ctx, func(ctx context.Context) error {
        // All operations in this function run in same transaction
        
        // 1. Create ledger entry
        // 2. Update user balance
        // 3. Update timer state
        // 4. Create audit log
        
        return nil // Commits on success, rolls back on error
    })
}
```

### Idempotency Pattern
```go
func (r *LedgerRepository) Create(ctx context.Context, entry *entity.LedgerEntry) error {
    query := `
        INSERT INTO ledger (id, user_id, currency_code, delta, reason, ref_id)
        VALUES (?, ?, ?, ?, ?, ?)
        ON CONFLICT (reason, ref_id) DO NOTHING
    `
    
    result, err := r.db.ExecContext(ctx, query, /* values */)
    if err != nil {
        return err
    }
    
    rows, _ := result.RowsAffected()
    if rows == 0 {
        // Entry already exists, fetch and return it
        return r.findByRefID(ctx, entry.RefID)
    }
    
    return nil
}
```

## Debugging Tips

1. **Use SQLite CLI for debugging**:
   ```bash
   sqlite3 adhd_bot.db
   .mode column
   .headers on
   SELECT * FROM timers WHERE state = 'running';
   ```

2. **Enable query logging**:
   ```go
   db.QueryRowContext(ctx, query, args...) // Log query and args
   ```

3. **Test transactions in isolation**:
   ```go
   // Use savepoints for nested testing
   tx.Exec("SAVEPOINT test")
   // ... test operations
   tx.Exec("ROLLBACK TO test")
   ```

## Performance Optimization

1. **Batch Operations**:
   ```go
   // Instead of multiple queries
   stmt, _ := db.Prepare("INSERT INTO ledger ...")
   for _, entry := range entries {
       stmt.Exec(entry.Values()...)
   }
   ```

2. **Connection Pooling**:
   ```go
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(5)
   db.SetConnMaxLifetime(5 * time.Minute)
   ```

3. **Prepared Statements**:
   ```go
   // Cache frequently used queries
   type Repository struct {
       db *sql.DB
       stmts map[string]*sql.Stmt
   }
   ```

## Checklist for Each Component

- [ ] Write failing test
- [ ] Implement minimal code to pass
- [ ] Refactor for clarity
- [ ] Add edge case tests
- [ ] Document public APIs
- [ ] Run with race detector
- [ ] Check test coverage
- [ ] Update Memory Bank

## Next Steps After Core Completion

1. Performance benchmarks
2. Stress testing
3. API documentation
4. Deployment scripts
5. Monitoring setup
6. Phase 2 planning (UI integration)