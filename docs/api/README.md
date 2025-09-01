# API Contracts Documentation

This directory contains documentation for the service layer APIs and domain contracts.

## Structure
- `services/` - Service layer interfaces and contracts
- `repositories/` - Repository interfaces
- `domain/` - Domain rules and invariants
- `validation/` - Input validation rules
- `examples/` - Usage examples

## Services Overview

### ShopService
Manages shop operations including:
- Item management
- Purchases
- Inventory

### TaskService
Handles task-related operations:
- Task creation/modification
- Timer management
- Reward calculation

## Domain Rules
Core business rules that must be maintained:
1. Single active timer per task/user
2. Atomic balance updates
3. Idempotent operations
4. Timezone-aware scheduling
5. Reward tier calculations

## Validation Rules
Input validation requirements for:
- Decimal precision
- String lengths
- Enum values
- Required fields
- State transitions