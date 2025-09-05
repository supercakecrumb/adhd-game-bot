package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/infra/inmemory"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
	"gopkg.in/telebot.v3"
)

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	// Initialize in-memory repositories
	userRepo := inmemory.NewUserRepository()
	shopItemRepo := inmemory.NewShopItemRepository()
	purchaseRepo := inmemory.NewPurchaseRepository()

	// Create test data
	ctx := context.Background()
	testUser := &entity.User{
		ID:      1,
		ChatID:  0,
		Balance: valueobject.NewDecimal("100.00"),
	}
	userRepo.Create(ctx, testUser)

	testItem := &entity.ShopItem{
		ID:       1,
		Name:     "Focus Booster",
		Code:     "FOCUS",
		Price:    valueobject.NewDecimal("10.00"),
		IsActive: true,
	}
	shopItemRepo.Create(ctx, testItem)

	// Initialize service with required dependencies
	shopService := usecase.NewShopServiceV2(
		shopItemRepo,
		purchaseRepo,
		userRepo,
		nil, // chatConfigRepo
		nil, // discountTierRepo
		nil, // uuidGen
		nil, // txManager
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
		return c.Send("🎮 Welcome to ADHD Game Bot!\n" +
			"Use /shop to see available items\n" +
			"Use /buy <code> to purchase items")
	})

	bot.Handle("/shop", func(c telebot.Context) error {
		items, err := shopService.GetShopItems(ctx, 0)
		if err != nil {
			return c.Send("❌ Error getting shop items")
		}

		message := "🛍️ Available Items:\n"
		for _, item := range items {
			message += fmt.Sprintf("- %s (%s): %s\n",
				item.Name, item.Code, item.Price)
		}
		return c.Send(message)
	})

	bot.Handle("/buy", func(c telebot.Context) error {
		if len(c.Args()) == 0 {
			return c.Send("Usage: /buy <item_code>")
		}

		itemCode := c.Args()[0]
		purchase, err := shopService.PurchaseItemWithIdempotency(ctx, 1, itemCode, 1, "")
		if err != nil {
			return c.Send(fmt.Sprintf("❌ Purchase failed: %v", err))
		}

		return c.Send(fmt.Sprintf("✅ Purchased %s for %s!",
			purchase.ItemName, purchase.TotalCost))
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
