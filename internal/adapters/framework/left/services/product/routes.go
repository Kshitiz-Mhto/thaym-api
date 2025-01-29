package product

import (
	"fmt"
	"net/http"
	"strconv"
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
	//public routes
	router.HandleFunc("/products", handler.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/product/{productId}", handler.handleGetProduct).Methods(http.MethodGet)
	// GET /selectiveproducts?ids=1,2,3
	router.HandleFunc("/selectiveproducts", handler.handleGetMultipleSelectiveProduct).Methods(http.MethodGet)

	//filtering-searching
	router.HandleFunc("/products/{category}", handler.handleFilteringByCategory).Methods(http.MethodGet)
	router.HandleFunc("/product/search/{queryStr}", handler.handleProductsSearching).Methods(http.MethodGet)
	router.HandleFunc("/products/{storeId}", auth.WithJWTAuth(handler.handleFilteringByStoreId, handler.userStore, "admin")).Methods(http.MethodGet)

	//inventory management
	router.HandleFunc("/product/quantity/{productId}/{num}", auth.WithJWTAuth(handler.handleUpdateProductQuatity, handler.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/product/stock/increase/{productId}/{num}", auth.WithJWTAuth(handler.handleIncrementProductStock, handler.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/product/stock/decrease/{productId}/{num}", auth.WithJWTAuth(handler.handleDecrementProductStock, handler.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/product/activate/{productId}", auth.WithJWTAuth(handler.handleProductActivation, handler.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/product/deactivate/{productId}", auth.WithJWTAuth(handler.handleProductDeactivation, handler.userStore)).Methods(http.MethodPut)

	//admin route
	router.HandleFunc("/products", auth.WithJWTAuth(handler.handleCreateProducts, handler.userStore, "admin", "storeowner")).Methods(http.MethodPost)
	router.HandleFunc("/delete_product/{productId}", auth.WithJWTAuth(handler.handleDeleteProduct, handler.userStore, "admin", "storeowner")).Methods(http.MethodDelete)
	router.HandleFunc("/update_product/{productId}", auth.WithJWTAuth(handler.handleUpdateProductById, handler.userStore, "admin", "storeowner")).Methods(http.MethodPut)
}

func (handler *ProductHandler) handleGetProducts(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	products, err := handler.store.GetAllProducts()

	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
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
		utils.WriteError(w, http.StatusNotFound, err)
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
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products, nil)
}

func (handler *ProductHandler) handleCreateProducts(w http.ResponseWriter, r *http.Request) {
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

	if r.Method != http.MethodDelete {
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

	utils.WriteJSON(w, http.StatusNoContent, productId, nil)
}

func (handler *ProductHandler) handleUpdateProductById(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}
	vars := mux.Vars(r)

	productId, ok := vars["productId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product ID"))
		return
	}

	_, err := handler.store.GetProductByID(productId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}

	err = handler.store.UpdateProductByID(productId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, productId, nil)
}

func (handler *ProductHandler) handleProductsSearching(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)
	query, ok := vars["queryStr"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing query string"))
		return
	}

	products, err := handler.store.SearchProducts(query)

	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products, nil)
}

func (handler *ProductHandler) handleFilteringByStoreId(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)
	storeId, ok := vars["storeId"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing storeId"))
		return
	}

	products, err := handler.store.GetProductsByStoreOwner(storeId)

	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products, nil)
}

func (handler *ProductHandler) handleFilteringByCategory(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)
	category, ok := vars["category"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing category"))
		return
	}

	products, err := handler.store.GetProductsByCategory(category)

	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products, nil)
}

func (handler *ProductHandler) handleProductDeactivation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}
	vars := mux.Vars(r)
	productId, ok := vars["productId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product id"))
		return
	}

	err := handler.store.DeactivateProduct(productId)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]bool{"success": true}, nil)

}

func (handler *ProductHandler) handleProductActivation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}
	vars := mux.Vars(r)
	productId, ok := vars["productId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product id"))
		return
	}

	err := handler.store.ActivateProduct(productId)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]bool{"success": true}, nil)

}

func (handler *ProductHandler) handleDecrementProductStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)
	productId, ok := vars["productId"]
	numStr, ko := vars["num"]

	if !ok || !ko {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing productId or quantity"))
		return
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("invalid quantity"))
		return
	}

	err = handler.store.DecreaseProductStock(productId, num)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]bool{"success": true}, nil)
}

func (handler *ProductHandler) handleIncrementProductStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)
	productId, ok := vars["productId"]
	numStr, ko := vars["num"]

	if !ok || !ko {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing productId or quantity"))
		return
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("invalid quantity"))
		return
	}

	err = handler.store.IncreaseProductStock(productId, num)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]bool{"success": true}, nil)
}

func (handler *ProductHandler) handleUpdateProductQuatity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)

	productId, ok := vars["productId"]
	numStr, ko := vars["num"]

	if !ok || !ko {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing productId or quantity"))
		return
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("invalid quantity"))
		return
	}

	err = handler.store.UpdateProductQuantity(productId, num)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]bool{"success": true}, nil)
}
