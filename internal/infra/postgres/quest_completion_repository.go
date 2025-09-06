package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type QuestCompletionRepository struct {
	db *sql.DB
}

func NewQuestCompletionRepository(db *sql.DB) *QuestCompletionRepository {
	return &QuestCompletionRepository{db: db}
}

func (r *QuestCompletionRepository) Insert(ctx context.Context, completion *entity.QuestCompletion) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO quest_completions (id, quest_id, user_id, dungeon_id, submitted_at, completion_ratio, minutes, awarded_points, idempotency_key)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			completion.ID, completion.QuestID, completion.UserID, completion.DungeonID, completion.SubmittedAt,
			completion.CompletionRatio, completion.Minutes, completion.AwardedPoints.String(), completion.IdempotencyKey)
		if err != nil {
			return fmt.Errorf("failed to insert quest completion: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO quest_completions (id, quest_id, user_id, dungeon_id, submitted_at, completion_ratio, minutes, awarded_points, idempotency_key)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			completion.ID, completion.QuestID, completion.UserID, completion.DungeonID, completion.SubmittedAt,
			completion.CompletionRatio, completion.Minutes, completion.AwardedPoints.String(), completion.IdempotencyKey)
		if err != nil {
			return fmt.Errorf("failed to insert quest completion: %w", err)
		}
	}

	return nil
}

func (r *QuestCompletionRepository) LastForUser(ctx context.Context, userID int64, questID string) (*entity.QuestCompletion, error) {
	var completion entity.QuestCompletion
	var awardedPointsStr string
	var submittedAt time.Time
	var completionRatio *float64
	var minutes *int

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT id, quest_id, user_id, dungeon_id, submitted_at, completion_ratio, minutes, awarded_points, idempotency_key
			FROM quest_completions 
			WHERE user_id = $1 AND quest_id = $2
			ORDER BY submitted_at DESC
			LIMIT 1`, userID, questID)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT id, quest_id, user_id, dungeon_id, submitted_at, completion_ratio, minutes, awarded_points, idempotency_key
			FROM quest_completions 
			WHERE user_id = $1 AND quest_id = $2
			ORDER BY submitted_at DESC
			LIMIT 1`, userID, questID)
	}

	err := row.Scan(&completion.ID, &completion.QuestID, &completion.UserID, &completion.DungeonID, &submittedAt, &completionRatio, &minutes, &awardedPointsStr, &completion.IdempotencyKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No completion found
		}
		return nil, fmt.Errorf("failed to query quest completion: %w", err)
	}

	// Parse the awarded points
	awardedPoints := valueobject.NewDecimal(awardedPointsStr)

	completion.SubmittedAt = submittedAt
	completion.CompletionRatio = completionRatio
	completion.Minutes = minutes
	completion.AwardedPoints = awardedPoints

	return &completion, nil
}

func (r *QuestCompletionRepository) SumAwardedForUserOnDay(ctx context.Context, userID int64, questID string, day time.Time, tz string) (valueobject.Decimal, error) {
	var sumStr string

	// Calculate the start and end of the day in the given timezone
	startOfDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT COALESCE(SUM(awarded_points), '0') as total
			FROM quest_completions 
			WHERE user_id = $1 AND quest_id = $2 AND submitted_at >= $3 AND submitted_at <= $4`,
			userID, questID, startOfDay, endOfDay)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT COALESCE(SUM(awarded_points), '0') as total
			FROM quest_completions 
			WHERE user_id = $1 AND quest_id = $2 AND submitted_at >= $3 AND submitted_at <= $4`,
			userID, questID, startOfDay, endOfDay)
	}

	err := row.Scan(&sumStr)
	if err != nil {
		return valueobject.NewDecimal("0"), fmt.Errorf("failed to sum awarded points: %w", err)
	}

	// Parse the sum
	sum := valueobject.NewDecimal(sumStr)

	return sum, nil
}
