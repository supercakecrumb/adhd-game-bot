package ports

import (
	"context"
	"errors"

	"github.com/yourusername/adhd-game-bot/internal/domain/entity"
	"github.com/yourusername/adhd-game-bot/internal/domain/valueobject"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	UpdateBalance(ctx context.Context, userID int64, currency string, delta valueobject.Decimal) error
	Delete(ctx context.Context, id int64) error
}
