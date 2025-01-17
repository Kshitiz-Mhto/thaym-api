package cart

import (
	"ecom-api/internal/adapters/framework/left/services/auth"
	"ecom-api/internal/application/core/types/entity/payloads"
	"ecom-api/internal/ports/right/rports"
	"ecom-api/utils"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type CartHandler struct {
	store      rports.ProductStore
	orderStore rports.OrderStore
	userStore  rports.UserStore
}

func NewCartHandler(store rports.ProductStore, orderStore rports.OrderStore, userStore rports.UserStore) *CartHandler {
	return &CartHandler{
		store:      store,
		orderStore: orderStore,
		userStore:  userStore,
	}
}

func (handler *CartHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(handler.handleCartCheckout, handler.userStore)).Methods(http.MethodPost)
}

func (handler *CartHandler) handleCartCheckout(w http.ResponseWriter, r *http.Request) {

	userID := auth.GetUserIDFromContext(r.Context())

	var cart payloads.CartCheckoutPayload

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

	orderId, subTotalPrice, totalPrice, err := handler.createOrder(products, cart.Items, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"total_price":                    totalPrice,
		"total_price_before_tax_and_dis": subTotalPrice,
		"order_id":                       orderId,
	}, nil)
}
