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
