package builders_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supercakecrumb/adhd-game-bot/test/fixtures/builders"
)

func TestUserBuilder(t *testing.T) {
	t.Run("WithDefaults creates valid user", func(t *testing.T) {
		user := builders.NewUserBuilder().
			WithDefaults().
			Build()

		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, int64(100), user.ChatID)
		assert.Equal(t, "UTC", user.Timezone)
		assert.Equal(t, "Test User", user.Username)
		assert.Equal(t, 100.0, user.Balance.Float64())
	})

	t.Run("Can override defaults", func(t *testing.T) {
		testTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		user := builders.NewUserBuilder().
			WithDefaults().
			WithID(2).
			WithChatID(200).
			WithUsername("Admin User").
			WithTimezone("Europe/Moscow").
			WithBalance("500.00").
			WithCreatedAt(testTime).
			WithUpdatedAt(testTime).
			Build()

		assert.Equal(t, int64(2), user.ID)
		assert.Equal(t, int64(200), user.ChatID)
		assert.Equal(t, "Europe/Moscow", user.Timezone)
		assert.Equal(t, "Admin User", user.Username)
		assert.Equal(t, 500.0, user.Balance.Float64())
		assert.Equal(t, testTime, user.CreatedAt)
		assert.Equal(t, testTime, user.UpdatedAt)
	})
}
