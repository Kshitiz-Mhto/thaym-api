package paymentrepo

import (
	"ecom-api/internal/application/core/types/entity/payloads"
	"fmt"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/paymentintent"
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

func (store *PaymentStore) CreatePaymentMethod(customerId string) (*stripe.PaymentMethod, error) {

	//test-mode
	params := &stripe.PaymentMethodParams{
		Type: stripe.String("card"),
		Card: &stripe.PaymentMethodCardParams{
			Token: stripe.String("tok_amex"), // Use a Stripe test token
		},
	}

	paymentMethod, err := paymentmethod.New(params)
	if err != nil {
		return nil, fmt.Errorf("unable to create payment method: %v", err)
	}

	_, err = AttachCustomerPaymentMethod(customerId, paymentMethod.ID)

	if err != nil {
		return nil, fmt.Errorf("unable to attach payment method to customer %s", customerId)
	}

	return paymentMethod, nil
}

func (store *PaymentStore) CreateStripeCharge(chargeParams *payloads.CustomerChargeRequest, customerId string) (*stripe.PaymentIntent, error) {
	testPaymentMethodID := "pm_card_amex" // Stripe's test Visa payment method

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(chargeParams.Amount),
		Currency:           stripe.String(chargeParams.Currency),
		ReceiptEmail:       stripe.String(chargeParams.ReceiptEmail),
		Description:        stripe.String(chargeParams.Description),
		Customer:           stripe.String(customerId),
		PaymentMethod:      stripe.String(testPaymentMethodID),
		Confirm:            stripe.Bool(true), // Automatically confirm in test mode
		ConfirmationMethod: stripe.String(string(stripe.PaymentIntentConfirmationMethodAutomatic)),
	}

	// Add test-specific metadata (optional)
	params.AddMetadata("environment", "test")
	params.AddMetadata("test_case", "successful_payment")

	charge, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("payment faild:%s", err)
	}

	return charge, nil
}

func (store *PaymentStore) CreateCheckoutSession(orderID string, totalPrice float64, email string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Currency: stripe.String("usd"),
				Amount:   stripe.Int64(int64(totalPrice * 100)),
				Quantity: stripe.Int64(1),
			},
		},
		CustomerEmail: stripe.String(email),
		Params: stripe.Params{
			Metadata: map[string]string{
				"order_id": orderID,
			},
		},
		SuccessURL: stripe.String("https://localhost:8080/order/" + orderID),
	}

	session, err := session.New(params)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (store *PaymentStore) DeleteCustomer(id string) (bool, error) {
	params := &stripe.CustomerParams{}
	_, err := customer.Del(id, params)

	if err != nil {
		return false, fmt.Errorf("uable to delete customer: %v", err)
	}

	return true, nil
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
	paymentAttach, err := paymentmethod.Attach(
		paymentMethodId,
		params,
	)

	if err != nil {
		return nil, fmt.Errorf("uable to attach: %v", err)
	}

	return paymentAttach, nil
}

func InProduction() {
	// // Use Stripe.js to create a token
	// const { token, error } = await stripe.createToken('card', {
	// 	number: '4242424242424242',
	// 	exp_month: '12',
	// 	exp_year: '2025',
	// 	cvc: '123',
	// });

	// Use stripe.createToken to convert information collected by card elements into a single-use Token that you safely pass to your server to use in an API call.

	// // Send the token (e.g., `tok_visa`) to your backend

	// params := &stripe.PaymentMethodParams{
	//     Type: stripe.String("card"),
	//     Card: &stripe.PaymentMethodCardParams{
	//         Token: stripe.String(token), // Use the client-side token
	//     },
	// }

	// paymentMethod, err := paymentmethod.New(params)
	// ... rest of the code ...
}
