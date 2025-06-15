package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCartsService_Get(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should return cart for existing user", func(t *testing.T) {
		rr := makeAuthenticatedRequest(t, server, "GET", "/carts/users/1", nil, token)
		assertStatus(t, rr, http.StatusOK)

		var cart map[string]any
		err := decodeJSON(rr, &cart)
		require.NoError(t, err)

		assert.Equal(t, "1", cart["userId"])
		assert.NotNil(t, cart["items"])
		items := cart["items"].([]any)
		assert.Len(t, items, 0) // Initially empty
		assert.NotEmpty(t, cart["id"])
		assert.NotEmpty(t, cart["createdAt"])
		assert.NotEmpty(t, cart["updatedAt"])
	})

	t.Run("should create empty cart for new user", func(t *testing.T) {
		rr := makeAuthenticatedRequest(t, server, "GET", "/carts/users/999", nil, token)
		assertStatus(t, rr, http.StatusOK)

		var cart map[string]any
		err := decodeJSON(rr, &cart)
		require.NoError(t, err)

		assert.Equal(t, "999", cart["userId"])
		items := cart["items"].([]any)
		assert.Len(t, items, 0)
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/carts/users/1", nil)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestCartsService_AddItem(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should add item to cart", func(t *testing.T) {
		addItem := map[string]any{
			"productId": "1",
			"quantity":  2,
		}

		rr := makeAuthenticatedRequest(t, server, "POST", "/carts/users/1/items", addItem, token)
		assertStatus(t, rr, http.StatusOK)

		var cart map[string]any
		err := decodeJSON(rr, &cart)
		require.NoError(t, err)

		items := cart["items"].([]any)
		assert.Len(t, items, 1)

		item := items[0].(map[string]any)
		assert.Equal(t, "1", item["productId"])
		assert.Equal(t, float64(2), item["quantity"])
	})

	t.Run("should increase quantity when adding existing item", func(t *testing.T) {
		// First, add item
		addItem := map[string]any{
			"productId": "1",
			"quantity":  2,
		}
		makeAuthenticatedRequest(t, server, "POST", "/carts/users/2/items", addItem, token)

		// Add same item again
		rr := makeAuthenticatedRequest(t, server, "POST", "/carts/users/2/items", addItem, token)
		assertStatus(t, rr, http.StatusOK)

		var cart map[string]any
		err := decodeJSON(rr, &cart)
		require.NoError(t, err)

		items := cart["items"].([]any)
		assert.Len(t, items, 1)

		item := items[0].(map[string]any)
		assert.Equal(t, float64(4), item["quantity"]) // 2 + 2
	})

	t.Run("should return 404 for non-existent product", func(t *testing.T) {
		addItem := map[string]any{
			"productId": "999",
			"quantity":  1,
		}

		rr := makeAuthenticatedRequest(t, server, "POST", "/carts/users/1/items", addItem, token)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})

	t.Run("should return 400 for insufficient stock", func(t *testing.T) {
		addItem := map[string]any{
			"productId": "1",
			"quantity":  100, // More than available stock (10)
		}

		rr := makeAuthenticatedRequest(t, server, "POST", "/carts/users/3/items", addItem, token)
		assertStatus(t, rr, http.StatusBadRequest)
		assertErrorResponse(t, rr, "INSUFFICIENT_STOCK")
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		addItem := map[string]any{
			"productId": "1",
			"quantity":  1,
		}

		rr := makeRequest(t, server, "POST", "/carts/users/1/items", addItem)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestCartsService_UpdateItem(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	// Setup: Add item to cart first
	setupCart := func(userID string) {
		addItem := map[string]any{
			"productId": "1",
			"quantity":  2,
		}
		makeAuthenticatedRequest(t, server, "POST", "/carts/users/"+userID+"/items", addItem, token)
	}

	t.Run("should update item quantity", func(t *testing.T) {
		setupCart("4")

		update := map[string]any{
			"quantity": 5,
		}

		rr := makeAuthenticatedRequest(t, server, "PATCH", "/carts/users/4/items/1", update, token)
		assertStatus(t, rr, http.StatusOK)

		var cart map[string]any
		err := decodeJSON(rr, &cart)
		require.NoError(t, err)

		items := cart["items"].([]any)
		item := items[0].(map[string]any)
		assert.Equal(t, float64(5), item["quantity"])
	})

	t.Run("should return 404 for item not in cart", func(t *testing.T) {
		update := map[string]any{
			"quantity": 5,
		}

		rr := makeAuthenticatedRequest(t, server, "PATCH", "/carts/users/5/items/999", update, token)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})

	t.Run("should return 400 for insufficient stock", func(t *testing.T) {
		setupCart("6")

		update := map[string]any{
			"quantity": 100,
		}

		rr := makeAuthenticatedRequest(t, server, "PATCH", "/carts/users/6/items/1", update, token)
		assertStatus(t, rr, http.StatusBadRequest)
		assertErrorResponse(t, rr, "INSUFFICIENT_STOCK")
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		update := map[string]any{
			"quantity": 5,
		}

		rr := makeRequest(t, server, "PATCH", "/carts/users/1/items/1", update)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestCartsService_RemoveItem(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	// Setup: Add items to cart
	setupCart := func(userID string) {
		// Add two different products
		makeAuthenticatedRequest(t, server, "POST", "/carts/users/"+userID+"/items", map[string]any{
			"productId": "1",
			"quantity":  2,
		}, token)
		makeAuthenticatedRequest(t, server, "POST", "/carts/users/"+userID+"/items", map[string]any{
			"productId": "2",
			"quantity":  1,
		}, token)
	}

	t.Run("should remove item from cart", func(t *testing.T) {
		setupCart("7")

		rr := makeAuthenticatedRequest(t, server, "DELETE", "/carts/users/7/items/1", nil, token)
		assertStatus(t, rr, http.StatusNoContent)

		// Verify item was removed
		cartRR := makeAuthenticatedRequest(t, server, "GET", "/carts/users/7", nil, token)
		var cart map[string]any
		err := decodeJSON(cartRR, &cart)
		require.NoError(t, err)

		items := cart["items"].([]any)
		assert.Len(t, items, 1)
		remainingItem := items[0].(map[string]any)
		assert.Equal(t, "2", remainingItem["productId"])
	})

	t.Run("should return 404 for item not in cart", func(t *testing.T) {
		rr := makeAuthenticatedRequest(t, server, "DELETE", "/carts/users/8/items/999", nil, token)
		assertStatus(t, rr, http.StatusNotFound)
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		rr := makeRequest(t, server, "DELETE", "/carts/users/1/items/1", nil)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestCartsService_Clear(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	// Setup: Add items to cart
	setupCart := func(userID string) {
		makeAuthenticatedRequest(t, server, "POST", "/carts/users/"+userID+"/items", map[string]any{
			"productId": "1",
			"quantity":  2,
		}, token)
		makeAuthenticatedRequest(t, server, "POST", "/carts/users/"+userID+"/items", map[string]any{
			"productId": "2",
			"quantity":  1,
		}, token)
	}

	t.Run("should clear all items from cart", func(t *testing.T) {
		setupCart("9")

		rr := makeAuthenticatedRequest(t, server, "DELETE", "/carts/users/9/items", nil, token)
		assertStatus(t, rr, http.StatusNoContent)

		// Verify cart is empty
		cartRR := makeAuthenticatedRequest(t, server, "GET", "/carts/users/9", nil, token)
		var cart map[string]any
		err := decodeJSON(cartRR, &cart)
		require.NoError(t, err)

		items := cart["items"].([]any)
		assert.Len(t, items, 0)
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		rr := makeRequest(t, server, "DELETE", "/carts/users/1/items", nil)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestCartsService_Integration(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should handle complete cart workflow", func(t *testing.T) {
		userID := "100"

		// 1. Start with empty cart
		getEmptyCart := makeAuthenticatedRequest(t, server, "GET", "/carts/users/"+userID, nil, token)
		assertStatus(t, getEmptyCart, http.StatusOK)
		var cart map[string]any
		decodeJSON(getEmptyCart, &cart)
		items := cart["items"].([]any)
		assert.Len(t, items, 0)

		// 2. Add items
		addItem1 := makeAuthenticatedRequest(t, server, "POST", "/carts/users/"+userID+"/items", map[string]any{
			"productId": "1",
			"quantity":  2,
		}, token)
		assertStatus(t, addItem1, http.StatusOK)

		addItem2 := makeAuthenticatedRequest(t, server, "POST", "/carts/users/"+userID+"/items", map[string]any{
			"productId": "2",
			"quantity":  1,
		}, token)
		assertStatus(t, addItem2, http.StatusOK)

		// 3. Update quantity
		updateItem := makeAuthenticatedRequest(t, server, "PATCH", "/carts/users/"+userID+"/items/1", map[string]any{
			"quantity": 3,
		}, token)
		assertStatus(t, updateItem, http.StatusOK)

		// 4. Remove one item
		removeItem := makeAuthenticatedRequest(t, server, "DELETE", "/carts/users/"+userID+"/items/2", nil, token)
		assertStatus(t, removeItem, http.StatusNoContent)

		// 5. Verify final state
		getFinalCart := makeAuthenticatedRequest(t, server, "GET", "/carts/users/"+userID, nil, token)
		assertStatus(t, getFinalCart, http.StatusOK)
		decodeJSON(getFinalCart, &cart)
		finalItems := cart["items"].([]any)
		assert.Len(t, finalItems, 1)
		finalItem := finalItems[0].(map[string]any)
		assert.Equal(t, "1", finalItem["productId"])
		assert.Equal(t, float64(3), finalItem["quantity"])

		// 6. Clear cart
		clearCart := makeAuthenticatedRequest(t, server, "DELETE", "/carts/users/"+userID+"/items", nil, token)
		assertStatus(t, clearCart, http.StatusNoContent)

		// 7. Verify empty
		getCleared := makeAuthenticatedRequest(t, server, "GET", "/carts/users/"+userID, nil, token)
		assertStatus(t, getCleared, http.StatusOK)
		decodeJSON(getCleared, &cart)
		clearedItems := cart["items"].([]any)
		assert.Len(t, clearedItems, 0)
	})
}