package builders_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supercakecrumb/adhd-game-bot/test/fixtures/builders"
)

func TestTaskBuilder(t *testing.T) {
	t.Run("WithDefaults creates valid task", func(t *testing.T) {
		task := builders.NewTaskBuilder().
			WithDefaults().
			Build()

		assert.Equal(t, "task-1", task.ID)
		assert.Equal(t, int64(100), task.ChatID)
		assert.Equal(t, "Sample Task", task.Title)
		assert.Equal(t, "Complete this sample task", task.Description)
		assert.Equal(t, "daily", task.Category)
		assert.Equal(t, "medium", task.Difficulty)
		assert.Equal(t, 30*60, task.BaseDuration)
		assert.Equal(t, "active", task.Status)
	})

	t.Run("Can override defaults", func(t *testing.T) {
		now := time.Now()
		task := builders.NewTaskBuilder().
			WithDefaults().
			WithID("task-2").
			WithChatID(200).
			WithTitle("Important Task").
			WithDescription("Must complete today").
			WithCategory("weekly").
			WithDifficulty("hard").
			WithBaseDuration(60 * 60). // 1 hour
			WithLastCompletedAt(now).
			WithTimeZone("America/New_York").
			WithStatus("inactive").
			Build()

		assert.Equal(t, "task-2", task.ID)
		assert.Equal(t, int64(200), task.ChatID)
		assert.Equal(t, "Important Task", task.Title)
		assert.Equal(t, "Must complete today", task.Description)
		assert.Equal(t, "weekly", task.Category)
		assert.Equal(t, "hard", task.Difficulty)
		assert.Equal(t, 60*60, task.BaseDuration)
		assert.Equal(t, "inactive", task.Status)
		assert.Equal(t, now, *task.LastCompletedAt)
		assert.Equal(t, "America/New_York", task.TimeZone)
	})
}
