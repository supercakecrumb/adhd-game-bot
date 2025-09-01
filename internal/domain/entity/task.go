package entity

import (
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type Decimal = valueobject.Decimal

type Task struct {
	ID              string
	ChatID          int64 // Telegram chat ID this task belongs to
	Title           string
	Description     string
	Category        string // daily, weekly, adhoc
	Difficulty      string // easy, medium, hard
	ScheduleJSON    string
	BaseDuration    int
	GracePeriod     int
	Cooldown        int
	RewardCurveJSON string
	PartialCredit   *Reward
	StreakEnabled   bool
	Status          string // inactive, active
	LastCompletedAt *time.Time
	StreakCount     int
}

type Reward struct {
	CurrencyID int64
	Amount     Decimal
}
