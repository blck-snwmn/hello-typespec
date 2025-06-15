package auth

import (
	"context"

	"github.com/blck-snwmn/hello-typespec/go/internal/storage"
)

// userKey is a custom type for context keys to avoid collisions
type userKey struct{}

// UserKey is the context key for storing user information
var UserKey = userKey{}

// WithUser adds a user to the context
func WithUser(ctx context.Context, user *storage.AuthUser) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

// GetUser retrieves a user from the context
func GetUser(ctx context.Context) (*storage.AuthUser, bool) {
	user, ok := ctx.Value(UserKey).(*storage.AuthUser)
	return user, ok
}
