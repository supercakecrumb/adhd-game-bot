package entity

import "errors"

var (
	ErrTaskNotFound         = errors.New("task not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrExchangeRateNotFound = errors.New("exchange rate not found")
	ErrCurrencyNotFound     = errors.New("currency not found")
)
