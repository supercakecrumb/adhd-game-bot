package builders_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/test/fixtures/builders"
)

func TestQuestBuilder(t *testing.T) {
	t.Run("WithDefaults creates valid quest", func(t *testing.T) {
		quest := builders.NewQuestBuilder().
			WithDefaults().
			Build()

		assert.Equal(t, "quest-1", quest.ID)
		assert.Equal(t, "Sample Quest", quest.Title)
		assert.Equal(t, "Complete this sample quest", quest.Description)
		assert.Equal(t, "daily", quest.Category)
		assert.Equal(t, "medium", quest.Difficulty)
		assert.Equal(t, "BINARY", quest.Mode)
		assert.Equal(t, valueobject.NewDecimal("10"), quest.PointsAward)
		assert.Equal(t, 0, quest.CooldownSec)
		assert.Equal(t, true, quest.StreakEnabled)
		assert.Equal(t, "active", quest.Status)
	})

	t.Run("Can override defaults", func(t *testing.T) {
		now := time.Now()
		points := valueobject.NewDecimal("25")

		quest := builders.NewQuestBuilder().
			WithDefaults().
			WithID("quest-2").
			WithTitle("Important Quest").
			WithDescription("Must complete today").
			WithCategory("weekly").
			WithDifficulty("hard").
			WithMode("PROGRESSIVE").
			WithPointsAward(points).
			WithCooldownSec(3600).
			WithStreakEnabled(false).
			WithLastCompletedAt(now).
			WithTimeZone("America/New_York").
			WithStatus("inactive").
			Build()

		assert.Equal(t, "quest-2", quest.ID)
		assert.Equal(t, "Important Quest", quest.Title)
		assert.Equal(t, "Must complete today", quest.Description)
		assert.Equal(t, "weekly", quest.Category)
		assert.Equal(t, "hard", quest.Difficulty)
		assert.Equal(t, "PROGRESSIVE", quest.Mode)
		assert.Equal(t, points, quest.PointsAward)
		assert.Equal(t, 3600, quest.CooldownSec)
		assert.Equal(t, false, quest.StreakEnabled)
		assert.Equal(t, "inactive", quest.Status)
		assert.Equal(t, now, *quest.LastCompletedAt)
		assert.Equal(t, "America/New_York", quest.TimeZone)
	})
}
