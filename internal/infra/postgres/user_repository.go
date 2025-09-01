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
		INSERT INTO users (id, role, timezone, display_name, preferences_json)
		VALUES ($1, $2, $3, $4, $5)`,
		user.ID, user.Role, user.TimeZone, user.DisplayName, prefs)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Initialize balances for all currencies
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO user_balances (user_id, currency_code)
		SELECT $1, code FROM currencies`,
		user.ID)
	if err != nil {
		return fmt.Errorf("failed to initialize balances: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	var prefsJSON []byte

	err := r.db.QueryRowContext(ctx, `
		SELECT id, role, timezone, display_name, preferences_json
		FROM users WHERE id = $1`, id).Scan(
		&user.ID, &user.Role, &user.TimeZone, &user.DisplayName, &prefsJSON)
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
		SELECT currency_code, amount FROM user_balances WHERE user_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query balances: %w", err)
	}
	defer rows.Close()

	user.Balances = make(map[string]valueobject.Decimal)
	for rows.Next() {
		var currency string
		var amount string
		if err := rows.Scan(&currency, &amount); err != nil {
			return nil, fmt.Errorf("failed to scan balance: %w", err)
		}
		user.Balances[currency] = valueobject.NewDecimal(amount)
	}

	return &user, nil
}

func (r *UserRepository) UpdateBalance(ctx context.Context, userID int64, currency string, delta valueobject.Decimal) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE user_balances
		SET amount = CAST(amount AS NUMERIC) + CAST($1 AS NUMERIC)
		WHERE user_id = $2 AND currency_code = $3`,
		delta.String(), userID, currency)
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
