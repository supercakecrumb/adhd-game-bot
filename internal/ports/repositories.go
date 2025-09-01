package ports

import (
	"context"
	"errors"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	UpdateBalance(ctx context.Context, userID int64, currencyID int64, delta entity.Decimal) error
	Delete(ctx context.Context, id int64) error
}

type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	FindByID(ctx context.Context, id string) (*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id string) error
	FindByUser(ctx context.Context, userID int64) ([]*entity.Task, error)
}

type CurrencyRepository interface {
	Create(ctx context.Context, currency *entity.Currency) error
	FindByID(ctx context.Context, id int64) (*entity.Currency, error)
	FindByCode(ctx context.Context, chatID int64, code string) (*entity.Currency, error)
	FindByChatID(ctx context.Context, chatID int64) ([]*entity.Currency, error)
	GetBaseCurrency(ctx context.Context, chatID int64) (*entity.Currency, error)
	Update(ctx context.Context, currency *entity.Currency) error
	Delete(ctx context.Context, id int64) error
}

type UUIDGenerator interface {
	New() string
}
