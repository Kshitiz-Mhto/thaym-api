package user

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	// store interfaces.UserStore
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (handler *UserHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", handler.handleLogin).Methods("POST")
	router.HandleFunc("/register", handler.handleRegister).Methods("POST")

	//admin routes
	// router.HandleFunc("/users/{userID}", auth.WithJWTAuth(handler.handleGetUser, handler.store)).Methods(http.MethodGet)
}

func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Print("jjjjj000000")
}

func (h *UserHandler) handleRegister(w http.ResponseWriter, r *http.Request) {

}
