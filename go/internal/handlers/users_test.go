package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersService_List(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should return all users with default pagination", func(t *testing.T) {
		rr := makeAuthenticatedRequest(t, server, "GET", "/users", nil, token)

		assertStatus(t, rr, http.StatusOK)
		response := assertPaginatedResponse(t, rr, 2, 20, 0) // 2 default users

		items := response["items"].([]any)
		assert.Len(t, items, 2)

		// Check first user structure
		user := items[0].(map[string]any)
		assert.Contains(t, user, "id")
		assert.Contains(t, user, "email")
		assert.Contains(t, user, "name")
		assert.Contains(t, user, "address")
		assert.Contains(t, user, "createdAt")
		assert.Contains(t, user, "updatedAt")
	})

	t.Run("should support pagination with limit and offset", func(t *testing.T) {
		// First page
		rr1 := makeAuthenticatedRequest(t, server, "GET", "/users?limit=1&offset=0", nil, token)
		assertStatus(t, rr1, http.StatusOK)
		response1 := assertPaginatedResponse(t, rr1, 2, 1, 0)
		items1 := response1["items"].([]any)
		assert.Len(t, items1, 1)

		// Second page
		rr2 := makeAuthenticatedRequest(t, server, "GET", "/users?limit=1&offset=1", nil, token)
		assertStatus(t, rr2, http.StatusOK)
		response2 := assertPaginatedResponse(t, rr2, 2, 1, 1)
		items2 := response2["items"].([]any)
		assert.Len(t, items2, 1)

		// Ensure different users
		user1 := items1[0].(map[string]any)
		user2 := items2[0].(map[string]any)
		assert.NotEqual(t, user1["id"], user2["id"])
	})

	t.Run("should handle empty results with offset beyond total", func(t *testing.T) {
		rr := makeAuthenticatedRequest(t, server, "GET", "/users?limit=10&offset=100", nil, token)
		assertStatus(t, rr, http.StatusOK)
		response := assertPaginatedResponse(t, rr, 2, 10, 100)
		items := response["items"].([]any)
		assert.Len(t, items, 0)
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/users", nil)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestUsersService_Get(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should return a user by id", func(t *testing.T) {
		rr := makeAuthenticatedRequest(t, server, "GET", "/users/1", nil, token)
		assertStatus(t, rr, http.StatusOK)

		var user map[string]any
		err := decodeJSON(rr, &user)
		require.NoError(t, err)

		assert.Equal(t, "1", user["id"])
		assert.Equal(t, "user1@example.com", user["email"])
		assert.Equal(t, "Test User 1", user["name"])
		assert.NotNil(t, user["address"])
	})

	t.Run("should return 404 for non-existent user", func(t *testing.T) {
		rr := makeAuthenticatedRequest(t, server, "GET", "/users/999", nil, token)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/users/1", nil)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestUsersService_Create(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should create a new user with address", func(t *testing.T) {
		newUser := map[string]any{
			"email": "newuser@example.com",
			"name":  "New User",
			"address": map[string]any{
				"street":     "789 New St",
				"city":       "New City",
				"state":      "NC",
				"postalCode": "67890",
				"country":    "USA",
			},
		}

		rr := makeAuthenticatedRequest(t, server, "POST", "/users", newUser, token)
		assertStatus(t, rr, http.StatusCreated)

		var user map[string]any
		err := decodeJSON(rr, &user)
		require.NoError(t, err)

		assert.NotEmpty(t, user["id"])
		assert.Equal(t, newUser["email"], user["email"])
		assert.Equal(t, newUser["name"], user["name"])
		assert.NotNil(t, user["address"])
		assert.NotEmpty(t, user["createdAt"])
		assert.NotEmpty(t, user["updatedAt"])

		// Verify cart was created for new user
		userID := user["id"].(string)
		cartRR := makeAuthenticatedRequest(t, server, "GET", "/carts/users/"+userID, nil, token)
		assertStatus(t, cartRR, http.StatusOK)

		var cart map[string]any
		err = decodeJSON(cartRR, &cart)
		require.NoError(t, err)
		assert.Equal(t, userID, cart["userId"])
		assert.Empty(t, cart["items"])
	})

	t.Run("should create a new user without address", func(t *testing.T) {
		newUser := map[string]any{
			"email": "minimal@example.com",
			"name":  "Minimal User",
		}

		rr := makeAuthenticatedRequest(t, server, "POST", "/users", newUser, token)
		assertStatus(t, rr, http.StatusCreated)

		var user map[string]any
		err := decodeJSON(rr, &user)
		require.NoError(t, err)

		assert.NotEmpty(t, user["id"])
		assert.Equal(t, newUser["email"], user["email"])
		assert.Equal(t, newUser["name"], user["name"])
		assert.Nil(t, user["address"])
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		newUser := map[string]any{
			"email": "unauth@example.com",
			"name":  "Unauthorized User",
		}

		rr := makeRequest(t, server, "POST", "/users", newUser)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})

	// TODO: Implement validation in Go handlers
	// t.Run("should return 400 for invalid user data", func(t *testing.T) {
	// 	invalidUser := map[string]any{
	// 		"email": "", // Empty email
	// 		"name":  "Invalid User",
	// 	}

	// 	rr := makeAuthenticatedRequest(t, server, "POST", "/users", invalidUser, token)
	// 	assertStatus(t, rr, http.StatusBadRequest)
	// 	assertErrorResponse(t, rr, "VALIDATION_ERROR")
	// })
}

func TestUsersService_Update(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should update user details", func(t *testing.T) {
		update := map[string]any{
			"name": "Updated Name",
		}

		rr := makeAuthenticatedRequest(t, server, "PATCH", "/users/1", update, token)
		assertStatus(t, rr, http.StatusOK)

		var user map[string]any
		err := decodeJSON(rr, &user)
		require.NoError(t, err)

		assert.Equal(t, "1", user["id"])
		assert.Equal(t, "Updated Name", user["name"])
		assert.Equal(t, "user1@example.com", user["email"]) // Email unchanged
	})

	t.Run("should update user address", func(t *testing.T) {
		update := map[string]any{
			"address": map[string]any{
				"street":     "999 Updated St",
				"city":       "Updated City",
				"state":      "UC",
				"postalCode": "99999",
				"country":    "USA",
			},
		}

		rr := makeAuthenticatedRequest(t, server, "PATCH", "/users/1", update, token)
		assertStatus(t, rr, http.StatusOK)

		var user map[string]any
		err := decodeJSON(rr, &user)
		require.NoError(t, err)

		address := user["address"].(map[string]any)
		assert.Equal(t, "999 Updated St", address["street"])
		assert.Equal(t, "Updated City", address["city"])
	})

	t.Run("should return 404 when updating non-existent user", func(t *testing.T) {
		update := map[string]any{
			"name": "Ghost User",
		}

		rr := makeAuthenticatedRequest(t, server, "PATCH", "/users/999", update, token)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		update := map[string]any{
			"name": "Unauthorized Update",
		}

		rr := makeRequest(t, server, "PATCH", "/users/1", update)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestUsersService_Delete(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should delete a user", func(t *testing.T) {
		// Create a new user to delete (note: createTestUser doesn't use auth)
		// We'll create using authenticated request instead
		newUser := map[string]any{
			"email": "delete@example.com",
			"name":  "Delete Me",
			"address": map[string]any{
				"street":     "123 Test St",
				"city":       "Test City",
				"state":      "TC",
				"postalCode": "12345",
				"country":    "Test Country",
			},
		}

		createRR := makeAuthenticatedRequest(t, server, "POST", "/users", newUser, token)
		assertStatus(t, createRR, http.StatusCreated)

		var user map[string]any
		err := decodeJSON(createRR, &user)
		require.NoError(t, err)
		userID := user["id"].(string)

		// Delete the user
		rr := makeAuthenticatedRequest(t, server, "DELETE", "/users/"+userID, nil, token)
		assertStatus(t, rr, http.StatusNoContent)

		// Verify user is deleted
		getRR := makeAuthenticatedRequest(t, server, "GET", "/users/"+userID, nil, token)
		assertStatus(t, getRR, http.StatusNotFound)
	})

	t.Run("should return 404 when deleting non-existent user", func(t *testing.T) {
		rr := makeAuthenticatedRequest(t, server, "DELETE", "/users/999", nil, token)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		rr := makeRequest(t, server, "DELETE", "/users/1", nil)
		assertStatus(t, rr, http.StatusUnauthorized)
		assertErrorResponse(t, rr, "UNAUTHORIZED")
	})
}

func TestUsersService_Integration(t *testing.T) {
	server, _, token := setupTestServerWithAuth(t)

	t.Run("should handle full user lifecycle", func(t *testing.T) {
		// Create user
		newUser := map[string]any{
			"email": "lifecycle@example.com",
			"name":  "Lifecycle User",
			"address": map[string]any{
				"street":     "123 Test St",
				"city":       "Test City",
				"state":      "TC",
				"postalCode": "12345",
				"country":    "Test Country",
			},
		}

		createRR := makeAuthenticatedRequest(t, server, "POST", "/users", newUser, token)
		assertStatus(t, createRR, http.StatusCreated)

		var user map[string]any
		err := decodeJSON(createRR, &user)
		require.NoError(t, err)
		userID := user["id"].(string)

		// Update user
		update := map[string]any{
			"name": "Updated Lifecycle User",
		}
		updateRR := makeAuthenticatedRequest(t, server, "PATCH", "/users/"+userID, update, token)
		assertStatus(t, updateRR, http.StatusOK)

		// Verify update
		getRR := makeAuthenticatedRequest(t, server, "GET", "/users/"+userID, nil, token)
		assertStatus(t, getRR, http.StatusOK)

		var updatedUser map[string]any
		err = decodeJSON(getRR, &updatedUser)
		require.NoError(t, err)
		assert.Equal(t, "Updated Lifecycle User", updatedUser["name"])

		// Delete user
		deleteRR := makeAuthenticatedRequest(t, server, "DELETE", "/users/"+userID, nil, token)
		assertStatus(t, deleteRR, http.StatusNoContent)

		// Verify deletion
		getDeletedRR := makeAuthenticatedRequest(t, server, "GET", "/users/"+userID, nil, token)
		assertStatus(t, getDeletedRR, http.StatusNotFound)
	})
}
