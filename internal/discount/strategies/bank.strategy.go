package strategies

import (
	"github.com/ahsmha/discounts/internal/models"
	"github.com/shopspring/decimal"
)

type BankDiscountStrategy struct{}

func (s *BankDiscountStrategy) IsApplicable(discount *models.Discount, cart []models.CartItem, customer models.CustomerProfile, payment *models.PaymentInfo) bool {
	if discount.Type != models.DiscountTypeBank || !discount.IsValid() || !discount.IsApplicableToCustomer(customer) {
		return false
	}

	if payment == nil || payment.Method != "CARD" {
		return false
	}

	if len(discount.ApplicableTo) > 0 && (payment.BankName == nil || !isInList(*payment.BankName, discount.ApplicableTo)) {
		return false
	}

	cartTotal := calculateCartTotal(cart)
	return !discount.MinAmount.IsZero() || cartTotal.GreaterThanOrEqual(discount.MinAmount)
}

func (s *BankDiscountStrategy) Calculate(discount *models.Discount, cart []models.CartItem, currentTotal decimal.Decimal) decimal.Decimal {
	eligibleAmount := currentTotal
	return calculateDiscountValue(discount, eligibleAmount)
}
