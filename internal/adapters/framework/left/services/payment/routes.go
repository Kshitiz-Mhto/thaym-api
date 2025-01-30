package payment

import (
	"fmt"
	"net/http"

	"ecom-api/internal/adapters/framework/left/services/auth"
	"ecom-api/internal/application/core/types/entity/payloads"
	"ecom-api/internal/ports/right/rports"
	"ecom-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type PaymentHandler struct {
	paymentStore rports.PaymentStore
	userStore    rports.UserStore
}

func NewPaymentHandler(paymentStore rports.PaymentStore, userStore rports.UserStore) *PaymentHandler {
	return &PaymentHandler{paymentStore: paymentStore, userStore: userStore}
}

func (handler *PaymentHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create/customer", auth.WithJWTAuth(handler.handleCustomerCreation, handler.userStore, "admin", "user")).Methods(http.MethodPost)
	router.HandleFunc("/customer/{id}", auth.WithJWTAuth(handler.handleGetCustomerById, handler.userStore, "admin")).Methods(http.MethodGet)
	router.HandleFunc("/customers", auth.WithJWTAuth(handler.handleGetCustomers, handler.userStore, "admin")).Methods(http.MethodGet)
	router.HandleFunc("customer/delete/{customerId}", auth.WithJWTAuth(handler.handleCustomerDeletion, handler.userStore, "admin", "storeowner")).Methods(http.MethodPost)

	router.HandleFunc("/payment_method/{customerId}", auth.WithJWTAuth(handler.handlePaymentMethodCreation, handler.userStore, "admin", "user")).Methods(http.MethodPost)
	router.HandleFunc("/charges/{customerId}", auth.WithJWTAuth(handler.handleCustomeChargeProcess, handler.userStore, "admin", "user")).Methods(http.MethodPost)
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

	utils.WriteJSON(w, http.StatusCreated, charge, nil)

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
