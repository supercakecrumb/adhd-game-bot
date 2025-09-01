package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/supercakecrumb/adhd-game-bot/internal/domain/entity"
	"github.com/supercakecrumb/adhd-game-bot/internal/domain/valueobject"
)

type CurrencyRepository struct {
	db *sql.DB
}

func NewCurrencyRepository(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r *CurrencyRepository) Create(ctx context.Context, currency *entity.Currency) error {
	rates, err := json.Marshal(currency.ExchangeRates)
	if err != nil {
		return fmt.Errorf("failed to marshal exchange rates: %w", err)
	}

	err = r.db.QueryRowContext(ctx, `
		INSERT INTO currencies (chat_id, code, name, decimals, is_base_currency, exchange_rates_json)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`,
		currency.ChatID, currency.Code, currency.Name, currency.Decimals,
		currency.IsBaseCurrency, rates).Scan(&currency.ID, &currency.CreatedAt, &currency.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create currency: %w", err)
	}

	return nil
}

func (r *CurrencyRepository) FindByID(ctx context.Context, id int64) (*entity.Currency, error) {
	var currency entity.Currency
	var ratesJSON []byte

	err := r.db.QueryRowContext(ctx, `
		SELECT id, chat_id, code, name, decimals, is_base_currency, exchange_rates_json, created_at, updated_at
		FROM currencies WHERE id = $1`, id).Scan(
		&currency.ID, &currency.ChatID, &currency.Code, &currency.Name, &currency.Decimals,
		&currency.IsBaseCurrency, &ratesJSON, &currency.CreatedAt, &currency.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrCurrencyNotFound
		}
		return nil, fmt.Errorf("failed to query currency: %w", err)
	}

	// Unmarshal exchange rates
	currency.ExchangeRates = make(map[string]valueobject.Decimal)
	if len(ratesJSON) > 0 {
		var rates map[string]string
		if err := json.Unmarshal(ratesJSON, &rates); err != nil {
			return nil, fmt.Errorf("failed to unmarshal exchange rates: %w", err)
		}
		for code, rate := range rates {
			currency.ExchangeRates[code] = valueobject.NewDecimal(rate)
		}
	}

	return &currency, nil
}

func (r *CurrencyRepository) FindByCode(ctx context.Context, chatID int64, code string) (*entity.Currency, error) {
	var currency entity.Currency
	var ratesJSON []byte

	err := r.db.QueryRowContext(ctx, `
		SELECT id, chat_id, code, name, decimals, is_base_currency, exchange_rates_json, created_at, updated_at
		FROM currencies WHERE chat_id = $1 AND code = $2`, chatID, code).Scan(
		&currency.ID, &currency.ChatID, &currency.Code, &currency.Name, &currency.Decimals,
		&currency.IsBaseCurrency, &ratesJSON, &currency.CreatedAt, &currency.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrCurrencyNotFound
		}
		return nil, fmt.Errorf("failed to query currency: %w", err)
	}

	// Unmarshal exchange rates
	currency.ExchangeRates = make(map[string]valueobject.Decimal)
	if len(ratesJSON) > 0 {
		var rates map[string]string
		if err := json.Unmarshal(ratesJSON, &rates); err != nil {
			return nil, fmt.Errorf("failed to unmarshal exchange rates: %w", err)
		}
		for code, rate := range rates {
			currency.ExchangeRates[code] = valueobject.NewDecimal(rate)
		}
	}

	return &currency, nil
}

func (r *CurrencyRepository) FindByChatID(ctx context.Context, chatID int64) ([]*entity.Currency, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, chat_id, code, name, decimals, is_base_currency, exchange_rates_json, created_at, updated_at
		FROM currencies WHERE chat_id = $1 ORDER BY is_base_currency DESC, code`, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query currencies: %w", err)
	}
	defer rows.Close()

	var currencies []*entity.Currency
	for rows.Next() {
		var currency entity.Currency
		var ratesJSON []byte

		if err := rows.Scan(&currency.ID, &currency.ChatID, &currency.Code, &currency.Name,
			&currency.Decimals, &currency.IsBaseCurrency, &ratesJSON,
			&currency.CreatedAt, &currency.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan currency: %w", err)
		}

		// Unmarshal exchange rates
		currency.ExchangeRates = make(map[string]valueobject.Decimal)
		if len(ratesJSON) > 0 {
			var rates map[string]string
			if err := json.Unmarshal(ratesJSON, &rates); err != nil {
				return nil, fmt.Errorf("failed to unmarshal exchange rates: %w", err)
			}
			for code, rate := range rates {
				currency.ExchangeRates[code] = valueobject.NewDecimal(rate)
			}
		}

		currencies = append(currencies, &currency)
	}

	return currencies, nil
}

func (r *CurrencyRepository) GetBaseCurrency(ctx context.Context, chatID int64) (*entity.Currency, error) {
	var currency entity.Currency
	var ratesJSON []byte

	err := r.db.QueryRowContext(ctx, `
		SELECT id, chat_id, code, name, decimals, is_base_currency, exchange_rates_json, created_at, updated_at
		FROM currencies WHERE chat_id = $1 AND is_base_currency = TRUE`, chatID).Scan(
		&currency.ID, &currency.ChatID, &currency.Code, &currency.Name, &currency.Decimals,
		&currency.IsBaseCurrency, &ratesJSON, &currency.CreatedAt, &currency.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrCurrencyNotFound
		}
		return nil, fmt.Errorf("failed to query base currency: %w", err)
	}

	// Unmarshal exchange rates
	currency.ExchangeRates = make(map[string]valueobject.Decimal)
	if len(ratesJSON) > 0 {
		var rates map[string]string
		if err := json.Unmarshal(ratesJSON, &rates); err != nil {
			return nil, fmt.Errorf("failed to unmarshal exchange rates: %w", err)
		}
		for code, rate := range rates {
			currency.ExchangeRates[code] = valueobject.NewDecimal(rate)
		}
	}

	return &currency, nil
}

func (r *CurrencyRepository) Update(ctx context.Context, currency *entity.Currency) error {
	rates, err := json.Marshal(currency.ExchangeRates)
	if err != nil {
		return fmt.Errorf("failed to marshal exchange rates: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE currencies 
		SET code = $2, name = $3, decimals = $4, is_base_currency = $5, 
		    exchange_rates_json = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`,
		currency.ID, currency.Code, currency.Name, currency.Decimals,
		currency.IsBaseCurrency, rates)

	if err != nil {
		return fmt.Errorf("failed to update currency: %w", err)
	}

	return nil
}

func (r *CurrencyRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM currencies WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete currency: %w", err)
	}
	return nil
}
