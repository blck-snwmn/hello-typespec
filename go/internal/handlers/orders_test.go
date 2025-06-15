package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrdersService_List(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should return all orders with default pagination", func(t *testing.T) {
		// Create some test orders first
		userID := createTestUser(t, server, "orders@example.com", "Order User")
		productID := createTestProduct(t, server, "Order Product", 50.00, 10)

		// Add to cart with auth
		addToCartAuth(t, server, userID, productID, 2, token)
		orderID := createOrderAuth(t, server, userID, token)

		rr := makeAuthenticatedRequest(t, server, "GET", "/orders", nil, token)
		assertStatus(t, rr, http.StatusOK)
		response := assertPaginatedResponse(t, rr, 1, 20, 0)

		items := response["items"].([]any)
		assert.Len(t, items, 1)

		// Check order structure
		order := items[0].(map[string]any)
		assert.Contains(t, order, "id")
		assert.Contains(t, order, "userId")
		assert.Contains(t, order, "items")
		assert.Contains(t, order, "totalAmount")
		assert.Contains(t, order, "status")
		assert.Contains(t, order, "createdAt")
		assert.Contains(t, order, "updatedAt")
		assert.Equal(t, orderID, order["id"])
	})

	t.Run("should filter by user", func(t *testing.T) {
		// Create orders for different users
		user1 := createTestUser(t, server, "user1@example.com", "User 1")
		user2 := createTestUser(t, server, "user2@example.com", "User 2")
		productID := createTestProduct(t, server, "Filter Product", 25.00, 20)

		addToCartAuth(t, server, user1, productID, 1, token)
		createOrderAuth(t, server, user1, token)

		addToCartAuth(t, server, user2, productID, 2, token)
		createOrderAuth(t, server, user2, token)

		// Filter by user1
		rr := makeAuthenticatedRequest(t, server, "GET", "/orders?userId="+user1, nil, token)
		assertStatus(t, rr, http.StatusOK)
		response := assertPaginatedResponse(t, rr, 1, 20, 0)

		items := response["items"].([]any)
		assert.Len(t, items, 1)
		assert.Equal(t, user1, items[0].(map[string]any)["userId"])
	})

	t.Run("should filter by status", func(t *testing.T) {
		// Create orders with different statuses
		userID := createTestUser(t, server, "status@example.com", "Status User")
		productID := createTestProduct(t, server, "Status Product", 30.00, 15)

		// Create multiple orders
		for i := 0; i < 3; i++ {
			addToCartAuth(t, server, userID, productID, 1, token)
			orderID := createOrderAuth(t, server, userID, token)

			// Update some to different status
			if i > 0 {
				updateOrderStatus(t, server, orderID, "processing", token)
			}
		}

		// Filter by pending status
		rr := makeAuthenticatedRequest(t, server, "GET", "/orders?status=pending", nil, token)
		assertStatus(t, rr, http.StatusOK)

		var response map[string]any
		err := decodeJSON(rr, &response)
		require.NoError(t, err)

		items := response["items"].([]any)
		// Should have at least 1 pending order (may have more from previous tests)
		assert.GreaterOrEqual(t, len(items), 1)
		// All returned orders should have pending status
		for _, item := range items {
			assert.Equal(t, "pending", item.(map[string]any)["status"])
		}
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/orders", nil)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestOrdersService_ListByUser(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should return orders for specific user", func(t *testing.T) {
		userID := createTestUser(t, server, "userorders@example.com", "User Orders")
		productID := createTestProduct(t, server, "User Order Product", 40.00, 10)

		// Create 2 orders
		for i := 0; i < 2; i++ {
			addToCartAuth(t, server, userID, productID, 1, token)
			createOrderAuth(t, server, userID, token)
		}

		rr := makeAuthenticatedRequest(t, server, "GET", "/orders/users/"+userID, nil, token)
		assertStatus(t, rr, http.StatusOK)
		response := assertPaginatedResponse(t, rr, 2, 20, 0)

		items := response["items"].([]any)
		assert.Len(t, items, 2)
		for _, item := range items {
			assert.Equal(t, userID, item.(map[string]any)["userId"])
		}
	})

	t.Run("should return empty list for user with no orders", func(t *testing.T) {
		userID := createTestUser(t, server, "noorders@example.com", "No Orders")

		rr := makeAuthenticatedRequest(t, server, "GET", "/orders/users/"+userID, nil, token)
		assertStatus(t, rr, http.StatusOK)
		response := assertPaginatedResponse(t, rr, 0, 20, 0)

		items := response["items"].([]any)
		assert.Len(t, items, 0)
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/orders/users/1", nil)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestOrdersService_Get(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should return an order by id", func(t *testing.T) {
		userID := createTestUser(t, server, "getorder@example.com", "Get Order")
		productID := createTestProduct(t, server, "Get Order Product", 75.00, 5)
		addToCartAuth(t, server, userID, productID, 2, token)
		orderID := createOrderAuth(t, server, userID, token)

		rr := makeAuthenticatedRequest(t, server, "GET", "/orders/"+orderID, nil, token)
		assertStatus(t, rr, http.StatusOK)

		var order map[string]any
		err := decodeJSON(rr, &order)
		require.NoError(t, err)

		assert.Equal(t, orderID, order["id"])
		assert.Equal(t, userID, order["userId"])
		assert.Equal(t, "pending", order["status"])
		assert.Equal(t, float64(150), order["totalAmount"]) // 75 * 2
		assert.NotNil(t, order["items"])
		assert.NotNil(t, order["shippingAddress"])
	})

	t.Run("should return 404 for non-existent order", func(t *testing.T) {
		rr := makeAuthenticatedRequest(t, server, "GET", "/orders/999", nil, token)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/orders/1", nil)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestOrdersService_Create(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should create order from cart", func(t *testing.T) {
		userID := createTestUser(t, server, "createorder@example.com", "Create Order")
		product1 := createTestProduct(t, server, "Product 1", 20.00, 10)
		product2 := createTestProduct(t, server, "Product 2", 30.00, 15)

		// Add items to cart
		addToCartAuth(t, server, userID, product1, 2, token)
		addToCartAuth(t, server, userID, product2, 1, token)

		// Create order with explicit items
		orderRequest := map[string]any{
			"items": []any{
				map[string]any{
					"productId": product1,
					"quantity":  2,
				},
				map[string]any{
					"productId": product2,
					"quantity":  1,
				},
			},
			"shippingAddress": map[string]any{
				"street":     "123 Order St",
				"city":       "Order City",
				"state":      "OC",
				"postalCode": "12345",
				"country":    "USA",
			},
		}

		rr := makeAuthenticatedRequest(t, server, "POST", "/orders/users/"+userID, orderRequest, token)
		assertStatus(t, rr, http.StatusCreated)

		var order map[string]any
		err := decodeJSON(rr, &order)
		require.NoError(t, err)

		assert.NotEmpty(t, order["id"])
		assert.Equal(t, userID, order["userId"])
		assert.Equal(t, "pending", order["status"])
		assert.Equal(t, float64(70), order["totalAmount"]) // (20*2) + (30*1)

		// Check items
		items := order["items"].([]any)
		assert.Len(t, items, 2)

		// Verify cart was cleared
		cartRR := makeAuthenticatedRequest(t, server, "GET", "/carts/users/"+userID, nil, token)
		var cart map[string]any
		decodeJSON(cartRR, &cart)
		cartItems := cart["items"].([]any)
		assert.Len(t, cartItems, 0)

		// Verify stock was reduced
		productRR := makeRequest(t, server, "GET", "/products/"+product1, nil)
		var product map[string]any
		decodeJSON(productRR, &product)
		assert.Equal(t, float64(8), product["stock"]) // 10 - 2
	})

	t.Run("should fail if insufficient stock", func(t *testing.T) {
		userID := createTestUser(t, server, "nostock@example.com", "No Stock")
		productID := createTestProduct(t, server, "Limited Product", 100.00, 2)

		// Try to add more than available stock - this should fail
		rr := makeAuthenticatedRequest(t, server, "POST", "/carts/users/"+userID+"/items", map[string]any{
			"productId": productID,
			"quantity":  3,
		}, token)
		assertStatus(t, rr, http.StatusBadRequest)

		// Now try to create order with items that exceed stock
		orderRequest := map[string]any{
			"items": []any{
				map[string]any{
					"productId": productID,
					"quantity":  3,
				},
			},
			"shippingAddress": map[string]any{
				"street":     "123 Test",
				"city":       "Test",
				"state":      "TS",
				"postalCode": "12345",
				"country":    "USA",
			},
		}

		orderRR := makeAuthenticatedRequest(t, server, "POST", "/orders/users/"+userID, orderRequest, token)
		// Order creation should fail due to insufficient stock
		assertStatus(t, orderRR, http.StatusBadRequest)
		assertErrorResponse(t, orderRR, "INSUFFICIENT_STOCK")
	})

	t.Run("should fail if product not found", func(t *testing.T) {
		userID := createTestUser(t, server, "notfound@example.com", "Not Found")

		// Add non-existent product to cart (this should fail at cart level)
		rr := makeAuthenticatedRequest(t, server, "POST", "/carts/users/"+userID+"/items", map[string]any{
			"productId": "999",
			"quantity":  1,
		}, token)
		assertStatus(t, rr, http.StatusNotFound)
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		orderRequest := map[string]any{
			"items": []any{},
			"shippingAddress": map[string]any{
				"street":     "123 Test",
				"city":       "Test",
				"state":      "TS",
				"postalCode": "12345",
				"country":    "USA",
			},
		}

		rr := makeRequest(t, server, "POST", "/orders/users/1", orderRequest)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestOrdersService_UpdateStatus(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should update order status", func(t *testing.T) {
		userID := createTestUser(t, server, "updatestatus@example.com", "Update Status")
		productID := createTestProduct(t, server, "Status Product", 50.00, 10)
		addToCartAuth(t, server, userID, productID, 1, token)
		orderID := createOrderAuth(t, server, userID, token)

		// Update status to processing
		statusUpdate := map[string]any{
			"status": "processing",
		}

		rr := makeAuthenticatedRequest(t, server, "PATCH", "/orders/status/"+orderID, statusUpdate, token)
		assertStatus(t, rr, http.StatusOK)

		var order map[string]any
		err := decodeJSON(rr, &order)
		require.NoError(t, err)

		assert.Equal(t, "processing", order["status"])
	})

	t.Run("should return 404 for non-existent order", func(t *testing.T) {
		statusUpdate := map[string]any{
			"status": "processing",
		}

		rr := makeAuthenticatedRequest(t, server, "PATCH", "/orders/status/999", statusUpdate, token)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		statusUpdate := map[string]any{
			"status": "processing",
		}

		rr := makeRequest(t, server, "PATCH", "/orders/status/1", statusUpdate)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestOrdersService_Cancel(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should cancel pending order and restore inventory", func(t *testing.T) {
		userID := createTestUser(t, server, "cancel@example.com", "Cancel Order")
		productID := createTestProduct(t, server, "Cancel Product", 60.00, 10)
		addToCartAuth(t, server, userID, productID, 3, token)
		orderID := createOrderAuth(t, server, userID, token)

		// Cancel the order
		rr := makeAuthenticatedRequest(t, server, "POST", "/orders/cancel/"+orderID, nil, token)
		assertStatus(t, rr, http.StatusOK)

		var order map[string]any
		err := decodeJSON(rr, &order)
		require.NoError(t, err)

		assert.Equal(t, "cancelled", order["status"])

		// Verify stock was restored
		productRR := makeRequest(t, server, "GET", "/products/"+productID, nil)
		var product map[string]any
		decodeJSON(productRR, &product)
		assert.Equal(t, float64(10), product["stock"]) // Back to original
	})

	t.Run("should not cancel non-pending order", func(t *testing.T) {
		userID := createTestUser(t, server, "nocancel@example.com", "No Cancel")
		productID := createTestProduct(t, server, "No Cancel Product", 70.00, 10)
		addToCartAuth(t, server, userID, productID, 1, token)
		orderID := createOrderAuth(t, server, userID, token)

		// Update to shipped first (shipped orders cannot be cancelled)
		updateOrderStatus(t, server, orderID, "processing", token)
		updateOrderStatus(t, server, orderID, "shipped", token)

		// Try to cancel
		rr := makeAuthenticatedRequest(t, server, "POST", "/orders/cancel/"+orderID, nil, token)
		assertStatus(t, rr, http.StatusBadRequest)
		assertErrorResponse(t, rr, "VALIDATION_ERROR")
	})

	t.Run("should return 404 for non-existent order", func(t *testing.T) {
		rr := makeAuthenticatedRequest(t, server, "POST", "/orders/cancel/999", nil, token)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		rr := makeRequest(t, server, "POST", "/orders/cancel/1", nil)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestOrdersService_Integration(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should handle complete order lifecycle", func(t *testing.T) {
		// Create user and products
		userID := createTestUser(t, server, "lifecycle@example.com", "Lifecycle User")
		product1 := createTestProduct(t, server, "Lifecycle Product 1", 25.00, 20)
		product2 := createTestProduct(t, server, "Lifecycle Product 2", 35.00, 15)

		// Add to cart
		addToCartAuth(t, server, userID, product1, 2, token)
		addToCartAuth(t, server, userID, product2, 1, token)

		// Create order with explicit items
		orderRequest := map[string]any{
			"items": []any{
				map[string]any{
					"productId": product1,
					"quantity":  2,
				},
				map[string]any{
					"productId": product2,
					"quantity":  1,
				},
			},
			"shippingAddress": map[string]any{
				"street":     "456 Lifecycle Ave",
				"city":       "Lifecycle City",
				"state":      "LC",
				"postalCode": "54321",
				"country":    "USA",
			},
		}

		createRR := makeAuthenticatedRequest(t, server, "POST", "/orders/users/"+userID, orderRequest, token)
		assertStatus(t, createRR, http.StatusCreated)

		var order map[string]any
		err := decodeJSON(createRR, &order)
		require.NoError(t, err)
		orderID := order["id"].(string)

		// Verify order details
		assert.Equal(t, float64(85), order["totalAmount"]) // (25*2) + (35*1)
		assert.Equal(t, "pending", order["status"])

		// Update status to processing
		updateRR := makeAuthenticatedRequest(t, server, "PATCH", "/orders/status/"+orderID, map[string]any{
			"status": "processing",
		}, token)
		assertStatus(t, updateRR, http.StatusOK)

		// Update to shipped
		updateRR2 := makeAuthenticatedRequest(t, server, "PATCH", "/orders/status/"+orderID, map[string]any{
			"status": "shipped",
		}, token)
		assertStatus(t, updateRR2, http.StatusOK)

		// Verify final state
		getRR := makeAuthenticatedRequest(t, server, "GET", "/orders/"+orderID, nil, token)
		assertStatus(t, getRR, http.StatusOK)
		decodeJSON(getRR, &order)
		assert.Equal(t, "shipped", order["status"])

		// List user's orders
		listRR := makeAuthenticatedRequest(t, server, "GET", "/orders/users/"+userID, nil, token)
		assertStatus(t, listRR, http.StatusOK)
		response := assertPaginatedResponse(t, listRR, 1, 20, 0)
		items := response["items"].([]any)
		assert.Len(t, items, 1)
	})
}

// Helper functions specific to orders that need auth
func addToCartAuth(t testing.TB, server *TestServer, userID, productID string, quantity int, token string) {
	t.Helper()

	item := map[string]any{
		"productId": productID,
		"quantity":  quantity,
	}

	rr := makeAuthenticatedRequest(t, server, "POST", "/carts/users/"+userID+"/items", item, token)
	require.Equal(t, http.StatusOK, rr.Code, "failed to add item to cart")
}

func createOrderAuth(t testing.TB, server *TestServer, userID string, token string) string {
	t.Helper()

	// Get cart items first
	cartRR := makeAuthenticatedRequest(t, server, "GET", "/carts/users/"+userID, nil, token)
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

	rr := makeAuthenticatedRequest(t, server, "POST", "/orders/users/"+userID, order, token)
	require.Equal(t, http.StatusCreated, rr.Code, "failed to create order")

	var response map[string]any
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	id, ok := response["id"].(string)
	require.True(t, ok, "response should contain id as string")

	return id
}

func updateOrderStatus(t testing.TB, server *TestServer, orderID string, status string, token string) {
	t.Helper()

	statusUpdate := map[string]any{
		"status": status,
	}

	rr := makeAuthenticatedRequest(t, server, "PATCH", "/orders/status/"+orderID, statusUpdate, token)
	require.Equal(t, http.StatusOK, rr.Code, "failed to update order status")
}
