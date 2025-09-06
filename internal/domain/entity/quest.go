package entity

import (
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type Quest struct {
	// Core identification
	ID          string
	DungeonID   string
	Title       string
	Description string

	// Categorization
	Category   string // "daily" | "weekly" | "adhoc"
	Difficulty string // "easy" | "medium" | "hard"

	// MVP Scoring Configuration
	Mode             string               // "BINARY" | "PARTIAL" | "PER_MINUTE"
	PointsAward      valueobject.Decimal  // Fixed award (BINARY) or max (PARTIAL)
	RatePointsPerMin *valueobject.Decimal // For PER_MINUTE mode
	MinMinutes       *int                 // Optional floor for PER_MINUTE
	MaxMinutes       *int                 // Optional cap for PER_MINUTE
	DailyPointsCap   *valueobject.Decimal // Optional anti-abuse limit

	// Behavioral Controls
	CooldownSec   int // Minimum seconds between completions
	StreakEnabled bool

	// Operational State
	Status          string // "active" | "paused" | "archived"
	LastCompletedAt *time.Time
	StreakCount     int
	TimeZone        string // IANA timezone for streak boundaries

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}
