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

	user, accessToken, refreshToken, err := h.service.Register(input)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Set HTTP-only cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true if using HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60, // 15 minutes
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true if using HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   30 * 24 * 60 * 60, // 30 days
	})

	// Return user data in response (without tokens in body if preferred)
	user.Password = "" // hide password
	response := map[string]interface{}{
		"user": user,
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