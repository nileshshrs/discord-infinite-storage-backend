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

	user, err := h.service.Register(input)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// return the created user without password
	user.Password = ""
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}



// func(u *User) Login(w http.ResponseWriter, r *http.Request){
// 	fmt.Println("Login endpoint hit")
// 	w.WriteHeader(http.StatusNotImplemented)
// }
// func(u *User) Update(w http.ResponseWriter, r *http.Request){
// 	fmt.Println("Login endpoint hit")
// 	w.WriteHeader(http.StatusNotImplemented)
// }