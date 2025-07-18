package interfaces

import (
	"context"

	"github.com/ahsmha/discounts/internal/models"
)

// IDiscountRepository interface defines methods for discount data operations
type IDiscountRepository interface {
	// GetActiveDiscounts retrieves all active discounts
	GetActiveDiscounts(ctx context.Context) ([]models.Discount, error)

	// GetDiscountByCode retrieves a discount by its code
	GetDiscountByCode(ctx context.Context, code string) (*models.Discount, error)

	// GetDiscountByID retrieves a discount by its ID
	GetDiscountByID(ctx context.Context, id string) (*models.Discount, error)

	// CreateDiscount creates a new discount
	CreateDiscount(ctx context.Context, discount *models.Discount) error

	// UpdateDiscount updates an existing discount
	UpdateDiscount(ctx context.Context, discount *models.Discount) error

	// DeleteDiscount deletes a discount by ID
	DeleteDiscount(ctx context.Context, id string) error

	// IncrementUsageCount increments the usage count for a discount
	IncrementUsageCount(ctx context.Context, id string) error
}
