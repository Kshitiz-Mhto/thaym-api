package entity

type CartCheckoutItem struct {
	ProductID string  `json:"productID" validate:"required,uuid"`  // Unique identifier for the product (UUID)
	Quantity  int     `json:"quantity" validate:"required,gt=0"`   // Quantity of the product being purchased (must be greater than 0)
	Currency  string  `json:"currency" validate:"required,len=3"`  // ISO 4217 currency code (e.g., USD, EUR)
	Discount  float64 `json:"discount,omitempty" validate:"gte=0"` // Discount applied to this item, optional but must be non-negative
	Tax       float64 `json:"tax,omitempty" validate:"gte=0"`      // Tax applied to this item, optional but must be non-negative
}
