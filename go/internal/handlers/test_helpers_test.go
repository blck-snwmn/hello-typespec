package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blck-snwmn/hello-typespec/go/generated"
	"github.com/blck-snwmn/hello-typespec/go/internal/handlers"
	"github.com/blck-snwmn/hello-typespec/go/internal/middleware"
	"github.com/blck-snwmn/hello-typespec/go/internal/storage"
	"github.com/blck-snwmn/hello-typespec/go/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServer wraps the server handler for testing
type TestServer struct {
	*httptest.Server
	handler     http.Handler
	authStorage *storage.AuthStore
	store       store.Store
}

// setupTestServer creates a test server with a memory store
func setupTestServer(t testing.TB) *TestServer {
	t.Helper()

	memStore := store.NewMemoryStore()
	authStorage := storage.NewAuthStore()
	server := handlers.NewServer(memStore, authStorage)
	
	// Create handler with auth middleware applied to protected routes
	authMiddleware := middleware.AuthMiddleware(authStorage)
	handler := handlers.CreateHandlerWithMiddleware(server, authMiddleware)

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	return &TestServer{
		Server:      ts,
		handler:     handler,
		authStorage: authStorage,
		store:       memStore,
	}
}


// makeRequest is a helper to make HTTP requests in tests
func makeRequest(t testing.TB, server *TestServer, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var req *http.Request
	var err error

	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		req, err = http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, path, nil)
		require.NoError(t, err)
	}

	rr := httptest.NewRecorder()
	server.handler.ServeHTTP(rr, req)

	return rr
}

// assertStatus checks if the response has the expected status code
func assertStatus(t testing.TB, rr *httptest.ResponseRecorder, expectedStatus int) {
	t.Helper()
	assert.Equal(t, expectedStatus, rr.Code, "unexpected status code")
}

// assertErrorResponse checks if the error response has the expected format
func assertErrorResponse(t testing.TB, rr *httptest.ResponseRecorder, expectedCode string) {
	t.Helper()

	var errorResp map[string]any
	err := json.NewDecoder(rr.Body).Decode(&errorResp)
	require.NoError(t, err)

	errorObj, ok := errorResp["error"].(map[string]any)
	require.True(t, ok, "response should contain error object")

	assert.Equal(t, expectedCode, errorObj["code"], "unexpected error code")
	assert.NotEmpty(t, errorObj["message"], "error message should not be empty")
}

// assertPaginatedResponse checks pagination response format
func assertPaginatedResponse(t testing.TB, rr *httptest.ResponseRecorder, expectedTotal, expectedLimit, expectedOffset int) map[string]any {
	t.Helper()

	var response map[string]any
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response, "items")
	assert.Contains(t, response, "total")
	assert.Contains(t, response, "limit")
	assert.Contains(t, response, "offset")

	assert.Equal(t, float64(expectedTotal), response["total"])
	assert.Equal(t, float64(expectedLimit), response["limit"])
	assert.Equal(t, float64(expectedOffset), response["offset"])

	items, ok := response["items"].([]any)
	require.True(t, ok, "items should be an array")
	assert.NotNil(t, items)

	return response
}

// Auth test helpers

// makeAuthenticatedRequest makes an HTTP request with authentication
func makeAuthenticatedRequest(t testing.TB, server *TestServer, method, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()

	var req *http.Request
	var err error

	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		req, err = http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, path, nil)
		require.NoError(t, err)
	}

	// Add authorization header
	req.Header.Set("Authorization", "Bearer "+token)

	// For /auth/me endpoint, we need to simulate the middleware behavior
	if path == "/auth/me" && token != "" {
		user, err := server.authStorage.ValidateToken(token)
		if err == nil && user != nil {
			ctx := context.WithValue(req.Context(), "user", user)
			req = req.WithContext(ctx)
		}
	}

	rr := httptest.NewRecorder()
	server.handler.ServeHTTP(rr, req)

	return rr
}

// loginTestUser logs in a test user and returns the token
func loginTestUser(t testing.TB, server *TestServer, email, password string) string {
	t.Helper()

	loginReq := generated.LoginRequest{
		Email:    email,
		Password: password,
	}

	rr := makeRequest(t, server, "POST", "/auth/login", loginReq)
	assertStatus(t, rr, http.StatusOK)

	var loginResp generated.LoginResponse
	err := json.NewDecoder(rr.Body).Decode(&loginResp)
	require.NoError(t, err)

	return loginResp.AccessToken
}

// setupTestServerWithAuth creates a test server and logs in a default user
func setupTestServerWithAuth(t testing.TB) (*TestServer, string, string) {
	t.Helper()
	
	server := setupTestServer(t)
	token := loginTestUser(t, server, "alice@example.com", "password123")
	
	return server, "550e8400-e29b-41d4-a716-446655440001", token
}

// Test helpers for creating test data

// createTestUser creates a test user and returns its ID
func createTestUser(t testing.TB, server *TestServer, email, name string) string {
	t.Helper()

	// Get a token for creating users
	token := loginTestUser(t, server, "alice@example.com", "password123")

	user := map[string]any{
		"email": email,
		"name":  name,
		"address": map[string]any{
			"street":     "123 Test St",
			"city":       "Test City",
			"state":      "TC",
			"postalCode": "12345",
			"country":    "Test Country",
		},
	}

	rr := makeAuthenticatedRequest(t, server, "POST", "/users", user, token)
	require.Equal(t, http.StatusCreated, rr.Code, "failed to create test user")

	var response map[string]any
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	id, ok := response["id"].(string)
	require.True(t, ok, "response should contain id as string")

	return id
}

// createTestProduct creates a test product and returns its ID
func createTestProduct(t testing.TB, server *TestServer, name string, price float64, stock int) string {
	t.Helper()

	// Get a token for creating products
	token := loginTestUser(t, server, "alice@example.com", "password123")

	product := map[string]any{
		"name":        name,
		"description": "Test product description",
		"price":       price,
		"stock":       stock,
		"categoryId":  "1", // Default category
		"imageUrls":   []string{},
	}

	rr := makeAuthenticatedRequest(t, server, "POST", "/products", product, token)
	require.Equal(t, http.StatusCreated, rr.Code, "failed to create test product")

	var response map[string]any
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	id, ok := response["id"].(string)
	require.True(t, ok, "response should contain id as string")

	return id
}

// createTestCategory creates a test category and returns its ID
func createTestCategory(t testing.TB, server *TestServer, name string, parentID *string) string {
	t.Helper()

	// Get a token for creating categories
	token := loginTestUser(t, server, "alice@example.com", "password123")

	category := map[string]any{
		"name":     name,
		"parentId": parentID,
	}

	rr := makeAuthenticatedRequest(t, server, "POST", "/categories", category, token)
	require.Equal(t, http.StatusCreated, rr.Code, "failed to create test category")

	var response map[string]any
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	id, ok := response["id"].(string)
	require.True(t, ok, "response should contain id as string")

	return id
}

// decodeJSON is a helper to decode JSON response
func decodeJSON(rr *httptest.ResponseRecorder, v any) error {
	return json.NewDecoder(rr.Body).Decode(v)
}

