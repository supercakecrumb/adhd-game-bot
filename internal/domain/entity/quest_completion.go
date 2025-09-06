package entity

import (
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type QuestCompletion struct {
	ID          string
	QuestID     string
	UserID      int64
	DungeonID   string
	SubmittedAt time.Time

	// Scoring inputs
	CompletionRatio *float64 // 0..1 for PARTIAL mode
	Minutes         *int     // for PER_MINUTE mode

	// Outcome
	AwardedPoints  valueobject.Decimal
	IdempotencyKey string
}
