package discount

import (
	"github.com/ahsmha/discounts/internal/discount/strategies"
	"github.com/ahsmha/discounts/internal/models"
)

// Factory holds a mapping from discount type to strategy instance.
type StrategyFactory struct {
	strategies map[models.DiscountType]DiscountStrategy
}

func NewStrategyFactory() *StrategyFactory {
	return &StrategyFactory{
		strategies: map[models.DiscountType]DiscountStrategy{
			models.DiscountTypeBrand:    &strategies.BrandDiscountStrategy{},
			models.DiscountTypeCategory: &strategies.CategoryDiscountStrategy{},
			models.DiscountTypeVoucher:  &strategies.VoucherDiscountStrategy{},
			models.DiscountTypeBank:     &strategies.BankDiscountStrategy{},
		},
	}
}

func (sf *StrategyFactory) Get(discountType models.DiscountType) DiscountStrategy {
	return sf.strategies[discountType]
}
