# ShopService API Contract

## Overview
The ShopService handles all shop-related operations including item management, purchases, and inventory tracking.

## Methods

### `PurchaseItem(userID int64, itemCode string, quantity int) (*Purchase, error)`
Purchases an item from the shop.

**Parameters:**
- `userID` - ID of the purchasing user
- `itemCode` - Unique code of the item to purchase  
- `quantity` - Number of items to purchase

**Returns:**
- `*Purchase` - Details of the completed purchase
- `error` - Possible errors:
  - `ErrItemNotFound`
  - `ErrInsufficientBalance` 
  - `ErrInsufficientStock`
  - `ErrInvalidQuantity`

**Preconditions:**
1. User must exist
2. Item must exist and be active
3. User must have sufficient balance
4. Item must have sufficient stock (if limited)

**Postconditions:**
1. User balance is reduced by total cost
2. Item stock is reduced (if limited)
3. Purchase record is created
4. Audit log entry is created

### `GetShopItems(chatID int64) ([]ShopItem, error)`
Retrieves available shop items for a chat.

**Parameters:**
- `chatID` - Chat ID to filter items (0 for global items)

**Returns:**
- `[]ShopItem` - List of available items
- `error` - Possible errors:
  - `ErrDatabaseError`

### `SetCurrencyName(chatID int64, name string) error`
Sets the currency name for a chat.

**Parameters:**
- `chatID` - Chat to configure
- `name` - New currency name (e.g. "Gold Coins")

**Returns:**
- `error` - Possible errors:
  - `ErrInvalidInput`
  - `ErrDatabaseError`

## Transaction Handling
All purchase operations are executed within a transaction to ensure:
- Atomic balance updates
- Consistent inventory tracking
- Reliable audit logging

## Error Handling
The service returns typed errors that clearly indicate failure reasons:
- Input validation failures
- Business rule violations
- System/DB errors