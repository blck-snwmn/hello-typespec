package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/blck-snwmn/hello-typespec/go/generated"
	"github.com/blck-snwmn/hello-typespec/go/internal/storage"
)

// AuthHandlers handles authentication endpoints
type AuthHandlers struct {
	authStore *storage.AuthStore
}

// NewAuthHandlers creates a new auth handlers instance
func NewAuthHandlers(authStore *storage.AuthStore) *AuthHandlers {
	return &AuthHandlers{
		authStore: authStore,
	}
}

// Login handles POST /auth/login
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req generated.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, generated.BADREQUEST, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		errorResponse(w, http.StatusBadRequest, generated.VALIDATIONERROR, "Email and password are required")
		return
	}

	session, err := h.authStore.Login(req.Email, req.Password)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, generated.UNAUTHORIZED, "Invalid email or password")
		return
	}

	response := generated.LoginResponse{
		AccessToken: session.Token,
		TokenType:   "Bearer",
		ExpiresIn:   86400, // 24 hours in seconds
		User: struct {
			Email string `json:"email"`
			Id    string `json:"id"`
			Name  string `json:"name"`
		}{
			Id:    session.User.ID,
			Email: session.User.Email,
			Name:  session.User.Name,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logout handles POST /auth/logout
func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token == "" {
		errorResponse(w, http.StatusUnauthorized, generated.UNAUTHORIZED, "Missing authorization token")
		return
	}

	h.authStore.Logout(token)

	response := generated.OkResponse{
		Message: "Logged out successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCurrentUser handles GET /auth/me
func (h *AuthHandlers) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// User is already validated by middleware
	user := r.Context().Value("user").(*storage.AuthUser)

	response := generated.AuthUser{
		Id:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// extractToken extracts the bearer token from Authorization header
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}