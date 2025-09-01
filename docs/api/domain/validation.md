# Validation Rules

## Core Validation Patterns

### Value Objects
1. **Decimal Values**
   - Must be positive
   - Precision limited to 4 decimal places
   - Parsed from string with locale support

2. **Identifiers**
   - UUID format validation
   - Null checks for required IDs
   - Referential integrity checks

### Entities

#### User
```go
type User struct {
    Username  string  // required, alphanumeric, 3-32 chars
    Email     string  // valid email format
    Timezone  string  // valid IANA timezone
    Balance   Decimal // non-negative
}
```

#### Task
```go
type Task struct {
    Title       string  // required, 1-100 chars
    Description string  // optional, max 500 chars
    Duration    int     // positive, in minutes
    Reward      Decimal // non-negative
    DueDate     time.Time // must be in future if set
}
```

### Business Rules

1. **Timer Start Validation**
   - Task must be active
   - No other active timers for user
   - Sufficient time remaining in task window

2. **Purchase Validation**
```go
func ValidatePurchase(user User, item ShopItem, quantity int) error {
    if quantity <= 0 {
        return ErrInvalidQuantity
    }
    if item.Limited && quantity > item.Stock {
        return ErrInsufficientStock
    }
    if user.Balance.LessThan(item.Price.Mul(quantity)) {
        return ErrInsufficientFunds
    }
    return nil
}
```

3. **Task Completion**
   - All prerequisites completed
   - Within valid completion window
   - Not already completed

## Error Hierarchy
```mermaid
classDiagram
    ValidationError <|-- DomainError
    ValidationError <|-- BusinessRuleError
    
    class ValidationError {
        +string Field
        +string Message
    }
    
    class BusinessRuleError {
        +string Rule
        +string Context
    }