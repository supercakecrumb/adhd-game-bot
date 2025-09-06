package inmemory

import (
	"context"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type InMemoryScheduler struct {
	scheduledTasks map[string]*entity.Quest
}

func NewInMemoryScheduler() *InMemoryScheduler {
	return &InMemoryScheduler{
		scheduledTasks: make(map[string]*entity.Quest),
	}
}

func (s *InMemoryScheduler) ScheduleRecurringTask(ctx context.Context, quest *entity.Quest) error {
	s.scheduledTasks[quest.ID] = quest
	return nil
}

func (s *InMemoryScheduler) CancelScheduledTask(ctx context.Context, questID string) error {
	delete(s.scheduledTasks, questID)
	return nil
}

func (s *InMemoryScheduler) GetNextOccurrence(ctx context.Context, questID string) (time.Time, error) {
	if _, exists := s.scheduledTasks[questID]; exists {
		// Simple implementation - assumes daily recurrence
		now := time.Now()
		next := now.Add(24 * time.Hour)
		return next, nil
	}
	return time.Time{}, ports.ErrTaskNotFound
}
