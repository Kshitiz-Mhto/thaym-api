package payloads

import (
	"ecom-api/internal/application/core/types/entity"
)

type CreateProductPayload struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Image       string   `json:"image" validate:"required,url"`
	Price       float64  `json:"price" validate:"required,gt=0"`
	Currency    string   `json:"currency" validate:"required,len=3"` // ISO 4217 currency code
	Quantity    int      `json:"quantity" validate:"required,gte=0"`
	Category    string   `json:"category" validate:"required"`
	Tags        []string `json:"tags"` // JSON array of tags
	IsActive    bool     `json:"isActive" validate:"required"`
}

type RegisterUserPayload struct {
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=5,max=130"`
	PasswordRe string `json:"passwordRe" validate:"required,min=5,max=130,eqfield=Password"`
}

type RegisterUserConfirmationPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=5,max=130"`
	Role      string `json:"role"`
	Token     string `json:"token" validate:"required"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CartCheckoutPayload struct {
	Items []entity.CartCheckoutItem `json:"items" validate:"required,dive,required"` // List of items in the cart, each item is required
}
type CustomerPayload struct {
	Email       string            `json:"email" validate:"required,email"`
	Name        string            `json:"name" validate:"required"`
	Balance     int64             `json:"balance" validate:"gte=0"`
	Phone       string            `json:"phone" validate:"required,e164"`
	Description string            `json:"description,omitempty" validat:"omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty" validat:"omitempty"`
	Address     Address           `json:"address" validate:"required"`
	Shipping    Shipping          `json:"shipping,omitempty" validate:"omitempty"`
}

type Address struct {
	City       string `json:"city,omitempty" validate:"omitempty"`
	Country    string `json:"country,omitempty" validate:"omitempty"`
	Line1      string `json:"line1,omitempty" validate:"omitempty"`
	Line2      string `json:"line2,omitempty" validate:"omitempty"`
	PostalCode string `json:"postal_code,omitempty" validate:"omitempty,numeric"`
	State      string `json:"state,omitempty" validate:"omitempty"`
}

type Shipping struct {
	Name    string  `json:"name" validate:"required"`
	Phone   string  `json:"phone" validate:"required,e164"`
	Address Address `json:"address" validate:"required"`
}
