package testdata

import (
	"time"

	"github.com/ahsmha/discount-service/internal/models"
	"github.com/shopspring/decimal"
)

// GetSampleProducts returns sample products for testing
func GetSampleProducts() []models.Product {
	return []models.Product{
		{
			ID:           "prod-001",
			Brand:        "PUMA",
			BrandTier:    models.BrandTierPremium,
			Category:     "T-shirts",
			BasePrice:    decimal.NewFromInt(1000),
			CurrentPrice: decimal.NewFromInt(600), // After 40% brand discount
		},
		{
			ID:           "prod-002",
			Brand:        "Nike",
			BrandTier:    models.BrandTierPremium,
			Category:     "Shoes",
			BasePrice:    decimal.NewFromInt(5000),
			CurrentPrice: decimal.NewFromInt(5000), // No discount applied yet
		},
		{
			ID:           "prod-003",
			Brand:        "Adidas",
			BrandTier:    models.BrandTierPremium,
			Category:     "T-shirts",
			BasePrice:    decimal.NewFromInt(800),
			CurrentPrice: decimal.NewFromInt(800), // No discount applied yet
		},
		{
			ID:           "prod-004",
			Brand:        "Zara",
			BrandTier:    models.BrandTierRegular,
			Category:     "Jeans",
			BasePrice:    decimal.NewFromInt(1200),
			CurrentPrice: decimal.NewFromInt(1200), // No discount applied yet
		},
	}
}

// GetSampleCartItems returns sample cart items for testing
func GetSampleCartItems() []models.CartItem {
	products := GetSampleProducts()
	return []models.CartItem{
		{
			Product:  products[0], // PUMA T-shirt
			Quantity: 2,
			Size:     "M",
		},
		{
			Product:  products[1], // Nike Shoes
			Quantity: 1,
			Size:     "42",
		},
		{
			Product:  products[2], // Adidas T-shirt
			Quantity: 1,
			Size:     "L",
		},
	}
}

// GetSampleCustomers returns sample customers for testing
func GetSampleCustomers() []models.CustomerProfile {
	return []models.CustomerProfile{
		{
			ID:   "cust-001",
			Tier: "premium",
		},
		{
			ID:   "cust-002",
			Tier: "regular",
		},
		{
			ID:   "cust-003",
			Tier: "bronze",
		},
	}
}

// GetSamplePaymentInfo returns sample payment information
func GetSamplePaymentInfo() []models.PaymentInfo {
	iciciBank := "ICICI"
	creditCard := "CREDIT"

	return []models.PaymentInfo{
		{
			Method:   "CARD",
			BankName: &iciciBank,
			CardType: &creditCard,
		},
		{
			Method: "UPI",
		},
		{
			Method: "COD",
		},
	}
}

// GetSampleDiscounts returns sample discounts for testing the multiple discount scenario
func GetSampleDiscounts() []models.Discount {
	now := time.Now()
	validFrom := now.Add(-24 * time.Hour)
	validTo := now.Add(30 * 24 * time.Hour)

	return []models.Discount{
		{
			ID:            "disc-001",
			Name:          "PUMA Brand Discount - Min 40% off",
			Type:          models.DiscountTypeBrand,
			Value:         decimal.NewFromInt(40),
			IsPercentage:  true,
			MinAmount:     decimal.NewFromInt(500),
			MaxAmount:     decimal.Zero, // No max limit
			ApplicableTo:  []string{"PUMA"},
			ExcludedItems: []string{},
			CustomerTiers: []string{},
			Code:          "",
			ValidFrom:     validFrom,
			ValidTo:       validTo,
			IsActive:      true,
			UsageLimit:    0, // No usage limit
			UsedCount:     0,
			Priority:      100,
		},
		{
			ID:            "disc-002",
			Name:          "T-shirts Category Discount - Extra 10% off",
			Type:          models.DiscountTypeCategory,
			Value:         decimal.NewFromInt(10),
			IsPercentage:  true,
			MinAmount:     decimal.Zero,
			MaxAmount:     decimal.NewFromInt(200), // Max 200 discount
			ApplicableTo:  []string{"T-shirts"},
			ExcludedItems: []string{},
			CustomerTiers: []string{},
			Code:          "",
			ValidFrom:     validFrom,
			ValidTo:       validTo,
			IsActive:      true,
			UsageLimit:    0,
			UsedCount:     0,
			Priority:      90,
		},
		{
			ID:            "disc-003",
			Name:          "ICICI Bank Offer - 10% instant discount",
			Type:          models.DiscountTypeBank,
			Value:         decimal.NewFromInt(10),
			IsPercentage:  true,
			MinAmount:     decimal.NewFromInt(1000),
			MaxAmount:     decimal.NewFromInt(500), // Max 500 discount
			ApplicableTo:  []string{"ICICI"},
			ExcludedItems: []string{},
			CustomerTiers: []string{},
			Code:          "",
			ValidFrom:     validFrom,
			ValidTo:       validTo,
			IsActive:      true,
			UsageLimit:    0,
			UsedCount:     0,
			Priority:      80,
		},
		{
			ID:            "disc-004",
			Name:          "SUPER69 Voucher - 69% off",
			Type:          models.DiscountTypeVoucher,
			Value:         decimal.NewFromInt(69),
			IsPercentage:  true,
			MinAmount:     decimal.NewFromInt(2000),
			MaxAmount:     decimal.NewFromInt(1000), // Max 1000 discount
			ApplicableTo:  []string{},
			ExcludedItems: []string{"Electronics", "Luxury"},
			CustomerTiers: []string{"premium"},
			Code:          "SUPER69",
			ValidFrom:     validFrom,
			ValidTo:       validTo,
			IsActive:      true,
			UsageLimit:    100,
			UsedCount:     0,
			Priority:      70,
		},
		{
			ID:            "disc-005",
			Name:          "Nike Brand Discount - 30% off",
			Type:          models.DiscountTypeBrand,
			Value:         decimal.NewFromInt(30),
			IsPercentage:  true,
			MinAmount:     decimal.NewFromInt(1000),
			MaxAmount:     decimal.Zero,
			ApplicableTo:  []string{"Nike"},
			ExcludedItems: []string{},
			CustomerTiers: []string{},
			Code:          "",
			ValidFrom:     validFrom,
			ValidTo:       validTo,
			IsActive:      true,
			UsageLimit:    0,
			UsedCount:     0,
			Priority:      95,
		},
		{
			ID:            "disc-006",
			Name:          "Premium Customer Discount - 15% off",
			Type:          models.DiscountTypeVoucher,
			Value:         decimal.NewFromInt(15),
			IsPercentage:  true,
			MinAmount:     decimal.NewFromInt(500),
			MaxAmount:     decimal.NewFromInt(300),
			ApplicableTo:  []string{},
			ExcludedItems: []string{},
			CustomerTiers: []string{"premium"},
			Code:          "PREMIUM15",
			ValidFrom:     validFrom,
			ValidTo:       validTo,
			IsActive:      true,
			UsageLimit:    0,
			UsedCount:     0,
			Priority:      60,
		},
	}
}

// GetMultipleDiscountScenario returns data for testing the multiple discount scenario
// PUMA T-shirt with "Min 40% off" + Additional 10% off on T-shirts category + ICICI bank offer of 10% instant discount
func GetMultipleDiscountScenario() ([]models.CartItem, models.CustomerProfile, *models.PaymentInfo) {
	products := GetSampleProducts()
	customers := GetSampleCustomers()
	paymentInfo := GetSamplePaymentInfo()

	// Cart with PUMA T-shirt
	cartItems := []models.CartItem{
		{
			Product:  products[0], // PUMA T-shirt - already has 40% off (CurrentPrice: 600)
			Quantity: 2,
			Size:     "M",
		},
	}

	// Premium customer
	customer := customers[0]

	// ICICI bank payment
	payment := &paymentInfo[0]

	return cartItems, customer, payment
}

// GetComplexDiscountScenario returns a more complex scenario with multiple items and discounts
func GetComplexDiscountScenario() ([]models.CartItem, models.CustomerProfile, *models.PaymentInfo) {
	products := GetSampleProducts()
	customers := GetSampleCustomers()
	paymentInfo := GetSamplePaymentInfo()

	// Cart with multiple items
	cartItems := []models.CartItem{
		{
			Product:  products[0], // PUMA T-shirt
			Quantity: 1,
			Size:     "M",
		},
		{
			Product:  products[1], // Nike Shoes
			Quantity: 1,
			Size:     "42",
		},
		{
			Product:  products[2], // Adidas T-shirt
			Quantity: 1,
			Size:     "L",
		},
	}

	// Premium customer
	customer := customers[0]

	// ICICI bank payment
	payment := &paymentInfo[0]

	return cartItems, customer, payment
}
