package strategies

import (
	"github.com/ahsmha/discounts/internal/models"
	"github.com/shopspring/decimal"
)

func calculateDiscountValue(discount *models.Discount, baseAmount decimal.Decimal) decimal.Decimal {
	if baseAmount.IsZero() {
		return decimal.Zero
	}

	var discountAmount decimal.Decimal
	if discount.IsPercentage {
		discountAmount = baseAmount.Mul(discount.Value).Div(decimal.NewFromInt(models.PercentageBase))
	} else {
		discountAmount = discount.Value
	}

	if !discount.MaxAmount.IsZero() && discountAmount.GreaterThan(discount.MaxAmount) {
		discountAmount = discount.MaxAmount
	}

	if discountAmount.GreaterThan(baseAmount) {
		return baseAmount
	}

	return discountAmount
}

func calculateCartTotal(cart []models.CartItem) decimal.Decimal {
	total := decimal.Zero
	for _, item := range cart {
		total = total.Add(item.GetTotalPrice())
	}
	return total
}

func isInList(item string, list []string) bool {
	for _, l := range list {
		if l == item {
			return true
		}
	}
	return false
}
