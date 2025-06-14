package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blck-snwmn/hello-typespec/go/generated"
	"github.com/blck-snwmn/hello-typespec/go/internal/handlers"
	"github.com/blck-snwmn/hello-typespec/go/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServer wraps the server handler for testing
type TestServer struct {
	*httptest.Server
	handler http.Handler
}

// setupTestServer creates a test server with a memory store
func setupTestServer(t testing.TB) *TestServer {
	t.Helper()

	memStore := store.NewMemoryStore()
	server := handlers.NewServer(memStore)
	handler := generated.Handler(server)

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	return &TestServer{
		Server:  ts,
		handler: handler,
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

// Test helpers for creating test data

// createTestUser creates a test user and returns its ID
func createTestUser(t testing.TB, server *TestServer, email, name string) string {
	t.Helper()

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

	rr := makeRequest(t, server, "POST", "/users", user)
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

	product := map[string]any{
		"name":        name,
		"description": "Test product description",
		"price":       price,
		"stock":       stock,
		"categoryId":  "1", // Default category
		"imageUrls":   []string{},
	}

	rr := makeRequest(t, server, "POST", "/products", product)
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

	category := map[string]any{
		"name":     name,
		"parentId": parentID,
	}

	rr := makeRequest(t, server, "POST", "/categories", category)
	require.Equal(t, http.StatusCreated, rr.Code, "failed to create test category")

	var response map[string]any
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	id, ok := response["id"].(string)
	require.True(t, ok, "response should contain id as string")

	return id
}

// addToCart adds an item to user's cart
func addToCart(t testing.TB, server *TestServer, userID, productID string, quantity int) {
	t.Helper()

	item := map[string]any{
		"productId": productID,
		"quantity":  quantity,
	}

	rr := makeRequest(t, server, "POST", "/carts/users/"+userID+"/items", item)
	require.Equal(t, http.StatusOK, rr.Code, "failed to add item to cart")
}

// createOrder creates an order for a user
func createOrder(t testing.TB, server *TestServer, userID string) string {
	t.Helper()

	// Get cart items first
	cartRR := makeRequest(t, server, "GET", "/carts/users/"+userID, nil)
	require.Equal(t, http.StatusOK, cartRR.Code, "failed to get cart")

	var cart map[string]any
	err := json.NewDecoder(cartRR.Body).Decode(&cart)
	require.NoError(t, err)

	cartItems := cart["items"].([]any)
	orderItems := make([]map[string]any, len(cartItems))
	for i, item := range cartItems {
		cartItem := item.(map[string]any)
		orderItems[i] = map[string]any{
			"productId": cartItem["productId"],
			"quantity":  cartItem["quantity"],
		}
	}

	order := map[string]any{
		"items": orderItems,
		"shippingAddress": map[string]any{
			"street":     "456 Order Ave",
			"city":       "Order City",
			"state":      "OC",
			"postalCode": "54321",
			"country":    "Order Country",
		},
	}

	rr := makeRequest(t, server, "POST", "/orders/users/"+userID, order)
	require.Equal(t, http.StatusCreated, rr.Code, "failed to create order")

	var response map[string]any
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	id, ok := response["id"].(string)
	require.True(t, ok, "response should contain id as string")

	return id
}

// decodeJSON is a helper to decode JSON response
func decodeJSON(rr *httptest.ResponseRecorder, v any) error {
	return json.NewDecoder(rr.Body).Decode(v)
}
