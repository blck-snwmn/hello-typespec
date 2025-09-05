package handlers

import (
	"net/http"

	"github.com/blck-snwmn/hello-typespec/go/generated"
	"github.com/blck-snwmn/hello-typespec/go/internal/storage"
	"github.com/blck-snwmn/hello-typespec/go/internal/store"
)

// Server implements the generated.ServerInterface
type Server struct {
	store       store.Store
	authHandler *AuthHandlers
}

// NewServer creates a new Server instance
func NewServer(store store.Store, authStore *storage.AuthStore) *Server {
	return &Server{
		store:       store,
		authHandler: NewAuthHandlers(authStore),
	}
}

// AuthServiceLogin handles user login
func (s *Server) AuthServiceLogin(w http.ResponseWriter, r *http.Request) {
	s.authHandler.Login(w, r)
}

// AuthServiceLogout handles user logout
func (s *Server) AuthServiceLogout(w http.ResponseWriter, r *http.Request) {
    s.authHandler.Logout(w, r)
}

// AuthServiceGetCurrentUser gets the current authenticated user
func (s *Server) AuthServiceGetCurrentUser(w http.ResponseWriter, r *http.Request) {
    s.authHandler.GetCurrentUser(w, r)
}

// Ensure Server implements generated.ServerInterface
var _ generated.ServerInterface = (*Server)(nil)
