package product

import (
	"fmt"
	"net/http"
	"strings"

	"ecom-api/internal/adapters/framework/left/services/auth"
	"ecom-api/internal/application/core/types/entity/payloads"
	"ecom-api/internal/ports/right/rports"
	"ecom-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	store     rports.ProductStore
	userStore rports.UserStore
}

func NewProductHandler(store rports.ProductStore, userStore rports.UserStore) *ProductHandler {
	return &ProductHandler{store: store, userStore: userStore}
}

func (handler *ProductHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", handler.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products/{productId}", handler.handleGetProduct).Methods(http.MethodGet)
	// GET /selectiveproducts?ids=1,2,3
	router.HandleFunc("/selectiveproducts", handler.handleGetMultipleSelectiveProduct).Methods(http.MethodGet)
	//admin route
	router.HandleFunc("/create_product", auth.WithJWTAuth(handler.handleCreateProduct, handler.userStore)).Methods(http.MethodPost)
}

func (handler *ProductHandler) handleGetProducts(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	products, err := handler.store.GetAllProducts()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products, nil)
}

func (handler *ProductHandler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}
	// Extract Path Parameter (productID) from the URL
	vars := mux.Vars(r)
	productId, ok := vars["productId"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product ID"))
		return
	}

	product, err := handler.store.GetProductByID(productId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, product, nil)

}

func (handler *ProductHandler) handleGetMultipleSelectiveProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	idsParam := r.URL.Query().Get("ids")

	if idsParam == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing 'ids' query parameter"))
		return
	}

	ids := strings.Split(idsParam, ",")

	products, err := handler.store.GetProductsByID(ids)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products, nil)
}

func (handler *ProductHandler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var product payloads.CreateProductPayload

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	if err := utils.ParseJSON(r, &product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(product); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors)
		return
	}

	err := handler.store.CreateProduct(product)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, product, nil)
}
