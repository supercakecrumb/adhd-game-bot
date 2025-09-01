package builders

import (
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type ShopItemBuilder struct {
	BaseBuilder[entity.ShopItem]
}

func NewShopItemBuilder() *ShopItemBuilder {
	return &ShopItemBuilder{
		BaseBuilder: BaseBuilder[entity.ShopItem]{
			Builder: *NewBuilder[entity.ShopItem](),
		},
	}
}

func (b *ShopItemBuilder) WithDefaults() *ShopItemBuilder {
	stock := 10
	return b.
		WithID(1).
		WithChatID(100).
		WithCode("BOOST").
		WithName("XP Boost").
		WithDescription("Gives 2x XP for 1 hour").
		WithPrice("50.00").
		WithCategory("rewards").
		WithIsActive(true).
		WithStock(&stock)
}

func (b *ShopItemBuilder) WithID(id int64) *ShopItemBuilder {
	b.With(func(i *entity.ShopItem) {
		i.ID = id
	})
	return b
}

func (b *ShopItemBuilder) WithChatID(chatID int64) *ShopItemBuilder {
	b.With(func(i *entity.ShopItem) {
		i.ChatID = chatID
	})
	return b
}

func (b *ShopItemBuilder) WithCode(code string) *ShopItemBuilder {
	b.With(func(i *entity.ShopItem) {
		i.Code = code
	})
	return b
}

func (b *ShopItemBuilder) WithName(name string) *ShopItemBuilder {
	b.With(func(i *entity.ShopItem) {
		i.Name = name
	})
	return b
}

func (b *ShopItemBuilder) WithDescription(description string) *ShopItemBuilder {
	b.With(func(i *entity.ShopItem) {
		i.Description = description
	})
	return b
}

func (b *ShopItemBuilder) WithPrice(price string) *ShopItemBuilder {
	b.With(func(i *entity.ShopItem) {
		i.Price = valueobject.NewDecimal(price)
	})
	return b
}

func (b *ShopItemBuilder) WithCategory(category string) *ShopItemBuilder {
	b.With(func(i *entity.ShopItem) {
		i.Category = category
	})
	return b
}

func (b *ShopItemBuilder) WithIsActive(isActive bool) *ShopItemBuilder {
	b.With(func(i *entity.ShopItem) {
		i.IsActive = isActive
	})
	return b
}

func (b *ShopItemBuilder) WithStock(stock *int) *ShopItemBuilder {
	b.With(func(i *entity.ShopItem) {
		i.Stock = stock
	})
	return b
}

func (b *ShopItemBuilder) WithCreatedAt(time time.Time) *ShopItemBuilder {
	b.With(func(i *entity.ShopItem) {
		i.CreatedAt = time
	})
	return b
}

func (b *ShopItemBuilder) WithUpdatedAt(time time.Time) *ShopItemBuilder {
	b.With(func(i *entity.ShopItem) {
		i.UpdatedAt = time
	})
	return b
}
