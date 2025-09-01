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
	_, err = db.Exec("DROP TABLE IF EXISTS users CASCADE")
	require.NoError(t, err)

	// Apply migrations
	_, err = db.Exec(`
		CREATE TABLE users (
			id BIGINT PRIMARY KEY,
			chat_id BIGINT NOT NULL,
			role VARCHAR(10) NOT NULL,
			timezone VARCHAR(50) NOT NULL,
			display_name VARCHAR(255) NOT NULL,
			preferences_json JSONB NOT NULL,
			balance NUMERIC(20, 8) NOT NULL DEFAULT 0
		);
	`)
	require.NoError(t, err)

	repo := postgres.NewUserRepository(db)

	t.Run("Create and find user", func(t *testing.T) {
		user := &entity.User{
			ID:          1,
			ChatID:      100,
			Role:        "member",
			TimeZone:    "UTC",
			DisplayName: "Test User",
			Balance:     valueobject.NewDecimal("10.5"),
			Preferences: entity.UserPreferences{
				EditIntervalSec: 300,
				NotifyOnAward:   true,
			},
		}

		err := repo.Create(context.Background(), user)
		require.NoError(t, err)

		found, err := repo.FindByID(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, user.ID, found.ID)
		require.Equal(t, user.ChatID, found.ChatID)
		require.Equal(t, user.DisplayName, found.DisplayName)
		require.Equal(t, 10.5, found.Balance.Float64())
	})

	t.Run("Update balance", func(t *testing.T) {
		err := repo.UpdateBalance(context.Background(), 1, valueobject.NewDecimal("5.25"))
		require.NoError(t, err)

		found, err := repo.FindByID(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, 15.75, found.Balance.Float64()) // 10.5 + 5.25
	})

	t.Run("Update balance negative", func(t *testing.T) {
		err := repo.UpdateBalance(context.Background(), 1, valueobject.NewDecimal("-5.25"))
		require.NoError(t, err)

		found, err := repo.FindByID(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, 10.5, found.Balance.Float64()) // 15.75 - 5.25
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
