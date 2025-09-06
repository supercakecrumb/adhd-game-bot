package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
)

type DungeonRepository struct {
	db *sql.DB
}

func NewDungeonRepository(db *sql.DB) *DungeonRepository {
	return &DungeonRepository{db: db}
}

func (r *DungeonRepository) Create(ctx context.Context, dungeon *entity.Dungeon) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO dungeons (id, title, admin_user_id, telegram_chat_id, created_at)
			VALUES ($1, $2, $3, $4, $5)`,
			dungeon.ID, dungeon.Title, dungeon.AdminUserID, dungeon.TelegramChatID, dungeon.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to create dungeon: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO dungeons (id, title, admin_user_id, telegram_chat_id, created_at)
			VALUES ($1, $2, $3, $4, $5)`,
			dungeon.ID, dungeon.Title, dungeon.AdminUserID, dungeon.TelegramChatID, dungeon.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to create dungeon: %w", err)
		}
	}

	return nil
}

func (r *DungeonRepository) GetByID(ctx context.Context, dungeonID string) (*entity.Dungeon, error) {
	var dungeon entity.Dungeon
	var telegramChatID *int64
	var createdAt time.Time

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT id, title, admin_user_id, telegram_chat_id, created_at
			FROM dungeons WHERE id = $1`, dungeonID)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT id, title, admin_user_id, telegram_chat_id, created_at
			FROM dungeons WHERE id = $1`, dungeonID)
	}

	err := row.Scan(&dungeon.ID, &dungeon.Title, &dungeon.AdminUserID, &telegramChatID, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("dungeon not found: %w", ErrDungeonNotFound)
		}
		return nil, fmt.Errorf("failed to query dungeon: %w", err)
	}

	dungeon.TelegramChatID = telegramChatID
	dungeon.CreatedAt = createdAt

	return &dungeon, nil
}

func (r *DungeonRepository) ListByAdmin(ctx context.Context, userID int64) ([]*entity.Dungeon, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, admin_user_id, telegram_chat_id, created_at
		FROM dungeons WHERE admin_user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query dungeons: %w", err)
	}
	defer rows.Close()

	var dungeons []*entity.Dungeon
	for rows.Next() {
		var dungeon entity.Dungeon
		var telegramChatID *int64
		var createdAt time.Time

		err := rows.Scan(&dungeon.ID, &dungeon.Title, &dungeon.AdminUserID, &telegramChatID, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan dungeon: %w", err)
		}

		dungeon.TelegramChatID = telegramChatID
		dungeon.CreatedAt = createdAt

		dungeons = append(dungeons, &dungeon)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over dungeon rows: %w", err)
	}

	return dungeons, nil
}
