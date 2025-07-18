# Technical Assumptions and Decisions 📋

## Overview

This document outlines the key technical assumptions and decisions made during the implementation of the Unifize Discount Service. These decisions were made to balance code quality, maintainability, performance, and extensibility within the 3-hour time constraint.

## 🏗️ Architecture Decisions

### 1. Clean Architecture Implementation

**Decision**: Implemented a simplified Clean Architecture with 3 main layers
- **Service Layer**: Business logic and discount calculation
- **Repository Layer**: Data access abstraction
- **Models Layer**: Domain entities and value objects

**Assumption**: The service will grow in complexity, so proper separation of concerns is essential from the start.

**Alternative Considered**: Flat structure for simplicity, but rejected due to extensibility requirements.

### 2. Dependency Injection Pattern

**Decision**: Used constructor-based dependency injection without a DI container
```go
func NewDiscountService(discountRepo repository.DiscountRepository) DiscountService {
    return &discountService{
        discountRepo: discountRepo,
        calculator:   NewDiscountCalculator(),
        validator:    NewDiscountValidator(),
    }
}
```

**Assumption**: Manual DI is sufficient for the current scope and maintains simplicity.

**Alternative Considered**: Wire or Dig framework, but manual DI was chosen to reduce dependencies.

### 3. Interface-First Design

**Decision**: Defined interfaces before implementations
```go
type DiscountService interface {
    CalculateCartDiscounts(ctx context.Context, ...) (*models.DiscountedPrice, error)
    ValidateDiscountCode(ctx context.Context, ...) (bool, error)
}
```

**Assumption**: Interfaces enable easy testing, mocking, and future implementations.

## 💾 Data Storage Decisions

### 1. In-Memory Repository

**Decision**: Implemented `MemoryDiscountRepository` for data storage
```go
type memoryDiscountRepository struct {
    discounts map[string]*models.Discount
    codeIndex map[string]string
    mu        sync.RWMutex
}
```

**Assumption**: For demo/testing purposes, in-memory storage is sufficient. Production would use persistent storage.

**Alternative Considered**: Mock repository, but functional in-memory repo provides better demonstration.

### 2. Thread Safety

**Decision**: Used `sync.RWMutex` for concurrent access
```go
func (r *memoryDiscountRepository) GetActiveDiscounts(ctx context.Context) ([]models.Discount, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    // ... implementation
}
```

**Assumption**: Service may handle concurrent requests, so thread safety is essential.

## 🧮 Discount Calculation Logic

### 1. Discount Stacking Order

**Decision**: Implemented specific discount application order:
1. Brand/Category discounts (highest priority)
2. Voucher codes (medium priority)
3. Bank offers (lowest priority)

**Assumption**: This order maximizes customer savings while being business-friendly.

**Business Logic**: 
```go
// Step 1: Apply brand/category discounts
// Step 2: Apply voucher codes
// Step 3: Apply bank offers
```

### 2. Decimal Precision

**Decision**: Used `shopspring/decimal` package for financial calculations
```go
type Product struct {
    BasePrice    decimal.Decimal `json:"base_price"`
    CurrentPrice decimal.Decimal `json:"current_price"`
}
```

**Assumption**: Financial calculations require exact precision to avoid rounding errors.

**Alternative Considered**: `float64` was rejected due to precision issues in financial calculations.

### 3. Discount Calculation Strategy

**Decision**: Percentage discounts calculated on current cart value, fixed discounts applied as-is
```go
if discount.IsPercentage {
    discountAmount = totalEligibleAmount.Mul(discount.Value).Div(decimal.NewFromInt(100))
} else {
    discountAmount = discount.Value
}
```

**Assumption**: This provides flexibility for both percentage and fixed-amount discounts.

## 🧪 Testing Decisions

### 1. Table-Driven Tests

**Decision**: Implemented comprehensive table-driven tests
```go
tests := []struct {
    name                string
    cartItems           []models.CartItem
    customer            models.CustomerProfile
    expectedFinalPrice  decimal.Decimal
    expectedDiscounts   int
    expectError         bool
}{
    // ... test cases
}
```

**Assumption**: Table-driven tests provide better coverage and maintainability.

### 2. Test Data Management

**Decision**: Created `testdata` package with realistic scenarios
```go
func GetMultipleDiscountScenario() ([]models.CartItem, models.CustomerProfile, *models.PaymentInfo)
```

**Assumption**: Realistic test data scenarios help validate business logic accurately.

### 3. Integration Testing

**Decision**: Included end-to-end integration tests that exercise the full discount calculation pipeline

**Assumption**: Integration tests catch issues that unit tests might miss.

## 🔧 Error Handling Decisions

### 1. Custom Error Types

**Decision**: Implemented custom error types for better error handling
```go
type ValidationError struct {
    Message string
}

func NewValidationError(message string) error {
    return ValidationError{Message: message}
}
```

**Assumption**: Custom error types enable better error categorization and handling.

### 2. Context Usage

**Decision**: Used `context.Context` throughout the service layer
```go
func (ds *discountService) CalculateCartDiscounts(ctx context.Context, ...)
```

**Assumption**: Context enables timeout handling, cancellation, and request tracing in production.

## 📦 Package Organization

### 1. Internal vs Pkg Structure

**Decision**: Used `internal/` for private packages, `pkg/` for potentially reusable components
```
internal/     # Private to this service
├── models/
├── service/
└── repository/

pkg/          # Potentially reusable
└── errors/
```

**Assumption**: Following Go module conventions improves maintainability.

### 2. Domain-Driven Package Naming

**Decision**: Named packages by domain responsibility (`service`, `repository`) not technical role (`handlers`, `database`)

**Assumption**: Domain-driven naming is more maintainable and follows Go idioms.

## 🎯 Business Logic Assumptions

### 1. Discount Priority System

**Decision**: Implemented priority-based discount ordering
```go
type Discount struct {
    Priority int `json:"priority"`
}

// Sort by priority (highest first)
sort.Slice(applicable, func(i, j int) bool {
    return applicable[i].Priority > applicable[j].Priority
})
```

**Assumption**: Business may want to control which discounts take precedence.

### 2. Customer Tier System

**Decision**: Implemented flexible customer tier validation
```go
func (d *Discount) IsApplicableToCustomer(customer CustomerProfile) bool {
    if len(d.CustomerTiers) == 0 {
        return true // No tier restrictions
    }
    return d.isInList(customer.Tier, d.CustomerTiers)
}
```

**Assumption**: Different customer tiers should have access to different discount levels.

### 3. Minimum/Maximum Discount Limits

**Decision**: Implemented both minimum order amounts and maximum discount caps
```go
type Discount struct {
    MinAmount decimal.Decimal `json:"min_amount"`
    MaxAmount decimal.Decimal `json:"max_amount"`
}
```

**Assumption**: Business needs control over discount limits to manage profitability.

## 🚀 Performance Considerations

### 1. Discount Filtering Strategy

**Decision**: Filter discounts by type before validation to reduce processing
```go
brandCategoryDiscounts := ds.filterDiscountsByType(allDiscounts, []models.DiscountType{
    models.DiscountTypeBrand,
    models.DiscountTypeCategory,
})
```

**Assumption**: Filtering early reduces unnecessary validation overhead.

### 2. Memory-Efficient Product Validation

**Decision**: Validate product applicability per discount without creating intermediate collections

**Assumption**: Large product catalogs require efficient validation logic.

## 🔮 Future Extensibility

### 1. Strategy Pattern for Calculators

**Decision**: Separated calculation logic into dedicated calculator interfaces
```go
type DiscountCalculator interface {
    CalculateDiscountAmount(discount *models.Discount, ...) decimal.Decimal
}
```

**Assumption**: Future discount types may require different calculation strategies.

### 2. Pluggable Validation Rules

**Decision**: Isolated validation logic in separate validator
```go
type DiscountValidator interface {
    ValidateDiscount(discount *models.Discount, ...) bool
}
```

**Assumption**: Validation rules may become more complex and require customization.

## 🎨 Code Quality Decisions

### 1. SOLID Principles Implementation

**Single Responsibility**: Each struct has one clear responsibility
**Open/Closed**: New discount types can be added without modifying existing code
**Liskov Substitution**: All repository implementations are interchangeable
**Interface Segregation**: Small, focused interfaces
**Dependency Inversion**: Depends on abstractions, not concretions

### 2. Go Best Practices

- Used Go modules for dependency management
- Followed Go naming conventions
- Implemented proper error handling
- Used interfaces for testability
- Applied Go formatting standards

## 📊 Monitoring and Observability

### 1. Usage Tracking

**Decision**: Implemented discount usage tracking
```go
func (ds *discountService) IncrementUsageCount(ctx context.Context, id string) error
```

**Assumption**: Business needs to track discount usage for analytics.

### 2. Detailed Error Messages

**Decision**: Provided detailed, contextual error messages
```go
return fmt.Errorf("failed to apply brand/category discounts: %w", err)
```

**Assumption**: Detailed errors help with debugging and monitoring.

## 🔄 Trade-offs Made

### 1. Simplicity vs Completeness

**Trade-off**: Chose implementation simplicity over complete feature coverage
**Rationale**: 3-hour time constraint required focusing on core functionality

### 2. In-Memory vs Persistent Storage

**Trade-off**: Used in-memory storage for demo purposes instead of database
**Rationale**: Faster implementation, sufficient for demonstration requirements

### 3. Manual DI vs Framework

**Trade-off**: Used manual dependency injection instead of DI framework
**Rationale**: Reduces complexity and external dependencies

## 🎯 Success Criteria Met

✅ **SOLID Principles**: All five principles implemented
✅ **Clean Architecture**: Clear separation of concerns
✅ **Extensibility**: Easy to add new discount types
✅ **Testability**: Comprehensive test coverage
✅ **Go Idioms**: Follows Go best practices
✅ **Business Logic**: Correctly implements discount stacking
✅ **Error Handling**: Proper error types and messages
✅ **Documentation**: Comprehensive documentation provided

---

These assumptions and decisions form the foundation of the Unifize Discount Service implementation, balancing code quality, maintainability, and business requirements within the given constraints.