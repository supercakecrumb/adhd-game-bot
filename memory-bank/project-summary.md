# ADHD Game Bot - Project Summary & Next Steps

## Project Overview

We're building a gamified routine/support tool for ADHD management. Phase 1 focuses on the core domain logic with no UI - just a rock-solid foundation with comprehensive tests.

## Architecture Summary

### Clean Architecture Layers
1. **Domain Layer**: Pure business logic (entities, value objects, domain errors)
2. **Application Layer**: Use cases and port interfaces
3. **Infrastructure Layer**: SQLite persistence, system clock, UUID generation
4. **Scheduler**: Timezone-aware task recurrence logic

### Key Design Principles
- **Test-Driven Development**: Write tests first, then implementation
- **Dependency Inversion**: Domain doesn't depend on infrastructure
- **Single Responsibility**: Each component has one reason to change
- **Idempotency**: Operations can be safely retried
- **Atomicity**: Financial operations in transactions

## Core Features

### 1. Task Management
- Categories: Daily, Weekly, Adhoc
- Schedules with timezone-aware windows
- Prerequisites and tags
- Streak tracking

### 2. Timer System
- Stretchy timers with reward tiers
- Snooze policies (demotion or decrement)
- Single active timer per task/user
- Crash-safe with resume capability

### 3. Multi-Currency Economy
- Decimal precision (no floating point)
- Atomic balance updates
- Complete audit trail
- Idempotent transactions

### 4. Store & Redemptions
- Item catalog with multi-currency pricing
- Two-phase redemption (create → fulfill)
- Admin approval workflow

### 5. Scheduling Engine
- Recurrence patterns (daily, weekly, custom)
- DST-aware calculations
- Materialized occurrences
- Window-based execution

## Technical Stack

- **Language**: Go
- **Database**: SQLite with WAL mode
- **Testing**: testify, property-based tests
- **Libraries**: 
  - shopspring/decimal (precision math)
  - google/uuid (ID generation)
  - mattn/go-sqlite3 (database driver)

## Implementation Roadmap

### Week 1: Foundation
- [x] Architecture design
- [x] Database schema
- [x] Test plan
- [ ] Project setup
- [ ] Domain entities
- [ ] Value objects
- [ ] Port interfaces

### Week 2: Core Logic
- [ ] Timer state machine
- [ ] Reward calculations
- [ ] Balance management
- [ ] Basic use cases
- [ ] In-memory test doubles

### Week 3: Persistence
- [ ] SQLite repositories
- [ ] Migration system
- [ ] Transaction manager
- [ ] Integration tests

### Week 4: Advanced Features
- [ ] Scheduler implementation
- [ ] Timezone handling
- [ ] Idempotency mechanisms
- [ ] Property-based tests

### Week 5: Polish & Documentation
- [ ] Acceptance test suite
- [ ] Performance optimization
- [ ] API documentation
- [ ] Deployment preparation

## Key Files Created

1. **memory-bank/architecture-plan.md**: System design with diagrams
2. **memory-bank/test-plan.md**: Comprehensive test strategy
3. **memory-bank/database-schema.md**: Complete SQLite schema
4. **memory-bank/implementation-guide.md**: Step-by-step coding guide

## Critical Path Items

1. **Decimal Arithmetic**: Must be precise for currency
2. **Timer Monotonic Time**: Wall clock can jump
3. **Transaction Atomicity**: No partial updates
4. **Idempotency**: Prevent double rewards
5. **Timezone Correctness**: DST transitions

## Success Criteria

- [ ] All tests pass (unit, integration, acceptance)
- [ ] 95%+ code coverage on domain layer
- [ ] No race conditions (tested with -race)
- [ ] Idempotent operations verified
- [ ] DST transitions handled correctly
- [ ] Crash recovery tested
- [ ] Performance benchmarks met

## Questions to Clarify

1. **Snooze Policy**: Should we implement both demotion and decrement, or choose one?
2. **Partial Credit**: Is this a fixed amount or percentage of lowest tier?
3. **Store Items**: Should prices be dynamic or fixed at redemption time?
4. **Audit Retention**: How long to keep audit logs?
5. **Backup Strategy**: Frequency and retention policy?

## Next Immediate Steps

1. **Switch to Code Mode** to begin implementation
2. **Start with project setup** (go mod init, directory structure)
3. **Create domain errors** as the foundation
4. **Write first failing test** for Decimal value object
5. **Follow TDD cycle** throughout

## Memory Bank Status

The Memory Bank has been initialized with:
- Product context and project overview
- Architectural decisions log
- Active development context
- Progress tracking

Remember to update the Memory Bank as you make progress and encounter new decisions or patterns.

## Ready to Start?

The planning phase is complete. We have:
- ✅ Clear architecture with clean separation of concerns
- ✅ Comprehensive test plan following TDD
- ✅ Complete database schema with migrations
- ✅ Step-by-step implementation guide
- ✅ All design decisions documented

The next step is to switch to Code mode and begin the implementation, starting with the project setup and first failing tests.