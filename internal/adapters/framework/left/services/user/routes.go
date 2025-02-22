package user

import (
	"fmt"
	"net/http"
	"time"

	"ecom-api/internal/adapters/framework/left/services/auth"
	"ecom-api/internal/adapters/framework/left/services/auth/token"
	"ecom-api/internal/application/core/types/entity"
	"ecom-api/internal/application/core/types/entity/payloads"
	"ecom-api/internal/ports/right/rports"
	"ecom-api/pkg/configs"
	"ecom-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	store      rports.UserStore
	tokenStore *token.TokenStore
}

func NewUserHandler(store rports.UserStore, tokenStore *token.TokenStore) *UserHandler {
	return &UserHandler{store: store, tokenStore: tokenStore}
}

func (handler *UserHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", handler.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", handler.handleRegister).Methods(http.MethodPost)
	router.HandleFunc("/register_confirm", handler.handleRegisterConfirmation).Methods(http.MethodPost)

	//admin routes
	router.HandleFunc("/users/{userID}", auth.WithJWTAuth(handler.handleGetUser, handler.store, "admin")).Methods(http.MethodGet)
	router.HandleFunc("/filter_users/{userRole}", auth.WithJWTAuth(handler.handleGetUsersByRole, handler.store, "admin")).Methods(http.MethodGet)
	router.HandleFunc("/user/lock/{userEmail}", auth.WithJWTAuth(handler.handleUserLocking, handler.store, "admin")).Methods(http.MethodPost)
}

func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user payloads.LoginUserPayload

	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

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

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid password"))
		return
	}

	secret := []byte(configs.Envs.JWTSecret)

	token, err := auth.CreateJWT(secret, u.ID, u.Role)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token}, nil)

}

func (h *UserHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user payloads.RegisterUserPayload

	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	// get json paload to check the user exist or not
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
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("user with email %s already exists", user.Email))
		return
	}

	// Generate a new token
	newToken, err := token.GenerateToken()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to generate token"))
		return
	}

	h.tokenStore.Set(user.Email, newToken, 7*time.Minute)

	// send verification code to user email
	HTMLTemplateEmailHandler(w, r, user.Email, map[string]string{"verification_code": newToken})

}

func (h *UserHandler) handleRegisterConfirmation(w http.ResponseWriter, r *http.Request) {
	var user payloads.RegisterUserConfirmationPayload

	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	// get json paload to check the user exist or not
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

	storedToken, exists := h.tokenStore.Get(user.Email)
	if !exists || storedToken != user.Token {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid or expired token"))
		return
	}

	h.tokenStore.Delete(user.Email)

	// hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(entity.User{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Password:   hashedPassword,
		Role:       user.Role,
		IsVerified: true,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"code": storedToken, "usertoken": user.Token}, nil)
}

func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)
	userID, ok := vars["userID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user, nil)
}

func (h *UserHandler) handleGetUsersByRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)

	userRole, ok := vars["userRole"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user role"))
		return
	}

	users, err := h.store.GetUsersByRole(userRole)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, users, nil)
}

func (h *UserHandler) handleUserLocking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	vars := mux.Vars(r)

	userEmail, ok := vars["userEmail"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user role"))
		return
	}

	userExists, err := h.store.GetUserByEmail(userEmail)
	if err != nil || userExists == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user with email %s not found", userEmail))
		return
	}

	err = h.store.SetUserLocking(userEmail, true)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to lock user: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "User locked successfully", "email": userEmail, "isLocked": true}, nil)
}
