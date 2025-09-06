# Phase 4: Quest Reward System

## Overview
Implement the core quest completion and point reward system. This is the heart of the gamification - users must actually earn points when they complete quests through both web and Telegram interfaces.

## Goals
- Implement point calculation based on quest mode (BINARY, PARTIAL, PER_MINUTE)
- Award points to user balance when quest is completed
- Create quest completion records for history tracking
- Handle streak bonuses and daily caps
- Ensure all operations are atomic and idempotent
- Support completion from both web and Telegram

## Tasks

### 4.1 Implement Point Calculation Logic
**File**: `internal/usecase/quest_reward_calculator.go` (new)

```go
package usecase

import (
    "errors"
    "fmt"
    "math"
    
    "github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
    "github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type QuestRewardCalculator struct{}

type RewardCalculation struct {
    BaseReward    valueobject.Decimal
    StreakBonus   valueobject.Decimal
    FinalReward   valueobject.Decimal
    StreakCount   int
    CappedByDaily bool
}

func NewQuestRewardCalculator() *QuestRewardCalculator {
    return &QuestRewardCalculator{}
}

func (c *QuestRewardCalculator) CalculateReward(quest *entity.Quest, input CompleteQuestInput, currentStreak int) (*RewardCalculation, error) {
    // Calculate base reward based on quest mode
    baseReward, err := c.calculateBaseReward(quest, input)
    if err != nil {
        return nil, err
    }
    
    // Calculate new streak count
    newStreakCount := currentStreak
    if quest.StreakEnabled {
        newStreakCount++
    }
    
    // Apply streak bonus
    streakBonus := c.calculateStreakBonus(baseReward, newStreakCount)
    finalReward := baseReward.Add(streakBonus)
    
    return &RewardCalculation{
        BaseReward:    baseReward,
        StreakBonus:   streakBonus,
        FinalReward:   finalReward,
        StreakCount:   newStreakCount,
        CappedByDaily: false,
    }, nil
}

func (c *QuestRewardCalculator) calculateBaseReward(quest *entity.Quest, input CompleteQuestInput) (valueobject.Decimal, error) {
    switch quest.Mode {
    case "BINARY":
        return quest.PointsAward, nil
        
    case "PARTIAL":
        if input.CompletionRatio == nil {
            return valueobject.NewDecimal("0"), errors.New("completion ratio required for PARTIAL mode")
        }
        
        ratio := *input.CompletionRatio
        if ratio < 0 || ratio > 1 {
            return valueobject.NewDecimal("0"), errors.New("completion ratio must be between 0 and 1")
        }
        
        ratioDecimal := valueobject.NewDecimal(fmt.Sprintf("%.4f", ratio))
        return quest.PointsAward.Multiply(ratioDecimal), nil
        
    case "PER_MINUTE":
        if input.Minutes == nil {
            return valueobject.NewDecimal("0"), errors.New("minutes required for PER_MINUTE mode")
        }
        if quest.RatePointsPerMin == nil {
            return valueobject.NewDecimal("0"), errors.New("rate_points_per_min not configured for PER_MINUTE quest")
        }
        
        minutes := *input.Minutes
        
        // Apply min/max constraints
        if quest.MinMinutes != nil && minutes < *quest.MinMinutes {
            minutes = *quest.MinMinutes
        }
        if quest.MaxMinutes != nil && minutes > *quest.MaxMinutes {
            minutes = *quest.MaxMinutes
        }
        
        minutesDecimal := valueobject.NewDecimal(fmt.Sprintf("%d", minutes))
        return quest.RatePointsPerMin.Multiply(minutesDecimal), nil
        
    default:
        return valueobject.NewDecimal("0"), fmt.Errorf("unknown quest mode: %s", quest.Mode)
    }
}

func (c *QuestRewardCalculator) calculateStreakBonus(baseReward valueobject.Decimal, streakCount int) valueobject.Decimal {
    if streakCount <= 1 {
        return valueobject.NewDecimal("0")
    }
    
    // Streak bonus: 2% per day, max 50% at 25 days
    bonusPercent := math.Min(float64(streakCount-1)*2, 50)
    
    if bonusPercent <= 0 {
        return valueobject.NewDecimal("0")
    }
    
    bonusMultiplier := valueobject.NewDecimal(fmt.Sprintf("%.2f", bonusPercent/100))
    return baseReward.Multiply(bonusMultiplier)
}

func (c *QuestRewardCalculator) ApplyDailyCap(reward valueobject.Decimal, quest *entity.Quest, todayTotal valueobject.Decimal) (valueobject.Decimal, bool) {
    if quest.DailyPointsCap == nil {
        return reward, false
    }
    
    newTotal := todayTotal.Add(reward)
    if newTotal.LessThanOrEqual(*quest.DailyPointsCap) {
        return reward, false
    }
    
    // Cap the reward
    cappedReward := quest.DailyPointsCap.Subtract(todayTotal)
    if cappedReward.LessThanOrEqual(valueobject.NewDecimal("0")) {
        return valueobject.NewDecimal("0"), true
    }
    
    return cappedReward, true
}
```

**Testing**:
- Test BINARY mode returns exact points
- Test PARTIAL mode with various completion ratios
- Test PER_MINUTE mode with min/max constraints
- Test streak bonus calculations
- Test daily cap enforcement

### 4.2 Update Quest Completion Service
**File**: `internal/usecase/quest_service.go` (modify existing CompleteQuest method)

```go
func (s *QuestService) CompleteQuest(ctx context.Context, userID int64, questID string, input CompleteQuestInput) (*entity.QuestCompletion, error) {
    // Create idempotency key if not provided
    if input.IdempotencyKey == "" {
        input.IdempotencyKey = fmt.Sprintf("quest_complete_%d_%s_%d", userID, questID, time.Now().UnixNano())
    }
    
    // Check idempotency
    idempKey := &entity.IdempotencyKey{
        Key:       input.IdempotencyKey,
        Operation: "quest_complete",
        UserID:    userID,
        Status:    "pending",
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }
    
    existingKey, err := s.idempotencyRepo.FindByKey(ctx, idempKey.Key)
    if err == nil && existingKey != nil {
        if existingKey.IsCompleted() {
            // Return existing completion
            return s.getCompletionFromIdempotencyResult(existingKey.Result)
        }
        if !existingKey.IsExpired() {
            return nil, errors.New("operation in progress")
        }
    }
    
    // Create idempotency key
    err = s.idempotencyRepo.Create(ctx, idempKey)
    if err != nil && err != ports.ErrIdempotencyKeyExists {
        return nil, err
    }
    
    var completion *entity.QuestCompletion
    
    // Execute the operation in a transaction
    err = s.txManager.WithTx(ctx, func(ctx context.Context) error {
        // Get quest and user
        quest, err := s.questRepo.GetByID(ctx, questID)
        if err != nil {
            return err
        }
        
        user, err := s.userRepo.FindByID(ctx, userID)
        if err != nil {
            return err
        }
        
        // Check if user has access to this quest (via dungeon membership)
        hasAccess, err := s.checkQuestAccess(ctx, userID, quest.DungeonID)
        if err != nil {
            return err
        }
        if !hasAccess {
            return errors.New("user does not have access to this quest")
        }
        
        // Check cooldown
        if quest.CooldownSec > 0 && quest.LastCompletedAt != nil {
            timeSinceLastCompletion := time.Since(*quest.LastCompletedAt)
            if timeSinceLastCompletion.Seconds() < float64(quest.CooldownSec) {
                return fmt.Errorf("quest is on cooldown for %d more seconds", 
                    quest.CooldownSec - int(timeSinceLastCompletion.Seconds()))
            }
        }
        
        // Calculate reward
        calculator := NewQuestRewardCalculator()
        rewardCalc, err := calculator.CalculateReward(quest, input, quest.StreakCount)
        if err != nil {
            return err
        }
        
        // Check daily cap if applicable
        finalReward := rewardCalc.FinalReward
        cappedByDaily := false
        
        if quest.DailyPointsCap != nil {
            todayCompletions, err := s.questCompletionRepo.GetTodayCompletions(ctx, userID, questID)
            if err != nil {
                return err
            }
            
            todayTotal := valueobject.NewDecimal("0")
            for _, comp := range todayCompletions {
                todayTotal = todayTotal.Add(comp.PointsAwarded)
            }
            
            finalReward, cappedByDaily = calculator.ApplyDailyCap(finalReward, quest, todayTotal)
            if finalReward.LessThanOrEqual(valueobject.NewDecimal("0")) {
                return errors.New("daily points cap reached for this quest")
            }
        }
        
        // Award points to user
        user.Balance = user.Balance.Add(finalReward)
        user.UpdatedAt = time.Now()
        err = s.userRepo.Update(ctx, user)
        if err != nil {
            return err
        }
        
        // Create completion record
        completion = &entity.QuestCompletion{
            QuestID:         questID,
            UserID:          userID,
            PointsAwarded:   finalReward,
            CompletionRatio: input.CompletionRatio,
            MinutesSpent:    input.Minutes,
            StreakCount:     rewardCalc.StreakCount,
            Notes:           input.Notes,
            CompletedAt:     time.Now(),
        }
        
        err = s.questCompletionRepo.Create(ctx, completion)
        if err != nil {
            return err
        }
        
        // Update quest state
        now := time.Now()
        quest.LastCompletedAt = &now
        quest.StreakCount = rewardCalc.StreakCount
        quest.UpdatedAt = now
        
        err = s.questRepo.Update(ctx, quest)
        if err != nil {
            return err
        }
        
        return nil
    })
    
    // Update idempotency key status
    completedAt := time.Now()
    idempKey.CompletedAt = &completedAt
    if err != nil {
        idempKey.Status = "failed"
        idempKey.Result = err.Error()
    } else {
        idempKey.Status = "completed"
        // Store completion ID for retrieval
        idempKey.Result = fmt.Sprintf("completion_id:%d", completion.ID)
    }
    
    updateErr := s.idempotencyRepo.Update(ctx, idempKey)
    if updateErr != nil {
        // Log error but don't fail the operation
        log.Printf("Failed to update idempotency key: %v", updateErr)
    }
    
    return completion, err
}

func (s *QuestService) checkQuestAccess(ctx context.Context, userID int64, dungeonID string) (bool, error) {
    // Check if user is a member of the dungeon
    members, err := s.dungeonMemberRepo.GetDungeonMembers(ctx, dungeonID)
    if err != nil {
        return false, err
    }
    
    for _, member := range members {
        if member.UserID == userID {
            return true, nil
        }
    }
    
    return false, nil
}

func (s *QuestService) getCompletionFromIdempotencyResult(result string) (*entity.QuestCompletion, error) {
    // Parse completion ID from result
    if strings.HasPrefix(result, "completion_id:") {
        completionIDStr := strings.TrimPrefix(result, "completion_id:")
        completionID, err := strconv.ParseInt(completionIDStr, 10, 64)
        if err != nil {
            return nil, err
        }
        
        // Get completion by ID
        return s.questCompletionRepo.GetByID(context.Background(), completionID)
    }
    
    return nil, errors.New("invalid idempotency result format")
}
```

**Testing**:
- Test point awarding for each quest mode
- Test streak bonus application
- Test daily cap enforcement
- Test cooldown enforcement
- Test quest access control
- Test idempotency prevents double rewards
- Test transaction rollback on errors

### 4.3 Update Quest Completion Repository
**File**: `internal/infra/postgres/quest_completion_repository.go` (modify existing)

```go
func (r *QuestCompletionRepository) GetTodayCompletions(ctx context.Context, userID int64, questID string) ([]*entity.QuestCompletion, error) {
    query := `
        SELECT id, quest_id, user_id, points_awarded, completion_ratio, minutes_spent, 
               streak_count, notes, completed_at
        FROM quest_completions 
        WHERE user_id = $1 AND quest_id = $2 
        AND completed_at >= CURRENT_DATE AT TIME ZONE 'UTC'
        AND completed_at < (CURRENT_DATE + INTERVAL '1 day') AT TIME ZONE 'UTC'
        ORDER BY completed_at DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, userID, questID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var completions []*entity.QuestCompletion
    for rows.Next() {
        completion := &entity.QuestCompletion{}
        err := rows.Scan(
            &completion.ID,
            &completion.QuestID,
            &completion.UserID,
            &completion.PointsAwarded,
            &completion.CompletionRatio,
            &completion.MinutesSpent,
            &completion.StreakCount,
            &completion.Notes,
            &completion.CompletedAt,
        )
        if err != nil {
            return nil, err
        }
        completions = append(completions, completion)
    }
    
    return completions, nil
}

func (r *QuestCompletionRepository) GetByID(ctx context.Context, id int64) (*entity.QuestCompletion, error) {
    query := `
        SELECT id, quest_id, user_id, points_awarded, completion_ratio, minutes_spent, 
               streak_count, notes, completed_at
        FROM quest_completions 
        WHERE id = $1
    `
    
    completion := &entity.QuestCompletion{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &completion.ID,
        &completion.QuestID,
        &completion.UserID,
        &completion.PointsAwarded,
        &completion.CompletionRatio,
        &completion.MinutesSpent,
        &completion.StreakCount,
        &completion.Notes,
        &completion.CompletedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    return completion, nil
}

func (r *QuestCompletionRepository) GetUserStats(ctx context.Context, userID int64) (*UserQuestStats, error) {
    query := `
        SELECT 
            COUNT(*) as total_completions,
            COALESCE(SUM(points_awarded), 0) as total_points_earned,
            COALESCE(MAX(streak_count), 0) as max_streak,
            COUNT(DISTINCT quest_id) as unique_quests_completed
        FROM quest_completions 
        WHERE user_id = $1
    `
    
    stats := &UserQuestStats{}
    err := r.db.QueryRowContext(ctx, query, userID).Scan(
        &stats.TotalCompletions,
        &stats.TotalPointsEarned,
        &stats.MaxStreak,
        &stats.UniqueQuestsCompleted,
    )
    
    return stats, err
}

func (r *QuestCompletionRepository) GetRecentCompletions(ctx context.Context, userID int64, limit int) ([]*entity.QuestCompletion, error) {
    query := `
        SELECT qc.id, qc.quest_id, qc.user_id, qc.points_awarded, qc.completion_ratio, 
               qc.minutes_spent, qc.streak_count, qc.notes, qc.completed_at,
               q.title as quest_title
        FROM quest_completions qc
        JOIN quests q ON qc.quest_id = q.id
        WHERE qc.user_id = $1 
        ORDER BY qc.completed_at DESC
        LIMIT $2
    `
    
    rows, err := r.db.QueryContext(ctx, query, userID, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var completions []*entity.QuestCompletion
    for rows.Next() {
        completion := &entity.QuestCompletion{}
        var questTitle string
        
        err := rows.Scan(
            &completion.ID,
            &completion.QuestID,
            &completion.UserID,
            &completion.PointsAwarded,
            &completion.CompletionRatio,
            &completion.MinutesSpent,
            &completion.StreakCount,
            &completion.Notes,
            &completion.CompletedAt,
            &questTitle,
        )
        if err != nil {
            return nil, err
        }
        
        // Add quest title for display purposes
        completion.QuestTitle = questTitle
        completions = append(completions, completion)
    }
    
    return completions, nil
}
```

**Testing**:
- Test today's completions filtering with timezone handling
- Test user stats aggregation
- Test recent completions with quest titles
- Test completion retrieval by ID

### 4.4 Add Missing Entity Fields
**File**: `internal/domain/entity/quest_completion.go` (modify existing)

```go
type QuestCompletion struct {
    ID              int64
    QuestID         string
    UserID          int64
    PointsAwarded   valueobject.Decimal
    CompletionRatio *float64 // For PARTIAL mode (0.0 to 1.0)
    MinutesSpent    *int     // For PER_MINUTE mode
    StreakCount     int
    Notes           string   // Optional completion notes
    CompletedAt     time.Time
    
    // Derived fields (not stored in DB)
    QuestTitle      string `json:"quest_title,omitempty"`
}

type UserQuestStats struct {
    TotalCompletions       int                     `json:"total_completions"`
    TotalPointsEarned      valueobject.Decimal     `json:"total_points_earned"`
    MaxStreak              int                     `json:"max_streak"`
    UniqueQuestsCompleted  int                     `json:"unique_quests_completed"`
}

// Add to CompleteQuestInput
type CompleteQuestInput struct {
    IdempotencyKey  string   `json:"idempotency_key,omitempty"`
    CompletionRatio *float64 `json:"completion_ratio,omitempty"` // For PARTIAL mode
    Minutes         *int     `json:"minutes,omitempty"`          // For PER_MINUTE mode
    Notes           string   `json:"notes,omitempty"`            // Optional notes
}
```

**Testing**:
- Test entity creation and validation
- Test decimal precision for points
- Test optional fields handling
- Test JSON serialization

## Testing Strategy

### 4.5 Unit Tests
**File**: `internal/usecase/quest_reward_calculator_test.go`

```go
func TestQuestRewardCalculator_BINARY(t *testing.T) {
    calculator := NewQuestRewardCalculator()
    
    quest := &entity.Quest{
        Mode:        "BINARY",
        PointsAward: valueobject.NewDecimal("10.00"),
    }
    
    input := CompleteQuestInput{}
    
    result, err := calculator.CalculateReward(quest, input, 0)
    assert.NoError(t, err)
    assert.Equal(t, "10.00", result.BaseReward.String())
    assert.Equal(t, "0.00", result.StreakBonus.String())
    assert.Equal(t, "10.00", result.FinalReward.String())
}

func TestQuestRewardCalculator_PARTIAL(t *testing.T) {
    calculator := NewQuestRewardCalculator()
    
    quest := &entity.Quest{
        Mode:        "PARTIAL",
        PointsAward: valueobject.NewDecimal("20.00"),
    }
    
    ratio := 0.75
    input := CompleteQuestInput{
        CompletionRatio: &ratio,
    }
    
    result, err := calculator.CalculateReward(quest, input, 0)
    assert.NoError(t, err)
    assert.Equal(t, "15.00", result.BaseReward.String())
}

func TestQuestRewardCalculator_PER_MINUTE(t *testing.T) {
    calculator := NewQuestRewardCalculator()
    
    ratePerMin := valueobject.NewDecimal("1.50")
    minMinutes := 10
    maxMinutes := 60
    
    quest := &entity.Quest{
        Mode:             "PER_MINUTE",
        RatePointsPerMin: &ratePerMin,
        MinMinutes:       &minMinutes,
        MaxMinutes:       &maxMinutes,
    }
    
    minutes := 30
    input := CompleteQuestInput{
        Minutes: &minutes,
    }
    
    result, err := calculator.CalculateReward(quest, input, 0)
    assert.NoError(t, err)
    assert.Equal(t, "45.00", result.BaseReward.String()) // 30 * 1.50
}

func TestStreakBonus(t *testing.T) {
    calculator := NewQuestRewardCalculator()
    
    baseReward := valueobject.NewDecimal("10.00")
    
    // 5-day streak = 8% bonus (4 * 2%)
    bonus := calculator.calculateStreakBonus(baseReward, 5)
    assert.Equal(t, "0.80", bonus.String()) // 10.00 * 0.08
}

func TestDailyCap(t *testing.T) {
    calculator := NewQuestRewardCalculator()
    
    dailyCap := valueobject.NewDecimal("50.00")
    quest := &entity.Quest{
        DailyPointsCap: &dailyCap,
    }
    
    reward := valueobject.NewDecimal("20.00")
    todayTotal := valueobject.NewDecimal("40.00")
    
    cappedReward, wasCapped := calculator.ApplyDailyCap(reward, quest, todayTotal)
    assert.True(t, wasCapped)
    assert.Equal(t, "10.00", cappedReward.String()) // 50 - 40 = 10
}
```

### 4.6 Integration Tests
**File**: `test/integration/quest_completion_test.go`

```go
func TestCompleteQuestFullWorkflow(t *testing.T) {
    // Setup test database and services
    // Create user with initial balance
    // Create dungeon and add user as member
    // Create quest
    // Complete quest
    // Verify points awarded
    // Verify completion record created
    // Verify quest state updated
    // Verify user balance increased
}

func TestQuestCompletionIdempotency(t *testing.T) {
    // Complete quest with idempotency key
    // Try to complete again with same key
    // Verify no double reward
    // Verify same completion returned
}

func TestQuestAccessControl(t *testing.T) {
    // Create quest in dungeon
    // Try to complete as non-member
    // Verify access denied
    // Add user to dungeon
    // Verify completion succeeds
}
```

## Verification Checklist

- [ ] Point calculation works for all quest modes (BINARY, PARTIAL, PER_MINUTE)
- [ ] Points are actually awarded to user balance
- [ ] Quest completion records are created with all details
- [ ] Streak bonuses are applied correctly
- [ ] Daily caps are enforced properly
- [ ] Cooldowns prevent rapid completion
- [ ] Quest access control works (dungeon membership)
- [ ] All operations are atomic (transaction rollback works)
- [ ] Idempotency prevents double rewards
- [ ] Unit tests cover all calculation logic
- [ ] Integration tests verify full workflow
- [ ] Both web and Telegram completion work

## Files to Create/Modify

1. `internal/usecase/quest_reward_calculator.go` - New reward calculation logic
2. `internal/usecase/quest_service.go` - Update CompleteQuest method
3. `internal/infra/postgres/quest_completion_repository.go` - Add missing methods
4. `internal/domain/entity/quest_completion.go` - Update entity with new fields
5. `internal/usecase/quest_reward_calculator_test.go` - Unit tests
6. `test/integration/quest_completion_test.go` - Integration tests

## Success Criteria

✅ Users earn points when completing quests via web or Telegram
✅ All quest modes (BINARY, PARTIAL, PER_MINUTE) work correctly
✅ Streak bonuses increase rewards appropriately
✅ Daily caps prevent abuse while allowing fair usage
✅ Cooldowns work as expected
✅ Quest access control prevents unauthorized completions
✅ All operations are atomic and idempotent
✅ Comprehensive test coverage validates all scenarios
✅ Performance is acceptable for concurrent completions

## Next Phase
Once quest rewards are working, move to **Phase 5: MVP Polish & Deployment** to add final touches and prepare for production deployment.