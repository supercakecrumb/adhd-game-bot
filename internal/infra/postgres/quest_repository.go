package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type QuestRepository struct {
	db *sql.DB
}

func NewQuestRepository(db *sql.DB) *QuestRepository {
	return &QuestRepository{db: db}
}

func (r *QuestRepository) Create(ctx context.Context, quest *entity.Quest) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO quests (id, dungeon_id, title, description, category, difficulty, mode, points_award, 
				rate_points_per_min, min_minutes, max_minutes, daily_points_cap, cooldown_sec, streak_enabled, 
				status, last_completed_at, streak_count, time_zone, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`,
			quest.ID, quest.DungeonID, quest.Title, quest.Description, quest.Category, quest.Difficulty, quest.Mode,
			quest.PointsAward.String(), quest.RatePointsPerMin, quest.MinMinutes, quest.MaxMinutes, quest.DailyPointsCap,
			quest.CooldownSec, quest.StreakEnabled, quest.Status, quest.LastCompletedAt, quest.StreakCount,
			quest.TimeZone, quest.CreatedAt, quest.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to create quest: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO quests (id, dungeon_id, title, description, category, difficulty, mode, points_award, 
				rate_points_per_min, min_minutes, max_minutes, daily_points_cap, cooldown_sec, streak_enabled, 
				status, last_completed_at, streak_count, time_zone, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`,
			quest.ID, quest.DungeonID, quest.Title, quest.Description, quest.Category, quest.Difficulty, quest.Mode,
			quest.PointsAward.String(), quest.RatePointsPerMin, quest.MinMinutes, quest.MaxMinutes, quest.DailyPointsCap,
			quest.CooldownSec, quest.StreakEnabled, quest.Status, quest.LastCompletedAt, quest.StreakCount,
			quest.TimeZone, quest.CreatedAt, quest.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to create quest: %w", err)
		}
	}

	return nil
}

func (r *QuestRepository) GetByID(ctx context.Context, questID string) (*entity.Quest, error) {
	var quest entity.Quest
	var pointsAwardStr string
	var lastCompletedAt *time.Time
	var createdAt time.Time
	var updatedAt time.Time

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT id, dungeon_id, title, description, category, difficulty, mode, points_award, 
				rate_points_per_min, min_minutes, max_minutes, daily_points_cap, cooldown_sec, streak_enabled, 
				status, last_completed_at, streak_count, time_zone, created_at, updated_at
			FROM quests WHERE id = $1`, questID)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT id, dungeon_id, title, description, category, difficulty, mode, points_award, 
				rate_points_per_min, min_minutes, max_minutes, daily_points_cap, cooldown_sec, streak_enabled, 
				status, last_completed_at, streak_count, time_zone, created_at, updated_at
			FROM quests WHERE id = $1`, questID)
	}

	err := row.Scan(&quest.ID, &quest.DungeonID, &quest.Title, &quest.Description, &quest.Category, &quest.Difficulty,
		&quest.Mode, &pointsAwardStr, &quest.RatePointsPerMin, &quest.MinMinutes, &quest.MaxMinutes, &quest.DailyPointsCap,
		&quest.CooldownSec, &quest.StreakEnabled, &quest.Status, &lastCompletedAt, &quest.StreakCount,
		&quest.TimeZone, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("quest not found: %w", ErrQuestNotFound)
		}
		return nil, fmt.Errorf("failed to query quest: %w", err)
	}

	// Parse the points award
	pointsAward := valueobject.NewDecimal(pointsAwardStr)

	quest.PointsAward = pointsAward
	quest.LastCompletedAt = lastCompletedAt
	quest.CreatedAt = createdAt
	quest.UpdatedAt = updatedAt

	return &quest, nil
}

func (r *QuestRepository) ListByDungeon(ctx context.Context, dungeonID string) ([]*entity.Quest, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, dungeon_id, title, description, category, difficulty, mode, points_award, 
			rate_points_per_min, min_minutes, max_minutes, daily_points_cap, cooldown_sec, streak_enabled, 
			status, last_completed_at, streak_count, time_zone, created_at, updated_at
		FROM quests WHERE dungeon_id = $1`, dungeonID)
	if err != nil {
		return nil, fmt.Errorf("failed to query quests: %w", err)
	}
	defer rows.Close()

	var quests []*entity.Quest
	for rows.Next() {
		var quest entity.Quest
		var pointsAwardStr string
		var lastCompletedAt *time.Time
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(&quest.ID, &quest.DungeonID, &quest.Title, &quest.Description, &quest.Category, &quest.Difficulty,
			&quest.Mode, &pointsAwardStr, &quest.RatePointsPerMin, &quest.MinMinutes, &quest.MaxMinutes, &quest.DailyPointsCap,
			&quest.CooldownSec, &quest.StreakEnabled, &quest.Status, &lastCompletedAt, &quest.StreakCount,
			&quest.TimeZone, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quest: %w", err)
		}

		// Parse the points award
		pointsAward := valueobject.NewDecimal(pointsAwardStr)

		quest.PointsAward = pointsAward
		quest.LastCompletedAt = lastCompletedAt
		quest.CreatedAt = createdAt
		quest.UpdatedAt = updatedAt

		quests = append(quests, &quest)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over quest rows: %w", err)
	}

	return quests, nil
}

func (r *QuestRepository) Update(ctx context.Context, quest *entity.Quest) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			UPDATE quests 
			SET title = $1, description = $2, category = $3, difficulty = $4, mode = $5, points_award = $6,
				rate_points_per_min = $7, min_minutes = $8, max_minutes = $9, daily_points_cap = $10,
				cooldown_sec = $11, streak_enabled = $12, status = $13, last_completed_at = $14,
				streak_count = $15, time_zone = $16, updated_at = $17
			WHERE id = $18`,
			quest.Title, quest.Description, quest.Category, quest.Difficulty, quest.Mode,
			quest.PointsAward.String(), quest.RatePointsPerMin, quest.MinMinutes, quest.MaxMinutes, quest.DailyPointsCap,
			quest.CooldownSec, quest.StreakEnabled, quest.Status, quest.LastCompletedAt,
			quest.StreakCount, quest.TimeZone, quest.UpdatedAt, quest.ID)
		if err != nil {
			return fmt.Errorf("failed to update quest: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			UPDATE quests 
			SET title = $1, description = $2, category = $3, difficulty = $4, mode = $5, points_award = $6,
				rate_points_per_min = $7, min_minutes = $8, max_minutes = $9, daily_points_cap = $10,
				cooldown_sec = $11, streak_enabled = $12, status = $13, last_completed_at = $14,
				streak_count = $15, time_zone = $16, updated_at = $17
			WHERE id = $18`,
			quest.Title, quest.Description, quest.Category, quest.Difficulty, quest.Mode,
			quest.PointsAward.String(), quest.RatePointsPerMin, quest.MinMinutes, quest.MaxMinutes, quest.DailyPointsCap,
			quest.CooldownSec, quest.StreakEnabled, quest.Status, quest.LastCompletedAt,
			quest.StreakCount, quest.TimeZone, quest.UpdatedAt, quest.ID)
		if err != nil {
			return fmt.Errorf("failed to update quest: %w", err)
		}
	}

	return nil
}

func (r *QuestRepository) Delete(ctx context.Context, questID string) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `DELETE FROM quests WHERE id = $1`, questID)
		if err != nil {
			return fmt.Errorf("failed to delete quest: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `DELETE FROM quests WHERE id = $1`, questID)
		if err != nil {
			return fmt.Errorf("failed to delete quest: %w", err)
		}
	}

	return nil
}
