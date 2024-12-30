package user

import (
	"fmt"
	"net/http"

	"ecom-api/internal/application/core/services/auth"
	"ecom-api/internal/application/core/types/entity"
	"ecom-api/internal/application/core/types/entity/payloads"
	"ecom-api/internal/application/core/types/interfaces"
	"ecom-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	store interfaces.UserStore
}

func NewUserHandler(store interfaces.UserStore) *UserHandler {
	return &UserHandler{store: store}
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
	// get json paload to check the user exist or not
	var user payloads.RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//Field Validation
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	//Check for Existing User
	_, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", user.Email))
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// if it dont exist
	err = h.store.CreateUser(entity.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil, nil)

}
