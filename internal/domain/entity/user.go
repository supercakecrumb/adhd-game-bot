package entity

import "github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"

type User struct {
	ID          int64
	ChatID      int64  // Telegram chat ID this user belongs to
	Role        string // "admin" or "member"
	TimeZone    string
	DisplayName string
	Balances    map[int64]valueobject.Decimal // currencyID -> amount
	Preferences UserPreferences
}

type UserPreferences struct {
	EditIntervalSec int
	NotifyOnAward   bool
}
