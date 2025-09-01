# ADHD Game Bot - Test Plan

## Test-Driven Development Strategy

This document outlines the comprehensive test plan following strict TDD principles. Each test will be written before the implementation.

## Test Categories

### 1. Domain Unit Tests

#### Reward Tier Selection Tests
```go
// TestTierSelection_Boundaries
// - 9:59 gets tier1 (15 MM)
// - 10:00 gets tier2 (10 MM)
// - 19:59 gets tier2 (10 MM)
// - 20:00 gets tier3 (5 MM)
// - 30:00 gets tier3 (5 MM)
// - 30:01 gets partial credit (3 MM)

// TestTierSelection_NoPartialCredit
// - When partial credit is nil, >30m gets 0 reward

// TestTierSelection_EmptyTiers
// - Should return error for empty tier list
```

#### Snooze Policy Tests
```go
// TestSnooze_DemotionPolicy
// - First snooze: tier1 → tier2
// - Second snooze: tier2 → tier3
// - Third snooze: tier3 → partial (if exists)
// - Fourth snooze: no further demotion

// TestSnooze_DecrementPolicy
// - Each snooze reduces reward by fixed amount
// - Cannot go below zero
// - Maintains current tier deadline

// TestSnooze_MaxExtension
// - Total extension cannot exceed configured maximum
// - Attempting to exceed returns error
```

#### Timer State Machine Tests
```go
// TestTimerState_ValidTransitions
// - Running → Paused → Running
// - Running → Completed
// - Running → Expired
// - Running → Cancelled

// TestTimerState_InvalidTransitions
// - Completed → Running (error)
// - Expired → Paused (error)
// - Cancelled → Completed (error)

// TestSingleActiveTimer_PerTaskUser
// - Starting timer when one exists returns error
// - Can start after previous completed/expired/cancelled
```

#### Balance and Decimal Math Tests
```go
// TestDecimal_Addition
// - 1.23 + 4.56 = 5.79
// - Handle different decimal places

// TestDecimal_Subtraction
// - 10.00 - 3.33 = 6.67
// - Insufficient balance detection

// TestDecimal_Multiplication
// - For currency conversion
// - Proper rounding rules

// TestBalance_ConcurrentUpdates
// - Race condition prevention
// - Atomic balance modifications
```

#### Streak Management Tests
```go
// TestStreak_IncrementWithinWindow
// - Complete task at 22:00, window 21:00-23:00 → increment
// - Streak goes from 5 to 6

// TestStreak_ResetOutsideWindow
// - Complete task at 20:00, window 21:00-23:00 → reset to 0
// - Previous streak of 10 → 0

// TestStreak_ConsecutiveDays
// - Track daily completion patterns
// - Handle timezone changes
```

### 2. Repository Interface Tests (with In-Memory Fakes)

#### User Repository Tests
```go
// TestUserRepo_CreateAndFind
// TestUserRepo_UpdateBalances_Atomic
// TestUserRepo_ConcurrentBalanceUpdates
```

#### Task Repository Tests
```go
// TestTaskRepo_CRUD
// TestTaskRepo_FindByCategory
// TestTaskRepo_UpdateStreakCount
// TestTaskRepo_ValidateScheduleWindows
```

#### Timer Repository Tests
```go
// TestTimerRepo_EnforceSingleActive
// TestTimerRepo_FindActiveByUser
// TestTimerRepo_UpdateState_Idempotent
```

#### Ledger Repository Tests
```go
// TestLedgerRepo_CreateEntry_Idempotent
// - Same RefID returns existing entry
// - No duplicate balance changes

// TestLedgerRepo_QueryByDateRange
// TestLedgerRepo_SumByUserCurrency
```

### 3. Use Case Integration Tests

#### Timer Lifecycle Tests
```go
// TestStartTimer_Success
// - Creates timer instance
// - Calculates tier deadlines
// - Enforces single active rule

// TestCompleteTimer_AwardCalculation
// - Correct tier selection
// - Ledger entry creation
// - Balance update
// - Audit log creation
// - All in single transaction

// TestExpireTimer_PartialCredit
// - Awards partial credit if configured
// - Updates timer state
// - Creates appropriate logs
```

#### Redemption Flow Tests
```go
// TestRedemption_FullFlow
// - Check balance sufficiency
// - Create pending redemption
// - Deduct balance atomically
// - Update redemption status
// - Audit trail complete

// TestRedemption_InsufficientBalance
// - Validation error
// - No state changes
// - Clear error message
```

### 4. SQLite Integration Tests

#### Schema and Migration Tests
```go
// TestMigrations_ApplyCleanly
// TestMigrations_Idempotent
// TestForeignKeys_Enforced
// TestIndices_Created
```

#### Transaction Tests
```go
// TestTransaction_Commit
// - All changes persisted

// TestTransaction_Rollback
// - No changes persisted
// - Clean state restoration

// TestTransaction_Nested
// - Savepoint support
```

#### Concurrency Tests
```go
// TestConcurrent_TimerCreation
// - Multiple goroutines
// - Only one succeeds

// TestConcurrent_BalanceUpdates
// - No lost updates
// - Correct final balance
```

### 5. Scheduler Tests

#### Recurrence Calculation Tests
```go
// TestDaily_NextOccurrence
// - Simple daily at 14:00
// - Handle today vs tomorrow

// TestWeekly_SpecificDays
// - Mon, Wed, Fri at 09:00
// - Skip to next valid day

// TestCustom_ComplexPattern
// - Every 3 days
// - Business days only
```

#### Timezone and DST Tests
```go
// TestScheduler_DSTSpringForward_Helsinki
// - Task at 02:30 on DST day
// - Clock jumps 02:00 → 03:00
// - Occurrence handled correctly

// TestScheduler_DSTFallBack_Helsinki
// - Task at 02:30 on DST day
// - Clock repeats 02:00-03:00
// - No duplicate occurrences

// TestScheduler_TimezoneChange
// - User moves Helsinki → Tokyo
// - Recalculate all schedules
```

### 6. Property-Based Tests

```go
// TestProperty_TierRewards_Deterministic
// - For any valid duration
// - Same input → same reward
// - Monotonic: longer time → less reward

// TestProperty_IdempotentOperations
// - Complete/Award/Fulfill
// - Multiple calls → single effect

// TestProperty_BalanceConsistency
// - Sum of ledger entries = current balance
// - No matter the operation sequence
```

### 7. Acceptance Tests

#### Complete User Journey Tests
```go
// TestJourney_MorningRoutine
// 1. Create "Brush Teeth" task (20m stretchy)
// 2. Start timer at 07:00
// 3. Complete at 07:08 (tier 1)
// 4. Receive 15 MM
// 5. Check balance = 15 MM
// 6. Verify ledger entry
// 7. Verify audit log

// TestJourney_SnoozeAndComplete
// 1. Start "Exercise" timer
// 2. Snooze at 15m (+10m)
// 3. Complete at 22m
// 4. Receive tier 2 reward (demotion policy)
// 5. Verify all records

// TestJourney_StoreRedemption
// 1. Accumulate 20 MM from tasks
// 2. Browse store items
// 3. Redeem "Massage 20 min"
// 4. Admin fulfills
// 5. Balance = 0 MM
// 6. Redemption marked fulfilled
```

#### Error Recovery Tests
```go
// TestRecovery_CrashDuringTimer
// 1. Start timer
// 2. Simulate crash
// 3. Restart system
// 4. Timer still active
// 5. Can complete normally
// 6. No double rewards

// TestRecovery_PartialTransaction
// 1. Start multi-step operation
// 2. Fail mid-transaction
// 3. Verify rollback complete
// 4. System in consistent state
```

## Test Data Fixtures

### Base Currencies
```go
var (
    MM = Currency{Code: "MM", Name: "Motivation Minutes", Decimals: 2}
    BS = Currency{Code: "BS", Name: "Bonus Stars", Decimals: 2}
)
```

### Sample Tasks
```go
var (
    BrushTeeth = Task{
        Title: "Brush Teeth",
        Category: Daily,
        BaseDuration: 20 * time.Minute,
        RewardCurve: RewardCurve{
            Tiers: []RewardTier{
                {DeadlineOffset: 10*time.Minute, Reward: map[string]Decimal{"MM": "15.00"}},
                {DeadlineOffset: 20*time.Minute, Reward: map[string]Decimal{"MM": "10.00"}},
                {DeadlineOffset: 30*time.Minute, Reward: map[string]Decimal{"MM": "5.00"}},
            },
            PartialCredit: &Reward{Currency: "MM", Amount: "3.00"},
        },
    }
    
    TakeMeds = Task{
        Title: "Take Medication",
        Category: Daily,
        Schedule: DailySchedule{
            WindowStart: time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC),
            WindowEnd:   time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC),
        },
        BaseDuration: 10 * time.Minute,
        StreakEnabled: true,
    }
)
```

### Test Users
```go
var (
    TestUser = User{
        ID: 1,
        DisplayName: "Test User",
        TimeZone: "Europe/Helsinki",
        Balances: map[string]Decimal{
            "MM": "0.00",
            "BS": "0.00",
        },
    }
    
    AdminUser = User{
        ID: 2,
        Role: Admin,
        DisplayName: "Admin",
        TimeZone: "UTC",
    }
)
```

## Test Execution Order

1. **Domain Unit Tests** - Pure logic, no dependencies
2. **Repository Interface Tests** - With in-memory fakes
3. **Use Case Tests** - Business logic with fakes
4. **SQLite Integration Tests** - Real database
5. **Scheduler Tests** - Time-sensitive logic
6. **Property Tests** - Invariant verification
7. **Acceptance Tests** - Full system validation

## Coverage Goals

- Domain Layer: 100% coverage
- Use Cases: 95%+ coverage
- Repositories: 90%+ coverage
- Critical Paths: 100% coverage
- Edge Cases: Comprehensive property tests

## Test Environment Setup

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test suite
go test ./internal/domain/...

# Run with race detection
go test -race ./...

# Run property tests (longer)
go test -tags=property ./...
```

## Continuous Testing

- Pre-commit: Unit tests only (fast)
- PR validation: All tests except property
- Nightly: Full suite including property tests
- Release: Full suite + performance benchmarks