package builders_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supercakecrumb/adhd-game-bot/test/fixtures/builders"
)

func TestShopItemBuilder(t *testing.T) {
	t.Run("WithDefaults creates valid shop item", func(t *testing.T) {
		item := builders.NewShopItemBuilder().
			WithDefaults().
			Build()

		assert.Equal(t, int64(1), item.ID)
		assert.Equal(t, int64(100), item.ChatID)
		assert.Equal(t, "BOOST", item.Code)
		assert.Equal(t, "XP Boost", item.Name)
		assert.Equal(t, "Gives 2x XP for 1 hour", item.Description)
		assert.Equal(t, 50.0, item.Price.Float64())
		assert.Equal(t, "rewards", item.Category)
		assert.True(t, item.IsActive)
		assert.Equal(t, 10, *item.Stock)
	})

	t.Run("Can override defaults", func(t *testing.T) {
		now := time.Now()
		stock := 5
		item := builders.NewShopItemBuilder().
			WithDefaults().
			WithID(2).
			WithChatID(200).
			WithCode("SUPER_BOOST").
			WithName("Super XP Boost").
			WithDescription("Gives 3x XP for 2 hours").
			WithPrice("100.00").
			WithCategory("premium").
			WithIsActive(false).
			WithStock(&stock).
			WithCreatedAt(now).
			WithUpdatedAt(now).
			Build()

		assert.Equal(t, int64(2), item.ID)
		assert.Equal(t, int64(200), item.ChatID)
		assert.Equal(t, "SUPER_BOOST", item.Code)
		assert.Equal(t, "Super XP Boost", item.Name)
		assert.Equal(t, "Gives 3x XP for 2 hours", item.Description)
		assert.Equal(t, 100.0, item.Price.Float64())
		assert.Equal(t, "premium", item.Category)
		assert.False(t, item.IsActive)
		assert.Equal(t, 5, *item.Stock)
		assert.Equal(t, now, item.CreatedAt)
		assert.Equal(t, now, item.UpdatedAt)
	})

	t.Run("Can set unlimited stock", func(t *testing.T) {
		item := builders.NewShopItemBuilder().
			WithStock(nil).
			Build()

		assert.Nil(t, item.Stock)
	})
}
