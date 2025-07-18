package service

import (
	"github.com/ahsmha/discount-service/internal/models"
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
		discountAmount = totalEligibleAmount.Mul(discount.Value).Div(decimal.NewFromInt(100))
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

// DiscountValidator interface defines methods for validating discount applicability
type DiscountValidator interface {
	ValidateDiscount(discount *models.Discount, cartItems []models.CartItem, customer models.CustomerProfile, paymentInfo *models.PaymentInfo) bool
}

// discountValidator implements the DiscountValidator interface
type discountValidator struct{}

// NewDiscountValidator creates a new instance of DiscountValidator
func NewDiscountValidator() DiscountValidator {
	return &discountValidator{}
}

// ValidateDiscount validates if a discount can be applied to the given cart and customer
func (dv *discountValidator) ValidateDiscount(discount *models.Discount, cartItems []models.CartItem,
	customer models.CustomerProfile, paymentInfo *models.PaymentInfo) bool {

	// Check if discount is valid (time, usage limits, etc.)
	if !discount.IsValid() {
		return false
	}

	// Check customer tier eligibility
	if !discount.IsApplicableToCustomer(customer) {
		return false
	}

	// Check minimum amount requirement
	if !discount.MinAmount.IsZero() {
		cartTotal := dv.calculateCartTotal(cartItems)
		if cartTotal.LessThan(discount.MinAmount) {
			return false
		}
	}

	// Validate based on discount type
	switch discount.Type {
	case models.DiscountTypeBrand, models.DiscountTypeCategory:
		return dv.validateProductDiscounts(discount, cartItems)
	case models.DiscountTypeBank:
		return dv.validateBankDiscount(discount, paymentInfo)
	case models.DiscountTypeVoucher:
		return dv.validateVoucherDiscount(discount, cartItems)
	default:
		return false
	}
}

// validateProductDiscounts validates brand/category discounts
func (dv *discountValidator) validateProductDiscounts(discount *models.Discount, cartItems []models.CartItem) bool {
	// Check if at least one product is eligible
	for _, item := range cartItems {
		if discount.IsApplicableToProduct(item.Product) {
			return true
		}
	}
	return false
}

// validateBankDiscount validates bank-specific discounts
func (dv *discountValidator) validateBankDiscount(discount *models.Discount, paymentInfo *models.PaymentInfo) bool {
	if paymentInfo == nil {
		return false
	}

	// Check if payment method is CARD
	if paymentInfo.Method != "CARD" {
		return false
	}

	// Check if bank name matches
	if len(discount.ApplicableTo) > 0 {
		if paymentInfo.BankName == nil {
			return false
		}
		return dv.isInList(*paymentInfo.BankName, discount.ApplicableTo)
	}

	return true
}

// validateVoucherDiscount validates voucher discounts
func (dv *discountValidator) validateVoucherDiscount(discount *models.Discount, cartItems []models.CartItem) bool {
	// Check if at least one product is not excluded
	for _, item := range cartItems {
		if discount.IsApplicableToProduct(item.Product) {
			return true
		}
	}
	return len(cartItems) > 0 && len(discount.ExcludedItems) == 0
}

// calculateCartTotal calculates the total cart value
func (dv *discountValidator) calculateCartTotal(cartItems []models.CartItem) decimal.Decimal {
	total := decimal.Zero
	for _, item := range cartItems {
		total = total.Add(item.GetTotalPrice())
	}
	return total
}

// isInList checks if a string is in a slice of strings
func (dv *discountValidator) isInList(item string, list []string) bool {
	for _, listItem := range list {
		if listItem == item {
			return true
		}
	}
	return false
}
