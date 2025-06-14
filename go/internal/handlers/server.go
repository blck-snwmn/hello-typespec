package handlers

import (
	"github.com/blck-snwmn/hello-typespec/go/generated"
	"github.com/blck-snwmn/hello-typespec/go/internal/store"
)

// Server implements the generated.ServerInterface
type Server struct {
	store store.Store
}

// NewServer creates a new Server instance
func NewServer(store store.Store) *Server {
	return &Server{
		store: store,
	}
}

// Ensure Server implements generated.ServerInterface
var _ generated.ServerInterface = (*Server)(nil)
