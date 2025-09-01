package scenarios

import (
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/test/fixtures/builders"
)

// NewUserWithPurchase creates a scenario with:
// - A new user with balance
// - An active shop item
// - A completed purchase
func NewUserWithPurchase() *ShopScenario {
	user := builders.NewUserBuilder().
		WithDefaults().
		WithBalance("100.00").
		Build()

	item := builders.NewShopItemBuilder().
		WithDefaults().
		Build()

	purchaseTime := time.Now()

	return &ShopScenario{
		User:     *user,
		ShopItem: *item,
		Purchase: entity.Purchase{
			UserID:      user.ID,
			ItemID:      item.ID,
			ItemName:    item.Name,
			ItemPrice:   item.Price,
			Quantity:    1,
			TotalCost:   item.Price,
			Status:      "completed",
			PurchasedAt: purchaseTime,
		},
	}
}

// NewUserWithInsufficientBalance creates a scenario with:
// - A user with low balance
// - An expensive shop item
func NewUserWithInsufficientBalance() *ShopScenario {
	user := builders.NewUserBuilder().
		WithDefaults().
		WithBalance("10.00").
		Build()

	item := builders.NewShopItemBuilder().
		WithDefaults().
		WithPrice("100.00").
		Build()

	return &ShopScenario{
		User:     *user,
		ShopItem: *item,
	}
}

// ShopScenario combines entities for shop-related test scenarios
type ShopScenario struct {
	User     entity.User
	ShopItem entity.ShopItem
	Purchase entity.Purchase
}
