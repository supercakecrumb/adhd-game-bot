# Dungeons & Quests Rename Implementation Plan

## Overview

This document outlines the complete plan for renaming "Task" to "Quest" throughout the codebase and introducing the "Dungeon" concept as a group container. This is a fresh start approach with no data migration from the old system.

## Branch Strategy

Create feature branch: `feat/dungeons-and-quests-rename`

## Implementation Order & Atomic Commits

### Stage 1: Domain Layer Rename

#### Commit 1: `refactor(domain): Task → Quest entities and interfaces`

**1. Create new file `internal/domain/entity/quest.go`:**

```go
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
    Mode             string                  // "BINARY" | "PARTIAL" | "PER_MINUTE"
    PointsAward      valueobject.Decimal    // Fixed award (BINARY) or max (PARTIAL)
    RatePointsPerMin *valueobject.Decimal   // For PER_MINUTE mode
    MinMinutes       *int                   // Optional floor for PER_MINUTE
    MaxMinutes       *int                   // Optional cap for PER_MINUTE
    DailyPointsCap   *valueobject.Decimal   // Optional anti-abuse limit
    
    // Behavioral Controls
    CooldownSec   int  // Minimum seconds between completions
    StreakEnabled bool
    
    // Operational State
    Status          string     // "active" | "paused" | "archived"
    LastCompletedAt *time.Time
    StreakCount     int
    TimeZone        string // IANA timezone for streak boundaries
    
    // Timestamps
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**2. Create `internal/domain/entity/dungeon.go`:**

```go
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
```

**3. Create `internal/domain/entity/quest_completion.go`:**

```go
package entity

import (
    "time"
    "github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type QuestCompletion struct {
    ID             string
    QuestID        string
    UserID         int64
    DungeonID      string
    SubmittedAt    time.Time
    
    // Scoring inputs
    CompletionRatio *float64 // 0..1 for PARTIAL mode
    Minutes         *int     // for PER_MINUTE mode
    
    // Outcome
    AwardedPoints  valueobject.Decimal
    IdempotencyKey string
}
```

**4. Update `internal/domain/entity/user.go`:**

```go
type User struct {
    ID        int64
    Username  string
    Balance   valueobject.Decimal
    TimeZone  string    // Renamed from Timezone
    ChatID    *int64    // Legacy field, kept temporarily
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**5. Update `internal/domain/entity/shop_item.go`:**

```go
type ShopItem struct {
    ID          int64
    DungeonID   *string // nil for global, string for scoped
    Code        string
    Name        string
    Description string
    Price       valueobject.Decimal
    Category    string
    IsActive    bool
    Stock       *int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Purchase struct {
    ID          int64
    UserID      int64
    ItemID      int64
    DungeonID   string // Denormalized for analytics
    ItemName    string
    ItemPrice   valueobject.Decimal
    Quantity    int
    TotalCost   valueobject.Decimal
    Status      string
    PurchasedAt time.Time
}
```

**6. Delete `internal/domain/entity/task.go`**

### Stage 2: Repository Layer Rename

#### Commit 2: `refactor(ports): TaskRepository → QuestRepository`

Update `internal/ports/repositories.go`:

```go
// Replace TaskRepository with:
type QuestRepository interface {
    Create(ctx context.Context, quest *entity.Quest) error
    GetByID(ctx context.Context, questID string) (*entity.Quest, error)
    ListByDungeon(ctx context.Context, dungeonID string) ([]*entity.Quest, error)
    Update(ctx context.Context, quest *entity.Quest) error
    Delete(ctx context.Context, questID string) error
}

type QuestCompletionRepository interface {
    Insert(ctx context.Context, completion *entity.QuestCompletion) error
    LastForUser(ctx context.Context, userID int64, questID string) (*entity.QuestCompletion, error)
    SumAwardedForUserOnDay(ctx context.Context, userID int64, questID string, day time.Time, tz string) (valueobject.Decimal, error)
}

type DungeonRepository interface {
    Create(ctx context.Context, dungeon *entity.Dungeon) error
    GetByID(ctx context.Context, dungeonID string) (*entity.Dungeon, error)
    ListByAdmin(ctx context.Context, userID int64) ([]*entity.Dungeon, error)
}

type DungeonMemberRepository interface {
    Add(ctx context.Context, dungeonID string, userID int64) error
    ListUsers(ctx context.Context, dungeonID string) ([]int64, error)
    IsMember(ctx context.Context, dungeonID string, userID int64) (bool, error)
}
```

#### Commit 3: `refactor(infra): rename repository implementations`

1. Rename files:
   - `internal/infra/postgres/task_repository.go` → `quest_repository.go`
   - `internal/infra/inmemory/task_repository.go` → `quest_repository.go` (if exists)

2. Update implementations to match new interfaces
3. Update struct names and method signatures

### Stage 3: Service Layer Rename

#### Commit 4: `refactor(usecase): TaskService → QuestService`

**1. Rename `internal/usecase/task_service.go` → `quest_service.go`**

**2. Update service structure:**

```go
type QuestService struct {
    questRepo       ports.QuestRepository
    userRepo        ports.UserRepository
    dungeonRepo     ports.DungeonRepository
    memberRepo      ports.DungeonMemberRepository
    completionRepo  ports.QuestCompletionRepository
    uuidGen         ports.UUIDGenerator
    idempotencyRepo ports.IdempotencyRepository
    txManager       ports.TxManager
}

func NewQuestService(
    questRepo ports.QuestRepository,
    userRepo ports.UserRepository,
    dungeonRepo ports.DungeonRepository,
    memberRepo ports.DungeonMemberRepository,
    completionRepo ports.QuestCompletionRepository,
    uuidGen ports.UUIDGenerator,
    idempotencyRepo ports.IdempotencyRepository,
    txManager ports.TxManager,
) *QuestService {
    return &QuestService{
        questRepo:       questRepo,
        userRepo:        userRepo,
        dungeonRepo:     dungeonRepo,
        memberRepo:      memberRepo,
        completionRepo:  completionRepo,
        uuidGen:         uuidGen,
        idempotencyRepo: idempotencyRepo,
        txManager:       txManager,
    }
}
```

**3. Update method signatures:**

```go
func (s *QuestService) CreateQuest(ctx context.Context, adminUserID int64, dungeonID string, input CreateQuestInput) (*entity.Quest, error)
func (s *QuestService) ListQuests(ctx context.Context, userID int64, dungeonID string) ([]*entity.Quest, error)
func (s *QuestService) CompleteQuest(ctx context.Context, userID int64, questID string, input CompleteQuestInput) (*entity.QuestCompletion, error)

type CreateQuestInput struct {
    Title            string
    Description      string
    Category         string
    Difficulty       string
    Mode             string
    PointsAward      valueobject.Decimal
    RatePointsPerMin *valueobject.Decimal
    MinMinutes       *int
    MaxMinutes       *int
    DailyPointsCap   *valueobject.Decimal
    CooldownSec      int
    StreakEnabled    bool
    Status           string
}

type CompleteQuestInput struct {
    IdempotencyKey  string
    CompletionRatio *float64 // For PARTIAL mode
    Minutes         *int     // For PER_MINUTE mode
}
```

**4. Rename test files:**
   - `internal/usecase/task_service_test.go` → `quest_service_test.go`
   - `internal/usecase/task_service_timezone_test.go` → `quest_service_timezone_test.go`

### Stage 4: HTTP Layer Rename

#### Commit 5: `refactor(infra/http): /tasks → /quests route names`

**1. Rename `internal/infra/http/task_handler.go` → `quest_handler.go`**

**2. Update `internal/infra/http/server.go`:**

```go
type Server struct {
    Router         *chi.Mux
    QuestService   *usecase.QuestService
    DungeonService *usecase.DungeonService
}

func (s *Server) setupRoutes() {
    s.Router.Route("/api/v1", func(r chi.Router) {
        // Quest routes
        r.Route("/dungeons/{dungeonId}/quests", func(r chi.Router) {
            r.Get("/", s.listQuestsHandler)
        })
        
        r.Route("/quests/{questId}", func(r chi.Router) {
            r.Post("/complete", s.completeQuestHandler)
        })
        
        // Admin routes
        r.Route("/dungeons", func(r chi.Router) {
            r.Post("/", s.createDungeonHandler)
            r.Route("/{dungeonId}", func(r chi.Router) {
                r.Post("/quests", s.createQuestHandler)
            })
        })
    })
}
```

**3. Update DTOs:**

```go
type QuestResponse struct {
    ID               string  `json:"id"`
    Title            string  `json:"title"`
    Description      string  `json:"description"`
    Category         string  `json:"category"`
    Difficulty       string  `json:"difficulty"`
    Mode             string  `json:"mode"`
    PointsAward      string  `json:"points_award"`
    RatePointsPerMin *string `json:"rate_points_per_min,omitempty"`
    MinMinutes       *int    `json:"min_minutes,omitempty"`
    MaxMinutes       *int    `json:"max_minutes,omitempty"`
    DailyPointsCap   *string `json:"daily_points_cap,omitempty"`
    CooldownSec      int     `json:"cooldown_sec"`
    StreakEnabled    bool    `json:"streak_enabled"`
    Status           string  `json:"status"`
}

type CreateQuestRequest struct {
    Title            string  `json:"title"`
    Description      string  `json:"description"`
    Category         string  `json:"category"`
    Difficulty       string  `json:"difficulty"`
    Mode             string  `json:"mode"`
    PointsAward      string  `json:"points_award"`
    RatePointsPerMin *string `json:"rate_points_per_min,omitempty"`
    MinMinutes       *int    `json:"min_minutes,omitempty"`
    MaxMinutes       *int    `json:"max_minutes,omitempty"`
    DailyPointsCap   *string `json:"daily_points_cap,omitempty"`
    CooldownSec      *int    `json:"cooldown_sec,omitempty"`
    StreakEnabled    *bool   `json:"streak_enabled,omitempty"`
    Status           *string `json:"status,omitempty"`
}

type CompleteQuestRequest struct {
    IdempotencyKey  string   `json:"idempotency_key"`
    CompletionRatio *float64 `json:"completion_ratio,omitempty"`
    Minutes         *int     `json:"minutes,omitempty"`
}

type CompleteQuestResponse struct {
    AwardedPoints string  `json:"awarded_points"`
    SubmittedAt   string  `json:"submitted_at"`
    StreakCount   *int    `json:"streak_count,omitempty"`
}
```

### Stage 5: Add New Domain Services

#### Commit 6: `feat(domain): add Dungeon, DungeonMember, QuestCompletion (types only)`

**1. Create `internal/usecase/dungeon_service.go`:**

```go
package usecase

import (
    "context"
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
    // Implementation stub
    return nil, nil
}

func (s *DungeonService) AddMember(ctx context.Context, adminUserID int64, dungeonID string, userID int64) error {
    // Implementation stub
    return nil
}

func (s *DungeonService) ListMembers(ctx context.Context, adminUserID int64, dungeonID string) ([]int64, error) {
    // Implementation stub
    return nil, nil
}
```

**2. Create repository implementations (stubs):**
   - `internal/infra/postgres/dungeon_repository.go`
   - `internal/infra/postgres/dungeon_member_repository.go`
   - `internal/infra/postgres/quest_completion_repository.go`

### Stage 6: Update Tests & Fixtures

#### Commit 7: `refactor(infra/tg): replace chat_id references with Dungeon plumbing`

This commit updates Telegram integration touchpoints (names only, no behavior changes).

#### Commit 8: `refactor(test): update fixtures and tests for quests`

**1. Rename `test/fixtures/builders/task_builder.go` → `quest_builder.go`:**

```go
package builders

import (
    "time"
    "github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
    "github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type QuestBuilder struct {
    BaseBuilder[entity.Quest]
}

func NewQuestBuilder() *QuestBuilder {
    return &QuestBuilder{
        BaseBuilder: BaseBuilder[entity.Quest]{
            Builder: *NewBuilder[entity.Quest](),
        },
    }
}

func (b *QuestBuilder) WithDefaults() *QuestBuilder {
    pointsAward := valueobject.NewFromInt(100)
    return b.
        WithID("quest-1").
        WithDungeonID("dungeon-1").
        WithTitle("Sample Quest").
        WithDescription("Complete this sample quest").
        WithCategory("daily").
        WithDifficulty("medium").
        WithMode("BINARY").
        WithPointsAward(pointsAward).
        WithCooldownSec(3600).
        WithStreakEnabled(true).
        WithStatus("active").
        WithTimeZone("UTC")
}

func (b *QuestBuilder) WithID(id string) *QuestBuilder {
    b.With(func(q *entity.Quest) {
        q.ID = id
    })
    return b
}

func (b *QuestBuilder) WithDungeonID(dungeonID string) *QuestBuilder {
    b.With(func(q *entity.Quest) {
        q.DungeonID = dungeonID
    })
    return b
}

// ... other builder methods
```

**2. Update test files:**
   - Replace all references to `Task` with `Quest`
   - Update test data to use new field names
   - Update assertions to check new fields

### Stage 7: Documentation Updates

#### Commit 9: `docs: update names in README/tech overview`

**1. Update `README.md`:**
   - Replace all instances of "Task" with "Quest"
   - Add section explaining Dungeons:
     ```markdown
     ## Core Concepts
     
     ### Dungeons
     A Dungeon is a group container with one admin (Dungeon Master) and one or more members. 
     Each dungeon can optionally be linked to a Telegram group via TelegramChatID.
     
     ### Quests
     Quests are activities that users can complete to earn points. Each quest belongs to a 
     specific dungeon and supports three scoring modes: BINARY, PARTIAL, and PER_MINUTE.
     ```

**2. Update `TECHNICAL_OVERVIEW.md`:**
   - Update architecture diagrams
   - Update domain model section
   - Add Dungeon/Quest relationship diagram

**3. Rename and update `docs/api/services/task_service.md` → `quest_service.md`**

**4. Update API documentation with new endpoints**

### Stage 8: Database Migration

#### Commit 10: `feat(db): add migration for dungeons and quests`

Create `internal/infra/postgres/migrations/007_dungeons_and_quests.sql`:

```sql
BEGIN;

-- Create dungeons table
CREATE TABLE dungeons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    admin_user_id BIGINT NOT NULL REFERENCES users(id),
    telegram_chat_id BIGINT UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create dungeon_members table
CREATE TABLE dungeon_members (
    dungeon_id UUID NOT NULL REFERENCES dungeons(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (dungeon_id, user_id)
);

-- Create quests table (fresh start, no migration from tasks)
CREATE TABLE quests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dungeon_id UUID NOT NULL REFERENCES dungeons(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(10) NOT NULL CHECK (category IN ('daily', 'weekly', 'adhoc')),
    difficulty VARCHAR(10) NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
    
    -- Scoring configuration
    mode VARCHAR(20) NOT NULL DEFAULT 'BINARY' CHECK (mode IN ('BINARY', 'PARTIAL', 'PER_MINUTE')),
    points_award NUMERIC(20,8) NOT NULL,
    rate_points_per_min NUMERIC(20,8),
    min_minutes INTEGER,
    max_minutes INTEGER,
    daily_points_cap NUMERIC(20,8),
    
    -- Behavioral controls
    cooldown_sec INTEGER NOT NULL DEFAULT 0,
    streak_enabled BOOLEAN NOT NULL DEFAULT true,
    
    -- State
    status VARCHAR(10) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paused', 'archived')),
    last_completed_at TIMESTAMP WITH TIME ZONE,
    streak_count INTEGER NOT NULL DEFAULT 0,
    time_zone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create quest_completions table
CREATE TABLE quest_completions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quest_id UUID NOT NULL REFERENCES quests(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    dungeon_id UUID NOT NULL REFERENCES dungeons(id),
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Scoring inputs
    completion_ratio NUMERIC(3,2) CHECK (completion_ratio >= 0 AND completion_ratio <= 1),
    minutes INTEGER CHECK (minutes > 0),
    
    -- Outcome
    awarded_points NUMERIC(20,8) NOT NULL,
    idempotency_key VARCHAR(255) NOT NULL UNIQUE,
    
    CONSTRAINT valid_scoring_input CHECK (
        (completion_ratio IS NOT NULL AND minutes IS NULL) OR
        (completion_ratio IS NULL AND minutes IS NOT NULL) OR
        (completion_ratio IS NULL AND minutes IS NULL)
    )
);

-- Update shop_items to use dungeon scope
ALTER TABLE shop_items ADD COLUMN dungeon_id UUID REFERENCES dungeons(id);

-- Update users table
ALTER TABLE users RENAME COLUMN timezone TO time_zone;
ALTER TABLE users ADD COLUMN chat_id BIGINT; -- Legacy field

-- Create indices
CREATE INDEX idx_quests_dungeon ON quests(dungeon_id);
CREATE INDEX idx_quest_completions_user ON quest_completions(user_id);
CREATE INDEX idx_quest_completions_quest ON quest_completions(quest_id);
CREATE INDEX idx_dungeon_members_user ON dungeon_members(user_id);

COMMIT;
```

## API Endpoint Changes

### Before (Old System)
```
POST   /api/tasks
GET    /api/tasks/{taskID}
PUT    /api/tasks/{taskID}
POST   /api/tasks/{taskID}/complete
GET    /api/users/{userID}/tasks
```

### After (New System)

#### User-facing Endpoints
```
GET    /api/v1/dungeons/{dungeonId}/quests
POST   /api/v1/quests/{questId}/complete
```

#### Admin Endpoints
```
POST   /api/v1/dungeons
POST   /api/v1/dungeons/{dungeonId}/quests
POST   /api/v1/dungeons/{dungeonId}/members
GET    /api/v1/dungeons/{dungeonId}/members
```

## JSON Response Examples

### Quest Response
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Morning Exercise",
  "description": "Complete 30 minutes of exercise",
  "category": "daily",
  "difficulty": "medium",
  "mode": "PER_MINUTE",
  "points_award": "100.00",
  "rate_points_per_min": "3.33",
  "min_minutes": 10,
  "max_minutes": 60,
  "daily_points_cap": "200.00",
  "cooldown_sec": 3600,
  "streak_enabled": true,
  "status": "active"
}
```

### Complete Quest Request
```json
{
  "idempotency_key": "user123-quest456-1234567890",
  "minutes": 45
}
```

### Complete Quest Response
```json
{
  "awarded_points": "149.85",
  "submitted_at": "2025-01-06T15:30:00Z",
  "streak_count": 5
}
```

## Acceptance Checklist

- [ ] All `Task` types/interfaces renamed to `Quest`
- [ ] New domain types exist: `Dungeon`, `DungeonMember`, `QuestCompletion`
- [ ] `Quest` contains all MVP scoring fields
- [ ] Repository interfaces updated with new methods
- [ ] Service layer uses new names and signatures
- [ ] API routes use `/quests` and `/dungeons`
- [ ] JSON responses use new field names
- [ ] Tests updated with new builders
- [ ] Documentation reflects new terminology
- [ ] Database migration ready (fresh start)
- [ ] Build passes with all tests green

## Notes

1. This is a fresh start - no data migration from the old task system
2. The old `ChatID` concept is replaced with `DungeonID`
3. Users must be members of a dungeon to see/complete its quests
4. Each dungeon has one admin who can create quests
5. Shop items can be scoped to specific dungeons or be global