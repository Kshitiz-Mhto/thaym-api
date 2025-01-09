package payloads

import (
	"ecom-api/internal/application/core/types/entity"
)

type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
}

type RegisterUserPayload struct {
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=5,max=130"`
	PasswordRe string `json:"passwordRe" validate:"required,min=5,max=130,eqfield=Password"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CartCheckoutPayload struct {
	Items []entity.CartCheckoutItem `json:"items" validate:"required"`
}

type EmailWithTemplateRequestBody struct {
	ToAddr   string            `json:"to_addr"`
	Template string            `json:"template"`
	Vars     map[string]string `json:"vars"`
}
