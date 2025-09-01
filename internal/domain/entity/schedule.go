package entity

import (
	"time"
)

// Schedule defines when and how often a task should recur
type Schedule struct {
	ID           string
	TaskID       string     // Associated task
	Type         string     // daily, weekly, monthly
	StartDate    time.Time  // When scheduling should begin
	EndDate      *time.Time // Optional end date
	TimeOfDay    string     // HH:MM format
	DaysOfWeek   []int      // For weekly schedules (0-6, Sunday=0)
	DaysOfMonth  []int      // For monthly schedules (1-31)
	SkipIfMissed bool       // Whether to skip if previous occurrence wasn't completed
	Timezone     string     // IANA timezone ID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ScheduleOccurrence tracks completion of scheduled tasks
type ScheduleOccurrence struct {
	ID           string
	ScheduleID   string
	ScheduledFor time.Time // When this was supposed to occur
	CompletedAt  *time.Time
	Skipped      bool
	Missed       bool
}
