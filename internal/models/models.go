package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type BrandTier string

const (
	BrandTierPremium BrandTier = "premium"
	BrandTierRegular BrandTier = "regular"
	BrandTierBudget  BrandTier = "budget"
)

type Brand struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Tier BrandTier `json:"tier"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID           string          `json:"id"`
	Brand        Brand           `json:"brand"`
	Category     Category        `json:"category"`
	BasePrice    decimal.Decimal `json:"base_price"`
	CurrentPrice decimal.Decimal `json:"current_price"` // After brand/category discount
}

type CartItem struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
	Size     string  `json:"size"`
}

func (ci *CartItem) GetTotalPrice() decimal.Decimal {
	return ci.Product.CurrentPrice.Mul(decimal.NewFromInt(int64(ci.Quantity)))
}

type CardType string

const (
	Credit CardType = "CREDIT"
	Debit  CardType = "DEBIT"
)

type PaymentMethod string

const (
	UPI  PaymentMethod = "UPI"
	Card PaymentMethod = "CARD"
)

type PaymentInfo struct {
	Method   PaymentMethod `json:"method"`
	BankName *string       `json:"bank_name"`
	CardType *CardType     `json:"card_type"`
}

type DiscountedPrice struct {
	OriginalPrice    decimal.Decimal            `json:"original_price"`
	FinalPrice       decimal.Decimal            `json:"final_price"`
	AppliedDiscounts map[string]decimal.Decimal `json:"applied_discounts"` // discount_id -> amount
	Message          string                     `json:"message"`
}

func (dp *DiscountedPrice) GetTotalDiscount() decimal.Decimal {
	total := decimal.Zero
	for _, discount := range dp.AppliedDiscounts {
		total = total.Add(discount)
	}
	return total
}

func (dp *DiscountedPrice) GetDiscountPercentage() decimal.Decimal {
	if dp.OriginalPrice.IsZero() {
		return decimal.Zero
	}
	return dp.GetTotalDiscount().Div(dp.OriginalPrice).Mul(decimal.NewFromInt(100))
}

type CustomerProfile struct {
	ID   string `json:"id"`
	Tier string `json:"tier"`
}

type DiscountType string

const (
	DiscountTypeBrand    DiscountType = "brand"
	DiscountTypeCategory DiscountType = "category"
	DiscountTypeBank     DiscountType = "bank"
	DiscountTypeVoucher  DiscountType = "voucher"
)

type Discount struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Type          DiscountType    `json:"type"`
	Value         decimal.Decimal `json:"value"`          // Percentage or fixed amount
	IsPercentage  bool            `json:"is_percentage"`  // True for percentage, false for fixed amount
	MinAmount     decimal.Decimal `json:"min_amount"`     // Minimum order amount
	MaxAmount     decimal.Decimal `json:"max_amount"`     // Maximum discount amount
	ApplicableTo  []string        `json:"applicable_to"`  // Brand names, categories, bank names, etc.
	ExcludedItems []string        `json:"excluded_items"` // Excluded brand ids, category ids, etc.
	CustomerTiers []string        `json:"customer_tiers"` // Applicable customer tiers
	Code          string          `json:"code"`           // Voucher code (for voucher discounts)
	ValidFrom     time.Time       `json:"valid_from"`
	ValidTo       time.Time       `json:"valid_to"`
	IsActive      bool            `json:"is_active"`
	UsageLimit    int             `json:"usage_limit"` // Maximum number of uses
	UsedCount     int             `json:"used_count"`  // Current usage count
	Priority      int             `json:"priority"`    // Higher number = higher priority
}

func (d *Discount) IsValid() bool {
	now := time.Now()
	return d.IsActive &&
		now.After(d.ValidFrom) &&
		now.Before(d.ValidTo) &&
		(d.UsageLimit == 0 || d.UsedCount < d.UsageLimit)
}

func (d *Discount) IsApplicableToProduct(product Product) bool {
	for _, excludedBrandIDOrCategoryID := range d.ExcludedItems {
		if excludedBrandIDOrCategoryID == product.Brand.ID || excludedBrandIDOrCategoryID == product.Category.ID {
			return false
		}
	}

	switch d.Type {
	case DiscountTypeBrand:
		return d.isInList(product.Brand.ID, d.ApplicableTo)
	case DiscountTypeCategory:
		return d.isInList(product.Category.ID, d.ApplicableTo)
	case DiscountTypeVoucher:
		return true // Vouchers are generally applicable to all products unless excluded
	default:
		return true
	}
}

func (d *Discount) IsApplicableToCustomer(customer CustomerProfile) bool {
	if len(d.CustomerTiers) == 0 {
		return true // No tier restrictions
	}
	return d.isInList(customer.Tier, d.CustomerTiers)
}

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
