package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var Validate = validator.New()

// v --> destination variable where the parsed JSON data will be stored.
func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil || r.ContentLength == 0 {
		return fmt.Errorf("missing request body")
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// convenience wrapper for sending error responses in JSON format.
func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()}, nil)
}

// writes a JSON-encoded HTTP response.
func WriteJSON(w http.ResponseWriter, status int, v any, headers map[string]string) error {
	w.Header().Add("Content-Type", "application/json")
	if headers != nil {
		for key, value := range headers {
			w.Header().Set(key, value)
		}
	}
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func GenerateRandomUniqueIdentifier() string {
	return uuid.New().String()
}
