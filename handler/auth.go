package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nileshshrs/infinite-storage/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// Register handles POST /sign-up
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input service.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
		return
	}

	// If userAgent not sent in body, take it from header
	if input.UserAgent == "" {
		input.UserAgent = r.UserAgent()
	}

	user, accessToken, refreshToken, err := h.service.Register(input)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}


// func(u *User) Login(w http.ResponseWriter, r *http.Request){
// 	fmt.Println("Login endpoint hit")
// 	w.WriteHeader(http.StatusNotImplemented)
// }
// func(u *User) Update(w http.ResponseWriter, r *http.Request){
// 	fmt.Println("Login endpoint hit")
// 	w.WriteHeader(http.StatusNotImplemented)
// }