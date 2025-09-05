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
	FindActiveByUser(ctx context.Context, userID int64) ([]*entity.Task, error)
	FindWithTimers(ctx context.Context, userID int64) ([]*entity.Task, error)
	FindWithSchedules(ctx context.Context, userID int64) ([]*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id string) error
	FindByUser(ctx context.Context, userID int64) ([]*entity.Task, error)
	BulkUpdate(ctx context.Context, tasks []*entity.Task) error
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

type TimerRepository interface {
	Create(ctx context.Context, timer *entity.Timer) error
	FindByID(ctx context.Context, id string) (*entity.Timer, error)
	FindByUser(ctx context.Context, userID int64) ([]*entity.Timer, error)
	FindByTask(ctx context.Context, taskID string) ([]*entity.Timer, error)
	FindActiveByUser(ctx context.Context, userID int64) ([]*entity.Timer, error)
	Update(ctx context.Context, timer *entity.Timer) error
	Delete(ctx context.Context, id string) error
	BulkUpdate(ctx context.Context, timers []*entity.Timer) error
}

type TimerEventRepository interface {
	Create(ctx context.Context, event *entity.TimerEvent) error
	FindByTimer(ctx context.Context, timerID string) ([]*entity.TimerEvent, error)
	DeleteOldEvents(ctx context.Context, olderThan time.Time) error
}

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *entity.Schedule) error
	FindByID(ctx context.Context, id string) (*entity.Schedule, error)
	FindByTask(ctx context.Context, taskID string) ([]*entity.Schedule, error)
	FindByUser(ctx context.Context, userID int64) ([]*entity.Schedule, error)
	Update(ctx context.Context, schedule *entity.Schedule) error
	Delete(ctx context.Context, id string) error
}

type RewardTierRepository interface {
	Create(ctx context.Context, tier *entity.RewardTier) error
	FindByID(ctx context.Context, id int64) (*entity.RewardTier, error)
	FindAll(ctx context.Context) ([]*entity.RewardTier, error)
	Update(ctx context.Context, tier *entity.RewardTier) error
	Delete(ctx context.Context, id int64) error
}

type DiscountTierRepository interface {
	Create(ctx context.Context, tier *entity.DiscountTier) error
	FindByID(ctx context.Context, id int64) (*entity.DiscountTier, error)
	FindAll(ctx context.Context) ([]*entity.DiscountTier, error)
	Update(ctx context.Context, tier *entity.DiscountTier) error
	Delete(ctx context.Context, id int64) error
}

type IdempotencyRepository interface {
	Create(ctx context.Context, key *entity.IdempotencyKey) error
	FindByKey(ctx context.Context, key string) (*entity.IdempotencyKey, error)
	Update(ctx context.Context, key *entity.IdempotencyKey) error
	DeleteExpired(ctx context.Context) error
	Purge(ctx context.Context, olderThan time.Time) error
}
