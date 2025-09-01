package builders

import (
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
		WithRole("member").
		WithTimeZone("UTC").
		WithDisplayName("Test User").
		WithBalance("100.00")
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

func (b *UserBuilder) WithRole(role string) *UserBuilder {
	b.With(func(u *entity.User) {
		u.Role = role
	})
	return b
}

func (b *UserBuilder) WithTimeZone(timeZone string) *UserBuilder {
	b.With(func(u *entity.User) {
		u.TimeZone = timeZone
	})
	return b
}

func (b *UserBuilder) WithDisplayName(name string) *UserBuilder {
	b.With(func(u *entity.User) {
		u.DisplayName = name
	})
	return b
}

func (b *UserBuilder) WithBalance(balance string) *UserBuilder {
	b.With(func(u *entity.User) {
		u.Balance = valueobject.NewDecimal(balance)
	})
	return b
}

func (b *UserBuilder) WithPreferences(prefs entity.UserPreferences) *UserBuilder {
	b.With(func(u *entity.User) {
		u.Preferences = prefs
	})
	return b
}
