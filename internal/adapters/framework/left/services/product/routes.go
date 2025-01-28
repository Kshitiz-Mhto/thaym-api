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
	router.HandleFunc("/delete_product/{productId}", auth.WithJWTAuth(handler.handleDeleteProduct, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/update_product/{productId}", auth.WithJWTAuth(handler.handleUpdateProductById, handler.userStore)).Methods(http.MethodPost)
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
	// decodedProductTags := utils.DecodeTheByteData(string(product.Tags))

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// log.Println(decodedProductTags)

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

	products, err := handler.store.GetProductsByIDs(ids)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products, nil)
}

func (handler *ProductHandler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var products []payloads.CreateProductPayload

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	if err := utils.ParseJSON(r, &products); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	for _, product := range products {
		if err := utils.Validate.Struct(product); err != nil {
			errors := err.(validator.ValidationErrors)
			utils.WriteError(w, http.StatusBadRequest, errors)
			return
		}
	}

	for _, product := range products {
		err := handler.store.CreateProduct(product)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}

	utils.WriteJSON(w, http.StatusCreated, products, nil)
}

func (handler *ProductHandler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)
	productId, ok := vars["productId"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product ID"))
		return
	}

	err := handler.store.DeleteProductByID(productId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, productId, nil)
}

func (handler *ProductHandler) handleUpdateProductById(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}
	vars := mux.Vars(r)

	productId, ok := vars["productId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product ID"))
		return
	}

	err := handler.store.UpdateProductByID(productId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, productId, nil)
}
