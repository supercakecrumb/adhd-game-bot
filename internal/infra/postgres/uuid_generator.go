package postgres

import (
	"github.com/google/uuid"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type UUIDGenerator struct{}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (u *UUIDGenerator) New() string {
	return uuid.New().String()
}

// Compile-time check that UUIDGenerator implements ports.UUIDGenerator
var _ ports.UUIDGenerator = (*UUIDGenerator)(nil)
