package rports

import (
	"ecom-api/internal/application/core/types/entity/payloads"

	"github.com/stripe/stripe-go"
)

type PaymentStore interface {
	//customer
	CreateStripeCustomer(customerParams *payloads.CustomerPayload) (*stripe.Customer, error)
	GetStripeCustomer(customerId string) (*stripe.Customer, error)
	GetAllStripeCustomers() ([]stripe.Customer, error)

	//payment method
	CreatePaymentMethod(customerId string, card *payloads.StripeCardPayload) (*stripe.PaymentMethod, error)

	//charge method
	CreateStripeCharge(chargeParams *payloads.CustomerChargeRequest, customerId string) (*stripe.Charge, error)
}
