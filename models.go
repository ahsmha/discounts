package models

import (
	"time"
	"github.com/shopspring/decimal"
)

// BrandTier represents different brand tiers
type BrandTier string

const (
	BrandTierPremium BrandTier = "premium"
	BrandTierRegular BrandTier = "regular"
	BrandTierBudget  BrandTier = "budget"
)

// Product represents a product with its pricing information
type Product struct {
	ID           string          `json:"id"`
	Brand        string          `json:"brand"`
	BrandTier    BrandTier       `json:"brand_tier"`
	Category     string          `json:"category"`
	BasePrice    decimal.Decimal `json:"base_price"`
	CurrentPrice decimal.Decimal `json:"current_price"` // After brand/category discount
}

// CartItem represents an item in the shopping cart
type CartItem struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
	Size     string  `json:"size"`
}

// GetTotalPrice calculates the total price for the cart item
func (ci *CartItem) GetTotalPrice() decimal.Decimal {
	return ci.Product.CurrentPrice.Mul(decimal.NewFromInt(int64(ci.Quantity)))
}

// PaymentInfo represents payment information
type PaymentInfo struct {
	Method   string  `json:"method"`    // CARD, UPI, etc
	BankName *string `json:"bank_name"` // Optional
	CardType *string `json:"card_type"` // Optional: CREDIT, DEBIT
}

// DiscountedPrice represents the final pricing after discounts
type DiscountedPrice struct {
	OriginalPrice    decimal.Decimal            `json:"original_price"`
	FinalPrice       decimal.Decimal            `json:"final_price"`
	AppliedDiscounts map[string]decimal.Decimal `json:"applied_discounts"` // discount_name -> amount
	Message          string                     `json:"message"`
}

// GetTotalDiscount calculates the total discount amount
func (dp *DiscountedPrice) GetTotalDiscount() decimal.Decimal {
	total := decimal.Zero
	for _, discount := range dp.AppliedDiscounts {
		total = total.Add(discount)
	}
	return total
}

// GetDiscountPercentage calculates the overall discount percentage
func (dp *DiscountedPrice) GetDiscountPercentage() decimal.Decimal {
	if dp.OriginalPrice.IsZero() {
		return decimal.Zero
	}
	return dp.GetTotalDiscount().Div(dp.OriginalPrice).Mul(decimal.NewFromInt(100))
}

// CustomerProfile represents customer information
type CustomerProfile struct {
	ID   string `json:"id"`
	Tier string `json:"tier"`
	// Add more fields as needed
}

// DiscountType represents different types of discounts
type DiscountType string

const (
	DiscountTypeBrand    DiscountType = "brand"
	DiscountTypeCategory DiscountType = "category"
	DiscountTypeBank     DiscountType = "bank"
	DiscountTypeVoucher  DiscountType = "voucher"
)

// Discount represents a discount rule
type Discount struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Type           DiscountType    `json:"type"`
	Value          decimal.Decimal `json:"value"`          // Percentage or fixed amount
	IsPercentage   bool            `json:"is_percentage"`   // True for percentage, false for fixed amount
	MinAmount      decimal.Decimal `json:"min_amount"`      // Minimum order amount
	MaxAmount      decimal.Decimal `json:"max_amount"`      // Maximum discount amount
	ApplicableTo   []string        `json:"applicable_to"`   // Brand names, categories, bank names, etc.
	ExcludedItems  []string        `json:"excluded_items"`  // Excluded brands, categories, etc.
	CustomerTiers  []string        `json:"customer_tiers"`  // Applicable customer tiers
	Code           string          `json:"code"`            // Voucher code (for voucher discounts)
	ValidFrom      time.Time       `json:"valid_from"`
	ValidTo        time.Time       `json:"valid_to"`
	IsActive       bool            `json:"is_active"`
	UsageLimit     int             `json:"usage_limit"`     // Maximum number of uses
	UsedCount      int             `json:"used_count"`      // Current usage count
	Priority       int             `json:"priority"`        // Higher number = higher priority
}

// IsValid checks if the discount is currently valid
func (d *Discount) IsValid() bool {
	now := time.Now()
	return d.IsActive && 
		   now.After(d.ValidFrom) && 
		   now.Before(d.ValidTo) && 
		   (d.UsageLimit == 0 || d.UsedCount < d.UsageLimit)
}

// IsApplicableToProduct checks if the discount is applicable to a product
func (d *Discount) IsApplicableToProduct(product Product) bool {
	// Check exclusions first
	for _, excluded := range d.ExcludedItems {
		if excluded == product.Brand || excluded == product.Category {
			return false
		}
	}

	// Check applicability based on discount type
	switch d.Type {
	case DiscountTypeBrand:
		return d.isInList(product.Brand, d.ApplicableTo)
	case DiscountTypeCategory:
		return d.isInList(product.Category, d.ApplicableTo)
	case DiscountTypeVoucher:
		return true // Vouchers are generally applicable to all products unless excluded
	default:
		return true
	}
}

// IsApplicableToCustomer checks if the discount is applicable to a customer
func (d *Discount) IsApplicableToCustomer(customer CustomerProfile) bool {
	if len(d.CustomerTiers) == 0 {
		return true // No tier restrictions
	}
	return d.isInList(customer.Tier, d.CustomerTiers)
}

// isInList checks if a string is in a slice of strings
func (d *Discount) isInList(item string, list []string) bool {
	if len(list) == 0 {
		return true // No restrictions
	}
	for _, listItem := range list {
		if listItem == item {
			return true
		}
	}
	return false
}

// CalculateDiscount calculates the discount amount for a given price
func (d *Discount) CalculateDiscount(price decimal.Decimal) decimal.Decimal {
	if d.IsPercentage {
		discount := price.Mul(d.Value).Div(decimal.NewFromInt(100))
		if !d.MaxAmount.IsZero() && discount.GreaterThan(d.MaxAmount) {
			return d.MaxAmount
		}
		return discount
	}
	return d.Value
}