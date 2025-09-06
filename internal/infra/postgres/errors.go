package postgres

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrDungeonNotFound      = errors.New("dungeon not found")
	ErrQuestNotFound        = errors.New("quest not found")
	ErrTimerNotFound        = errors.New("timer not found")
	ErrScheduleNotFound     = errors.New("schedule not found")
	ErrItemNotFound         = errors.New("item not found")
	ErrPurchaseNotFound     = errors.New("purchase not found")
	ErrRewardTierNotFound   = errors.New("reward tier not found")
	ErrDiscountTierNotFound = errors.New("discount tier not found")
)
