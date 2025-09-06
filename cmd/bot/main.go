package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/infra/postgres"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
	"gopkg.in/telebot.v3"
)

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	// Database connection
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "user=postgres dbname=adhd_bot sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize PostgreSQL repositories
	userRepo := postgres.NewUserRepository(db)
	shopItemRepo := postgres.NewShopItemRepository(db)
	purchaseRepo := postgres.NewPurchaseRepository(db)
	chatConfigRepo := postgres.NewChatConfigRepository(db)
	txManager := postgres.NewTxManager(db)

	// Initialize service with required dependencies
	shopService := usecase.NewShopServiceV2(
		shopItemRepo,
		purchaseRepo,
		userRepo,
		chatConfigRepo,
		nil, // discountTierRepo
		nil, // uuidGen
		txManager,
		nil, // idempotencyRepo
	)

	// Create bot
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Command handlers
	bot.Handle("/start", func(c telebot.Context) error {
		// Register user if not exists
		userID := c.Sender().ID
		chatID := c.Chat().ID

		ctx := context.Background()
		_, err := userRepo.FindByID(ctx, userID)
		if err != nil {
			// User not found, create new user
			newUser := &entity.User{
				ID:       userID,
				ChatID:   chatID,
				Username: c.Sender().FirstName,
				Balance:  valueobject.NewDecimal("0.00"),
				TimeZone: "UTC",
			}
			if err := userRepo.Create(ctx, newUser); err != nil {
				log.Printf("Failed to create user: %v", err)
				return c.Send("❌ Failed to register user")
			}
		}

		return c.Send("🎮 Welcome to ADHD Game Bot!\n" +
			"Use /shop to see available items\n" +
			"Use /buy <code> to purchase items\n" +
			"Use /balance to check your balance")
	})

	bot.Handle("/shop", func(c telebot.Context) error {
		userID := c.Sender().ID
		chatID := c.Chat().ID

		// Register user if not exists
		ctx := context.Background()
		_, err := userRepo.FindByID(ctx, userID)
		if err != nil {
			// User not found, create new user
			newUser := &entity.User{
				ID:       userID,
				ChatID:   chatID,
				Username: c.Sender().FirstName,
				Balance:  valueobject.NewDecimal("0.00"),
				TimeZone: "UTC",
			}
			if err := userRepo.Create(ctx, newUser); err != nil {
				log.Printf("Failed to create user: %v", err)
				return c.Send("❌ Failed to register user")
			}
		}

		items, err := shopService.GetShopItems(ctx, chatID)
		if err != nil {
			return c.Send("❌ Error getting shop items")
		}

		if len(items) == 0 {
			return c.Send("🛒 No items available in the shop right now.")
		}

		message := "🛍️ Available Items:\n"
		for _, item := range items {
			stockInfo := ""
			if item.Stock != nil {
				stockInfo = fmt.Sprintf(" (Stock: %d)", *item.Stock)
			}
			message += fmt.Sprintf("- %s (%s): %s%s\n",
				item.Name, item.Code, item.Price, stockInfo)
		}
		return c.Send(message)
	})

	bot.Handle("/buy", func(c telebot.Context) error {
		userID := c.Sender().ID
		chatID := c.Chat().ID

		if len(c.Args()) == 0 {
			return c.Send("Usage: /buy <item_code>")
		}

		// Register user if not exists
		ctx := context.Background()
		_, err := userRepo.FindByID(ctx, userID)
		if err != nil {
			// User not found, create new user
			newUser := &entity.User{
				ID:       userID,
				ChatID:   chatID,
				Username: c.Sender().FirstName,
				Balance:  valueobject.NewDecimal("0.00"),
				TimeZone: "UTC",
			}
			if err := userRepo.Create(ctx, newUser); err != nil {
				log.Printf("Failed to create user: %v", err)
				return c.Send("❌ Failed to register user")
			}
		}

		itemCode := c.Args()[0]
		purchase, err := shopService.PurchaseItemWithIdempotency(ctx, userID, itemCode, 1, "")
		if err != nil {
			return c.Send(fmt.Sprintf("❌ Purchase failed: %v", err))
		}

		return c.Send(fmt.Sprintf("✅ Purchased %s for %s!",
			purchase.ItemName, purchase.TotalCost))
	})

	bot.Handle("/balance", func(c telebot.Context) error {
		userID := c.Sender().ID
		chatID := c.Chat().ID

		// Register user if not exists
		ctx := context.Background()
		user, err := userRepo.FindByID(ctx, userID)
		if err != nil {
			// User not found, create new user
			newUser := &entity.User{
				ID:       userID,
				ChatID:   chatID,
				Username: c.Sender().FirstName,
				Balance:  valueobject.NewDecimal("0.00"),
				TimeZone: "UTC",
			}
			if err := userRepo.Create(ctx, newUser); err != nil {
				log.Printf("Failed to create user: %v", err)
				return c.Send("❌ Failed to register user")
			}
			user = newUser
		}

		// Get currency name
		currencyName := "Points"
		config, err := shopService.GetCurrencyName(ctx, chatID)
		if err == nil {
			currencyName = config
		}

		return c.Send(fmt.Sprintf("💰 Your balance: %s %s", user.Balance, currencyName))
	})

	// Start bot
	log.Println("Bot starting...")
	go bot.Start()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	bot.Stop()
}
