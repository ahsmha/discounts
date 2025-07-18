package interfaces

import (
	"context"

	"github.com/ahsmha/discounts/internal/models"
	"github.com/shopspring/decimal"
)

// IDiscountService interface defines the contract for discount service operations
type IDiscountService interface {
	// CalculateCartDiscounts calculates final price after applying discount logic:
	// - First apply brand/category discounts
	// - Then apply coupon codes
	// - Then apply bank offers
	CalculateCartDiscounts(ctx context.Context, cartItems []models.CartItem,
		customer models.CustomerProfile, paymentInfo *models.PaymentInfo) (*models.DiscountedPrice, error)

	// ValidateDiscountCode validates if a discount code can be applied.
	// Handle specific cases like:
	// - Brand exclusions
	// - Category restrictions
	// - Customer tier requirements
	ValidateDiscountCode(ctx context.Context, code string, cartItems []models.CartItem,
		customer models.CustomerProfile) (bool, error)
}

type IValidator interface {
	ValidateDiscount(discount *models.Discount, cartItems []models.CartItem, customer models.CustomerProfile, paymentInfo *models.PaymentInfo) bool
}

type ICalculator interface {
	CalculateDiscountAmount(discount *models.Discount, cartItems []models.CartItem, currentPrice decimal.Decimal) decimal.Decimal
}
