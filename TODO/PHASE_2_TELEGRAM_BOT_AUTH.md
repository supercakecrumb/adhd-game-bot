# Phase 2: Telegram Bot Authentication & Quick Actions

## Overview
Implement Telegram bot as an authentication gateway and quick action interface. Users authenticate via Telegram to access the web interface, and can perform quick actions like completing quests and viewing lists through bot commands.

## Goals
- Implement Telegram authentication flow for web access
- Add quick action buttons for quest completion
- Add simple list commands for viewing quests
- Create session management between bot and web
- Send quest reminders and notifications

## Tasks

### 2.1 Telegram Authentication Flow
**File**: `internal/usecase/telegram_auth_service.go` (new)

```go
type TelegramAuthService struct {
    userRepo        ports.UserRepository
    sessionRepo     ports.TelegramSessionRepository
    uuidGen         ports.UUIDGenerator
    txManager       ports.TxManager
}

type AuthResponse struct {
    SessionToken string    `json:"session_token"`
    ExpiresAt    time.Time `json:"expires_at"`
    User         *entity.User `json:"user"`
}

func (s *TelegramAuthService) AuthenticateUser(ctx context.Context, telegramUserID int64, username string) (*AuthResponse, error) {
    return s.txManager.WithTxResult(ctx, func(ctx context.Context) (*AuthResponse, error) {
        // Find or create user
        user, err := s.userRepo.FindByTelegramID(ctx, telegramUserID)
        if err != nil {
            // Create new user
            user = &entity.User{
                TelegramUserID: &telegramUserID,
                Username:       username,
                Balance:        valueobject.NewDecimal("0.00"),
                TimeZone:       "UTC",
            }
            err = s.userRepo.Create(ctx, user)
            if err != nil {
                return nil, err
            }
        }
        
        // Generate session token
        sessionToken := s.uuidGen.New()
        expiresAt := time.Now().Add(24 * time.Hour) // 24 hour sessions
        
        // Create or update session
        session := &entity.TelegramSession{
            TelegramUserID: telegramUserID,
            UserID:         user.ID,
            SessionToken:   sessionToken,
            ExpiresAt:      expiresAt,
        }
        
        err = s.sessionRepo.CreateOrUpdate(ctx, session)
        if err != nil {
            return nil, err
        }
        
        return &AuthResponse{
            SessionToken: sessionToken,
            ExpiresAt:    expiresAt,
            User:         user,
        }, nil
    })
}

func (s *TelegramAuthService) ValidateSession(ctx context.Context, sessionToken string) (*entity.User, error) {
    session, err := s.sessionRepo.FindByToken(ctx, sessionToken)
    if err != nil {
        return nil, err
    }
    
    if session.IsExpired() {
        return nil, errors.New("session expired")
    }
    
    return s.userRepo.FindByID(ctx, session.UserID)
}

func (s *TelegramAuthService) GenerateWebLoginURL(ctx context.Context, telegramUserID int64) (string, error) {
    // Generate temporary auth token for web login
    authToken := s.uuidGen.New()
    
    // Store temporary token (expires in 5 minutes)
    tempSession := &entity.TelegramSession{
        TelegramUserID: telegramUserID,
        SessionToken:   authToken,
        ExpiresAt:      time.Now().Add(5 * time.Minute),
    }
    
    err := s.sessionRepo.CreateOrUpdate(ctx, tempSession)
    if err != nil {
        return "", err
    }
    
    // Return web URL with auth token
    webURL := os.Getenv("WEB_URL")
    if webURL == "" {
        webURL = "http://localhost:3000"
    }
    
    return fmt.Sprintf("%s/auth/telegram?token=%s", webURL, authToken), nil
}
```

**Testing**:
- Test user creation and authentication
- Test session token generation and validation
- Test session expiration
- Test web login URL generation

### 2.2 Update Telegram Bot with Auth Commands
**File**: `cmd/bot/main.go` (modify existing)

```go
// Add auth service initialization
authService := usecase.NewTelegramAuthService(userRepo, sessionRepo, uuidGen, txManager)

// Update /start command
bot.Handle("/start", func(c telebot.Context) error {
    telegramUserID := c.Sender().ID
    username := c.Sender().FirstName
    if c.Sender().Username != "" {
        username = c.Sender().Username
    }
    
    ctx := context.Background()
    
    // Authenticate user
    authResp, err := authService.AuthenticateUser(ctx, telegramUserID, username)
    if err != nil {
        log.Printf("Failed to authenticate user: %v", err)
        return c.Send("❌ Authentication failed")
    }
    
    // Generate web login URL
    webURL, err := authService.GenerateWebLoginURL(ctx, telegramUserID)
    if err != nil {
        log.Printf("Failed to generate web URL: %v", err)
        webURL = "Check the web interface"
    }
    
    // Create inline keyboard for web access
    webButton := &telebot.InlineButton{
        Text: "🌐 Open Web Interface",
        URL:  webURL,
    }
    
    keyboard := &telebot.ReplyMarkup{}
    keyboard.Inline(keyboard.Row(webButton))
    
    message := fmt.Sprintf(
        "🎮 Welcome to ADHD Game Bot!\n\n" +
        "💰 Balance: %s points\n\n" +
        "Quick commands:\n" +
        "• /quests - View your active quests\n" +
        "• /complete - Complete a quest\n" +
        "• /balance - Check your balance\n\n" +
        "For full quest management, use the web interface:",
        authResp.User.Balance,
    )
    
    return c.Send(message, keyboard)
})

// Add /web command for getting login URL
bot.Handle("/web", func(c telebot.Context) error {
    telegramUserID := c.Sender().ID
    ctx := context.Background()
    
    webURL, err := authService.GenerateWebLoginURL(ctx, telegramUserID)
    if err != nil {
        return c.Send("❌ Failed to generate web login URL")
    }
    
    webButton := &telebot.InlineButton{
        Text: "🌐 Open Web Interface",
        URL:  webURL,
    }
    
    keyboard := &telebot.ReplyMarkup{}
    keyboard.Inline(keyboard.Row(webButton))
    
    return c.Send("Click the button below to access the web interface:", keyboard)
})
```

**Testing**:
- Test /start command creates user and session
- Test /web command generates valid login URLs
- Test web login URL expires correctly
- Test inline keyboard buttons work

### 2.3 Quick Quest Actions
**File**: `cmd/bot/main.go` (add to existing)

```go
// Add /quests command with inline buttons
bot.Handle("/quests", func(c telebot.Context) error {
    telegramUserID := c.Sender().ID
    ctx := context.Background()
    
    // Get user
    user, err := userRepo.FindByTelegramID(ctx, telegramUserID)
    if err != nil {
        return c.Send("❌ User not found. Use /start first.")
    }
    
    // Get user's dungeons
    dungeons, err := dungeonService.GetUserDungeons(ctx, user.ID)
    if err != nil || len(dungeons) == 0 {
        return c.Send("📝 No quests available. Join a dungeon or create quests on the web interface.")
    }
    
    // Get active quests from first dungeon (for simplicity)
    quests, err := questService.ListQuests(ctx, user.ID, dungeons[0].ID)
    if err != nil {
        return c.Send("❌ Failed to load quests")
    }
    
    if len(quests) == 0 {
        return c.Send("📝 No active quests found.")
    }
    
    // Create message with quest list and completion buttons
    message := "📋 Your Active Quests:\n\n"
    keyboard := &telebot.ReplyMarkup{}
    var buttons []telebot.Row
    
    for i, quest := range quests {
        if i >= 5 { // Limit to 5 quests for readability
            break
        }
        
        // Quest info
        message += fmt.Sprintf("🎯 %s\n", quest.Title)
        message += fmt.Sprintf("   💎 %s points • %s\n\n", quest.PointsAward, quest.Difficulty)
        
        // Complete button
        completeBtn := &telebot.InlineButton{
            Text: fmt.Sprintf("✅ Complete: %s", quest.Title),
            Data: fmt.Sprintf("complete_%s", quest.ID),
        }
        buttons = append(buttons, keyboard.Row(completeBtn))
    }
    
    if len(quests) > 5 {
        message += fmt.Sprintf("... and %d more quests\n", len(quests)-5)
    }
    
    // Add web interface button
    webURL, _ := authService.GenerateWebLoginURL(ctx, telegramUserID)
    webBtn := &telebot.InlineButton{
        Text: "🌐 View All Quests",
        URL:  webURL,
    }
    buttons = append(buttons, keyboard.Row(webBtn))
    
    keyboard.Inline(buttons...)
    return c.Send(message, keyboard)
})

// Handle quest completion callbacks
bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
    data := c.Callback().Data
    
    if strings.HasPrefix(data, "complete_") {
        questID := strings.TrimPrefix(data, "complete_")
        telegramUserID := c.Sender().ID
        ctx := context.Background()
        
        // Get user
        user, err := userRepo.FindByTelegramID(ctx, telegramUserID)
        if err != nil {
            return c.Respond(&telebot.CallbackResponse{Text: "❌ User not found"})
        }
        
        // Complete quest (BINARY mode by default)
        idempotencyKey := fmt.Sprintf("tg_complete_%d_%s_%d", telegramUserID, questID, time.Now().Unix())
        completion, err := questService.CompleteQuest(ctx, user.ID, questID, usecase.CompleteQuestInput{
            IdempotencyKey: idempotencyKey,
        })
        
        if err != nil {
            return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("❌ %v", err)})
        }
        
        // Success response
        responseText := fmt.Sprintf("🎉 Quest completed! +%s points", completion.PointsAwarded)
        
        // Update the message to show completion
        c.Edit(c.Message().Text + "\n\n✅ Quest completed via Telegram!")
        
        return c.Respond(&telebot.CallbackResponse{Text: responseText, ShowAlert: true})
    }
    
    return c.Respond(&telebot.CallbackResponse{Text: "Unknown action"})
})
```

**Testing**:
- Test /quests command shows active quests
- Test completion buttons work correctly
- Test idempotency prevents double completion
- Test error handling for invalid quests

### 2.4 Quest Reminders System
**File**: `internal/usecase/reminder_service.go` (new)

```go
type ReminderService struct {
    questRepo       ports.QuestRepository
    userRepo        ports.UserRepository
    dungeonRepo     ports.DungeonRepository
    bot             *telebot.Bot
}

func (s *ReminderService) SendDailyReminders(ctx context.Context) error {
    // Get all active daily quests
    quests, err := s.questRepo.GetActiveQuestsByCategory(ctx, "daily")
    if err != nil {
        return err
    }
    
    for _, quest := range quests {
        // Check if quest needs reminder (not completed today)
        if s.shouldSendReminder(quest) {
            err := s.sendQuestReminder(ctx, quest)
            if err != nil {
                log.Printf("Failed to send reminder for quest %s: %v", quest.ID, err)
            }
        }
    }
    
    return nil
}

func (s *ReminderService) shouldSendReminder(quest *entity.Quest) bool {
    if quest.LastCompletedAt == nil {
        return true // Never completed
    }
    
    // Check if completed today in quest's timezone
    loc, err := time.LoadLocation(quest.TimeZone)
    if err != nil {
        loc = time.UTC
    }
    
    now := time.Now().In(loc)
    lastCompleted := quest.LastCompletedAt.In(loc)
    
    // Not completed today
    return !isSameDay(now, lastCompleted)
}

func (s *ReminderService) sendQuestReminder(ctx context.Context, quest *entity.Quest) error {
    // Get dungeon
    dungeon, err := s.dungeonRepo.GetByID(ctx, quest.DungeonID)
    if err != nil {
        return err
    }
    
    // Only send to Telegram-linked dungeons
    if dungeon.TelegramChatID == nil {
        return nil
    }
    
    // Create reminder message with completion button
    message := fmt.Sprintf(
        "⏰ Daily Quest Reminder\n\n" +
        "🎯 %s\n" +
        "💎 %s points\n\n" +
        "Ready to complete it?",
        quest.Title,
        quest.PointsAward,
    )
    
    completeBtn := &telebot.InlineButton{
        Text: "✅ Complete Now",
        Data: fmt.Sprintf("complete_%s", quest.ID),
    }
    
    keyboard := &telebot.ReplyMarkup{}
    keyboard.Inline(keyboard.Row(completeBtn))
    
    chat := &telebot.Chat{ID: *dungeon.TelegramChatID}
    _, err = s.bot.Send(chat, message, keyboard)
    
    return err
}

func isSameDay(t1, t2 time.Time) bool {
    y1, m1, d1 := t1.Date()
    y2, m2, d2 := t2.Date()
    return y1 == y2 && m1 == m2 && d1 == d2
}
```

**Testing**:
- Test daily reminder logic
- Test timezone handling for reminders
- Test reminder messages are sent correctly
- Test completion buttons in reminders work

### 2.5 Web Authentication Endpoint
**File**: `internal/infra/http/auth_handler.go` (new)

```go
type AuthHandler struct {
    authService *usecase.TelegramAuthService
}

func (h *AuthHandler) TelegramAuth(w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")
    if token == "" {
        http.Error(w, "Missing auth token", http.StatusBadRequest)
        return
    }
    
    ctx := r.Context()
    
    // Validate session token
    user, err := h.authService.ValidateSession(ctx, token)
    if err != nil {
        http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
        return
    }
    
    // Create long-term session for web
    authResp, err := h.authService.AuthenticateUser(ctx, *user.TelegramUserID, user.Username)
    if err != nil {
        http.Error(w, "Authentication failed", http.StatusInternalServerError)
        return
    }
    
    // Set session cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    authResp.SessionToken,
        Expires:  authResp.ExpiresAt,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
    })
    
    // Redirect to dashboard
    http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func (h *AuthHandler) ValidateSession(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        http.Error(w, "No session", http.StatusUnauthorized)
        return
    }
    
    ctx := r.Context()
    user, err := h.authService.ValidateSession(ctx, cookie.Value)
    if err != nil {
        http.Error(w, "Invalid session", http.StatusUnauthorized)
        return
    }
    
    // Return user info
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

**Testing**:
- Test Telegram auth flow from web
- Test session cookie creation
- Test session validation endpoint
- Test redirect after successful auth

## Testing Strategy

### 2.6 Integration Tests
**File**: `test/integration/telegram_auth_test.go`

```go
func TestTelegramAuthFlow(t *testing.T) {
    // Test complete auth flow
    // 1. User starts bot
    // 2. Gets web login URL
    // 3. Accesses web with token
    // 4. Gets session cookie
    // 5. Can access protected endpoints
}

func TestQuickActions(t *testing.T) {
    // Test quest completion via bot
    // Test quest listing
    // Test reminder system
}
```

## Verification Checklist

- [ ] Telegram authentication creates user sessions
- [ ] Web login URLs work and expire correctly
- [ ] Quest completion buttons work in Telegram
- [ ] Quest reminders are sent at appropriate times
- [ ] Session management works between bot and web
- [ ] All bot commands handle errors gracefully
- [ ] Web authentication endpoints work correctly

## Files to Create/Modify

1. `internal/usecase/telegram_auth_service.go` - New auth service
2. `internal/usecase/reminder_service.go` - New reminder service
3. `internal/infra/http/auth_handler.go` - New web auth endpoints
4. `cmd/bot/main.go` - Update with auth and quick actions
5. `test/integration/telegram_auth_test.go` - Integration tests

## Success Criteria

✅ Users can authenticate via Telegram for web access
✅ Quick quest completion works through bot buttons
✅ Quest reminders are sent automatically
✅ Session management works seamlessly
✅ Bot provides good UX for quick actions
✅ Web interface can authenticate Telegram users

## Next Phase
Once Telegram auth is working, move to **Phase 3: Web Interface Integration** to connect the frontend with the backend APIs.