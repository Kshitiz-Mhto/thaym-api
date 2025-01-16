package entity

import (
	"time"
)

type OrderItem struct {
	ID          int       `json:"id"`                        // Unique identifier for the order item
	OrderID     int       `json:"orderId"`                   // Foreign key to associate with the order
	ProductID   int       `json:"productId"`                 // Foreign key to associate with the product
	ProductName string    `json:"productName"`               // Cached product name to prevent dependency on product table
	Quantity    int       `json:"quantity" validate:"gte=1"` // Quantity ordered, minimum of 1
	Price       float64   `json:"price" validate:"gt=0"`     // Price per unit
	TotalPrice  float64   `json:"totalPrice"`                // Calculated total price (Quantity * Price)
	Currency    string    `json:"currency" validate:"len=3"` // ISO 4217 currency code
	Discount    float64   `json:"discount"`                  // Discount applied to this item
	Tax         float64   `json:"tax"`                       // Tax applied to this item
	CreatedAt   time.Time `json:"createdAt"`                 // Timestamp for when the item was created
	UpdatedAt   time.Time `json:"updatedAt"`                 // Timestamp for when the item was last updated
}
