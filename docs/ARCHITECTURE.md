# Architecture Documentation 🏗️

## System Overview

The Unifize Discount Service is designed as a modular, extensible system following Clean Architecture principles. It handles complex discount calculations for e-commerce platforms with support for multiple discount types and sophisticated business rules.

## 🎯 Core Architecture Principles

### 1. Clean Architecture Layers

```
┌─────────────────────────────────────────────┐
│                Presentation                 │  (Future HTTP/gRPC API)
├─────────────────────────────────────────────┤
│               Application                   │  cmd/server
├─────────────────────────────────────────────┤
│              Business Logic                 │  internal/service
├─────────────────────────────────────────────┤
│               Data Access                   │  internal/repository
├─────────────────────────────────────────────┤
│              Infrastructure                 │  pkg/errors, external deps
└─────────────────────────────────────────────┘
```

### 2. SOLID Principles Implementation

#### Single Responsibility Principle (SRP)
- `DiscountService`: Orchestrates discount calculations
- `DiscountCalculator`: Handles discount amount calculations
- `DiscountValidator`: Manages discount validation rules
- `DiscountRepository`: Manages discount data access

#### Open/Closed Principle (OCP)
- New discount types can be added without modifying existing code
- New validation rules can be plugged in through interfaces
- New calculation strategies can be implemented independently

#### Liskov Substitution Principle (LSP)
- All repository implementations are interchangeable
- Different calculator implementations can be substituted
- Validator implementations can be swapped without breaking functionality

#### Interface Segregation Principle (ISP)
- Small, focused interfaces (DiscountCalculator, DiscountValidator)
- Clients only depend on methods they actually use
- No fat interfaces with unused methods

#### Dependency Inversion Principle (DIP)
- High-level modules depend on abstractions
- Concrete implementations depend on interfaces
- Dependencies are injected rather than created internally

## 📦 Package Structure

```
unifize-discount-service/
│
├── cmd/
│   └── server/                 # Application entry points
│       └── main.go            # Main application bootstrap
│
├── internal/                   # Private application packages
│   ├── models/                # Domain entities and value objects
│   │   └── discount.go       # Core domain models
│   ├── service/              # Business logic layer
│   │   ├── discount_service.go      # Main service orchestrator
│   │   ├── discount_calculator.go   # Calculation strategies
│   │   └── discount_validator.go    # Validation logic
│   ├── repository/           # Data access layer
│   │   ├── discount_repository.go       # Repository interface
│   │   └── memory_discount_repository.go # In-memory implementation
│   └── config/              # Configuration management
│
├── pkg/                      # Public, reusable packages
│   └── errors/              # Custom error types and utilities
│       └── errors.go
│
├── testdata/                 # Test data and scenarios
│   └── fake_data.go         # Sample data for testing
│
├── test/                    # Test utilities and integration tests
│   └── integration/         # End-to-end integration tests
│
└── docs/                    # Documentation
    ├── README.md
    ├── ASSUMPTIONS.md
    └── ARCHITECTURE.md      # This file
```

## 🔄 Data Flow Architecture

### Discount Calculation Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Request   │───▶│   Service   │───▶│ Repository  │
│  (Cart +    │    │             │    │             │
│  Customer + │    │             │    │             │
│  Payment)   │    │             │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
                          │
                          ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  Response   │◀───│ Calculator  │───▶│  Validator  │
│ (Discounted │    │             │    │             │
│   Price)    │    │             │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
```

### Detailed Processing Steps

1. **Input Validation**: Validate cart items, customer profile, and payment info
2. **Discount Retrieval**: Fetch all active discounts from repository
3. **Filtering**: Filter discounts by type and applicability
4. **Validation**: Validate each discount against business rules
5. **Prioritization**: Sort discounts by priority and type
6. **Calculation**: Apply discounts in the correct stacking order
7. **Response Formation**: Build final discounted price response

## 🎨 Design Patterns Used

### 1. Repository Pattern

```go
type DiscountRepository interface {
    GetActiveDiscounts(ctx context.Context) ([]models.Discount, error)
    GetDiscountByCode(ctx context.Context, code string) (*models.Discount, error)
    // ... other methods
}
```

**Benefits**:
- Abstracts data access layer
- Enables easy testing with mock repositories
- Supports multiple storage backends

### 2. Strategy Pattern

```go
type DiscountCalculator interface {
    CalculateDiscountAmount(discount *models.Discount, ...) decimal.Decimal
}
```

**Benefits**:
- Different calculation strategies for different discount types
- Easy to extend with new calculation methods
- Testable in isolation

### 3. Factory Pattern

```go
func NewDiscountService(discountRepo repository.DiscountRepository) DiscountService {
    return &discountService{
        discountRepo: discountRepo,
        calculator:   NewDiscountCalculator(),
        validator:    NewDiscountValidator(),
    }
}
```

**Benefits**:
- Centralized object creation
- Dependency injection setup
- Easy to modify construction logic

### 4. Template Method Pattern

```go
// Discount application follows a template:
// 1. Apply brand/category discounts
// 2. Apply voucher codes  
// 3. Apply bank offers
```

**Benefits**:
- Consistent discount processing order
- Easy to modify individual steps
- Clear business logic flow

## 🏛️ Domain Model

### Core Entities

```go
// Product represents items in the e-commerce catalog
type Product struct {
    ID           string
    Brand        string
    BrandTier    BrandTier
    Category     string
    BasePrice    decimal.Decimal
    CurrentPrice decimal.Decimal
}

// Discount represents a discount rule with business logic
type Discount struct {
    ID             string
    Name           string
    Type           DiscountType
    Value          decimal.Decimal
    IsPercentage   bool
    MinAmount      decimal.Decimal
    MaxAmount      decimal.Decimal
    ApplicableTo   []string
    ExcludedItems  []string
    CustomerTiers  []string
    // ... other fields
}

// CartItem represents an item in the shopping cart
type CartItem struct {
    Product  Product
    Quantity int
    Size     string
}
```

### Value Objects

```go
// DiscountedPrice represents the final pricing result
type DiscountedPrice struct {
    OriginalPrice    decimal.Decimal
    FinalPrice       decimal.Decimal
    AppliedDiscounts map[string]decimal.Decimal
    Message          string
}
```

## 🔧 Configuration Architecture

### Environment-Based Configuration

```go
type Config struct {
    Database DatabaseConfig
    Cache    CacheConfig
    API      APIConfig
    Logging  LoggingConfig
}
```

### Configuration Loading Strategy

1. **Default values**: Sensible defaults for all configurations
2. **Environment variables**: Override defaults with env vars
3. **Configuration files**: Support for JSON/YAML config files
4. **Command line flags**: Final override for specific deployments

## 🚀 Scalability Considerations

### Horizontal Scaling

- **Stateless service design**: No server-side state storage
- **Database connection pooling**: Efficient database resource usage
- **Cache-friendly**: Discount rules can be cached effectively
- **Load balancer ready**: Service can run multiple instances

### Performance Optimizations

- **Early filtering**: Filter discounts before expensive operations
- **Efficient data structures**: Use maps for O(1) lookups where possible
- **Minimal allocations**: Reuse objects and avoid unnecessary allocations
- **Concurrent-safe**: Thread-safe operations for multi-user access

### Caching Strategy

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    Client   │───▶│   Service   │───▶│    Cache    │
└─────────────┘    └─────────────┘    └─────────────┘
                          │                   │
                          ▼                   ▼
                   ┌─────────────┐    ┌─────────────┐
                   │  Database   │    │   Metrics   │
                   └─────────────┘    └─────────────┘
```

## 🔐 Security Architecture

### Input Validation

- **Request validation**: Validate all incoming requests
- **SQL injection prevention**: Use parameterized queries
- **XSS prevention**: Sanitize all user inputs
- **Rate limiting**: Prevent abuse and DoS attacks

### Authorization

- **Customer verification**: Verify customer identity
- **Discount code validation**: Prevent unauthorized code usage
- **Usage limits**: Enforce discount usage limits
- **Audit logging**: Track all discount applications

## 📊 Monitoring and Observability

### Metrics Collection

```go
type Metrics struct {
    DiscountApplications  counter
    DiscountValidations   counter
    CalculationDuration   histogram
    ErrorRates           counter
}
```

### Logging Strategy

- **Structured logging**: JSON format for easy parsing
- **Correlation IDs**: Track requests across services
- **Performance logging**: Log calculation times
- **Business metrics**: Track discount usage patterns

### Health Checks

```go
type HealthCheck struct {
    DatabaseConnection bool
    CacheConnection    bool
    ServiceHealth      bool
}
```

## 🧪 Testing Architecture

### Test Categories

1. **Unit Tests**: Individual component testing
2. **Integration Tests**: Component interaction testing
3. **End-to-End Tests**: Full workflow testing
4. **Performance Tests**: Load and stress testing
5. **Contract Tests**: API contract validation

### Test Data Management

```go
// testdata package provides realistic scenarios
func GetMultipleDiscountScenario() ([]models.CartItem, models.CustomerProfile, *models.PaymentInfo)
func GetComplexDiscountScenario() ([]models.CartItem, models.CustomerProfile, *models.PaymentInfo)
```

### Mocking Strategy

- **Repository mocking**: In-memory implementations for testing
- **External service mocking**: Mock external dependencies
- **Time mocking**: Control time-dependent discount validations

## 🔄 Future Extensions

### Planned Enhancements

1. **HTTP/gRPC API**: RESTful and gRPC interfaces
2. **Database Integration**: PostgreSQL/MySQL repository implementations
3. **Redis Caching**: Distributed caching for performance
4. **Event Sourcing**: Track discount application history
5. **A/B Testing**: Support for discount strategy testing
6. **Machine Learning**: Intelligent discount recommendations

### Extension Points

```go
// New discount types can be added easily
type CustomDiscountCalculator struct{}
func (c *CustomDiscountCalculator) CalculateDiscountAmount(...) decimal.Decimal

// New validation rules can be plugged in
type CustomDiscountValidator struct{}
func (v *CustomDiscountValidator) ValidateDiscount(...) bool
```

## 📋 Deployment Architecture

### Container Strategy

```dockerfile
# Multi-stage build for minimal image size
FROM golang:1.21-alpine AS builder
# ... build steps

FROM alpine:latest
# ... runtime setup
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: discount-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: discount-service
  template:
    # ... pod specification
```

### Service Mesh Integration

- **Istio compatibility**: Ready for service mesh deployment
- **Circuit breakers**: Fault tolerance patterns
- **Retry policies**: Automatic retry configurations
- **Load balancing**: Intelligent traffic distribution

## 🎯 Quality Attributes

### Maintainability
- Clear separation of concerns
- Comprehensive documentation
- Consistent coding standards
- Automated testing

### Reliability  
- Error handling at all layers
- Graceful degradation
- Circuit breaker patterns
- Health checks

### Performance
- Efficient algorithms
- Minimal memory allocation
- Database query optimization
- Caching strategies

### Security
- Input validation
- Authentication support
- Authorization checks
- Audit logging

### Scalability
- Horizontal scaling support
- Database partitioning ready
- Stateless design
- Microservice architecture

---

This architecture provides a solid foundation for building a comprehensive, scalable, and maintainable discount service that can evolve with business needs while maintaining code quality and performance.