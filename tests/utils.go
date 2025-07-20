package tests

import "github.com/ahsmha/discounts/internal/models"

func ResetCartPricesToBase(cart []models.CartItem) {
	for i := range cart {
		cart[i].Product.CurrentPrice = cart[i].Product.BasePrice
	}
}
