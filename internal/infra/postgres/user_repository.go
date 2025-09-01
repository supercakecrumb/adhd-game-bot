package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
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
	prefs, err := json.Marshal(user.Preferences)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %w", err)
	}

	// Check if we're in a transaction
	if tx, ok := GetTx(ctx); ok {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO users (id, chat_id, role, timezone, display_name, preferences_json, balance)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			user.ID, user.ChatID, user.Role, user.TimeZone, user.DisplayName, prefs, user.Balance.String())
	} else {
		_, err = r.db.ExecContext(ctx, `
			INSERT INTO users (id, chat_id, role, timezone, display_name, preferences_json, balance)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			user.ID, user.ChatID, user.Role, user.TimeZone, user.DisplayName, prefs, user.Balance.String())
	}

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	var prefsJSON []byte
	var balanceStr string

	var row *sql.Row
	if tx, ok := GetTx(ctx); ok {
		row = tx.QueryRowContext(ctx, `
			SELECT id, chat_id, role, timezone, display_name, preferences_json, balance
			FROM users WHERE id = $1`, id)
	} else {
		row = r.db.QueryRowContext(ctx, `
			SELECT id, chat_id, role, timezone, display_name, preferences_json, balance
			FROM users WHERE id = $1`, id)
	}

	err := row.Scan(&user.ID, &user.ChatID, &user.Role, &user.TimeZone, &user.DisplayName, &prefsJSON, &balanceStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// Unmarshal preferences
	if err := json.Unmarshal(prefsJSON, &user.Preferences); err != nil {
		return nil, fmt.Errorf("failed to unmarshal preferences: %w", err)
	}

	// Parse balance
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
