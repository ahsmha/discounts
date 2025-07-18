package services

import (
	"github.com/ahsmha/discounts/internal/models"
	"github.com/shopspring/decimal"
)

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
