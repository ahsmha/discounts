package service

import (
	"context"
	"fmt"
	"sort"
	"github.com/unifize/discount-service/internal/models"
	"github.com/unifize/discount-service/internal/repository"
	"github.com/unifize/discount-service/pkg/errors"
	"github.com/shopspring/decimal"
)

// DiscountService interface defines the contract for discount service operations
type DiscountService interface {
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

// discountService implements the DiscountService interface
type discountService struct {
	discountRepo repository.DiscountRepository
	calculator   DiscountCalculator
	validator    DiscountValidator
}

// NewDiscountService creates a new instance of DiscountService
func NewDiscountService(discountRepo repository.DiscountRepository) DiscountService {
	return &discountService{
		discountRepo: discountRepo,
		calculator:   NewDiscountCalculator(),
		validator:    NewDiscountValidator(),
	}
}

// CalculateCartDiscounts implements the main discount calculation logic
func (ds *discountService) CalculateCartDiscounts(ctx context.Context, cartItems []models.CartItem,
	customer models.CustomerProfile, paymentInfo *models.PaymentInfo) (*models.DiscountedPrice, error) {
	
	if len(cartItems) == 0 {
		return nil, errors.NewValidationError("cart is empty")
	}

	// Calculate original price
	originalPrice := decimal.Zero
	for _, item := range cartItems {
		originalPrice = originalPrice.Add(item.GetTotalPrice())
	}

	// Get all applicable discounts
	allDiscounts, err := ds.discountRepo.GetActiveDiscounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active discounts: %w", err)
	}

	// Initialize result
	result := &models.DiscountedPrice{
		OriginalPrice:    originalPrice,
		FinalPrice:       originalPrice,
		AppliedDiscounts: make(map[string]decimal.Decimal),
		Message:          "No discounts applied",
	}

	// Step 1: Apply brand/category discounts
	brandCategoryDiscounts := ds.filterDiscountsByType(allDiscounts, []models.DiscountType{
		models.DiscountTypeBrand,
		models.DiscountTypeCategory,
	})
	
	applicableBrandCategoryDiscounts := ds.filterApplicableDiscounts(brandCategoryDiscounts, cartItems, customer, paymentInfo)
	if len(applicableBrandCategoryDiscounts) > 0 {
		err := ds.applyDiscounts(ctx, applicableBrandCategoryDiscounts, cartItems, result)
		if err != nil {
			return nil, fmt.Errorf("failed to apply brand/category discounts: %w", err)
		}
	}

	// Step 2: Apply voucher codes (if any)
	voucherDiscounts := ds.filterDiscountsByType(allDiscounts, []models.DiscountType{models.DiscountTypeVoucher})
	applicableVoucherDiscounts := ds.filterApplicableDiscounts(voucherDiscounts, cartItems, customer, paymentInfo)
	if len(applicableVoucherDiscounts) > 0 {
		err := ds.applyDiscounts(ctx, applicableVoucherDiscounts, cartItems, result)
		if err != nil {
			return nil, fmt.Errorf("failed to apply voucher discounts: %w", err)
		}
	}

	// Step 3: Apply bank offers
	bankDiscounts := ds.filterDiscountsByType(allDiscounts, []models.DiscountType{models.DiscountTypeBank})
	applicableBankDiscounts := ds.filterApplicableDiscounts(bankDiscounts, cartItems, customer, paymentInfo)
	if len(applicableBankDiscounts) > 0 {
		err := ds.applyDiscounts(ctx, applicableBankDiscounts, cartItems, result)
		if err != nil {
			return nil, fmt.Errorf("failed to apply bank discounts: %w", err)
		}
	}

	// Update final message
	if len(result.AppliedDiscounts) > 0 {
		result.Message = fmt.Sprintf("Applied %d discount(s) - Total savings: %s", 
			len(result.AppliedDiscounts), result.GetTotalDiscount().String())
	}

	return result, nil
}

// ValidateDiscountCode validates if a discount code can be applied
func (ds *discountService) ValidateDiscountCode(ctx context.Context, code string, cartItems []models.CartItem,
	customer models.CustomerProfile) (bool, error) {
	
	if code == "" {
		return false, errors.NewValidationError("discount code cannot be empty")
	}

	discount, err := ds.discountRepo.GetDiscountByCode(ctx, code)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get discount by code: %w", err)
	}

	return ds.validator.ValidateDiscount(discount, cartItems, customer, nil), nil
}

// filterDiscountsByType filters discounts by their types
func (ds *discountService) filterDiscountsByType(discounts []models.Discount, types []models.DiscountType) []models.Discount {
	var filtered []models.Discount
	typeMap := make(map[models.DiscountType]bool)
	for _, t := range types {
		typeMap[t] = true
	}

	for _, discount := range discounts {
		if typeMap[discount.Type] {
			filtered = append(filtered, discount)
		}
	}
	return filtered
}

// filterApplicableDiscounts filters discounts that are applicable to the given cart and customer
func (ds *discountService) filterApplicableDiscounts(discounts []models.Discount, cartItems []models.CartItem, 
	customer models.CustomerProfile, paymentInfo *models.PaymentInfo) []models.Discount {
	
	var applicable []models.Discount
	for _, discount := range discounts {
		if ds.validator.ValidateDiscount(&discount, cartItems, customer, paymentInfo) {
			applicable = append(applicable, discount)
		}
	}

	// Sort by priority (highest first)
	sort.Slice(applicable, func(i, j int) bool {
		return applicable[i].Priority > applicable[j].Priority
	})

	return applicable
}

// applyDiscounts applies a list of discounts to the cart and updates the result
func (ds *discountService) applyDiscounts(ctx context.Context, discounts []models.Discount, 
	cartItems []models.CartItem, result *models.DiscountedPrice) error {
	
	for _, discount := range discounts {
		discountAmount := ds.calculator.CalculateDiscountAmount(&discount, cartItems, result.FinalPrice)
		
		if discountAmount.GreaterThan(decimal.Zero) {
			result.FinalPrice = result.FinalPrice.Sub(discountAmount)
			result.AppliedDiscounts[discount.Name] = discountAmount
			
			// Update usage count
			err := ds.discountRepo.IncrementUsageCount(ctx, discount.ID)
			if err != nil {
				return fmt.Errorf("failed to increment usage count for discount %s: %w", discount.ID, err)
			}
		}
	}
	
	return nil
}