package entity

import "github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"

type User struct {
	ID          int64
	Role        string // "admin" or "member"
	TimeZone    string
	DisplayName string
	Balances    map[string]valueobject.Decimal // currency -> amount
	Preferences UserPreferences
}

type UserPreferences struct {
	EditIntervalSec int
	NotifyOnAward   bool
}
