package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCartsService_Get(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should return cart for existing user", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/carts/users/1", nil)
		assertStatus(t, rr, http.StatusOK)

		var cart map[string]interface{}
		err := decodeJSON(rr, &cart)
		require.NoError(t, err)

		assert.Equal(t, "1", cart["userId"])
		assert.NotNil(t, cart["items"])
		items := cart["items"].([]interface{})
		assert.Len(t, items, 0) // Initially empty
		assert.NotEmpty(t, cart["id"])
		assert.NotEmpty(t, cart["createdAt"])
		assert.NotEmpty(t, cart["updatedAt"])
	})

	t.Run("should create empty cart for new user", func(t *testing.T) {
		// Request cart for non-existent user
		rr := makeRequest(t, server, "GET", "/carts/users/999", nil)
		assertStatus(t, rr, http.StatusOK)

		var cart map[string]interface{}
		err := decodeJSON(rr, &cart)
		require.NoError(t, err)

		assert.Equal(t, "999", cart["userId"])
		items := cart["items"].([]interface{})
		assert.Len(t, items, 0)
	})
}

func TestCartsService_AddItem(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should add item to cart", func(t *testing.T) {
		// Create a test product first
		productID := createTestProduct(t, server, "Test Product", 99.99, 10)

		// Add item to cart
		item := map[string]interface{}{
			"productId": productID,
			"quantity":  2,
		}

		rr := makeRequest(t, server, "POST", "/carts/users/1/items", item)
		assertStatus(t, rr, http.StatusOK)

		var cart map[string]interface{}
		err := decodeJSON(rr, &cart)
		require.NoError(t, err)

		items := cart["items"].([]interface{})
		assert.Len(t, items, 1)

		cartItem := items[0].(map[string]interface{})
		assert.Equal(t, productID, cartItem["productId"])
		assert.Equal(t, float64(2), cartItem["quantity"])
		// TODO: Go implementation doesn't include product details in cart items
	})

	t.Run("should increase quantity when adding existing item", func(t *testing.T) {
		// Create a test product
		productID := createTestProduct(t, server, "Test Product 2", 49.99, 20)

		// Add item first time
		item := map[string]interface{}{
			"productId": productID,
			"quantity":  1,
		}
		makeRequest(t, server, "POST", "/carts/users/2/items", item)

		// Add same item again
		item["quantity"] = 2
		rr := makeRequest(t, server, "POST", "/carts/users/2/items", item)
		assertStatus(t, rr, http.StatusOK)

		var cart map[string]interface{}
		err := decodeJSON(rr, &cart)
		require.NoError(t, err)

		items := cart["items"].([]interface{})
		assert.Len(t, items, 1)

		cartItem := items[0].(map[string]interface{})
		assert.Equal(t, float64(3), cartItem["quantity"]) // 1 + 2 = 3
	})

	t.Run("should return 404 for non-existent product", func(t *testing.T) {
		item := map[string]interface{}{
			"productId": "999999",
			"quantity":  1,
		}

		rr := makeRequest(t, server, "POST", "/carts/users/1/items", item)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})

	t.Run("should return 400 for insufficient stock", func(t *testing.T) {
		// Create product with limited stock
		productID := createTestProduct(t, server, "Limited Product", 19.99, 2)

		item := map[string]interface{}{
			"productId": productID,
			"quantity":  5, // More than available stock
		}

		rr := makeRequest(t, server, "POST", "/carts/users/1/items", item)
		assertStatus(t, rr, http.StatusBadRequest)
		assertErrorResponse(t, rr, "INSUFFICIENT_STOCK")
	})
}

func TestCartsService_UpdateItem(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should update item quantity", func(t *testing.T) {
		// Setup: Create product and add to cart
		productID := createTestProduct(t, server, "Update Product", 29.99, 15)
		addToCart(t, server, "3", productID, 2)

		// Update quantity
		update := map[string]interface{}{
			"quantity": 5,
		}

		rr := makeRequest(t, server, "PATCH", "/carts/users/3/items/"+productID, update)
		assertStatus(t, rr, http.StatusOK)

		var cart map[string]interface{}
		err := decodeJSON(rr, &cart)
		require.NoError(t, err)

		items := cart["items"].([]interface{})
		cartItem := items[0].(map[string]interface{})
		assert.Equal(t, float64(5), cartItem["quantity"])
	})

	t.Run("should return 404 for item not in cart", func(t *testing.T) {
		update := map[string]interface{}{
			"quantity": 1,
		}

		rr := makeRequest(t, server, "PATCH", "/carts/users/1/items/nonexistent", update)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})

	t.Run("should return 400 for insufficient stock", func(t *testing.T) {
		// Setup: Create product with limited stock and add to cart
		productID := createTestProduct(t, server, "Limited Update Product", 9.99, 3)
		addToCart(t, server, "4", productID, 1)

		// Try to update to quantity exceeding stock
		update := map[string]interface{}{
			"quantity": 10,
		}

		rr := makeRequest(t, server, "PATCH", "/carts/users/4/items/"+productID, update)
		assertStatus(t, rr, http.StatusBadRequest)
		assertErrorResponse(t, rr, "INSUFFICIENT_STOCK")
	})
}

func TestCartsService_RemoveItem(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should remove item from cart", func(t *testing.T) {
		// Setup: Create products and add to cart
		productID1 := createTestProduct(t, server, "Remove Product 1", 10.00, 5)
		productID2 := createTestProduct(t, server, "Remove Product 2", 20.00, 5)
		addToCart(t, server, "5", productID1, 1)
		addToCart(t, server, "5", productID2, 2)

		// Remove first product
		rr := makeRequest(t, server, "DELETE", "/carts/users/5/items/"+productID1, nil)
		assertStatus(t, rr, http.StatusNoContent)

		// Verify cart contents
		cartRR := makeRequest(t, server, "GET", "/carts/users/5", nil)
		var cart map[string]interface{}
		err := decodeJSON(cartRR, &cart)
		require.NoError(t, err)

		items := cart["items"].([]interface{})
		assert.Len(t, items, 1) // Only one item left

		remainingItem := items[0].(map[string]interface{})
		assert.Equal(t, productID2, remainingItem["productId"])
	})

	t.Run("should return 404 for item not in cart", func(t *testing.T) {
		rr := makeRequest(t, server, "DELETE", "/carts/users/1/items/nonexistent", nil)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})
}

func TestCartsService_Clear(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should clear all items from cart", func(t *testing.T) {
		// Setup: Add multiple items to cart
		productID1 := createTestProduct(t, server, "Clear Product 1", 15.00, 10)
		productID2 := createTestProduct(t, server, "Clear Product 2", 25.00, 10)
		addToCart(t, server, "6", productID1, 3)
		addToCart(t, server, "6", productID2, 2)

		// Clear cart
		rr := makeRequest(t, server, "DELETE", "/carts/users/6/items", nil)
		assertStatus(t, rr, http.StatusNoContent)

		// Verify cart is empty
		cartRR := makeRequest(t, server, "GET", "/carts/users/6", nil)
		var cart map[string]interface{}
		err := decodeJSON(cartRR, &cart)
		require.NoError(t, err)

		items := cart["items"].([]interface{})
		assert.Len(t, items, 0)
	})
}

func TestCartsService_Integration(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should handle complete cart workflow", func(t *testing.T) {
		userID := "integration-user"
		
		// Create products
		productID1 := createTestProduct(t, server, "Integration Product 1", 50.00, 20)
		productID2 := createTestProduct(t, server, "Integration Product 2", 75.00, 15)

		// Add first product
		addToCart(t, server, userID, productID1, 2)

		// Add second product
		addToCart(t, server, userID, productID2, 1)

		// Update first product quantity
		update := map[string]interface{}{
			"quantity": 3,
		}
		updateRR := makeRequest(t, server, "PATCH", "/carts/users/"+userID+"/items/"+productID1, update)
		assertStatus(t, updateRR, http.StatusOK)

		// Remove second product
		deleteRR := makeRequest(t, server, "DELETE", "/carts/users/"+userID+"/items/"+productID2, nil)
		assertStatus(t, deleteRR, http.StatusNoContent)

		// Verify final cart state
		cartRR := makeRequest(t, server, "GET", "/carts/users/"+userID, nil)
		var cart map[string]interface{}
		err := decodeJSON(cartRR, &cart)
		require.NoError(t, err)

		items := cart["items"].([]interface{})
		assert.Len(t, items, 1)

		item := items[0].(map[string]interface{})
		assert.Equal(t, productID1, item["productId"])
		assert.Equal(t, float64(3), item["quantity"])
	})
}