package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nileshshrs/infinite-storage/utils"
)

// contextKey is a private type to avoid collisions in context
type contextKey string

const (
	UserIDKey    contextKey = "userID"
	SessionIDKey contextKey = "sessionID"
)

type errorResponse struct {
	Error string `json:"error"`
}

// Authenticate middleware checks for accessToken in cookies and verifies it
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookies
		cookie, err := r.Cookie("accessToken")
		if err != nil || cookie.Value == "" {
			writeJSONError(w, "No access token provided", http.StatusUnauthorized)
			return
		}

		// Verify token
		claims, err := utils.VerifyAccessToken(cookie.Value)
		if err != nil {
			writeJSONError(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		userID, _ := claims["userID"].(string)
		sessionID, _ := claims["sessionID"].(string)

		if userID == "" || sessionID == "" {
			writeJSONError(w, "Invalid token payload", http.StatusUnauthorized)
			return
		}

		// Put into context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, SessionIDKey, sessionID)

		// Continue with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// helper to write JSON errors
func writeJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: message})
}
