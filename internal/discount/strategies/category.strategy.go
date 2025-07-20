package strategies

import (
	"github.com/ahsmha/discounts/internal/models"
	"github.com/shopspring/decimal"
)

type CategoryDiscountStrategy struct{}

func (s *CategoryDiscountStrategy) IsApplicable(discount *models.Discount, cart []models.CartItem, customer models.CustomerProfile, payment *models.PaymentInfo) bool {
	if discount.Type != models.DiscountTypeCategory || !discount.IsValid() || !discount.IsApplicableToCustomer(customer) {
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
	return false
}

func (s *CategoryDiscountStrategy) Calculate(discount *models.Discount, cart []models.CartItem, currentTotal decimal.Decimal) decimal.Decimal {
	var amount decimal.Decimal
	for _, item := range cart {
		if discount.MatchesProduct(item.Product) {
			amount = amount.Add(item.GetTotalPrice())
		}
	}

	return calculateDiscountValue(discount, amount)
}
