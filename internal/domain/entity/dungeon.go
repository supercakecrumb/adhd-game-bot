package entity

import "time"

type Dungeon struct {
	ID             string
	Title          string
	AdminUserID    int64
	TelegramChatID *int64 // Optional link to Telegram group
	CreatedAt      time.Time
}

type DungeonMember struct {
	DungeonID string
	UserID    int64
	JoinedAt  time.Time
}
