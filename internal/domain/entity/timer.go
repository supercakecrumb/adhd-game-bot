package entity

import (
	"time"
)

// Timer represents an active countdown or stopwatch for a task
type Timer struct {
	ID             string
	TaskID         string     // Reference to associated task
	UserID         int64      // User who started the timer
	Type           string     // "countdown" or "stopwatch"
	StartTime      time.Time  // When the timer was started
	Duration       int        // Duration in seconds (for countdown)
	CurrentValue   int        // Current value in seconds
	LastTick       *time.Time // Last time the timer was updated
	Status         string     // "running", "paused", "completed"
	Timezone       string     // IANA timezone for scheduling
	NotificationID *string    // ID of scheduled notification
}

// TimerEvent represents a state change in a timer
type TimerEvent struct {
	ID        string
	TimerID   string
	EventType string // "start", "pause", "resume", "complete", "cancel"
	Timestamp time.Time
	Metadata  map[string]interface{}
}

// RewardTier defines achievement levels for task completion
type RewardTier struct {
	ID          string
	Name        string
	StreakCount int     // Required streak count
	Reward      Reward  // Base reward
	Multiplier  Decimal // Optional multiplier
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
