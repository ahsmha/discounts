package services

import (
	"github.com/ahsmha/discounts/internal/models"
	"github.com/shopspring/decimal"
)

// DiscountCalculator interface defines methods for calculating discount amounts
type DiscountCalculator interface {
	CalculateDiscountAmount(discount *models.Discount, cartItems []models.CartItem, currentPrice decimal.Decimal) decimal.Decimal
}

// discountCalculator implements the DiscountCalculator interface
type discountCalculator struct{}

// NewDiscountCalculator creates a new instance of DiscountCalculator
func NewDiscountCalculator() DiscountCalculator {
	return &discountCalculator{}
}

// CalculateDiscountAmount calculates the discount amount based on the discount rules
func (dc *discountCalculator) CalculateDiscountAmount(discount *models.Discount, cartItems []models.CartItem, currentPrice decimal.Decimal) decimal.Decimal {
	// Check minimum amount requirement
	if !discount.MinAmount.IsZero() && currentPrice.LessThan(discount.MinAmount) {
		return decimal.Zero
	}

	var totalEligibleAmount decimal.Decimal

	switch discount.Type {
	case models.DiscountTypeBrand, models.DiscountTypeCategory:
		// Calculate eligible amount based on applicable products
		totalEligibleAmount = dc.calculateEligibleAmount(discount, cartItems)
	case models.DiscountTypeVoucher, models.DiscountTypeBank:
		// Apply to entire cart amount
		totalEligibleAmount = currentPrice
	default:
		totalEligibleAmount = currentPrice
	}

	if totalEligibleAmount.IsZero() {
		return decimal.Zero
	}

	// Calculate discount amount
	var discountAmount decimal.Decimal
	if discount.IsPercentage {
		discountAmount = totalEligibleAmount.Mul(discount.Value).Div(decimal.NewFromInt(models.PercentageBase))
	} else {
		discountAmount = discount.Value
	}

	// Apply maximum discount limit
	if !discount.MaxAmount.IsZero() && discountAmount.GreaterThan(discount.MaxAmount) {
		discountAmount = discount.MaxAmount
	}

	// Ensure discount doesn't exceed the eligible amount
	if discountAmount.GreaterThan(totalEligibleAmount) {
		discountAmount = totalEligibleAmount
	}

	return discountAmount
}

// calculateEligibleAmount calculates the total amount eligible for brand/category discounts
func (dc *discountCalculator) calculateEligibleAmount(discount *models.Discount, cartItems []models.CartItem) decimal.Decimal {
	eligibleAmount := decimal.Zero

	for _, item := range cartItems {
		if discount.IsApplicableToProduct(item.Product) {
			eligibleAmount = eligibleAmount.Add(item.GetTotalPrice())
		}
	}

	return eligibleAmount
}
