package payment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"ecom-api/internal/adapters/framework/left/services/auth"
	"ecom-api/internal/application/core/types/entity/payloads"
	"ecom-api/internal/ports/right/rports"
	"ecom-api/pkg/configs"
	"ecom-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
)

type PaymentHandler struct {
	paymentStore rports.PaymentStore
	userStore    rports.UserStore
	orderStore   rports.OrderStore
}

func NewPaymentHandler(paymentStore rports.PaymentStore, userStore rports.UserStore, orderStore rports.OrderStore) *PaymentHandler {
	return &PaymentHandler{paymentStore: paymentStore, userStore: userStore, orderStore: orderStore}
}

func (handler *PaymentHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create/customer", auth.WithJWTAuth(handler.handleCustomerCreation, handler.userStore, "admin", "user")).Methods(http.MethodPost)
	router.HandleFunc("/customer/{id}", auth.WithJWTAuth(handler.handleGetCustomerById, handler.userStore, "admin")).Methods(http.MethodGet)
	router.HandleFunc("/customers", auth.WithJWTAuth(handler.handleGetCustomers, handler.userStore, "admin")).Methods(http.MethodGet)
	router.HandleFunc("customer/delete/{customerId}", auth.WithJWTAuth(handler.handleCustomerDeletion, handler.userStore, "admin", "storeowner")).Methods(http.MethodPost)

	router.HandleFunc("/payment_method/{customerId}", auth.WithJWTAuth(handler.handlePaymentMethodCreation, handler.userStore, "admin", "user")).Methods(http.MethodPost)
	router.HandleFunc("/charges/{customerId}", auth.WithJWTAuth(handler.handleCustomeChargeProcess, handler.userStore, "admin", "user")).Methods(http.MethodPost)

	router.HandleFunc("/payment/webhook", handler.handlePaymentLiveUpdateThroughWebhook).Methods(http.MethodPost)
}

func (handler *PaymentHandler) handleCustomeChargeProcess(w http.ResponseWriter, r *http.Request) {
	var chargeParams payloads.CustomerChargeRequest
	vars := mux.Vars(r)

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	customerId, ok := vars["customerId"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing customer ID"))
		return
	}

	if err := utils.ParseJSON(r, &chargeParams); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(chargeParams); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors)
		return
	}

	charge, err := handler.paymentStore.CreateStripeCharge(&chargeParams, customerId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userID := auth.GetUserIDFromContext(r.Context())

	user, err := handler.userStore.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	HTMLTemplateEmailHandler(w, r, user.Email, map[string]string{
		"username": user.FirstName + " " + user.LastName,
		"email":    configs.Envs.FromEmail,
		"address":  "USA, New York, 123 wall street, Apt:10032, 143B",
	})

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{"charge": charge.Amount}, nil)

}

func (handler *PaymentHandler) handleCustomerCreation(w http.ResponseWriter, r *http.Request) {
	var customer payloads.CustomerPayload

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	if err := utils.ParseJSON(r, &customer); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(customer); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors)
		return
	}

	cus, err := handler.paymentStore.CreateStripeCustomer(&customer)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, cus, nil)
}

func (handler *PaymentHandler) handleGetCustomerById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)

	customerId, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing customer ID"))
		return
	}

	customer, err := handler.paymentStore.GetStripeCustomer(customerId)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, customer, nil)
}

func (handler *PaymentHandler) handleGetCustomers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	customerList, err := handler.paymentStore.GetAllStripeCustomers()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, customerList, nil)
}

func (handler *PaymentHandler) handlePaymentMethodCreation(w http.ResponseWriter, r *http.Request) {
	// var stripeCard payloads.StripeCardPayload
	vars := mux.Vars(r)

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	customerId, ok := vars["customerId"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing customer ID"))
		return
	}

	// if err := utils.ParseJSON(r, &stripeCard); err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err)
	// 	return
	// }

	// if err := utils.Validate.Struct(stripeCard); err != nil {
	// 	errors := err.(validator.ValidationErrors)
	// 	utils.WriteError(w, http.StatusBadRequest, errors)
	// 	return
	// }

	// if !luhn.Valid(stripeCard.Number) {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("provide valid card number"))
	// 	return
	// }

	// expMonth, err := strconv.Atoi(stripeCard.ExpMonth)
	// if err != nil || expMonth < 1 || expMonth > 12 {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("provide a valid expiration month (01-12)"))
	// 	return
	// }

	// expYear, err := strconv.Atoi(stripeCard.ExpYear)
	// if err != nil || expYear < time.Now().Year() {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("expiration year must be this year or later"))
	// 	return
	// }

	// currentYear, currentMonth, _ := time.Now().Date()
	// if expYear == currentYear && expMonth < int(currentMonth) {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("expiration date must be in the future"))
	// 	return
	// }

	paymentMethod, err := handler.paymentStore.CreatePaymentMethod(customerId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, paymentMethod, nil)

}

func (handler *PaymentHandler) handleCustomerDeletion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	customerId, ok := vars["customerId"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing customer ID"))
		return
	}

	isDeleted, err := handler.paymentStore.DeleteCustomer(customerId)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]bool{"isDeleted": isDeleted}, nil)
}

func (handler *PaymentHandler) handlePaymentLiveUpdateThroughWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error reading request body: %v", err))
		return
	}

	sigHeader := r.Header.Get("stripe-signature")
	if sigHeader == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error: signature is empty"))
		return
	}

	event, err := webhook.ConstructEvent(payload, sigHeader, configs.Envs.StripeWebhookSecret)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error verifying webhook signature: %v", err))
		return
	}

	switch event.Type {
	case "customer.created":
		var customer stripe.Customer
		if err := json.Unmarshal(event.Data.Raw, &customer); err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("webhook error: %v", err))
			return
		}
		log.Printf("New customer created: %s", customer.ID)

	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		if err := json.Unmarshal(event.Data.Raw, &paymentMethod); err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("webhook error: %v", err))
			return
		}
		log.Printf("Payment method attached: %s", paymentMethod.ID)

	case "payment_intent.created":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("webhook error: %v", err))
			return
		}
		log.Printf("PaymentIntent created: %s", paymentIntent.ID)

	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("webhook error: %v", err))
			return
		}
		log.Printf("PaymentIntent succeeded: %s", paymentIntent.ID)

	case "charge.succeeded":
		var charge stripe.Charge
		if err := json.Unmarshal(event.Data.Raw, &charge); err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("webhook error: %v", err))
			return
		}
		log.Printf("Charge succeeded for charge ID: %s, amount: %d", charge.ID, charge.Amount)

	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("webhook error: %v", err))
			return
		}
		log.Printf("PaymentIntent failed: %s, error: %v", paymentIntent.ID, paymentIntent.LastPaymentError)

	case "payment_method.detached":
		var paymentMethod stripe.PaymentMethod
		if err := json.Unmarshal(event.Data.Raw, &paymentMethod); err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("webhook error: %v", err))
			return
		}
		log.Printf("Payment method detached: %s", paymentMethod.ID)

	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("webhook error: %v", err))
			return
		}
		log.Printf("Checkout session completed: %s", session.ID)
		// You could fulfill the order, update database, and send a confirmation email.

		// -------- checkout session completion actions --------------
		/*
			orderID := session.Metadata["order_id"] // Extract order ID

			if orderID == "" {
				utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("order ID missing in metadata"))
				return
			}

			err := handler.orderStore.UpdateOrderPaymentStatus(orderID, configs.Envs.PaymentStatusPaid)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update order status: %v", err))
				return
			}

			userID, ok := r.Context().Value("userID").(string)
			if !ok {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user ID not found in context"))
				return
			}

			user, err := handler.userStore.GetUserByID(userID)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			HTMLTemplateEmailHandler(w, r, user.Email, map[string]string{
				"username": user.FirstName + user.LastName,
				"email":    user.Email,
				"address":  "USA, New York, 123 wall street, Apt:10032, 143B",
			})
		*/

	case "checkout.session.expired":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("webhook error: %v", err))
			return
		}
		log.Printf("Checkout session expired: %s", session.ID)
		// Handle expired sessions, such as notifying the customer.

	case "checkout.session.payment_failed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("webhook error: %v", err))
			return
		}
		log.Printf("Payment failed for session: %s", session.ID)
		// Handle failed payments, like sending an email to retry payment

	default:
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("unhandled event type: %s", event.Type))
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil, nil)
}
