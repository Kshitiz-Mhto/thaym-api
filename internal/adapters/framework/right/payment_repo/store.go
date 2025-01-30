package paymentrepo

import (
	"ecom-api/internal/application/core/types/entity/payloads"
	"fmt"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/paymentmethod"
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

func (store *PaymentStore) CreatePaymentMethod(customerId string, card *payloads.StripeCardPayload) (*stripe.PaymentMethod, error) {
	var params *stripe.PaymentMethodParams

	if card.Type == "card" {
		params = &stripe.PaymentMethodParams{
			Type: stripe.String(card.Type),
			Card: &stripe.PaymentMethodCardParams{
				Number:   stripe.String(card.Number),
				ExpMonth: stripe.String(card.ExpMonth),
				ExpYear:  stripe.String(card.ExpYear),
				CVC:      stripe.String(card.CVC),
			},
		}
	} else {
		return nil, fmt.Errorf("unsupported payment method type %s", card.Type)
	}

	paymentMethod, err := paymentmethod.New(params)
	if err != nil {
		return nil, fmt.Errorf("unable to create payment method")
	}

	_, error := AttachCustomerPaymentMethod(customerId, paymentMethod.ID)

	if error != nil {
		return nil, fmt.Errorf("unable to attach payment method to customer %s", customerId)
	}

	return paymentMethod, nil
}

func (store *PaymentStore) CreateStripeCharge(chargeParams *payloads.CustomerChargeRequest, customerId string) (*stripe.Charge, error) {
	params := &stripe.ChargeParams{
		Amount:       stripe.Int64(chargeParams.Amount),
		Currency:     stripe.String(chargeParams.Currency),
		ReceiptEmail: stripe.String(chargeParams.ReceiptEmail),
		Description:  stripe.String(chargeParams.Description),
		Source:       &stripe.SourceParams{Token: stripe.String("tok_visa")},
		Customer:     stripe.String(customerId),
	}

	charge, err := charge.New(params)
	if err != nil {
		return nil, fmt.Errorf("unable to create charge")
	}

	return charge, nil
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

func AttachCustomerPaymentMethod(customerId string, paymentMethodId string) (*stripe.PaymentMethod, error) {

	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerId),
	}
	paymentmethod, err := paymentmethod.Attach(
		paymentMethodId,
		params,
	)

	if err != nil {
		return nil, err
	}

	fmt.Println(paymentmethod)

	return paymentmethod, nil
}
