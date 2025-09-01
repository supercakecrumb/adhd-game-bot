package entity

import (
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type User struct {
	ID        int64
	ChatID    int64  // Chat/group this user belongs to
	Username  string // User's display name
	Balance   valueobject.Decimal
	Timezone  string // IANA timezone (e.g. "America/New_York")
	CreatedAt time.Time
	UpdatedAt time.Time
}
