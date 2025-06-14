package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersService_List(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should return all users with default pagination", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/users", nil)

		assertStatus(t, rr, http.StatusOK)
		response := assertPaginatedResponse(t, rr, 2, 20, 0) // 2 default users

		items := response["items"].([]interface{})
		assert.Len(t, items, 2)

		// Check first user structure
		user := items[0].(map[string]interface{})
		assert.Contains(t, user, "id")
		assert.Contains(t, user, "email")
		assert.Contains(t, user, "name")
		assert.Contains(t, user, "address")
		assert.Contains(t, user, "createdAt")
		assert.Contains(t, user, "updatedAt")
	})

	t.Run("should support pagination with limit and offset", func(t *testing.T) {
		// First page
		rr1 := makeRequest(t, server, "GET", "/users?limit=1&offset=0", nil)
		assertStatus(t, rr1, http.StatusOK)
		response1 := assertPaginatedResponse(t, rr1, 2, 1, 0)
		items1 := response1["items"].([]interface{})
		assert.Len(t, items1, 1)

		// Second page
		rr2 := makeRequest(t, server, "GET", "/users?limit=1&offset=1", nil)
		assertStatus(t, rr2, http.StatusOK)
		response2 := assertPaginatedResponse(t, rr2, 2, 1, 1)
		items2 := response2["items"].([]interface{})
		assert.Len(t, items2, 1)

		// Ensure different users
		user1 := items1[0].(map[string]interface{})
		user2 := items2[0].(map[string]interface{})
		assert.NotEqual(t, user1["id"], user2["id"])
	})

	t.Run("should handle empty results with offset beyond total", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/users?limit=10&offset=100", nil)
		assertStatus(t, rr, http.StatusOK)
		response := assertPaginatedResponse(t, rr, 2, 10, 100)
		items := response["items"].([]interface{})
		assert.Len(t, items, 0)
	})
}

func TestUsersService_Get(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should return a user by id", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/users/1", nil)
		assertStatus(t, rr, http.StatusOK)

		var user map[string]interface{}
		err := decodeJSON(rr, &user)
		require.NoError(t, err)

		assert.Equal(t, "1", user["id"])
		assert.Equal(t, "user1@example.com", user["email"])
		assert.Equal(t, "Test User 1", user["name"])
		assert.NotNil(t, user["address"])
	})

	t.Run("should return 404 for non-existent user", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/users/999", nil)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})
}

func TestUsersService_Create(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should create a new user with address", func(t *testing.T) {
		newUser := map[string]interface{}{
			"email": "newuser@example.com",
			"name":  "New User",
			"address": map[string]interface{}{
				"street":     "789 New St",
				"city":       "New City",
				"state":      "NC",
				"postalCode": "67890",
				"country":    "USA",
			},
		}

		rr := makeRequest(t, server, "POST", "/users", newUser)
		assertStatus(t, rr, http.StatusCreated)

		var user map[string]interface{}
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
		cartRR := makeRequest(t, server, "GET", "/carts/users/"+userID, nil)
		assertStatus(t, cartRR, http.StatusOK)

		var cart map[string]interface{}
		err = decodeJSON(cartRR, &cart)
		require.NoError(t, err)
		assert.Equal(t, userID, cart["userId"])
		assert.Empty(t, cart["items"])
	})

	t.Run("should create a new user without address", func(t *testing.T) {
		newUser := map[string]interface{}{
			"email": "minimal@example.com",
			"name":  "Minimal User",
		}

		rr := makeRequest(t, server, "POST", "/users", newUser)
		assertStatus(t, rr, http.StatusCreated)

		var user map[string]interface{}
		err := decodeJSON(rr, &user)
		require.NoError(t, err)

		assert.NotEmpty(t, user["id"])
		assert.Equal(t, newUser["email"], user["email"])
		assert.Equal(t, newUser["name"], user["name"])
		assert.Nil(t, user["address"])
	})

	// TODO: Implement validation in Go handlers
	// t.Run("should return 400 for invalid user data", func(t *testing.T) {
	// 	invalidUser := map[string]interface{}{
	// 		"email": "", // Empty email
	// 		"name":  "Invalid User",
	// 	}

	// 	rr := makeRequest(t, server, "POST", "/users", invalidUser)
	// 	assertStatus(t, rr, http.StatusBadRequest)
	// 	assertErrorResponse(t, rr, "VALIDATION_ERROR")
	// })
}

func TestUsersService_Update(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should update user details", func(t *testing.T) {
		update := map[string]interface{}{
			"name": "Updated Name",
		}

		rr := makeRequest(t, server, "PATCH", "/users/1", update)
		assertStatus(t, rr, http.StatusOK)

		var user map[string]interface{}
		err := decodeJSON(rr, &user)
		require.NoError(t, err)

		assert.Equal(t, "1", user["id"])
		assert.Equal(t, "Updated Name", user["name"])
		assert.Equal(t, "user1@example.com", user["email"]) // Email unchanged
	})

	t.Run("should update user address", func(t *testing.T) {
		update := map[string]interface{}{
			"address": map[string]interface{}{
				"street":     "999 Updated St",
				"city":       "Updated City",
				"state":      "UC",
				"postalCode": "99999",
				"country":    "USA",
			},
		}

		rr := makeRequest(t, server, "PATCH", "/users/1", update)
		assertStatus(t, rr, http.StatusOK)

		var user map[string]interface{}
		err := decodeJSON(rr, &user)
		require.NoError(t, err)

		address := user["address"].(map[string]interface{})
		assert.Equal(t, "999 Updated St", address["street"])
		assert.Equal(t, "Updated City", address["city"])
	})

	t.Run("should return 404 when updating non-existent user", func(t *testing.T) {
		update := map[string]interface{}{
			"name": "Ghost User",
		}

		rr := makeRequest(t, server, "PATCH", "/users/999", update)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})
}

func TestUsersService_Delete(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should delete a user", func(t *testing.T) {
		// Create a new user to delete
		userID := createTestUser(t, server, "delete@example.com", "Delete Me")

		// Delete the user
		rr := makeRequest(t, server, "DELETE", "/users/"+userID, nil)
		assertStatus(t, rr, http.StatusNoContent)

		// Verify user is deleted
		getRR := makeRequest(t, server, "GET", "/users/"+userID, nil)
		assertStatus(t, getRR, http.StatusNotFound)
	})

	t.Run("should return 404 when deleting non-existent user", func(t *testing.T) {
		rr := makeRequest(t, server, "DELETE", "/users/999", nil)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})
}

func TestUsersService_Integration(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should handle full user lifecycle", func(t *testing.T) {
		// Create user
		userID := createTestUser(t, server, "lifecycle@example.com", "Lifecycle User")

		// Update user
		update := map[string]interface{}{
			"name": "Updated Lifecycle User",
		}
		updateRR := makeRequest(t, server, "PATCH", "/users/"+userID, update)
		assertStatus(t, updateRR, http.StatusOK)

		// Verify update
		getRR := makeRequest(t, server, "GET", "/users/"+userID, nil)
		assertStatus(t, getRR, http.StatusOK)

		var user map[string]interface{}
		err := decodeJSON(getRR, &user)
		require.NoError(t, err)
		assert.Equal(t, "Updated Lifecycle User", user["name"])

		// Delete user
		deleteRR := makeRequest(t, server, "DELETE", "/users/"+userID, nil)
		assertStatus(t, deleteRR, http.StatusNoContent)

		// Verify deletion
		getDeletedRR := makeRequest(t, server, "GET", "/users/"+userID, nil)
		assertStatus(t, getDeletedRR, http.StatusNotFound)
	})
}

// decodeJSON is a helper to decode JSON response
func decodeJSON(rr *httptest.ResponseRecorder, v interface{}) error {
	return json.NewDecoder(rr.Body).Decode(v)
}