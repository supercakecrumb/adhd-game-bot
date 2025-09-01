# Domain Rules and Invariants

## Core Business Rules

### 1. Timer Management
- **Single Active Timer**: A user can have only one active timer per task
- **State Transitions**: 
  - Running → Paused/Completed/Expired/Cancelled
  - No reverse transitions allowed
- **Reward Calculation**:
  - Based on elapsed time and reward tiers
  - Snooze penalties applied according to policy

### 2. Task Completion
- **Streak Tracking**:
  - Reset if completed outside time window
  - Incremented only for on-time completion
- **Prerequisites**: All prerequisite tasks must be completed first

### 3. Currency and Economy
- **Atomic Transactions**: Balance updates and ledger entries must be atomic
- **No Negative Balances**: Rejected at domain level
- **Idempotent Operations**: Duplicate operations have no effect

### 4. Shop System
- **Stock Management**: 
  - Limited items can't go below zero stock
  - Unlimited items ignore stock checks
- **Price Lock**: Item price is fixed at purchase time

## Invariants

### User
1. Balance must equal sum of all ledger entries
2. Timezone must be valid IANA timezone

### Task
1. Schedule must have valid window (start < end)
2. Base duration must be positive
3. Reward curve must have at least one tier

### Timer
1. Started time must be before completion time
2. Total extended duration can't exceed max extension
3. Snooze count can't exceed max snoozes

### Purchase
1. Total cost must equal price × quantity
2. Status transitions are one-way (pending→completed→refunded)