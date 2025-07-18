# Unifize Discount Service 🛍️

A comprehensive, extensible discount calculation service built in Go following SOLID principles and clean architecture patterns.

## 🎯 Overview

This service handles e-commerce discount calculations for fashion retail, supporting multiple discount types with proper stacking order:

1. **Brand-specific discounts** (e.g., "Min 40% off on PUMA")
2. **Category-specific deals** (e.g., "Extra 10% off on T-shirts")
3. **Voucher codes** (e.g., "SUPER69" for 69% off)
4. **Bank card offers** (e.g., "10% instant discount on ICICI Bank cards")

## 🏗️ Architecture

The service follows **Clean Architecture** principles with clear separation of concerns:

```
cmd/
├── server/           # Application entry point
internal/
├── models/           # Domain entities and value objects
├── service/          # Business logic layer
├── repository/       # Data access layer
├── config/           # Configuration management
pkg/
├── errors/           # Custom error types
testdata/
├── fake_data.go      # Test data scenarios
```

### Key Design Patterns

- **Dependency Injection**: Constructor-based DI for loose coupling
- **Strategy Pattern**: Pluggable discount calculation strategies
- **Repository Pattern**: Abstracted data access layer
- **Factory Pattern**: Discount creation and validation
- **Interface Segregation**: Small, focused interfaces

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Make (optional, for convenience commands)

### Installation

```bash
# Clone the repository
git clone https://github.com/unifize/discount-service.git
cd discount-service

# Install dependencies
go mod download

# Run the application
make run
```

### Alternative Installation

```bash
# Direct run without Make
go run cmd/server/main.go

# Build and run
go build -o bin/discount-service cmd/server/main.go
./bin/discount-service
```

## 🧪 Testing

### Run All Tests

```bash
make test
```

### Run Tests with Coverage

```bash
make test-coverage
```

### Run Specific Test Suite

```bash
# Unit tests
go test ./internal/service/... -v

# Integration tests
go test ./test/integration/... -v

# Table-driven tests for discount calculations
go test -run TestDiscountService_CalculateCartDiscounts -v
```

### Test Coverage Report

```bash
make view-coverage
```

## 📋 Usage Examples

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/unifize/discount-service/internal/service"
    "github.com/unifize/discount-service/internal/repository"
    "github.com/unifize/discount-service/internal/models"
)

func main() {
    // Initialize repository and service
    repo := repository.NewMemoryDiscountRepository()
    discountService := service.NewDiscountService(repo)
    
    // Prepare cart items
    cartItems := []models.CartItem{
        {
            Product: models.Product{
                ID:           "prod-001",
                Brand:        "PUMA",
                Category:     "T-shirts",
                BasePrice:    decimal.NewFromInt(1000),
                CurrentPrice: decimal.NewFromInt(1000),
            },
            Quantity: 2,
            Size:     "M",
        },
    }
    
    // Customer information
    customer := models.CustomerProfile{
        ID:   "cust-001",
        Tier: "premium",
    }
    
    // Payment information
    bankName := "ICICI"
    paymentInfo := &models.PaymentInfo{
        Method:   "CARD",
        BankName: &bankName,
    }
    
    // Calculate discounts
    result, err := discountService.CalculateCartDiscounts(
        context.Background(),
        cartItems,
        customer,
        paymentInfo,
    )
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Original Price: ₹%s\n", result.OriginalPrice.String())
    fmt.Printf("Final Price: ₹%s\n", result.FinalPrice.String())
    fmt.Printf("Total Savings: ₹%s\n", result.GetTotalDiscount().String())
}
```

### Discount Code Validation

```go
// Validate a discount code
isValid, err := discountService.ValidateDiscountCode(
    context.Background(),
    "SUPER69",
    cartItems,
    customer,
)

if err != nil {
    fmt.Printf("Validation error: %v\n", err)
} else if isValid {
    fmt.Println("Discount code is valid!")
} else {
    fmt.Println("Invalid discount code")
}
```

## 🔧 Configuration

### Discount Types

The service supports four types of discounts:

1. **Brand Discounts** (`DiscountTypeBrand`)
2. **Category Discounts** (`DiscountTypeCategory`)
3. **Voucher Discounts** (`DiscountTypeVoucher`)
4. **Bank Discounts** (`DiscountTypeBank`)

### Discount Stacking Order

Discounts are applied in the following order:

1. **Brand/Category discounts** (applied first)
2. **Voucher codes** (applied second)
3. **Bank offers** (applied last)

### Sample Data

The service includes comprehensive test data scenarios:

```go
// Multiple discount scenario
cartItems, customer, paymentInfo := testdata.GetMultipleDiscountScenario()

// Complex scenario with multiple items
cartItems, customer, paymentInfo := testdata.GetComplexDiscountScenario()
```

## 🛠️ Development

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run all quality checks
make pre-commit
```

### Adding New Discount Types

1. **Define the discount type** in `models/discount.go`
2. **Update the calculator** in `service/discount_calculator.go`
3. **Update the validator** in `service/discount_validator.go`
4. **Add test cases** in the test files

### Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── models/
│   │   └── discount.go          # Domain models
│   ├── service/
│   │   ├── discount_service.go  # Main service interface
│   │   ├── discount_calculator.go # Discount calculation logic
│   │   └── discount_validator.go   # Validation logic
│   └── repository/
│       ├── discount_repository.go    # Repository interface
│       └── memory_discount_repository.go # In-memory implementation
├── pkg/
│   └── errors/
│       └── errors.go            # Custom error types
├── testdata/
│   └── fake_data.go            # Test data scenarios
├── test/
│   └── integration/
│       └── discount_service_test.go # Integration tests
├── docs/
│   ├── README.md               # This file
│   ├── ASSUMPTIONS.md          # Technical assumptions
│   └── ARCHITECTURE.md         # Architecture documentation
├── Makefile                    # Build and development commands
├── go.mod                      # Go module definition
└── .golangci.yml              # Linter configuration
```

## 📊 Performance Considerations

- **In-memory storage** for fast access during testing
- **Concurrent-safe** repository implementation
- **Efficient discount filtering** and validation
- **Minimal heap allocations** using decimal package
- **Table-driven tests** for comprehensive coverage

## 🔍 Monitoring and Logging

The service includes:

- **Structured error handling** with custom error types
- **Comprehensive logging** for discount applications
- **Usage tracking** for discount codes
- **Validation metrics** for business insights

## 🚀 Deployment

### Docker Support

```bash
# Build Docker image
docker build -t unifize/discount-service:latest .

# Run container
docker run -p 8080:8080 unifize/discount-service:latest
```

### Production Considerations

- Replace in-memory repository with persistent storage (Redis, PostgreSQL)
- Add HTTP/gRPC API layer for external access
- Implement caching for frequently accessed discounts
- Add circuit breaker patterns for external dependencies
- Configure proper logging and monitoring

## 📈 Metrics and Analytics

The service tracks:

- **Discount usage statistics**
- **Customer tier analysis**
- **Popular discount combinations**
- **Revenue impact calculations**

## 🤝 Contributing

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Write tests** for your changes
4. **Run quality checks**: `make pre-commit`
5. **Commit your changes**: `git commit -m 'Add amazing feature'`
6. **Push to the branch**: `git push origin feature/amazing-feature`
7. **Open a Pull Request**

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Robert C. Martin** for Clean Architecture principles
- **Go team** for excellent language design
- **Unifize team** for the challenging assignment
- **SOLID principles** for guiding the architecture

---

**Built with ❤️ using Go and Clean Architecture principles**