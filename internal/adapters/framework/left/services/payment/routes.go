package payment

import (
	"fmt"
	"net/http"

	"ecom-api/internal/adapters/framework/left/services/auth"
	"ecom-api/internal/application/core/types/entity/payloads"
	"ecom-api/internal/ports/right/rports"
	"ecom-api/utils"

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

	router.HandleFunc("/create/customer", auth.WithJWTAuth(handler.handleCustomerCreation, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/customer/{id}", auth.WithJWTAuth(handler.handleGetCustomerById, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/customers", auth.WithJWTAuth(handler.handleGetCustomers, handler.userStore)).Methods(http.MethodPost)

}

func (handler *PaymentHandler) handleCustomerCreation(w http.ResponseWriter, r *http.Request) {
	var customer payloads.CustomerPayload

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
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
	if r.Method != http.MethodPost {
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
	if r.Method != http.MethodPost {
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
