package builders

import (
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
)

type TaskBuilder struct {
	BaseBuilder[entity.Task]
}

func NewTaskBuilder() *TaskBuilder {
	return &TaskBuilder{
		BaseBuilder: BaseBuilder[entity.Task]{
			Builder: *NewBuilder[entity.Task](),
		},
	}
}

func (b *TaskBuilder) WithDefaults() *TaskBuilder {
	return b.
		WithID("task-1").
		WithChatID(100).
		WithTitle("Sample Task").
		WithDescription("Complete this sample task").
		WithCategory("daily").
		WithDifficulty("medium").
		WithBaseDuration(30 * 60). // 30 minutes
		WithStatus("active")
}

func (b *TaskBuilder) WithID(id string) *TaskBuilder {
	b.With(func(t *entity.Task) {
		t.ID = id
	})
	return b
}

func (b *TaskBuilder) WithChatID(chatID int64) *TaskBuilder {
	b.With(func(t *entity.Task) {
		t.ChatID = chatID
	})
	return b
}

func (b *TaskBuilder) WithTitle(title string) *TaskBuilder {
	b.With(func(t *entity.Task) {
		t.Title = title
	})
	return b
}

func (b *TaskBuilder) WithDescription(description string) *TaskBuilder {
	b.With(func(t *entity.Task) {
		t.Description = description
	})
	return b
}

func (b *TaskBuilder) WithCategory(category string) *TaskBuilder {
	b.With(func(t *entity.Task) {
		t.Category = category
	})
	return b
}

func (b *TaskBuilder) WithDifficulty(difficulty string) *TaskBuilder {
	b.With(func(t *entity.Task) {
		t.Difficulty = difficulty
	})
	return b
}

func (b *TaskBuilder) WithBaseDuration(seconds int) *TaskBuilder {
	b.With(func(t *entity.Task) {
		t.BaseDuration = seconds
	})
	return b
}

func (b *TaskBuilder) WithStatus(status string) *TaskBuilder {
	b.With(func(t *entity.Task) {
		t.Status = status
	})
	return b
}

func (b *TaskBuilder) WithLastCompletedAt(time time.Time) *TaskBuilder {
	b.With(func(t *entity.Task) {
		t.LastCompletedAt = &time
	})
	return b
}

func (b *TaskBuilder) WithTimeZone(timeZone string) *TaskBuilder {
	b.With(func(t *entity.Task) {
		t.TimeZone = timeZone
	})
	return b
}
