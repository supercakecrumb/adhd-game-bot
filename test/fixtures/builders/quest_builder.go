package builders

import (
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type QuestBuilder struct {
	BaseBuilder[entity.Quest]
}

func NewQuestBuilder() *QuestBuilder {
	return &QuestBuilder{
		BaseBuilder: BaseBuilder[entity.Quest]{
			Builder: *NewBuilder[entity.Quest](),
		},
	}
}

func (b *QuestBuilder) WithDefaults() *QuestBuilder {
	points := valueobject.NewDecimal("10")

	return b.
		WithID("quest-1").
		WithTitle("Sample Quest").
		WithDescription("Complete this sample quest").
		WithCategory("daily").
		WithDifficulty("medium").
		WithMode("BINARY").
		WithPointsAward(points).
		WithCooldownSec(0).
		WithStreakEnabled(true).
		WithStatus("active")
}

func (b *QuestBuilder) WithID(id string) *QuestBuilder {
	b.With(func(t *entity.Quest) {
		t.ID = id
	})
	return b
}

func (b *QuestBuilder) WithTitle(title string) *QuestBuilder {
	b.With(func(t *entity.Quest) {
		t.Title = title
	})
	return b
}

func (b *QuestBuilder) WithDescription(description string) *QuestBuilder {
	b.With(func(t *entity.Quest) {
		t.Description = description
	})
	return b
}

func (b *QuestBuilder) WithCategory(category string) *QuestBuilder {
	b.With(func(t *entity.Quest) {
		t.Category = category
	})
	return b
}

func (b *QuestBuilder) WithDifficulty(difficulty string) *QuestBuilder {
	b.With(func(t *entity.Quest) {
		t.Difficulty = difficulty
	})
	return b
}

func (b *QuestBuilder) WithMode(mode string) *QuestBuilder {
	b.With(func(t *entity.Quest) {
		t.Mode = mode
	})
	return b
}

func (b *QuestBuilder) WithPointsAward(points valueobject.Decimal) *QuestBuilder {
	b.With(func(t *entity.Quest) {
		t.PointsAward = points
	})
	return b
}

func (b *QuestBuilder) WithCooldownSec(seconds int) *QuestBuilder {
	b.With(func(t *entity.Quest) {
		t.CooldownSec = seconds
	})
	return b
}

func (b *QuestBuilder) WithStreakEnabled(enabled bool) *QuestBuilder {
	b.With(func(t *entity.Quest) {
		t.StreakEnabled = enabled
	})
	return b
}

func (b *QuestBuilder) WithStatus(status string) *QuestBuilder {
	b.With(func(t *entity.Quest) {
		t.Status = status
	})
	return b
}

func (b *QuestBuilder) WithLastCompletedAt(time time.Time) *QuestBuilder {
	b.With(func(t *entity.Quest) {
		t.LastCompletedAt = &time
	})
	return b
}

func (b *QuestBuilder) WithTimeZone(timeZone string) *QuestBuilder {
	b.With(func(t *entity.Quest) {
		t.TimeZone = timeZone
	})
	return b
}
