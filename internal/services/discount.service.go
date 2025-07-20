package services

import (
	"context"
	"fmt"
	"sort"

	"github.com/ahsmha/discounts/internal/discount"
	"github.com/ahsmha/discounts/internal/interfaces"
	"github.com/ahsmha/discounts/internal/models"
	"github.com/ahsmha/discounts/pkg/errors"
	"github.com/shopspring/decimal"
)

type discountService struct {
	discountRepo    interfaces.IDiscountRepository
	strategyFactory *discount.StrategyFactory
}

func NewDiscountService(discountRepo interfaces.IDiscountRepository) interfaces.IDiscountService {
	return &discountService{
		discountRepo:    discountRepo,
		strategyFactory: discount.NewStrategyFactory(),
	}
}

func (ds *discountService) CalculateCartDiscounts(ctx context.Context, cartItems []models.CartItem,
	customer models.CustomerProfile, paymentInfo *models.PaymentInfo) (*models.DiscountedPrice, error) {

	if len(cartItems) == 0 {
		return nil, errors.NewValidationError("cart is empty")
	}

	originalPrice := decimal.Zero
	for _, item := range cartItems {
		originalPrice = originalPrice.Add(item.GetTotalPrice())
	}

	allDiscounts, err := ds.discountRepo.GetActiveDiscounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get discounts: %w", err)
	}

	result := &models.DiscountedPrice{
		OriginalPrice:    originalPrice,
		FinalPrice:       originalPrice,
		AppliedDiscounts: make(map[string]decimal.Decimal),
		Message:          "No discounts applied",
	}

	// Sort by priority
	sort.Slice(allDiscounts, func(i, j int) bool {
		return allDiscounts[i].Priority > allDiscounts[j].Priority
	})

	for _, discount := range allDiscounts {
		strategy := ds.strategyFactory.Get(discount.Type)
		if strategy == nil {
			continue
		}

		applicable := strategy.IsApplicable(&discount, cartItems, customer, paymentInfo)
		if !applicable {
			continue
		}

		amount := strategy.Calculate(&discount, cartItems, result.FinalPrice)
		if amount.GreaterThan(decimal.Zero) {
			result.FinalPrice = result.FinalPrice.Sub(amount)
			result.AppliedDiscounts[discount.Name] = amount

			// Track usage
			err := ds.discountRepo.IncrementUsageCount(ctx, discount.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to increment usage: %w", err)
			}
		}
	}

	if len(result.AppliedDiscounts) > 0 {
		result.Message = fmt.Sprintf("Applied %d discount(s) - Savings: %s",
			len(result.AppliedDiscounts), result.GetTotalDiscount().String())
	}

	return result, nil
}

func (ds *discountService) ValidateDiscountCode(ctx context.Context, code string, cartItems []models.CartItem,
	customer models.CustomerProfile) (bool, error) {

	if code == "" {
		return false, errors.NewValidationError("code cannot be empty")
	}

	discount, err := ds.discountRepo.GetDiscountByCode(ctx, code)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return false, nil
		}
		return false, fmt.Errorf("repo error: %w", err)
	}

	strat := ds.strategyFactory.Get(discount.Type)
	if strat == nil {
		return false, nil
	}

	return strat.IsApplicable(discount, cartItems, customer, nil), nil
}
