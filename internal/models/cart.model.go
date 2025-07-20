package models

import "github.com/shopspring/decimal"

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

// CardType represents the type of card payment.
type CardType string

const (
	// Credit indicates payment by credit card.
	Credit CardType = "CREDIT"
	// Debit indicates payment by debit card.
	Debit CardType = "DEBIT"
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
