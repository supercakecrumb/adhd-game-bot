package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
)

type mockCurrencyRepo struct {
	currencies map[int64]*entity.Currency
	nextID     int64
}

func newMockCurrencyRepo() *mockCurrencyRepo {
	return &mockCurrencyRepo{
		currencies: make(map[int64]*entity.Currency),
		nextID:     1,
	}
}

func (m *mockCurrencyRepo) Create(ctx context.Context, currency *entity.Currency) error {
	currency.ID = m.nextID
	m.nextID++
	m.currencies[currency.ID] = currency
	return nil
}

func (m *mockCurrencyRepo) FindByID(ctx context.Context, id int64) (*entity.Currency, error) {
	currency, exists := m.currencies[id]
	if !exists {
		return nil, entity.ErrCurrencyNotFound
	}
	return currency, nil
}

func (m *mockCurrencyRepo) FindByCode(ctx context.Context, chatID int64, code string) (*entity.Currency, error) {
	for _, currency := range m.currencies {
		if currency.ChatID == chatID && currency.Code == code {
			return currency, nil
		}
	}
	return nil, entity.ErrCurrencyNotFound
}

func (m *mockCurrencyRepo) FindByChatID(ctx context.Context, chatID int64) ([]*entity.Currency, error) {
	var result []*entity.Currency
	for _, currency := range m.currencies {
		if currency.ChatID == chatID {
			result = append(result, currency)
		}
	}
	return result, nil
}

func (m *mockCurrencyRepo) GetBaseCurrency(ctx context.Context, chatID int64) (*entity.Currency, error) {
	for _, currency := range m.currencies {
		if currency.ChatID == chatID && currency.IsBaseCurrency {
			return currency, nil
		}
	}
	return nil, entity.ErrCurrencyNotFound
}

func (m *mockCurrencyRepo) Update(ctx context.Context, currency *entity.Currency) error {
	m.currencies[currency.ID] = currency
	return nil
}

func (m *mockCurrencyRepo) Delete(ctx context.Context, id int64) error {
	delete(m.currencies, id)
	return nil
}

// mockUserRepoForCurrency implements the UserRepository interface for currency tests
type mockUserRepoForCurrency struct {
	users map[int64]*entity.User
}

func (m *mockUserRepoForCurrency) Create(ctx context.Context, user *entity.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepoForCurrency) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, entity.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepoForCurrency) UpdateBalance(ctx context.Context, userID int64, currencyID int64, delta entity.Decimal) error {
	user, exists := m.users[userID]
	if !exists {
		return entity.ErrUserNotFound
	}
	if user.Balances == nil {
		user.Balances = make(map[int64]entity.Decimal)
	}
	current := user.Balances[currencyID]
	user.Balances[currencyID] = current.Add(delta)
	return nil
}

func (m *mockUserRepoForCurrency) Delete(ctx context.Context, id int64) error {
	delete(m.users, id)
	return nil
}

func TestCurrencyService(t *testing.T) {
	ctx := context.Background()
	currencyRepo := newMockCurrencyRepo()
	userRepo := &mockUserRepoForCurrency{users: make(map[int64]*entity.User)}

	service := usecase.NewCurrencyService(currencyRepo, userRepo)

	chatID := int64(12345) // Example Telegram chat ID

	t.Run("CreateFirstCurrency_BecomesBase", func(t *testing.T) {
		currency := &entity.Currency{
			Code:     "MM",
			Name:     "Motivation Minutes",
			Decimals: 2,
		}

		err := service.CreateCurrency(ctx, chatID, currency)
		require.NoError(t, err)
		require.True(t, currency.IsBaseCurrency)
		require.Equal(t, chatID, currency.ChatID)
	})

	t.Run("CreateAdditionalCurrencies", func(t *testing.T) {
		// Create second currency
		currency2 := &entity.Currency{
			Code:     "BS",
			Name:     "Bonus Stars",
			Decimals: 0,
		}

		err := service.CreateCurrency(ctx, chatID, currency2)
		require.NoError(t, err)
		require.False(t, currency2.IsBaseCurrency)

		// Create third currency
		currency3 := &entity.Currency{
			Code:     "XP",
			Name:     "Experience Points",
			Decimals: 0,
		}

		err = service.CreateCurrency(ctx, chatID, currency3)
		require.NoError(t, err)
		require.False(t, currency3.IsBaseCurrency)
	})

	t.Run("SetExchangeRates", func(t *testing.T) {
		// Get base currency (MM)
		baseCurrency, err := service.GetBaseCurrency(ctx, chatID)
		require.NoError(t, err)
		require.Equal(t, "MM", baseCurrency.Code)

		// Set exchange rates from MM to other currencies
		// 1 MM = 0.1 BS
		err = service.SetExchangeRate(ctx, baseCurrency.ID, "BS", valueobject.NewDecimal("0.1"))
		require.NoError(t, err)

		// 1 MM = 10 XP
		err = service.SetExchangeRate(ctx, baseCurrency.ID, "XP", valueobject.NewDecimal("10"))
		require.NoError(t, err)
	})

	t.Run("ConvertCurrency_Direct", func(t *testing.T) {
		// Get currencies
		baseCurrency, _ := service.GetBaseCurrency(ctx, chatID)
		bsCurrency, _ := currencyRepo.FindByCode(ctx, chatID, "BS")

		// Convert 100 MM to BS (should be 10 BS)
		amount := valueobject.NewDecimal("100")
		converted, err := service.ConvertCurrency(ctx, amount, baseCurrency.ID, bsCurrency.ID)
		require.NoError(t, err)
		require.Equal(t, "10", converted.String())
	})

	t.Run("ConvertCurrency_ThroughBase", func(t *testing.T) {
		// Get currencies
		bsCurrency, _ := currencyRepo.FindByCode(ctx, chatID, "BS")
		xpCurrency, _ := currencyRepo.FindByCode(ctx, chatID, "XP")

		// Set up exchange rate from BS to MM (inverse of MM to BS)
		// If 1 MM = 0.1 BS, then 1 BS = 10 MM
		// This is handled by the conversion logic

		// Convert 5 BS to XP
		// 5 BS = 50 MM (5 * 10)
		// 50 MM = 500 XP (50 * 10)
		amount := valueobject.NewDecimal("5")
		converted, err := service.ConvertCurrency(ctx, amount, bsCurrency.ID, xpCurrency.ID)
		require.NoError(t, err)
		require.Equal(t, "500", converted.String())
	})

	t.Run("GetCurrencies", func(t *testing.T) {
		currencies, err := service.GetCurrencies(ctx, chatID)
		require.NoError(t, err)
		require.Len(t, currencies, 3)

		// Verify we have all three currencies
		codes := make(map[string]bool)
		for _, c := range currencies {
			codes[c.Code] = true
		}
		require.True(t, codes["MM"])
		require.True(t, codes["BS"])
		require.True(t, codes["XP"])
	})
}
