package inmemory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/adhd-game-bot/internal/domain/entity"
	"github.com/yourusername/adhd-game-bot/internal/domain/valueobject"
	"github.com/yourusername/adhd-game-bot/internal/infra/inmemory"
)

func TestUserRepository(t *testing.T) {
	ctx := context.Background()
	repo := inmemory.NewUserRepository()

	// Test Create
	t.Run("Create new user", func(t *testing.T) {
		user := &entity.User{
			ID:          1,
			DisplayName: "Test User",
		}

		err := repo.Create(ctx, user)
		assert.NoError(t, err)
	})

	t.Run("Create duplicate user", func(t *testing.T) {
		user := &entity.User{
			ID:          1,
			DisplayName: "Duplicate User",
		}

		err := repo.Create(ctx, user)
		assert.Error(t, err)
	})

	// Test FindByID
	t.Run("Find existing user", func(t *testing.T) {
		user, err := repo.FindByID(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, "Test User", user.DisplayName)
	})

	t.Run("Find non-existent user", func(t *testing.T) {
		_, err := repo.FindByID(ctx, 999)
		assert.Error(t, err)
	})

	// Test UpdateBalance
	t.Run("Update balance", func(t *testing.T) {
		err := repo.UpdateBalance(ctx, 1, "MM", valueobject.NewDecimal("10.50"))
		assert.NoError(t, err)

		user, err := repo.FindByID(ctx, 1)
		assert.NoError(t, err)
		// shopspring/decimal's String() doesn't preserve trailing zeros
		val := user.Balances["MM"].Float64()
		assert.Equal(t, 10.5, val)
	})

	t.Run("Update balance nonexistent user", func(t *testing.T) {
		err := repo.UpdateBalance(ctx, 999, "MM", valueobject.NewDecimal("10.50"))
		assert.Error(t, err)
	})

	// Test Delete
	t.Run("Delete user", func(t *testing.T) {
		err := repo.Delete(ctx, 1)
		assert.NoError(t, err)

		_, err = repo.FindByID(ctx, 1)
		assert.Error(t, err)
	})
}
