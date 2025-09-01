package builders_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supercakecrumb/adhd-game-bot/test/fixtures/builders"
)

func TestTimerBuilder(t *testing.T) {
	t.Run("WithDefaults creates valid timer", func(t *testing.T) {
		timer := builders.NewTimerBuilder().
			WithDefaults().
			Build()

		assert.Equal(t, "timer-1", timer.ID)
		assert.Equal(t, "task-1", timer.TaskID)
		assert.Equal(t, int64(1), timer.UserID)
		assert.NotZero(t, timer.StartedAt)
		assert.Equal(t, 1800, timer.InitialDuration)
		assert.Equal(t, "running", timer.State)
	})

	t.Run("Can override defaults", func(t *testing.T) {
		now := time.Now()
		timer := builders.NewTimerBuilder().
			WithDefaults().
			WithID("timer-2").
			WithTaskID("task-2").
			WithUserID(2).
			WithStartedAt(now).
			WithInitialDuration(3600). // 1 hour
			WithSnoozeCount(3).
			WithState("paused").
			WithTotalExtended(900). // 15 minutes
			Build()

		assert.Equal(t, "timer-2", timer.ID)
		assert.Equal(t, "task-2", timer.TaskID)
		assert.Equal(t, int64(2), timer.UserID)
		assert.Equal(t, now, timer.StartedAt)
		assert.Equal(t, 3600, timer.InitialDuration)
		assert.Equal(t, 3, timer.SnoozeCount)
		assert.Equal(t, "paused", timer.State)
		assert.Equal(t, 900, timer.TotalExtended)
	})

	t.Run("Can set last tick time", func(t *testing.T) {
		now := time.Now()
		timer := builders.NewTimerBuilder().
			WithLastTickAt(now).
			Build()

		assert.Equal(t, now, *timer.LastTickAt)
	})
}
