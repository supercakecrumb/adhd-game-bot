package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO users (id, chat_id, timezone, display_name, balance, role, preferences_json)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			user.ID, user.ChatID, user.TimeZone, user.Username, user.Balance.String(), "member", "{}")
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO users (id, chat_id, timezone, display_name, balance, role, preferences_json)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			user.ID, user.ChatID, user.TimeZone, user.Username, user.Balance.String(), "member", "{}")
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	var balanceStr string
	var role, timezone, displayName, preferencesJSON string
	var createdAt, updatedAt interface{} // We'll ignore these for now

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT id, chat_id, role, timezone, display_name, preferences_json, balance, created_at, updated_at
			FROM users WHERE id = $1`, id)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT id, chat_id, role, timezone, display_name, preferences_json, balance, created_at, updated_at
			FROM users WHERE id = $1`, id)
	}

	err := row.Scan(&user.ID, &user.ChatID, &role, &timezone, &displayName, &preferencesJSON, &balanceStr, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// Map fields to user entity
	user.TimeZone = timezone
	user.Username = displayName
	user.Balance = valueobject.NewDecimal(balanceStr)

	return &user, nil
}

func (r *UserRepository) UpdateBalance(ctx context.Context, userID int64, delta valueobject.Decimal) error {
	var err error
	if tx, ok := GetTx(ctx); ok {
		_, err = tx.ExecContext(ctx, `
			UPDATE users
			SET balance = CAST(balance AS NUMERIC) + CAST($1 AS NUMERIC)
			WHERE id = $2`,
			delta.String(), userID)
	} else {
		_, err = r.db.ExecContext(ctx, `
			UPDATE users
			SET balance = CAST(balance AS NUMERIC) + CAST($1 AS NUMERIC)
			WHERE id = $2`,
			delta.String(), userID)
	}

	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	var err error
	if tx, ok := GetTx(ctx); ok {
		_, err = tx.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	} else {
		_, err = r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	}
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByChatID(ctx context.Context, chatID int64) ([]*entity.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, chat_id, role, timezone, display_name, preferences_json, balance, created_at, updated_at
		FROM users WHERE chat_id = $1`, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query users by chat_id: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		var balanceStr string
		var role, timezone, displayName, preferencesJSON string
		var createdAt, updatedAt interface{} // We'll ignore these for now
		err := rows.Scan(&user.ID, &user.ChatID, &role, &timezone, &displayName, &preferencesJSON, &balanceStr, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		user.TimeZone = timezone
		user.Username = displayName
		user.Balance = valueobject.NewDecimal(balanceStr)
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user rows: %w", err)
	}

	return users, nil
}
