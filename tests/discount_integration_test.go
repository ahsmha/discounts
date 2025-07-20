package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	repository "github.com/ahsmha/discounts/internal/repositories"
	"github.com/ahsmha/discounts/internal/services"
	"github.com/ahsmha/discounts/testdata"
)

// Integration test that fixes the import issue and demonstrates the working system
func TestDiscountService_Integration_WorkingDemo(t *testing.T) {
	// Setup repository and service
	repo := repository.NewInMemoryDiscountRepository()

	// Since we can't access the private SeedDiscounts method in this test,
	// we'll manually create and add discounts through the public interface
	ctx := context.Background()

	// Create sample discounts manually
	sampleDiscounts := testdata.GetSampleDiscounts()
	for _, discount := range sampleDiscounts {
		err := repo.CreateDiscount(ctx, &discount)
		require.NoError(t, err)
	}

	service := services.NewDiscountService(repo)

	// Test the multiple discount scenario
	cartItems, customer, paymentInfo := testdata.GetMultipleDiscountScenario()

	// Ensure we start with the base price for proper calculation
	ResetCartPricesToBase(cartItems)

	t.Logf("Testing with cart items: %+v", cartItems)
	t.Logf("Customer: %+v", customer)
	t.Logf("Payment info: %+v", paymentInfo)

	// Calculate discounts
	result, err := service.CalculateCartDiscounts(ctx, cartItems, customer, paymentInfo)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify results
	assert.True(t, result.FinalPrice.LessThan(result.OriginalPrice),
		"Final price should be less than original price")
	assert.True(t, len(result.AppliedDiscounts) > 0,
		"At least one discount should be applied")

	// Log the results for verification
	t.Logf("Original Price: %s", result.OriginalPrice.String())
	t.Logf("Final Price: %s", result.FinalPrice.String())
	t.Logf("Applied Discounts: %+v", result.AppliedDiscounts)
	t.Logf("Message: %s", result.Message)

	// Test discount code validation
	validationTests := []struct {
		code     string
		expected bool
	}{
		{"SUPER69", true},   // Valid for premium customers
		{"PREMIUM15", true}, // Valid for premium customers
		{"INVALID", false},  // Invalid code
	}

	for _, vt := range validationTests {
		t.Run("validate_"+vt.code, func(t *testing.T) {
			isValid, err := service.ValidateDiscountCode(ctx, vt.code, cartItems, customer)
			require.NoError(t, err)
			assert.Equal(t, vt.expected, isValid)
		})
	}
}
