package paymentrepo

import (
	"ecom-api/internal/application/core/types/entity/payloads"
	"fmt"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
)

type PaymentStore struct{}

func NewPaymentStore() *PaymentStore {
	return &PaymentStore{}
}

func (store *PaymentStore) CreateStripeCustomer(customerParams *payloads.CustomerPayload) (*stripe.Customer, error) {

	address := mapAddressToStripe(customerParams.Address)
	shipping := mapShippingToStripe(customerParams.Shipping, address)

	params := &stripe.CustomerParams{
		Name:        stripe.String(customerParams.Name),
		Email:       stripe.String(customerParams.Email),
		Phone:       stripe.String(customerParams.Phone),
		Description: stripe.String(customerParams.Description),
		Balance:     stripe.Int64(customerParams.Balance),
		Address:     address,
		Shipping:    shipping,
	}

	newCustomer, err := customer.New(params)

	if err != nil {
		return nil, err
	}

	return newCustomer, err
}

func (store *PaymentStore) GetStripeCustomer(customerId string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{}

	customer, err := customer.Get(customerId, params)

	if err != nil {
		return nil, fmt.Errorf("customer detail retrival failed :%v", err)
	}

	return customer, nil
}

func (store *PaymentStore) GetAllStripeCustomers() ([]stripe.Customer, error) {
	params := &stripe.CustomerListParams{}
	params.Limit = stripe.Int64(100)

	var customers []stripe.Customer

	iter := customer.List(params)

	for iter.Next() {
		customer := iter.Customer()
		customers = append(customers, *customer)
	}

	if iter.Err() != nil {
		return nil, fmt.Errorf("failed to list customers: %v", iter.Err())
	}

	return customers, nil
}

func mapAddressToStripe(address payloads.Address) *stripe.AddressParams {

	return &stripe.AddressParams{
		Line1:      stripe.String(address.Line1),
		Line2:      stripe.String(address.Line2),
		City:       stripe.String(address.City),
		State:      stripe.String(address.State),
		PostalCode: stripe.String(address.PostalCode),
		Country:    stripe.String(address.Country),
	}
}

func mapShippingToStripe(shipping payloads.Shipping, address *stripe.AddressParams) *stripe.CustomerShippingDetailsParams {
	return &stripe.CustomerShippingDetailsParams{
		Name:    stripe.String(shipping.Name),
		Phone:   stripe.String(shipping.Phone),
		Address: address,
	}
}
