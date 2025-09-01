# Product Context

This file provides a high-level overviewbased on project brief:

# ADHD Game Bot - Core System

## Overview
A gamified routine/support tool for ADHD management with tasks, timers, rewards, and streaks.

## Phase 1 Scope (Core Only)
- Domain modeling with clean architecture
- SQLite persistence layer
- Test-driven development
- NO UI (no Telegram, HTTP, CLI)
- Focus on domain structs, interfaces, methods, and tests

## Key Features
- Tasks (manual and scheduled)
- Timers with stretchy reward tiers
- Multi-currency economy
- Store/redemptions system
- Ledger/audit trail
- Streak tracking
- Idempotency guarantees

## Architecture
- Clean layering: domain, ports, adapters
- Deterministic time math
- Crash-safe timers with resume semantics
- SQLite with WAL mode
- Interfaces for Clock, UUID, TxManager

## Core Entities
- User (with balances and preferences)
- Currency (MM, BS with decimal support)
- Task (with schedule, duration, reward curve)
- TimerInstance (state machine)
- LedgerEntry (atomic financial records)
- Redemption (store purchases)
- AuditLog (all actions)
- ScheduleOccurrence (materialized runs)

## Key Behaviors
- One active timer per task/user
- Tier-based rewards with snooze policies
- Atomic ledger updates
- Idempotent operations
- Timezone-aware scheduling
- Streak management

...

*