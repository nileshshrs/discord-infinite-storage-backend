package middleware

import (
	"context"
	"net/http"

	"github.com/nileshshrs/infinite-storage/utils"
)

// Keys for context
type contextKey string

const (
	UserIDKey    contextKey = "userID"
	SessionIDKey contextKey = "sessionID"
)

// Authenticate middleware checks for accessToken in cookies and verifies it
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookies
		cookie, err := r.Cookie("accessToken")
		if err != nil || cookie.Value == "" {
			http.Error(w, `{"error":"No access token provided"}`, http.StatusUnauthorized)
			return
		}

		// Verify token
		claims, err := utils.VerifyAccessToken(cookie.Value)
		if err != nil {
			http.Error(w, `{"error":"Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// Extract claims
		userID, _ := claims["userID"].(string)
		sessionID, _ := claims["sessionID"].(string)

		// Put into context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, SessionIDKey, sessionID)

		// Continue with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
