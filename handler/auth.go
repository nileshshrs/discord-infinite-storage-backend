package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nileshshrs/infinite-storage/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// setAuthCookies sets HTTP-only cookies for access and refresh tokens
func setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true if using HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(15 * time.Minute / time.Second), // 15 minutes
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true if using HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(30 * 24 * time.Hour / time.Second), // 30 days
	})
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

	setAuthCookies(w, accessToken, refreshToken)

	user.Password = ""
	response := map[string]interface{}{
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles POST /login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input service.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
		return
	}

	user, accessToken, refreshToken, err := h.service.Login(input)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	setAuthCookies(w, accessToken, refreshToken)

	user.Password = ""
	response := map[string]interface{}{
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RefreshToken handles POST /refresh-token
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Get refresh token from cookie
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		http.Error(w, `{"error":"missing refresh token"}`, http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value

	// 2. Call service to refresh tokens
	accessToken, newRefreshToken, err := h.service.RefreshTokens(refreshToken)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	// 3. Set new cookies
	setAuthCookies(w, accessToken, newRefreshToken)

	// 4. Return success message
	response := map[string]interface{}{
		"message":      "access token has been refreshed",
		"accessToken":  accessToken,
		"refreshToken": newRefreshToken, // could be empty if not refreshed
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
