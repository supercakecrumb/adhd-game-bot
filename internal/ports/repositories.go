package ports

import (
	"context"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	FindByChatID(ctx context.Context, chatID int64) ([]*entity.User, error)
	UpdateBalance(ctx context.Context, userID int64, delta valueobject.Decimal) error
	Delete(ctx context.Context, id int64) error
}

type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	FindByID(ctx context.Context, id string) (*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id string) error
	FindByUser(ctx context.Context, userID int64) ([]*entity.Task, error)
}

type ChatConfigRepository interface {
	Create(ctx context.Context, config *entity.ChatConfig) error
	FindByChatID(ctx context.Context, chatID int64) (*entity.ChatConfig, error)
	Update(ctx context.Context, config *entity.ChatConfig) error
}

type ShopItemRepository interface {
	Create(ctx context.Context, item *entity.ShopItem) error
	FindByID(ctx context.Context, id int64) (*entity.ShopItem, error)
	FindByCode(ctx context.Context, chatID int64, code string) (*entity.ShopItem, error)
	FindByChatID(ctx context.Context, chatID int64) ([]*entity.ShopItem, error)
	Update(ctx context.Context, item *entity.ShopItem) error
	Delete(ctx context.Context, id int64) error
}

type PurchaseRepository interface {
	Create(ctx context.Context, purchase *entity.Purchase) error
	FindByID(ctx context.Context, id int64) (*entity.Purchase, error)
	FindByUserID(ctx context.Context, userID int64) ([]*entity.Purchase, error)
	FindByItemID(ctx context.Context, itemID int64) ([]*entity.Purchase, error)
}

type UUIDGenerator interface {
	New() string
}

type Scheduler interface {
	ScheduleRecurringTask(ctx context.Context, task *entity.Task) error
	CancelScheduledTask(ctx context.Context, taskID string) error
	GetNextOccurrence(ctx context.Context, taskID string) (time.Time, error)
}

type IdempotencyRepository interface {
	Create(ctx context.Context, key *entity.IdempotencyKey) error
	FindByKey(ctx context.Context, key string) (*entity.IdempotencyKey, error)
	Update(ctx context.Context, key *entity.IdempotencyKey) error
	DeleteExpired(ctx context.Context) error
}
