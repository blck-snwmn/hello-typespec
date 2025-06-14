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

// TestServer wraps the server with test utilities
type TestServer struct {
	server  *handlers.Server
	handler http.Handler
	store   store.Store
}

// setupTestServer creates a new test server with a clean memory store
func setupTestServer(t *testing.T) *TestServer {
	t.Helper()

	memStore := store.NewMemoryStore()
	server := handlers.NewServer(memStore)
	handler := generated.Handler(server)

	return &TestServer{
		server:  server,
		handler: handler,
		store:   memStore,
	}
}

// makeRequest executes an HTTP request against the test server
func makeRequest(t *testing.T, server *TestServer, method, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		require.NoError(t, err, "failed to marshal request body")
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	rr := httptest.NewRecorder()
	server.handler.ServeHTTP(rr, req)

	return rr
}

// assertStatus checks if the response status code matches the expected value
func assertStatus(t *testing.T, rr *httptest.ResponseRecorder, want int) {
	t.Helper()
	assert.Equal(t, want, rr.Code, "unexpected status code")
}

// assertErrorResponse validates an error response
func assertErrorResponse(t *testing.T, rr *httptest.ResponseRecorder, wantCode string) {
	t.Helper()

	var errResp generated.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	require.NoError(t, err, "failed to decode error response")

	assert.Equal(t, wantCode, errResp.Error.Code, "unexpected error code")
	assert.NotEmpty(t, errResp.Error.Message, "error message should not be empty")
}

// assertPaginatedResponse validates a paginated response structure
func assertPaginatedResponse(t *testing.T, rr *httptest.ResponseRecorder, wantTotal, wantLimit, wantOffset int) map[string]interface{} {
	t.Helper()

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err, "failed to decode response")

	assert.Contains(t, response, "items", "response should contain items")
	assert.Contains(t, response, "total", "response should contain total")
	assert.Contains(t, response, "limit", "response should contain limit")
	assert.Contains(t, response, "offset", "response should contain offset")

	items, ok := response["items"].([]interface{})
	require.True(t, ok, "items should be an array")

	total := int(response["total"].(float64))
	limit := int(response["limit"].(float64))
	offset := int(response["offset"].(float64))

	assert.Equal(t, wantTotal, total, "unexpected total")
	assert.Equal(t, wantLimit, limit, "unexpected limit")
	assert.Equal(t, wantOffset, offset, "unexpected offset")

	// Ensure items count doesn't exceed limit
	assert.LessOrEqual(t, len(items), limit, "items count should not exceed limit")

	return response
}

// Test data creation helpers

// createTestUser creates a test user and returns its ID
func createTestUser(t *testing.T, server *TestServer, email, name string) string {
	t.Helper()

	user := map[string]interface{}{
		"email": email,
		"name":  name,
		"address": map[string]interface{}{
			"street":     "123 Test St",
			"city":       "Test City",
			"state":      "TS",
			"postalCode": "12345",
			"country":    "Test Country",
		},
	}

	rr := makeRequest(t, server, "POST", "/users", user)
	require.Equal(t, http.StatusCreated, rr.Code, "failed to create test user")

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	id, ok := response["id"].(string)
	require.True(t, ok, "response should contain id as string")

	return id
}

// createTestProduct creates a test product and returns its ID
func createTestProduct(t *testing.T, server *TestServer, name string, price float64, stock int) string {
	t.Helper()

	product := map[string]interface{}{
		"name":        name,
		"description": "Test product description",
		"price":       price,
		"stock":       stock,
		"categoryId":  "1", // Default category
		"imageUrls":   []string{"https://example.com/image.jpg"},
	}

	rr := makeRequest(t, server, "POST", "/products", product)
	require.Equal(t, http.StatusCreated, rr.Code, "failed to create test product")

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	id, ok := response["id"].(string)
	require.True(t, ok, "response should contain id as string")

	return id
}

// createTestCategory creates a test category and returns its ID
func createTestCategory(t *testing.T, server *TestServer, name string, parentID *string) string {
	t.Helper()

	category := map[string]interface{}{
		"name":     name,
		"parentId": parentID,
	}

	rr := makeRequest(t, server, "POST", "/categories", category)
	require.Equal(t, http.StatusCreated, rr.Code, "failed to create test category")

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	id, ok := response["id"].(string)
	require.True(t, ok, "response should contain id as string")

	return id
}

// addToCart adds an item to user's cart
func addToCart(t *testing.T, server *TestServer, userID, productID string, quantity int) {
	t.Helper()

	item := map[string]interface{}{
		"productId": productID,
		"quantity":  quantity,
	}

	rr := makeRequest(t, server, "POST", "/carts/users/"+userID+"/items", item)
	require.Equal(t, http.StatusOK, rr.Code, "failed to add item to cart")
}

// createOrder creates an order for a user
func createOrder(t *testing.T, server *TestServer, userID string) string {
	t.Helper()

	// Get cart items first
	cartRR := makeRequest(t, server, "GET", "/carts/users/"+userID, nil)
	require.Equal(t, http.StatusOK, cartRR.Code, "failed to get cart")

	var cart map[string]interface{}
	err := json.NewDecoder(cartRR.Body).Decode(&cart)
	require.NoError(t, err)

	cartItems := cart["items"].([]interface{})
	orderItems := make([]map[string]interface{}, len(cartItems))
	for i, item := range cartItems {
		cartItem := item.(map[string]interface{})
		orderItems[i] = map[string]interface{}{
			"productId": cartItem["productId"],
			"quantity":  cartItem["quantity"],
		}
	}

	order := map[string]interface{}{
		"items": orderItems,
		"shippingAddress": map[string]interface{}{
			"street":     "456 Order Ave",
			"city":       "Order City",
			"state":      "OC",
			"postalCode": "54321",
			"country":    "Order Country",
		},
	}

	rr := makeRequest(t, server, "POST", "/orders/users/"+userID, order)
	require.Equal(t, http.StatusCreated, rr.Code, "failed to create order")

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	id, ok := response["id"].(string)
	require.True(t, ok, "response should contain id as string")

	return id
}

// assertJSONEqual compares two JSON responses
func assertJSONEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()

	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)

	actualJSON, err := json.Marshal(actual)
	require.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}