# API Usage Examples

## Basic Workflows

### 1. Task Creation and Timer Management
```go
// Create a new task
task := builders.NewTaskBuilder().
    WithTitle("Study for exam").
    WithDuration(45). // minutes
    WithReward(5.0). // currency
    Build()

// Start timer for task  
timer, err := taskService.StartTimer(userID, task.ID)
if err != nil {
    // Handle validation errors
}

// Pause/resume timer
err = timerService.Pause(timer.ID)
err = timerService.Resume(timer.ID)

// Complete task
err = taskService.Complete(userID, task.ID)
```

### 2. Shop Purchases
```go
// Get available items
items := shopService.ListItems()

// Purchase an item
receipt, err := shopService.Purchase(userID, items[0].ID, 1)
if err != nil {
    switch {
    case errors.Is(err, ErrInsufficientFunds):
        // Handle insufficient balance
    case errors.Is(err, ErrOutOfStock):
        // Handle out of stock
    }
}
```

## Integration Patterns

### Database Transactions
```go
err := txManager.Transaction(ctx, func(tx ports.Transaction) error {
    // Withdraw funds
    err := userRepo.UpdateBalance(tx, userID, -item.Price)
    if err != nil {
        return err
    }

    // Record purchase
    return purchaseRepo.Create(tx, purchase)
})
```

### Event-Driven Workflows
```go
// Subscribe to task completion events
eventBus.Subscribe("task.completed", func(e Event) {
    // Award XP
    userService.AddXP(e.UserID, e.Task.Reward)
    
    // Trigger next task if in a sequence
    if e.Task.NextTaskID != nil {
        taskService.Activate(e.UserID, *e.Task.NextTaskID)
    }
})