package cart

import (
	"ecom-api/internal/adapters/framework/left/services/auth"
	"ecom-api/internal/application/core/types/entity"
	"ecom-api/internal/application/core/types/entity/payloads"
	"ecom-api/internal/ports/right/rports"
	"ecom-api/utils"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type CartHandler struct {
	store        rports.ProductStore
	orderStore   rports.OrderStore
	userStore    rports.UserStore
	paymentStore rports.PaymentStore
}

func NewCartHandler(store rports.ProductStore, orderStore rports.OrderStore, userStore rports.UserStore, paymentStore rports.PaymentStore) *CartHandler {
	return &CartHandler{
		store:        store,
		orderStore:   orderStore,
		userStore:    userStore,
		paymentStore: paymentStore,
	}
}

func (handler *CartHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(handler.handleCartCheckout, handler.userStore, "admin", "storeowner", "user")).Methods(http.MethodPost)
	router.HandleFunc("/cart/delete_orderitem/{orderItemId}", auth.WithJWTAuth(handler.handleOrderItemDeletion, handler.userStore, "admin", "storeowner", "user")).Methods(http.MethodDelete)

	router.HandleFunc("/order/delete/{orderId}", auth.WithJWTAuth(handler.handleOrderDeletion, handler.userStore, "admin", "storeowner", "user")).Methods(http.MethodDelete)
	router.HandleFunc("/order/{orderId}", auth.WithJWTAuth(handler.handleGetOrderById, handler.userStore, "admin", "storeowner", "user")).Methods(http.MethodGet)
	router.HandleFunc("/order/user/{userId}", auth.WithJWTAuth(handler.handleGetOrderByUserId, handler.userStore, "admin", "storeowner", "user")).Methods(http.MethodGet)
	router.HandleFunc("/order/update/paymentstatus/{orderId}", auth.WithJWTAuth(handler.handlerUpdateOrderPaymentStatus, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/order/update/status/{orderId}", auth.WithJWTAuth(handler.handlerUpdateOrderStatus, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/orderitem/{orderId}", auth.WithJWTAuth(handler.handleGetOrderItemByOrderId, handler.userStore, "admin", "storeowner", "user")).Methods(http.MethodGet)

}

func (handler *CartHandler) handleCartCheckout(w http.ResponseWriter, r *http.Request) {

	userID := auth.GetUserIDFromContext(r.Context())

	var cart payloads.CartCheckoutPayload
	// var storeCustomer *entity.User

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors)
		return
	}

	productIDs, err := getCartItemsIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	products, err := handler.store.GetProductsByIDs(productIDs)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderId, subTotal, totalPrice, err := handler.createOrder(products, cart.Items, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// -------------------- for checkout session -------------------
	// storeCustomer, err = handler.userStore.GetUserByID(userID)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusInternalServerError, err)
	// 	return
	// }

	// session, err := handler.paymentStore.CreateCheckoutSession(orderId, totalPrice, storeCustomer.Email)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error creating checkout session: %v", err))
	// 	return
	// }

	// utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
	// 	"paid":     session.PaymentIntent.Amount,
	// 	"customer": session.Customer,
	// 	"items":    session.DisplayItems,
	// 	"balance":  session.Customer.Balance,
	// 	"currency": session.PaymentIntent.Currency,
	// }, nil)

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"total":    totalPrice,
		"subTotal": subTotal,
		"userId":   userID,
		"orderId":  orderId,
	}, nil)
}

func (handler *CartHandler) handleOrderDeletion(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)

	orderId, ok := vars["orderId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing order ID"))
		return
	}

	err := handler.orderStore.DeleteOrder(orderId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"orderId": orderId}, nil)

}

func (handler *CartHandler) handleOrderItemDeletion(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)

	orderItemId, ok := vars["orderItemId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing order ID"))
		return
	}

	err := handler.orderStore.DeleteOrderItem(orderItemId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"orderItemId": orderItemId}, nil)
}

func (handler *CartHandler) handleGetOrderById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)
	orderId, ok := vars["orderId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing order ID"))
		return
	}

	order, err := handler.orderStore.GetOrderByID(orderId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, order, nil)
}

func (handler *CartHandler) handleGetOrderByUserId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)
	userId, ok := vars["userId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	var orders []*entity.Order
	var err error

	orders, err = handler.orderStore.GetOrdersByUserID(userId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, orders, nil)
}

func (handler *CartHandler) handleGetOrderItemByOrderId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)
	orderId, ok := vars["orderId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing order ID"))
		return
	}

	var orderitems []*entity.OrderItem
	var err error

	orderitems, err = handler.orderStore.GetOrderItemsByOrderId(orderId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, orderitems, nil)
}

func (handler *CartHandler) handlerUpdateOrderPaymentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)
	var status string

	orderId, ok := vars["orderId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing order ID"))
		return
	}

	err := handler.orderStore.UpdateOrderPaymentStatus(orderId, status)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"orderId": orderId, "paymentSatus": status}, nil)
}

func (handler *CartHandler) handlerUpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)
	var status string

	orderId, ok := vars["orderId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing order ID"))
		return
	}

	err := handler.orderStore.UpdateOrderPaymentStatus(orderId, status)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"orderId": orderId, "orderSatus": status}, nil)
}
