package postgres_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/infra/postgres"
)

func TestUserRepository(t *testing.T) {
	// Get connection string from env or use default
	connStr := os.Getenv("TEST_DB_CONN")
	if connStr == "" {
		connStr = "user=postgres dbname=adhd_bot_test sslmode=disable"
	}

	// Connect to test database
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	// Clean database before tests
	_, err = db.Exec("DROP TABLE IF EXISTS users, user_balances CASCADE")
	require.NoError(t, err)

	// Apply migrations
	_, err = db.Exec(`
		CREATE TABLE users (
			id BIGSERIAL PRIMARY KEY,
			role VARCHAR(10) NOT NULL,
			timezone VARCHAR(50) NOT NULL,
			display_name VARCHAR(255) NOT NULL,
			preferences_json JSONB NOT NULL
		);
		
		CREATE TABLE user_balances (
			user_id BIGINT NOT NULL,
			currency_code VARCHAR(10) NOT NULL,
			amount NUMERIC(20, 8) NOT NULL,
			PRIMARY KEY (user_id, currency_code),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	require.NoError(t, err)

	repo := postgres.NewUserRepository(db)

	t.Run("Create and find user", func(t *testing.T) {
		user := &entity.User{
			ID:          1,
			Role:        "member",
			TimeZone:    "UTC",
			DisplayName: "Test User",
			Preferences: entity.UserPreferences{
				EditIntervalSec: 300,
				NotifyOnAward:   true,
			},
		}

		err := repo.Create(context.Background(), user)
		require.NoError(t, err)

		// Initialize some balances for testing
		_, err = db.Exec(`
			INSERT INTO user_balances (user_id, currency_code, amount)
			VALUES (1, 'MM', 10.5), (1, 'BS', 5.25)
		`)
		require.NoError(t, err)

		found, err := repo.FindByID(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, user.ID, found.ID)
		require.Equal(t, user.DisplayName, found.DisplayName)
		require.Equal(t, 10.5, found.Balances["MM"].Float64())
		require.Equal(t, 5.25, found.Balances["BS"].Float64())
	})

	t.Run("Update balance", func(t *testing.T) {
		err := repo.UpdateBalance(context.Background(), 1, "MM", valueobject.NewDecimal("5.25"))
		require.NoError(t, err)

		var balance string
		err = db.QueryRow("SELECT amount FROM user_balances WHERE user_id = 1 AND currency_code = 'MM'").
			Scan(&balance)
		require.NoError(t, err)
		require.Equal(t, "15.75000000", balance) // 10.5 + 5.25
	})

	t.Run("Delete user", func(t *testing.T) {
		err := repo.Delete(context.Background(), 1)
		require.NoError(t, err)

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = 1").Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}
