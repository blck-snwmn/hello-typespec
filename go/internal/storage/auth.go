package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("not found")
)

// AuthUser represents an authenticated user
type AuthUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// AuthSession represents a user session with token
type AuthSession struct {
	Token     string
	User      AuthUser
	ExpiresAt time.Time
}

// AuthStore handles authentication token storage
type AuthStore struct {
	sessions sync.Map // map[string]AuthSession
	users    map[string]struct {
		password string
		user     AuthUser
	}
}

// NewAuthStore creates a new authentication store
func NewAuthStore() *AuthStore {
	// Initialize with mock users
	return &AuthStore{
		users: map[string]struct {
			password string
			user     AuthUser
		}{
			"alice@example.com": {
				password: "password123", // In production, this would be hashed
				user: AuthUser{
					ID:    "550e8400-e29b-41d4-a716-446655440001",
					Email: "alice@example.com",
					Name:  "Alice Johnson",
				},
			},
			"bob@example.com": {
				password: "password456", // In production, this would be hashed
				user: AuthUser{
					ID:    "550e8400-e29b-41d4-a716-446655440002",
					Email: "bob@example.com",
					Name:  "Bob Smith",
				},
			},
		},
	}
}

// Login authenticates a user and returns a token
func (s *AuthStore) Login(email, password string) (*AuthSession, error) {
	user, exists := s.users[email]
	if !exists || user.password != password {
		return nil, ErrNotFound
	}

	// Generate token
	token := uuid.New().String()
	session := AuthSession{
		Token:     token,
		User:      user.user,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Store session
	s.sessions.Store(token, session)

	return &session, nil
}

// Logout removes a session
func (s *AuthStore) Logout(token string) error {
	s.sessions.Delete(token)
	return nil
}

// ValidateToken checks if a token is valid and returns the user
func (s *AuthStore) ValidateToken(token string) (*AuthUser, error) {
	value, ok := s.sessions.Load(token)
	if !ok {
		return nil, ErrNotFound
	}

	session := value.(AuthSession)

	// Check if token is expired
	if time.Now().After(session.ExpiresAt) {
		s.sessions.Delete(token)
		return nil, ErrNotFound
	}

	return &session.User, nil
}

// CleanupExpiredTokens removes expired sessions
func (s *AuthStore) CleanupExpiredTokens() {
	now := time.Now()
	s.sessions.Range(func(key, value any) bool {
		session := value.(AuthSession)
		if now.After(session.ExpiresAt) {
			s.sessions.Delete(key)
		}
		return true
	})
}