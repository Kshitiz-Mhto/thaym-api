package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ProductId   string          `json:"productId"`   // Unique identifier for the product
	Name        string          `json:"name"`        // Product name
	Description string          `json:"description"` // Product description
	Image       string          `json:"image"`       // URL or path to the product's image
	Price       decimal.Decimal `json:"price"`       // Price of the product
	Currency    string          `json:"currency"`    // Currency code (e.g., USD, EUR)
	Quantity    int             `json:"quantity"`    // Inventory count
	Category    string          `json:"category"`    // Product category
	Tags        []string        `json:"tags"`        // List of tags for filtering or searching
	IsActive    bool            `json:"isActive"`    // Indicates if the product is active/available for purchase
	CreatedAt   time.Time       `json:"createdAt"`   // Timestamp for when the product was created
	UpdatedAt   time.Time       `json:"updatedAt"`   // Timestamp for when the product was last updated
}
