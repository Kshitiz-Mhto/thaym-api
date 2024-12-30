package api

import (
	"database/sql"
	"log"
	"net/http"

	"ecom-api/internal/adapters/framework/left/services/user"
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

	userStore := user_repo.NewStore(api.db)
	userHandler := user.NewUserHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	log.Println("Listening to ", api.addr)
	return http.ListenAndServe(api.addr, router)

}
