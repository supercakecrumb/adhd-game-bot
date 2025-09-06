package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type DungeonMemberRepository struct {
	db *sql.DB
}

func NewDungeonMemberRepository(db *sql.DB) *DungeonMemberRepository {
	return &DungeonMemberRepository{db: db}
}

func (r *DungeonMemberRepository) Add(ctx context.Context, dungeonID string, userID int64) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO dungeon_members (dungeon_id, user_id)
			VALUES ($1, $2)
			ON CONFLICT (dungeon_id, user_id) DO NOTHING`,
			dungeonID, userID)
		if err != nil {
			return fmt.Errorf("failed to add dungeon member: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO dungeon_members (dungeon_id, user_id)
			VALUES ($1, $2)
			ON CONFLICT (dungeon_id, user_id) DO NOTHING`,
			dungeonID, userID)
		if err != nil {
			return fmt.Errorf("failed to add dungeon member: %w", err)
		}
	}

	return nil
}

func (r *DungeonMemberRepository) ListUsers(ctx context.Context, dungeonID string) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT user_id
		FROM dungeon_members 
		WHERE dungeon_id = $1`, dungeonID)
	if err != nil {
		return nil, fmt.Errorf("failed to query dungeon members: %w", err)
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var userID int64
		err := rows.Scan(&userID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user ID: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over dungeon member rows: %w", err)
	}

	return userIDs, nil
}

func (r *DungeonMemberRepository) IsMember(ctx context.Context, dungeonID string, userID int64) (bool, error) {
	var count int
	var row *sql.Row

	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT COUNT(*) 
			FROM dungeon_members 
			WHERE dungeon_id = $1 AND user_id = $2`, dungeonID, userID)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT COUNT(*) 
			FROM dungeon_members 
			WHERE dungeon_id = $1 AND user_id = $2`, dungeonID, userID)
	}

	err := row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check dungeon membership: %w", err)
	}

	return count > 0, nil
}
