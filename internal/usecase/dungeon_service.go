package usecase

import (
	"context"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type DungeonService struct {
	dungeonRepo ports.DungeonRepository
	memberRepo  ports.DungeonMemberRepository
	userRepo    ports.UserRepository
	uuidGen     ports.UUIDGenerator
	txManager   ports.TxManager
}

func NewDungeonService(
	dungeonRepo ports.DungeonRepository,
	memberRepo ports.DungeonMemberRepository,
	userRepo ports.UserRepository,
	uuidGen ports.UUIDGenerator,
	txManager ports.TxManager,
) *DungeonService {
	return &DungeonService{
		dungeonRepo: dungeonRepo,
		memberRepo:  memberRepo,
		userRepo:    userRepo,
		uuidGen:     uuidGen,
		txManager:   txManager,
	}
}

func (s *DungeonService) CreateDungeon(ctx context.Context, adminUserID int64, title string, telegramChatID *int64) (*entity.Dungeon, error) {
	// Verify admin user exists
	_, err := s.userRepo.FindByID(ctx, adminUserID)
	if err != nil {
		return nil, err
	}

	// Create dungeon entity
	dungeon := &entity.Dungeon{
		ID:             s.uuidGen.New(),
		Title:          title,
		AdminUserID:    adminUserID,
		TelegramChatID: telegramChatID,
		CreatedAt:      time.Now(),
	}

	// Create dungeon
	err = s.dungeonRepo.Create(ctx, dungeon)
	if err != nil {
		return nil, err
	}

	// Add admin as first member
	err = s.memberRepo.Add(ctx, dungeon.ID, adminUserID)
	if err != nil {
		return nil, err
	}

	return dungeon, nil
}

func (s *DungeonService) AddMember(ctx context.Context, adminUserID int64, dungeonID string, userID int64) error {
	// Verify admin user exists and is the dungeon admin
	dungeon, err := s.dungeonRepo.GetByID(ctx, dungeonID)
	if err != nil {
		return err
	}

	if dungeon.AdminUserID != adminUserID {
		return ports.ErrUserNotFound // Using this error for simplicity
	}

	// Verify user exists
	_, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Add member
	return s.memberRepo.Add(ctx, dungeonID, userID)
}

func (s *DungeonService) ListMembers(ctx context.Context, adminUserID int64, dungeonID string) ([]int64, error) {
	// Verify admin user exists and is the dungeon admin
	dungeon, err := s.dungeonRepo.GetByID(ctx, dungeonID)
	if err != nil {
		return nil, err
	}

	if dungeon.AdminUserID != adminUserID {
		return nil, ports.ErrUserNotFound // Using this error for simplicity
	}

	// List members
	return s.memberRepo.ListUsers(ctx, dungeonID)
}
