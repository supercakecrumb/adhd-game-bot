package builders

import (
	"time"
)

// Timer represents a timer instance for testing purposes
type Timer struct {
	ID              string
	TaskID          string
	UserID          int64
	StartedAt       time.Time
	InitialDuration int // seconds
	State           string
	LastTickAt      *time.Time
	SnoozeCount     int
	TotalExtended   int // seconds
}

type TimerBuilder struct {
	BaseBuilder[Timer]
}

func NewTimerBuilder() *TimerBuilder {
	return &TimerBuilder{
		BaseBuilder: BaseBuilder[Timer]{
			Builder: *NewBuilder[Timer](),
		},
	}
}

func (b *TimerBuilder) WithDefaults() *TimerBuilder {
	now := time.Now()
	return b.
		WithID("timer-1").
		WithTaskID("task-1").
		WithUserID(1).
		WithStartedAt(now).
		WithInitialDuration(1800). // 30 minutes
		WithState("running")
}

func (b *TimerBuilder) WithID(id string) *TimerBuilder {
	b.With(func(t *Timer) {
		t.ID = id
	})
	return b
}

func (b *TimerBuilder) WithTaskID(taskID string) *TimerBuilder {
	b.With(func(t *Timer) {
		t.TaskID = taskID
	})
	return b
}

func (b *TimerBuilder) WithUserID(userID int64) *TimerBuilder {
	b.With(func(t *Timer) {
		t.UserID = userID
	})
	return b
}

func (b *TimerBuilder) WithStartedAt(time time.Time) *TimerBuilder {
	b.With(func(t *Timer) {
		t.StartedAt = time
	})
	return b
}

func (b *TimerBuilder) WithInitialDuration(seconds int) *TimerBuilder {
	b.With(func(t *Timer) {
		t.InitialDuration = seconds
	})
	return b
}

func (b *TimerBuilder) WithState(state string) *TimerBuilder {
	b.With(func(t *Timer) {
		t.State = state
	})
	return b
}

func (b *TimerBuilder) WithLastTickAt(time time.Time) *TimerBuilder {
	b.With(func(t *Timer) {
		t.LastTickAt = &time
	})
	return b
}

func (b *TimerBuilder) WithSnoozeCount(count int) *TimerBuilder {
	b.With(func(t *Timer) {
		t.SnoozeCount = count
	})
	return b
}

func (b *TimerBuilder) WithTotalExtended(seconds int) *TimerBuilder {
	b.With(func(t *Timer) {
		t.TotalExtended = seconds
	})
	return b
}
