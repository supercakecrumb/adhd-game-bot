package inmemory

import (
	"context"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type InMemoryScheduler struct {
	scheduledTasks map[string]*entity.Task
}

func NewInMemoryScheduler() *InMemoryScheduler {
	return &InMemoryScheduler{
		scheduledTasks: make(map[string]*entity.Task),
	}
}

func (s *InMemoryScheduler) ScheduleRecurringTask(ctx context.Context, task *entity.Task) error {
	s.scheduledTasks[task.ID] = task
	return nil
}

func (s *InMemoryScheduler) CancelScheduledTask(ctx context.Context, taskID string) error {
	delete(s.scheduledTasks, taskID)
	return nil
}

func (s *InMemoryScheduler) GetNextOccurrence(ctx context.Context, taskID string) (time.Time, error) {
	if _, exists := s.scheduledTasks[taskID]; exists {
		// Simple implementation - assumes daily recurrence
		now := time.Now()
		next := now.Add(24 * time.Hour)
		return next, nil
	}
	return time.Time{}, ports.ErrTaskNotFound
}
