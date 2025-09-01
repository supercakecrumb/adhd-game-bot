# Decision Log

This file records architectural and implementation decisions...

*
[2025-09-01 13:51:48] - ## Architecture Decisions

### Clean Architecture Pattern
- **Decision**: Use Clean Architecture with clear separation between domain, application, and infrastructure layers
- **Rationale**: Ensures testability, maintainability, and independence from external frameworks
- **Trade-offs**: More initial setup complexity, but better long-term maintainability

### SQLite for Persistence
- **Decision**: Use SQLite with WAL mode as the primary database
- **Rationale**: Simple deployment, excellent performance for single-instance apps, built-in transaction support
- **Trade-offs**: Limited concurrent write performance, but sufficient for our use case

### Decimal Storage as Strings
- **Decision**: Store currency amounts as TEXT in database
- **Rationale**: Avoid floating-point precision issues, maintain exact decimal arithmetic
- **Implementation**: Use shopspring/decimal library in Go

### Monotonic Time for Timers
- **Decision**: Track both wall clock and monotonic time for timers
- **Rationale**: Wall clock can jump (NTP, DST), monotonic time ensures accurate duration measurement
- **Implementation**: Store started_at_mono as nanoseconds since arbitrary point

### Single Active Timer Constraint
- **Decision**: Enforce one active timer per task/user at database level
- **Rationale**: Prevents accidental double-starts, simplifies state management
- **Implementation**: Unique partial index on (task_id, user_id) WHERE state = 'running'
