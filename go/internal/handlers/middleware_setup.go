package handlers

import (
	"net/http"

	"github.com/blck-snwmn/hello-typespec/go/generated"
)

// ProtectedRoutes defines which routes require authentication
var ProtectedRoutes = map[string]bool{
	"/carts":       true,
	"/orders":      true,
	"/users":       true,
	"/auth/me":     true,
	"/auth/logout": true,
}

// CreateHandlerWithMiddleware creates an HTTP handler with authentication middleware applied to protected routes
func CreateHandlerWithMiddleware(server generated.ServerInterface, authMiddleware func(http.Handler) http.Handler) http.Handler {
	// Use the generated handler as the base
	baseHandler := generated.Handler(server)

	// Wrap the handler to apply middleware selectively
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the route requires authentication
		requiresAuth := false
		for route := range ProtectedRoutes {
			if r.URL.Path == route ||
				(len(r.URL.Path) > len(route) && r.URL.Path[:len(route)+1] == route+"/") {
				requiresAuth = true
				break
			}
		}

		if requiresAuth {
			// Apply auth middleware
			authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				baseHandler.ServeHTTP(w, r)
			})).ServeHTTP(w, r)
		} else {
			// No auth required
			baseHandler.ServeHTTP(w, r)
		}
	})
}
