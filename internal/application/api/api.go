package api

import (
	"database/sql"
	"log"
	"net/http"

	"ecom-api/internal/adapters/framework/left/services/auth/token"
	"ecom-api/internal/adapters/framework/left/services/cart"
	"ecom-api/internal/adapters/framework/left/services/payment"
	"ecom-api/internal/adapters/framework/left/services/product"
	"ecom-api/internal/adapters/framework/left/services/user"
	order "ecom-api/internal/adapters/framework/right/order_repo"
	paymentrepo "ecom-api/internal/adapters/framework/right/payment_repo"
	"ecom-api/internal/adapters/framework/right/product_repo"
	"ecom-api/internal/adapters/framework/right/user_repo"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (api *APIServer) Run() error {

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	tokenStore := token.NewTokenStore()
	userStore := user_repo.NewStore(api.db)
	userHandler := user.NewUserHandler(userStore, tokenStore)
	userHandler.RegisterRoutes(subrouter)

	productStore := product_repo.NewStore(api.db)
	productHandler := product.NewProductHandler(productStore, userStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(api.db)

	cartHandler := cart.NewCartHandler(productStore, orderStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	paymentStore := paymentrepo.NewPaymentStore()
	paymentHandler := payment.NewPaymentHandler(paymentStore, userStore)
	paymentHandler.RegisterRoutes(subrouter)

	log.Println("Listening to ", api.addr)
	return http.ListenAndServe(api.addr, router)

}
