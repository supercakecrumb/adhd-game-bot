# Phase 2: Quest Reward System

## Overview
Implement the core quest completion and point reward system. This is the heart of the gamification - users must actually earn points when they complete quests.

## Goals
- Implement point calculation based on quest mode (BINARY, PARTIAL, PER_MINUTE)
- Award points to user balance when quest is completed
- Create quest completion records for history tracking
- Handle streak bonuses and daily caps
- Ensure all operations are atomic and idempotent

## Tasks

### 2.1 Implement Point Calculation Logic
**File**: `internal/usecase/quest_reward_calculator.go`

Create a dedicated service for calculating quest rewards:

```go
type QuestRewardCalculator struct{}

func (c *QuestRewardCalculator) CalculateReward(quest *entity.Quest, input CompleteQuestInput) (valueobject.Decimal, error) {
    switch quest.Mode {
    case "BINARY":
        return quest.PointsAward, nil
    case "PARTIAL":
        if input.CompletionRatio == nil {
            return valueobject.NewDecimal("0"), errors.New("completion ratio required for PARTIAL mode")
        }
        ratio := valueobject.NewDecimal(fmt.Sprintf("%.4f", *input.CompletionRatio))
        return quest.PointsAward.Multiply(ratio), nil
    case "PER_MINUTE":
        if input.Minutes == nil {
            return valueobject.NewDecimal("0"), errors.New("minutes required for PER_MINUTE mode")
        }
        // Apply min/max constraints
        minutes := *input.Minutes
        if quest.MinMinutes != nil && minutes < *quest.MinMinutes {
            minutes = *quest.MinMinutes
        }
        if quest.MaxMinutes != nil && minutes > *quest.MaxMinutes {
            minutes = *quest.MaxMinutes
        }
        
        minutesDecimal := valueobject.NewDecimal(fmt.Sprintf("%d", minutes))
        return quest.RatePointsPerMin.Multiply(minutesDecimal), nil
    default:
        return valueobject.NewDecimal("0"), errors.New("unknown quest mode")
    }
}

func (c *QuestRewardCalculator) ApplyStreakBonus(baseReward valueobject.Decimal, streakCount int) valueobject.Decimal {
    // Apply streak multiplier: 1% bonus per streak day, max 50%
    bonusPercent := min(streakCount, 50)
    if bonusPercent == 0 {
        return baseReward
    }
    
    multiplier := valueobject.NewDecimal(fmt.Sprintf("1.%02d", bonusPercent))
    return baseReward.Multiply(multiplier)
}
```

**Testing**:
- Unit tests for each quest mode
- Test min/max constraints for PER_MINUTE
- Test streak bonus calculations
- Test edge cases (zero values, nil pointers)

### 2.2 Update Quest Completion Service
**File**: `internal/usecase/quest_service.go` (modify existing)

Update the `CompleteQuest` method to actually award points:

```go
func (s *QuestService) CompleteQuest(ctx context.Context, userID int64, questID string, input CompleteQuestInput) (*entity.QuestCompletion, error) {
    // ... existing idempotency logic ...

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
        
        // Check cooldown
        if quest.LastCompletedAt != nil {
            timeSinceLastCompletion := time.Since(*quest.LastCompletedAt)
            if timeSinceLastCompletion.Seconds() < float64(quest.CooldownSec) {
                return errors.New("quest is on cooldown")
            }
        }
        
        // Calculate base reward
        calculator := &QuestRewardCalculator{}
        baseReward, err := calculator.CalculateReward(quest, input)
        if err != nil {
            return err
        }
        
        // Apply streak bonus if enabled
        finalReward := baseReward
        newStreakCount := quest.StreakCount
        if quest.StreakEnabled {
            newStreakCount++
            finalReward = calculator.ApplyStreakBonus(baseReward, newStreakCount)
        }
        
        // Check daily cap
        if quest.DailyPointsCap != nil {
            todayCompletions, err := s.questCompletionRepo.GetTodayCompletions(ctx, userID, questID)
            if err != nil {
                return err
            }
            
            totalTodayPoints := valueobject.NewDecimal("0")
            for _, comp := range todayCompletions {
                totalTodayPoints = totalTodayPoints.Add(comp.PointsAwarded)
            }
            
            if totalTodayPoints.Add(finalReward).GreaterThan(*quest.DailyPointsCap) {
                finalReward = quest.DailyPointsCap.Subtract(totalTodayPoints)
                if finalReward.LessThanOrEqual(valueobject.NewDecimal("0")) {
                    return errors.New("daily points cap reached")
                }
            }
        }
        
        // Award points to user
        user.Balance = user.Balance.Add(finalReward)
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
            StreakCount:     newStreakCount,
            CompletedAt:     time.Now(),
        }
        
        err = s.questCompletionRepo.Create(ctx, completion)
        if err != nil {
            return err
        }
        
        // Update quest state
        now := time.Now()
        quest.LastCompletedAt = &now
        quest.StreakCount = newStreakCount
        quest.UpdatedAt = now
        
        err = s.questRepo.Update(ctx, quest)
        if err != nil {
            return err
        }
        
        return nil
    })
    
    // ... existing idempotency completion logic ...
    
    return completion, err
}
```

**Testing**:
- Test point awarding for each quest mode
- Test streak bonus application
- Test daily cap enforcement
- Test cooldown enforcement
- Test transaction rollback on errors

### 2.3 Create Quest Completion Repository
**File**: `internal/infra/postgres/quest_completion_repository.go` (modify existing)

Ensure the repository has all needed methods:

```go
func (r *QuestCompletionRepository) GetTodayCompletions(ctx context.Context, userID int64, questID string) ([]*entity.QuestCompletion, error) {
    query := `
        SELECT id, quest_id, user_id, points_awarded, completion_ratio, minutes_spent, streak_count, completed_at
        FROM quest_completions 
        WHERE user_id = $1 AND quest_id = $2 
        AND completed_at >= CURRENT_DATE 
        AND completed_at < CURRENT_DATE + INTERVAL '1 day'
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
            &completion.CompletedAt,
        )
        if err != nil {
            return nil, err
        }
        completions = append(completions, completion)
    }
    
    return completions, nil
}

func (r *QuestCompletionRepository) GetUserStats(ctx context.Context, userID int64) (*UserQuestStats, error) {
    query := `
        SELECT 
            COUNT(*) as total_completions,
            SUM(points_awarded) as total_points_earned,
            MAX(streak_count) as max_streak,
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
```

**Testing**:
- Test today's completions filtering
- Test user stats aggregation
- Test timezone handling for "today"

### 2.4 Add Quest Completion Entity
**File**: `internal/domain/entity/quest_completion.go` (modify existing)

Ensure the entity matches the database schema:

```go
type QuestCompletion struct {
    ID              int64
    QuestID         string
    UserID          int64
    PointsAwarded   valueobject.Decimal
    CompletionRatio *float64 // For PARTIAL mode
    MinutesSpent    *int     // For PER_MINUTE mode
    StreakCount     int
    CompletedAt     time.Time
}

type UserQuestStats struct {
    TotalCompletions       int
    TotalPointsEarned      valueobject.Decimal
    MaxStreak              int
    UniqueQuestsCompleted  int
}
```

**Testing**:
- Test entity creation and validation
- Test decimal precision for points
- Test optional fields handling

## Testing Strategy

### 2.5 Unit Tests
Create test file: `internal/usecase/quest_reward_calculator_test.go`

```go
func TestQuestRewardCalculator_BINARY(t *testing.T) {
    // Test binary mode returns exact points
}

func TestQuestRewardCalculator_PARTIAL(t *testing.T) {
    // Test partial completion calculations
    // Test 0%, 50%, 100% completion
}

func TestQuestRewardCalculator_PER_MINUTE(t *testing.T) {
    // Test per-minute calculations
    // Test min/max constraints
}

func TestStreakBonus(t *testing.T) {
    // Test streak bonus calculations
    // Test max bonus cap
}
```

### 2.6 Integration Tests
Create test file: `test/integration/quest_completion_test.go`

```go
func TestCompleteQuestFullWorkflow(t *testing.T) {
    // Create user with initial balance
    // Create quest
    // Complete quest
    // Verify points awarded
    // Verify completion record created
    // Verify quest state updated
}

func TestQuestCooldown(t *testing.T) {
    // Complete quest
    // Try to complete again immediately
    // Verify cooldown error
    // Wait for cooldown
    // Complete successfully
}

func TestDailyPointsCap(t *testing.T) {
    // Set quest with daily cap
    // Complete multiple times
    // Verify cap is enforced
}
```

## Verification Checklist

- [ ] Point calculation works for all quest modes
- [ ] Points are actually awarded to user balance
- [ ] Quest completion records are created
- [ ] Streak bonuses are applied correctly
- [ ] Daily caps are enforced
- [ ] Cooldowns prevent rapid completion
- [ ] All operations are atomic (transaction rollback works)
- [ ] Idempotency prevents double rewards
- [ ] Unit tests cover all calculation logic
- [ ] Integration tests verify full workflow

## Files to Create/Modify

1. `internal/usecase/quest_reward_calculator.go` - New reward calculation logic
2. `internal/usecase/quest_service.go` - Update CompleteQuest method
3. `internal/infra/postgres/quest_completion_repository.go` - Add missing methods
4. `internal/domain/entity/quest_completion.go` - Ensure entity matches schema
5. `internal/usecase/quest_reward_calculator_test.go` - Unit tests
6. `test/integration/quest_completion_test.go` - Integration tests

## Success Criteria

✅ Users earn points when completing quests
✅ All quest modes (BINARY, PARTIAL, PER_MINUTE) work correctly
✅ Streak bonuses increase rewards
✅ Daily caps prevent abuse
✅ Cooldowns work as expected
✅ All operations are atomic and idempotent
✅ Comprehensive test coverage

## Next Phase
Once quest rewards are working, move to **Phase 3: Bot Quest Commands** to add user-facing quest management commands.