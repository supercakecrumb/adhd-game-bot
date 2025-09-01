# System Patterns *Optional*

This file documents recurring patterns...

*
[2025-09-01 14:21:04] - ## Currency System Architecture

### Multi-Chat Currency Model
- Each chat/group has its own set of currencies
- Currencies are identified by ID (not code) for better isolation
- User balances map currency IDs to amounts
- Base currency concept: first/primary currency per chat

### Exchange Rate System
- Exchange rates stored in Currency entity
- Rates are from this currency to others
- ConvertTo method handles conversions
- Rates stored as JSONB in PostgreSQL

### Database Changes
- currencies table: added chat_id, is_base_currency
- user_balances: changed from currency_code to currency_id
- users table: added chat_id
- Unique constraint: one base currency per chat

### Repository Pattern
- CurrencyRepository: full CRUD + chat-specific queries
- UserRepository: updated to use currency IDs
- Proper error handling for currency not found
