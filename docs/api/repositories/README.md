# Repository Interfaces

This document describes the persistence layer contracts for all repositories.

## Common Repository Pattern

All repositories follow these principles:
- **Create/Read/Update** operations
- **Idempotency** for critical operations
- **Transaction** support via TxManager
- **Error types** for domain-specific failures

## UserRepository

```go
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id int64) (*entity.User, error)
    UpdateBalance(ctx context.Context, userID int64, delta valueobject.Decimal) error
}
```

**Idempotency Notes:**
- Balance updates use atomic operations
- Concurrent updates are handled via transactions

## TaskRepository

```go
type TaskRepository interface {
    Create(ctx context.Context, task *entity.Task) error
    FindByID(ctx context.Context, id string) (*entity.Task, error)
    FindActiveByUser(ctx context.Context, userID int64) ([]*entity.Task, error)
}
```

**Constraints:**
- Active tasks must have valid schedules
- Task titles are unique per chat

## TimerRepository

```go
type TimerRepository interface {
    Create(ctx context.Context, timer *entity.Timer) error
    FindActiveByUserTask(ctx context.Context, userID int64, taskID string) (*entity.Timer, error)
    UpdateState(ctx context.Context, timerID string, state string) error
}
```

**Invariants:**
- Single active timer per task/user
- State transitions are validated

## ShopItemRepository

```go
type ShopItemRepository interface {
    Create(ctx context.Context, item *entity.ShopItem) error
    FindByCode(ctx context.Context, chatID int64, code string) (*entity.ShopItem, error)
    UpdateStock(ctx context.Context, itemID int64, delta int) error
}
```

**Business Rules:**
- Stock cannot go negative
- Items are unique per chat by code