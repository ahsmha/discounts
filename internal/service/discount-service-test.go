package service

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ahsmha/discount-service/internal/models"
	"github.com/ahsmha/discount-service/internal/repository"
	"github.com/ahsmha/discount-service/testdata"
)

func TestDiscountService_CalculateCartDiscounts(t *testing.T) {
	// Setup
	repo := repository.NewMemoryDiscountRepository()
	memoryRepo := repo.(*repository.MemoryDiscountRepository)
	err := memoryRepo.SeedDiscounts(testdata.GetSampleDiscounts())
	require.NoError(t, err)

	service := NewDiscountService(repo)
	ctx := context.Background()

	tests := []struct {
		name               string
		cartItems          []models.CartItem
		customer           models.CustomerProfile
		paymentInfo        *models.PaymentInfo
		expectedFinalPrice decimal.Decimal
		expectedDiscounts  int
		expectError        bool
		errorMessage       string
	}{
		{
			name:         "Empty cart should return error",
			cartItems:    []models.CartItem{},
			customer:     testdata.GetSampleCustomers()[0],
			paymentInfo:  &testdata.GetSamplePaymentInfo()[0],
			expectError:  true,
			errorMessage: "cart is empty",
		},
		{
			name: "Multiple discount scenario - PUMA T-shirt with brand, category, and bank discounts",
			cartItems: []models.CartItem{
				{
					Product: models.Product{
						ID:           "prod-001",
						Brand:        "PUMA",
						BrandTier:    models.BrandTierPremium,
						Category:     "T-shirts",
						BasePrice:    decimal.NewFromInt(1000),
						CurrentPrice: decimal.NewFromInt(1000), // Start with base price
					},
					Quantity: 2,
					Size:     "M",
				},
			},
			customer:           testdata.GetSampleCustomers()[0],    // premium customer
			paymentInfo:        &testdata.GetSamplePaymentInfo()[0], // ICICI card
			expectedFinalPrice: decimal.NewFromFloat(972),           // Expected after all discounts
			expectedDiscounts:  3,                                   // Brand (40%) + Category (10%) + Bank (10%)
			expectError:        false,
		},
		{
			name: "Nike shoes with brand and bank discount only",
			cartItems: []models.CartItem{
				{
					Product: models.Product{
						ID:           "prod-002",
						Brand:        "Nike",
						BrandTier:    models.BrandTierPremium,
						Category:     "Shoes",
						BasePrice:    decimal.NewFromInt(5000),
						CurrentPrice: decimal.NewFromInt(5000),
					},
					Quantity: 1,
					Size:     "42",
				},
			},
			customer:           testdata.GetSampleCustomers()[0],    // premium customer
			paymentInfo:        &testdata.GetSamplePaymentInfo()[0], // ICICI card
			expectedFinalPrice: decimal.NewFromInt(3150),            // 5000 - 30%(1500) - 10%(350) = 3150
			expectedDiscounts:  2,                                   // Brand (30%) + Bank (10%)
			expectError:        false,
		},
		{
			name: "No payment info - should not apply bank discounts",
			cartItems: []models.CartItem{
				{
					Product: models.Product{
						ID:           "prod-001",
						Brand:        "PUMA",
						BrandTier:    models.BrandTierPremium,
						Category:     "T-shirts",
						BasePrice:    decimal.NewFromInt(1000),
						CurrentPrice: decimal.NewFromInt(1000),
					},
					Quantity: 1,
					Size:     "M",
				},
			},
			customer:           testdata.GetSampleCustomers()[0],
			paymentInfo:        nil,                     // No payment info
			expectedFinalPrice: decimal.NewFromInt(540), // 1000 - 40%(400) - 10%(60) = 540
			expectedDiscounts:  2,                       // Only brand and category discounts
			expectError:        false,
		},
		{
			name: "Regular customer with non-premium voucher",
			cartItems: []models.CartItem{
				{
					Product: models.Product{
						ID:           "prod-004",
						Brand:        "Zara",
						BrandTier:    models.BrandTierRegular,
						Category:     "Jeans",
						BasePrice:    decimal.NewFromInt(1200),
						CurrentPrice: decimal.NewFromInt(1200),
					},
					Quantity: 1,
					Size:     "32",
				},
			},
			customer:           testdata.GetSampleCustomers()[1],    // regular customer
			paymentInfo:        &testdata.GetSamplePaymentInfo()[0], // ICICI card
			expectedFinalPrice: decimal.NewFromInt(1080),            // Only bank discount applies: 1200 - 10%(120) = 1080
			expectedDiscounts:  1,                                   // Only bank discount
			expectError:        false,
		},
		{
			name: "Cart value below minimum for bank discount",
			cartItems: []models.CartItem{
				{
					Product: models.Product{
						ID:           "prod-003",
						Brand:        "Adidas",
						BrandTier:    models.BrandTierPremium,
						Category:     "T-shirts",
						BasePrice:    decimal.NewFromInt(500),
						CurrentPrice: decimal.NewFromInt(500),
					},
					Quantity: 1,
					Size:     "S",
				},
			},
			customer:           testdata.GetSampleCustomers()[0],
			paymentInfo:        &testdata.GetSamplePaymentInfo()[0], // ICICI card - but min amount is 1000
			expectedFinalPrice: decimal.NewFromInt(450),             // Only category discount: 500 - 10%(50) = 450
			expectedDiscounts:  1,                                   // Only category discount (bank discount requires min 1000)
			expectError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CalculateCartDiscounts(ctx, tt.cartItems, tt.customer, tt.paymentInfo)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectedDiscounts, len(result.AppliedDiscounts),
				"Expected %d discounts but got %d", tt.expectedDiscounts, len(result.AppliedDiscounts))

			assert.True(t, tt.expectedFinalPrice.Equal(result.FinalPrice),
				"Expected final price %s but got %s", tt.expectedFinalPrice.String(), result.FinalPrice.String())

			// Verify that discounts were applied correctly
			if len(result.AppliedDiscounts) > 0 {
				totalDiscount := result.GetTotalDiscount()
				expectedTotal := result.OriginalPrice.Sub(result.FinalPrice)
				assert.True(t, totalDiscount.Equal(expectedTotal),
					"Total discount %s does not match price difference %s",
					totalDiscount.String(), expectedTotal.String())
			}
		})
	}
}

func TestDiscountService_ValidateDiscountCode(t *testing.T) {
	// Setup
	repo := repository.NewMemoryDiscountRepository()
	memoryRepo := repo.(*repository.MemoryDiscountRepository)
	err := memoryRepo.SeedDiscounts(testdata.GetSampleDiscounts())
	require.NoError(t, err)

	service := NewDiscountService(repo)
	ctx := context.Background()

	cartItems := []models.CartItem{
		{
			Product:  testdata.GetSampleProducts()[0], // PUMA T-shirt
			Quantity: 1,
			Size:     "M",
		},
	}

	tests := []struct {
		name          string
		code          string
		cartItems     []models.CartItem
		customer      models.CustomerProfile
		expectedValid bool
		expectError   bool
		errorMessage  string
	}{
		{
			name:          "Valid voucher code for premium customer",
			code:          "SUPER69",
			cartItems:     cartItems,
			customer:      testdata.GetSampleCustomers()[0], // premium
			expectedValid: true,
			expectError:   false,
		},
		{
			name:          "Valid voucher code for premium customer with sufficient cart value",
			code:          "PREMIUM15",
			cartItems:     cartItems,
			customer:      testdata.GetSampleCustomers()[0], // premium
			expectedValid: true,
			expectError:   false,
		},
		{
			name:          "Invalid - voucher code for regular customer (premium required)",
			code:          "SUPER69",
			cartItems:     cartItems,
			customer:      testdata.GetSampleCustomers()[1], // regular
			expectedValid: false,
			expectError:   false,
		},
		{
			name:          "Invalid - non-existent voucher code",
			code:          "INVALID123",
			cartItems:     cartItems,
			customer:      testdata.GetSampleCustomers()[0],
			expectedValid: false,
			expectError:   false,
		},
		{
			name:          "Error - empty voucher code",
			code:          "",
			cartItems:     cartItems,
			customer:      testdata.GetSampleCustomers()[0],
			expectedValid: false,
			expectError:   true,
			errorMessage:  "discount code cannot be empty",
		},
		{
			name: "Invalid - cart value below minimum for SUPER69 (requires 2000)",
			code: "SUPER69",
			cartItems: []models.CartItem{
				{
					Product: models.Product{
						ID:           "prod-small",
						Brand:        "TestBrand",
						Category:     "TestCategory",
						BasePrice:    decimal.NewFromInt(500),
						CurrentPrice: decimal.NewFromInt(500),
					},
					Quantity: 1,
					Size:     "S",
				},
			},
			customer:      testdata.GetSampleCustomers()[0], // premium
			expectedValid: false,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid, err := service.ValidateDiscountCode(ctx, tt.code, tt.cartItems, tt.customer)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedValid, isValid,
				"Expected validation result %v but got %v", tt.expectedValid, isValid)
		})
	}
}

func TestDiscountService_Integration_MultipleDiscountScenario(t *testing.T) {
	// This test validates the exact scenario mentioned in the assignment:
	// PUMA T-shirt with "Min 40% off" + Additional 10% off on T-shirts category + ICICI bank offer of 10% instant discount

	repo := repository.NewMemoryDiscountRepository()
	memoryRepo := repo.(*repository.MemoryDiscountRepository)
	err := memoryRepo.SeedDiscounts(testdata.GetSampleDiscounts())
	require.NoError(t, err)

	service := NewDiscountService(repo)
	ctx := context.Background()

	// Get the multiple discount scenario data
	cartItems, customer, paymentInfo := testdata.GetMultipleDiscountScenario()

	// The PUMA T-shirt should start with base price for this test
	cartItems[0].Product.CurrentPrice = cartItems[0].Product.BasePrice

	result, err := service.CalculateCartDiscounts(ctx, cartItems, customer, paymentInfo)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify original price: PUMA T-shirt (1000) * 2 = 2000
	expectedOriginalPrice := decimal.NewFromInt(2000)
	assert.True(t, expectedOriginalPrice.Equal(result.OriginalPrice),
		"Expected original price %s but got %s", expectedOriginalPrice.String(), result.OriginalPrice.String())

	// Verify that all three types of discounts were applied
	assert.GreaterOrEqual(t, len(result.AppliedDiscounts), 3,
		"Expected at least 3 discounts but got %d", len(result.AppliedDiscounts))

	// Check for specific discount types
	hasGoodBrandDiscount := false
	hasCategoryDiscount := false
	hasBankDiscount := false

	for discountName := range result.AppliedDiscounts {
		if discountName == "PUMA Brand Discount - Min 40% off" {
			hasGoodBrandDiscount = true
		}
		if discountName == "T-shirts Category Discount - Extra 10% off" {
			hasCategoryDiscount = true
		}
		if discountName == "ICICI Bank Offer - 10% instant discount" {
			hasBankDiscount = true
		}
	}

	assert.True(t, hasGoodBrandDiscount, "PUMA brand discount should be applied")
	assert.True(t, hasCategoryDiscount, "T-shirts category discount should be applied")
	assert.True(t, hasBankDiscount, "ICICI bank discount should be applied")

	// Verify significant savings were achieved
	totalDiscount := result.GetTotalDiscount()
	assert.True(t, totalDiscount.GreaterThan(decimal.NewFromInt(800)),
		"Expected significant discount but got %s", totalDiscount.String())

	// Verify final price is reasonable
	assert.True(t, result.FinalPrice.LessThan(result.OriginalPrice),
		"Final price should be less than original price")
	assert.True(t, result.FinalPrice.GreaterThan(decimal.Zero),
		"Final price should be positive")

	t.Logf("Original Price: %s", result.OriginalPrice.String())
	t.Logf("Final Price: %s", result.FinalPrice.String())
	t.Logf("Total Discount: %s", totalDiscount.String())
	t.Logf("Discount Percentage: %s%%", result.GetDiscountPercentage().String())
	t.Logf("Applied Discounts: %v", result.AppliedDiscounts)
	t.Logf("Message: %s", result.Message)
}
