package builders

import (
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type UserBuilder struct {
	BaseBuilder[entity.User]
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		BaseBuilder: BaseBuilder[entity.User]{
			Builder: *NewBuilder[entity.User](),
		},
	}
}

func (b *UserBuilder) WithDefaults() *UserBuilder {
	return b.
		WithID(1).
		WithChatID(100).
		WithUsername("Test User").
		WithTimezone("UTC").
		WithBalance("100.00").
		WithCreatedAt(time.Now()).
		WithUpdatedAt(time.Now())
}

func (b *UserBuilder) WithID(id int64) *UserBuilder {
	b.With(func(u *entity.User) {
		u.ID = id
	})
	return b
}

func (b *UserBuilder) WithChatID(chatID int64) *UserBuilder {
	b.With(func(u *entity.User) {
		u.ChatID = chatID
	})
	return b
}

func (b *UserBuilder) WithUsername(username string) *UserBuilder {
	b.With(func(u *entity.User) {
		u.Username = username
	})
	return b
}

func (b *UserBuilder) WithTimezone(timezone string) *UserBuilder {
	b.With(func(u *entity.User) {
		u.Timezone = timezone
	})
	return b
}

func (b *UserBuilder) WithBalance(balance string) *UserBuilder {
	b.With(func(u *entity.User) {
		u.Balance = valueobject.NewDecimal(balance)
	})
	return b
}

func (b *UserBuilder) WithCreatedAt(t time.Time) *UserBuilder {
	b.With(func(u *entity.User) {
		u.CreatedAt = t
	})
	return b
}

func (b *UserBuilder) WithUpdatedAt(t time.Time) *UserBuilder {
	b.With(func(u *entity.User) {
		u.UpdatedAt = t
	})
	return b
}
