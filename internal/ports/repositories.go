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

type QuestRepository interface {
	Create(ctx context.Context, quest *entity.Quest) error
	GetByID(ctx context.Context, questID string) (*entity.Quest, error)
	ListByDungeon(ctx context.Context, dungeonID string) ([]*entity.Quest, error)
	Update(ctx context.Context, quest *entity.Quest) error
	Delete(ctx context.Context, questID string) error
}

type QuestCompletionRepository interface {
	Insert(ctx context.Context, completion *entity.QuestCompletion) error
	LastForUser(ctx context.Context, userID int64, questID string) (*entity.QuestCompletion, error)
	SumAwardedForUserOnDay(ctx context.Context, userID int64, questID string, day time.Time, tz string) (valueobject.Decimal, error)
}

type DungeonRepository interface {
	Create(ctx context.Context, dungeon *entity.Dungeon) error
	GetByID(ctx context.Context, dungeonID string) (*entity.Dungeon, error)
	ListByAdmin(ctx context.Context, userID int64) ([]*entity.Dungeon, error)
}

type DungeonMemberRepository interface {
	Add(ctx context.Context, dungeonID string, userID int64) error
	ListUsers(ctx context.Context, dungeonID string) ([]int64, error)
	IsMember(ctx context.Context, dungeonID string, userID int64) (bool, error)
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
	ScheduleRecurringTask(ctx context.Context, task *entity.Quest) error
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
