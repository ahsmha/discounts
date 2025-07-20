package strategies

import (
	"github.com/ahsmha/discounts/internal/models"
	"github.com/shopspring/decimal"
)

type VoucherDiscountStrategy struct{}

func (s *VoucherDiscountStrategy) IsApplicable(discount *models.Discount, cart []models.CartItem, customer models.CustomerProfile, payment *models.PaymentInfo) bool {

	if discount.Type != models.DiscountTypeVoucher || !discount.IsValid() || !discount.IsApplicableToCustomer(customer) {
		return false
	}

	total := calculateCartTotal(cart)
	if !discount.MinAmount.IsZero() && total.LessThan(discount.MinAmount) {
		return false
	}

	for _, item := range cart {
		if discount.MatchesProduct(item.Product) {
			return true
		}
	}

	return len(discount.ExcludedItems) == 0 && len(cart) > 0
}

func (s *VoucherDiscountStrategy) Calculate(discount *models.Discount, cart []models.CartItem, currentTotal decimal.Decimal) decimal.Decimal {
	return calculateDiscountValue(discount, currentTotal)
}
