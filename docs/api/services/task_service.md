# TaskService API Contract

## Overview
The TaskService manages all task-related operations including creation, scheduling, and reward calculation.

## Methods

### `CreateTask(task *Task) (*Task, error)`
Creates a new task.

**Parameters:**
- `task` - Task object with required fields

**Returns:**
- `*Task` - Created task with populated ID
- `error` - Possible errors:
  - `ErrInvalidInput`
  - `ErrDatabaseError`

**Validation Rules:**
- Title required (3-100 chars)
- Valid category (daily/weekly/adhoc)
- Base duration > 0
- Valid reward curve

### `StartTimer(userID int64, taskID string) (*Timer, error)`
Starts a timer for a task.

**Parameters:**
- `userID` - ID of user starting timer
- `taskID` - ID of task to time

**Returns:**
- `*Timer` - Created timer instance
- `error` - Possible errors:
  - `ErrTaskNotFound`
  - `ErrTimerAlreadyActive`
  - `ErrInvalidState`

**Preconditions:**
1. Task must exist
2. No active timer for this task/user
3. Task must be active

### `CompleteTimer(timerID string) (*Reward, error)`
Completes an active timer and awards rewards.

**Parameters:**
- `timerID` - ID of timer to complete

**Returns:**
- `*Reward` - Details of awarded amount
- `error` - Possible errors:
  - `ErrTimerNotFound`
  - `ErrTimerNotActive`
  - `ErrInsufficientBalance`

**Postconditions:**
1. Timer marked completed
2. User balance updated
3. Ledger entry created
4. Audit log created

## State Management Rules
- Only one active timer per task/user
- Timers can't be restarted after completion
- Rewards calculated based on elapsed time