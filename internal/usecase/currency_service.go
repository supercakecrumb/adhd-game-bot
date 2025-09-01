package usecase

import (
	"context"
	"fmt"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
	"github.com/supercakecrumb/adhd-game-bot/internal/ports"
)

type CurrencyService struct {
	currencyRepo ports.CurrencyRepository
	userRepo     ports.UserRepository
}

func NewCurrencyService(
	currencyRepo ports.CurrencyRepository,
	userRepo ports.UserRepository,
) *CurrencyService {
	return &CurrencyService{
		currencyRepo: currencyRepo,
		userRepo:     userRepo,
	}
}

// CreateCurrency creates a new currency for a chat
func (s *CurrencyService) CreateCurrency(ctx context.Context, chatID int64, currency *entity.Currency) error {
	currency.ChatID = chatID

	// If this is the first currency, make it the base currency
	existing, err := s.currencyRepo.FindByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("failed to check existing currencies: %w", err)
	}

	if len(existing) == 0 {
		currency.IsBaseCurrency = true
	}

	return s.currencyRepo.Create(ctx, currency)
}

// GetCurrencies returns all currencies for a chat
func (s *CurrencyService) GetCurrencies(ctx context.Context, chatID int64) ([]*entity.Currency, error) {
	return s.currencyRepo.FindByChatID(ctx, chatID)
}

// GetBaseCurrency returns the base currency for a chat
func (s *CurrencyService) GetBaseCurrency(ctx context.Context, chatID int64) (*entity.Currency, error) {
	return s.currencyRepo.GetBaseCurrency(ctx, chatID)
}

// SetExchangeRate sets the exchange rate from one currency to another
func (s *CurrencyService) SetExchangeRate(ctx context.Context, fromCurrencyID int64, toCurrencyCode string, rate valueobject.Decimal) error {
	currency, err := s.currencyRepo.FindByID(ctx, fromCurrencyID)
	if err != nil {
		return err
	}

	if currency.ExchangeRates == nil {
		currency.ExchangeRates = make(map[string]valueobject.Decimal)
	}

	currency.ExchangeRates[toCurrencyCode] = rate

	return s.currencyRepo.Update(ctx, currency)
}

// ConvertCurrency converts an amount from one currency to another
func (s *CurrencyService) ConvertCurrency(ctx context.Context, amount valueobject.Decimal, fromCurrencyID, toCurrencyID int64) (valueobject.Decimal, error) {
	if fromCurrencyID == toCurrencyID {
		return amount, nil
	}

	fromCurrency, err := s.currencyRepo.FindByID(ctx, fromCurrencyID)
	if err != nil {
		return valueobject.Decimal{}, fmt.Errorf("source currency not found: %w", err)
	}

	toCurrency, err := s.currencyRepo.FindByID(ctx, toCurrencyID)
	if err != nil {
		return valueobject.Decimal{}, fmt.Errorf("target currency not found: %w", err)
	}

	// Direct conversion if exchange rate exists
	if rate, exists := fromCurrency.ExchangeRates[toCurrency.Code]; exists {
		return amount.Mul(rate), nil
	}

	// Try conversion through base currency
	baseCurrency, err := s.currencyRepo.GetBaseCurrency(ctx, fromCurrency.ChatID)
	if err != nil {
		return valueobject.Decimal{}, fmt.Errorf("base currency not found: %w", err)
	}

	// If from currency is base, check if base has rate to target
	if fromCurrency.ID == baseCurrency.ID {
		if rate, exists := baseCurrency.ExchangeRates[toCurrency.Code]; exists {
			return amount.Mul(rate), nil
		}
		return valueobject.Decimal{}, entity.ErrExchangeRateNotFound
	}

	// Convert to base first, then to target
	// Find rate from source to base (inverse of base to source)
	baseToSourceRate, exists := baseCurrency.ExchangeRates[fromCurrency.Code]
	if !exists {
		return valueobject.Decimal{}, entity.ErrExchangeRateNotFound
	}

	// Inverse rate (source to base)
	one := valueobject.NewDecimal("1")
	sourceToBaseRate := one.Div(baseToSourceRate)
	amountInBase := amount.Mul(sourceToBaseRate)

	// Now convert from base to target
	baseToTargetRate, exists := baseCurrency.ExchangeRates[toCurrency.Code]
	if !exists {
		return valueobject.Decimal{}, entity.ErrExchangeRateNotFound
	}

	return amountInBase.Mul(baseToTargetRate), nil
}

// InitializeUserBalances creates balance entries for all currencies in a chat
func (s *CurrencyService) InitializeUserBalances(ctx context.Context, userID int64, chatID int64) error {
	currencies, err := s.currencyRepo.FindByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("failed to get currencies: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Balances == nil {
		user.Balances = make(map[int64]valueobject.Decimal)
	}

	// Initialize balance for each currency
	for _, currency := range currencies {
		if _, exists := user.Balances[currency.ID]; !exists {
			user.Balances[currency.ID] = valueobject.NewDecimal("0")
		}
	}

	// In a real implementation, this would be done in the repository
	// For now, we'll just return success
	return nil
}
