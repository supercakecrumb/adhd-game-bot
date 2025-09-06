package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *entity.Task) error {
	// Convert PartialCredit to JSON if it exists
	var partialCreditJSON *string
	if task.PartialCredit != nil {
		// In a real implementation, we would serialize this to JSON
		// For now, we'll use a placeholder
		partialCreditStr := fmt.Sprintf(`{"currency_id": %d, "amount": "%s"}`, task.PartialCredit.CurrencyID, task.PartialCredit.Amount.String())
		partialCreditJSON = &partialCreditStr
	}

	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO tasks (id, title, description, category, difficulty, schedule_json, 
				base_duration, grace_period, cooldown, reward_curve_json, partial_credit_json, 
				streak_enabled, status, time_zone, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
			task.ID, task.Title, task.Description, task.Category, task.Difficulty,
			task.ScheduleJSON, task.BaseDuration, task.GracePeriod, task.Cooldown,
			task.RewardCurveJSON, partialCreditJSON, task.StreakEnabled, task.Status,
			task.TimeZone, time.Now(), time.Now())
		if err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO tasks (id, title, description, category, difficulty, schedule_json, 
				base_duration, grace_period, cooldown, reward_curve_json, partial_credit_json, 
				streak_enabled, status, time_zone, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
			task.ID, task.Title, task.Description, task.Category, task.Difficulty,
			task.ScheduleJSON, task.BaseDuration, task.GracePeriod, task.Cooldown,
			task.RewardCurveJSON, partialCreditJSON, task.StreakEnabled, task.Status,
			task.TimeZone, time.Now(), time.Now())
		if err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}
	}

	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id string) (*entity.Task, error) {
	var task entity.Task
	var partialCreditJSON *string
	var lastCompletedAt *time.Time
	var createdAt, updatedAt time.Time

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT id, title, description, category, difficulty, schedule_json, 
				base_duration, grace_period, cooldown, reward_curve_json, partial_credit_json, 
				streak_enabled, status, last_completed_at, streak_count, time_zone, created_at, updated_at
			FROM tasks WHERE id = $1`, id)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT id, title, description, category, difficulty, schedule_json, 
				base_duration, grace_period, cooldown, reward_curve_json, partial_credit_json, 
				streak_enabled, status, last_completed_at, streak_count, time_zone, created_at, updated_at
			FROM tasks WHERE id = $1`, id)
	}

	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Category, &task.Difficulty,
		&task.ScheduleJSON, &task.BaseDuration, &task.GracePeriod, &task.Cooldown,
		&task.RewardCurveJSON, &partialCreditJSON, &task.StreakEnabled, &task.Status,
		&lastCompletedAt, &task.StreakCount, &task.TimeZone, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found: %w", ports.ErrTaskNotFound)
		}
		return nil, fmt.Errorf("failed to query task: %w", err)
	}

	// Set the last completed at time
	task.LastCompletedAt = lastCompletedAt

	// In a real implementation, we would deserialize partialCreditJSON to a Reward struct
	// For now, we'll leave it as nil

	return &task, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *entity.Task) error {
	// Convert PartialCredit to JSON if it exists
	var partialCreditJSON *string
	if task.PartialCredit != nil {
		// In a real implementation, we would serialize this to JSON
		// For now, we'll use a placeholder
		partialCreditStr := fmt.Sprintf(`{"currency_id": %d, "amount": "%s"}`, task.PartialCredit.CurrencyID, task.PartialCredit.Amount.String())
		partialCreditJSON = &partialCreditStr
	}

	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			UPDATE tasks SET 
				title = $1, description = $2, category = $3, difficulty = $4, schedule_json = $5,
				base_duration = $6, grace_period = $7, cooldown = $8, reward_curve_json = $9, 
				partial_credit_json = $10, streak_enabled = $11, status = $12, 
				last_completed_at = $13, streak_count = $14, time_zone = $15, updated_at = $16
			WHERE id = $17`,
			task.Title, task.Description, task.Category, task.Difficulty,
			task.ScheduleJSON, task.BaseDuration, task.GracePeriod, task.Cooldown,
			task.RewardCurveJSON, partialCreditJSON, task.StreakEnabled, task.Status,
			task.LastCompletedAt, task.StreakCount, task.TimeZone, time.Now(), task.ID)
		if err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			UPDATE tasks SET 
				title = $1, description = $2, category = $3, difficulty = $4, schedule_json = $5,
				base_duration = $6, grace_period = $7, cooldown = $8, reward_curve_json = $9, 
				partial_credit_json = $10, streak_enabled = $11, status = $12, 
				last_completed_at = $13, streak_count = $14, time_zone = $15, updated_at = $16
			WHERE id = $17`,
			task.Title, task.Description, task.Category, task.Difficulty,
			task.ScheduleJSON, task.BaseDuration, task.GracePeriod, task.Cooldown,
			task.RewardCurveJSON, partialCreditJSON, task.StreakEnabled, task.Status,
			task.LastCompletedAt, task.StreakCount, task.TimeZone, time.Now(), task.ID)
		if err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	var err error
	if tx, ok := GetTx(ctx); ok {
		_, err = tx.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id)
	} else {
		_, err = r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id)
	}
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

func (r *TaskRepository) FindByUser(ctx context.Context, userID int64) ([]*entity.Task, error) {
	// This is a simplified implementation since we don't have a direct user-task relationship in the schema
	// In a real implementation, we would need to join with a user_tasks table or similar
	// For now, we'll return all tasks (which is not correct but serves as a placeholder)

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, description, category, difficulty, schedule_json, 
			base_duration, grace_period, cooldown, reward_curve_json, partial_credit_json, 
			streak_enabled, status, last_completed_at, streak_count, time_zone, created_at, updated_at
		FROM tasks`)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*entity.Task
	for rows.Next() {
		var task entity.Task
		var partialCreditJSON *string
		var lastCompletedAt *time.Time
		var createdAt, updatedAt time.Time

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Category, &task.Difficulty,
			&task.ScheduleJSON, &task.BaseDuration, &task.GracePeriod, &task.Cooldown,
			&task.RewardCurveJSON, &partialCreditJSON, &task.StreakEnabled, &task.Status,
			&lastCompletedAt, &task.StreakCount, &task.TimeZone, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		// Set the last completed at time
		task.LastCompletedAt = lastCompletedAt

		// In a real implementation, we would deserialize partialCreditJSON to a Reward struct
		// For now, we'll leave it as nil

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over task rows: %w", err)
	}

	return tasks, nil
}

// Additional methods from the interface that need implementation
func (r *TaskRepository) FindActiveByUser(ctx context.Context, userID int64) ([]*entity.Task, error) {
	// Simplified implementation
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, description, category, difficulty, schedule_json, 
			base_duration, grace_period, cooldown, reward_curve_json, partial_credit_json, 
			streak_enabled, status, last_completed_at, streak_count, time_zone, created_at, updated_at
		FROM tasks WHERE status = 'active'`)
	if err != nil {
		return nil, fmt.Errorf("failed to query active tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*entity.Task
	for rows.Next() {
		var task entity.Task
		var partialCreditJSON *string
		var lastCompletedAt *time.Time
		var createdAt, updatedAt time.Time

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Category, &task.Difficulty,
			&task.ScheduleJSON, &task.BaseDuration, &task.GracePeriod, &task.Cooldown,
			&task.RewardCurveJSON, &partialCreditJSON, &task.StreakEnabled, &task.Status,
			&lastCompletedAt, &task.StreakCount, &task.TimeZone, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		// Set the last completed at time
		task.LastCompletedAt = lastCompletedAt

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over task rows: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) FindWithTimers(ctx context.Context, userID int64) ([]*entity.Task, error) {
	// Simplified implementation - same as FindByUser for now
	return r.FindByUser(ctx, userID)
}

func (r *TaskRepository) FindWithSchedules(ctx context.Context, userID int64) ([]*entity.Task, error) {
	// Simplified implementation - same as FindByUser for now
	return r.FindByUser(ctx, userID)
}

func (r *TaskRepository) BulkUpdate(ctx context.Context, tasks []*entity.Task) error {
	// Begin a transaction for bulk update
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare statement for update
	stmt, err := tx.PrepareContext(ctx, `
		UPDATE tasks SET 
			title = $1, description = $2, category = $3, difficulty = $4, schedule_json = $5,
			base_duration = $6, grace_period = $7, cooldown = $8, reward_curve_json = $9, 
			partial_credit_json = $10, streak_enabled = $11, status = $12, 
			last_completed_at = $13, streak_count = $14, time_zone = $15, updated_at = $16
		WHERE id = $17`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute updates
	for _, task := range tasks {
		// Convert PartialCredit to JSON if it exists
		var partialCreditJSON *string
		if task.PartialCredit != nil {
			partialCreditStr := fmt.Sprintf(`{"currency_id": %d, "amount": "%s"}`, task.PartialCredit.CurrencyID, task.PartialCredit.Amount.String())
			partialCreditJSON = &partialCreditStr
		}

		_, err := stmt.ExecContext(ctx,
			task.Title, task.Description, task.Category, task.Difficulty,
			task.ScheduleJSON, task.BaseDuration, task.GracePeriod, task.Cooldown,
			task.RewardCurveJSON, partialCreditJSON, task.StreakEnabled, task.Status,
			task.LastCompletedAt, task.StreakCount, task.TimeZone, time.Now(), task.ID)
		if err != nil {
			return fmt.Errorf("failed to update task %s: %w", task.ID, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
