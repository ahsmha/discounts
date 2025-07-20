package discount

import (
	"github.com/ahsmha/discounts/internal/models"
	"github.com/shopspring/decimal"
)

type DiscountStrategy interface {
	IsApplicable(discount *models.Discount, cart []models.CartItem, customer models.CustomerProfile, payment *models.PaymentInfo) bool
	Calculate(discount *models.Discount, cart []models.CartItem, currentTotal decimal.Decimal) decimal.Decimal
}
