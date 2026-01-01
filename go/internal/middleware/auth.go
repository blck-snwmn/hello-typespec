package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/blck-snwmn/hello-typespec/go/generated"
	"github.com/blck-snwmn/hello-typespec/go/internal/auth"
	"github.com/blck-snwmn/hello-typespec/go/internal/storage"
)

// AuthMiddleware validates Bearer tokens and adds user to context
func AuthMiddleware(authStore *storage.AuthStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				errorResponse(w, http.StatusUnauthorized, generated.UNAUTHORIZED, "Missing Authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				errorResponse(w, http.StatusUnauthorized, generated.UNAUTHORIZED, "Invalid Authorization header format")
				return
			}

			token := parts[1]

			// Validate token
			user, err := authStore.ValidateToken(token)
			if err != nil {
				errorResponse(w, http.StatusUnauthorized, generated.UNAUTHORIZED, "Invalid or expired token")
				return
			}

			// Add user to request context
			ctx := auth.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// errorResponse sends a standardized error response
func errorResponse(w http.ResponseWriter, statusCode int, code generated.ErrorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := generated.ErrorResponse{
		Error: struct {
			Code    generated.ErrorCode `json:"code"`
			Details interface{}         `json:"details,omitempty"`
			Message string              `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	}
	json.NewEncoder(w).Encode(response)
}
