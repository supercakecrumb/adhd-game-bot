package builders_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/test/fixtures/builders"
)

func TestUserBuilder(t *testing.T) {
	t.Run("WithDefaults creates valid user", func(t *testing.T) {
		user := builders.NewUserBuilder().
			WithDefaults().
			Build()

		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, int64(100), user.ChatID)
		assert.Equal(t, "member", user.Role)
		assert.Equal(t, "UTC", user.TimeZone)
		assert.Equal(t, "Test User", user.DisplayName)
		assert.Equal(t, "100.00", user.Balance.String())
	})

	t.Run("Can override defaults", func(t *testing.T) {
		user := builders.NewUserBuilder().
			WithDefaults().
			WithID(2).
			WithChatID(200).
			WithRole("admin").
			WithTimeZone("Europe/Moscow").
			WithDisplayName("Admin User").
			WithBalance("500.00").
			Build()

		assert.Equal(t, int64(2), user.ID)
		assert.Equal(t, int64(200), user.ChatID)
		assert.Equal(t, "admin", user.Role)
		assert.Equal(t, "Europe/Moscow", user.TimeZone)
		assert.Equal(t, "Admin User", user.DisplayName)
		assert.Equal(t, "500.00", user.Balance.String())
	})

	t.Run("WithPreferences sets user preferences", func(t *testing.T) {
		prefs := entity.UserPreferences{
			EditIntervalSec: 60,
			NotifyOnAward:   true,
		}

		user := builders.NewUserBuilder().
			WithPreferences(prefs).
			Build()

		assert.Equal(t, 60, user.Preferences.EditIntervalSec)
		assert.True(t, user.Preferences.NotifyOnAward)
	})
}
