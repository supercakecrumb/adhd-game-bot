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

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO users (id, chat_id, role, timezone, display_name, preferences_json)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, user.ChatID, user.Role, user.TimeZone, user.DisplayName, prefs)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Initialize balances for all currencies in the user's chat
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO user_balances (user_id, currency_id, amount)
		SELECT $1, id, '0' FROM currencies WHERE chat_id = $2`,
		user.ID, user.ChatID)
	if err != nil {
		return fmt.Errorf("failed to initialize balances: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	var prefsJSON []byte

	err := r.db.QueryRowContext(ctx, `
		SELECT id, chat_id, role, timezone, display_name, preferences_json
		FROM users WHERE id = $1`, id).Scan(
		&user.ID, &user.ChatID, &user.Role, &user.TimeZone, &user.DisplayName, &prefsJSON)
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

	// Get balances
	rows, err := r.db.QueryContext(ctx, `
		SELECT currency_id, amount FROM user_balances WHERE user_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query balances: %w", err)
	}
	defer rows.Close()

	user.Balances = make(map[int64]valueobject.Decimal)
	for rows.Next() {
		var currencyID int64
		var amount string
		if err := rows.Scan(&currencyID, &amount); err != nil {
			return nil, fmt.Errorf("failed to scan balance: %w", err)
		}
		user.Balances[currencyID] = valueobject.NewDecimal(amount)
	}

	return &user, nil
}

func (r *UserRepository) UpdateBalance(ctx context.Context, userID int64, currencyID int64, delta valueobject.Decimal) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE user_balances
		SET amount = CAST(amount AS NUMERIC) + CAST($1 AS NUMERIC)
		WHERE user_id = $2 AND currency_id = $3`,
		delta.String(), userID, currencyID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
