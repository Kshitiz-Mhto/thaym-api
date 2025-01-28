package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"ecom-api/internal/application/core/types/entity"
	"ecom-api/internal/application/core/types/entity/payloads"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewUserHandler(userStore, nil)

	t.Run("should fail if user payload is invalid", func(t *testing.T) {
		payload := payloads.RegisterUserPayload{
			FirstName: "user",
			LastName:  "123",
			Email:     "invalid",
			Password:  "asd",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatalf("error requesting %v", err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

	})

	t.Run("should pass if user payload is valid", func(t *testing.T) {
		payload := payloads.RegisterUserPayload{
			FirstName: "user",
			LastName:  "x00",
			Email:     "valid@mail.com",
			Password:  "asd87665757",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatalf("error requesting %v", err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}

	})

	t.Run("should fail if the user ID is not a number", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/user/abc", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}", handler.handleGetUser).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should handle get user by ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/user/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}", handler.handleGetUser).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type mockUserStore struct{}

// SetUserLocking implements rports.UserStore.
func (m *mockUserStore) SetUserLocking(email string, isLocked bool) error {
	panic("unimplemented")
}

// GetUsersByRole implements rports.UserStore.
func (m *mockUserStore) GetUsersByRole(role string) ([]*entity.User, error) {
	panic("unimplemented")
}

func (m *mockUserStore) GetUserByID(id string) (*entity.User, error) {
	return &entity.User{}, nil
}

func (m *mockUserStore) GetUserByEmail(email string) (*entity.User, error) {
	return &entity.User{}, nil
}
func (m *mockUserStore) CreateUser(u entity.User) error {
	return nil
}
func (m *mockUserStore) UpdateUser(u entity.User) error {
	return nil
}
