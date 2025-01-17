package entity

import (
	"time"
)

type Order struct {
	ID            string    `json:"id"`                        // Unique identifier for the order
	UserID        int       `json:"userID"`                    // Foreign key to associate with the user
	Total         float64   `json:"total"`                     // Total amount for the order
	Subtotal      float64   `json:"subtotal"`                  // Subtotal before tax and discounts
	Status        string    `json:"status"`                    // Order status (e.g., "Pending", "Shipped", "Delivered", "Cancelled")
	PaymentStatus string    `json:"paymentStatus"`             // Payment status (e.g., "Paid", "Pending", "Refunded")
	PaymentMethod string    `json:"paymentMethod"`             // Payment method used (e.g., "Credit Card", "PayPal")
	Address       string    `json:"address"`                   // Shipping address
	Currency      string    `json:"currency" validate:"len=3"` // ISO 4217 currency code
	CreatedAt     time.Time `json:"createdAt"`                 // Timestamp for when the order was created
	UpdatedAt     time.Time `json:"updatedAt"`                 // Timestamp for when the order was last updated
}
