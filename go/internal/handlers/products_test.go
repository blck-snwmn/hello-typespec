package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductsService_List(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should return all products with default pagination", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/products", nil)
		assertStatus(t, rr, http.StatusOK)
		response := assertPaginatedResponse(t, rr, 3, 20, 0) // 3 default products

		items := response["items"].([]any)
		assert.Len(t, items, 3)

		// Check first product structure
		product := items[0].(map[string]any)
		assert.Contains(t, product, "id")
		assert.Contains(t, product, "name")
		assert.Contains(t, product, "description")
		assert.Contains(t, product, "price")
		assert.Contains(t, product, "stock")
		assert.Contains(t, product, "categoryId")
		assert.Contains(t, product, "imageUrls")
		assert.Contains(t, product, "createdAt")
		assert.Contains(t, product, "updatedAt")
	})

	t.Run("should support pagination", func(t *testing.T) {
		// Create additional products
		for i := 0; i < 5; i++ {
			createTestProduct(t, server, "Extra Product", 10.00, 10)
		}

		// Test pagination
		rr := makeRequest(t, server, "GET", "/products?limit=5&offset=2", nil)
		assertStatus(t, rr, http.StatusOK)
		response := assertPaginatedResponse(t, rr, 8, 5, 2)
		items := response["items"].([]any)
		assert.Len(t, items, 5)
	})

	t.Run("should filter by name", func(t *testing.T) {
		// Create product with specific name
		createTestProduct(t, server, "Unique Search Product", 50.00, 5)

		rr := makeRequest(t, server, "GET", "/products?name=Unique+Search", nil)
		assertStatus(t, rr, http.StatusOK)

		var response map[string]any
		err := decodeJSON(rr, &response)
		require.NoError(t, err)

		items := response["items"].([]any)
		assert.Len(t, items, 1)
		product := items[0].(map[string]any)
		assert.Contains(t, product["name"], "Unique Search")
	})

	t.Run("should filter by category", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/products?categoryId=2", nil)
		assertStatus(t, rr, http.StatusOK)

		var response map[string]any
		err := decodeJSON(rr, &response)
		require.NoError(t, err)

		items := response["items"].([]any)
		// Should return products in category 2 (Computers)
		for _, item := range items {
			product := item.(map[string]any)
			assert.Equal(t, "2", product["categoryId"])
		}
	})

	t.Run("should filter by price range", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/products?minPrice=100&maxPrice=1000", nil)
		assertStatus(t, rr, http.StatusOK)

		var response map[string]any
		err := decodeJSON(rr, &response)
		require.NoError(t, err)

		items := response["items"].([]any)
		for _, item := range items {
			product := item.(map[string]any)
			price := product["price"].(float64)
			assert.GreaterOrEqual(t, price, float64(100))
			assert.LessOrEqual(t, price, float64(1000))
		}
	})

	t.Run("should sort by price", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/products?sortBy=price&order=asc", nil)
		assertStatus(t, rr, http.StatusOK)

		var response map[string]any
		err := decodeJSON(rr, &response)
		require.NoError(t, err)

		items := response["items"].([]any)
		assert.Greater(t, len(items), 1)

		// Verify ascending order
		for i := 1; i < len(items); i++ {
			prevPrice := items[i-1].(map[string]any)["price"].(float64)
			currPrice := items[i].(map[string]any)["price"].(float64)
			assert.LessOrEqual(t, prevPrice, currPrice)
		}
	})

	t.Run("should sort by name", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/products?sortBy=name&order=desc", nil)
		assertStatus(t, rr, http.StatusOK)

		var response map[string]any
		err := decodeJSON(rr, &response)
		require.NoError(t, err)

		items := response["items"].([]any)
		assert.Greater(t, len(items), 1)

		// Verify descending order
		for i := 1; i < len(items); i++ {
			prevName := items[i-1].(map[string]any)["name"].(string)
			currName := items[i].(map[string]any)["name"].(string)
			assert.GreaterOrEqual(t, prevName, currName)
		}
	})
}

func TestProductsService_Get(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should return a product by id", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/products/1", nil)
		assertStatus(t, rr, http.StatusOK)

		var product map[string]any
		err := decodeJSON(rr, &product)
		require.NoError(t, err)

		assert.Equal(t, "1", product["id"])
		assert.Equal(t, "MacBook Pro 16\"", product["name"])
		assert.Equal(t, float64(2499.99), product["price"])
		assert.Equal(t, "2", product["categoryId"])
	})

	t.Run("should return 404 for non-existent product", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/products/999", nil)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})
}

func TestProductsService_Create(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should create a new product", func(t *testing.T) {
		newProduct := map[string]any{
			"name":        "New Test Product",
			"description": "A test product description",
			"price":       199.99,
			"stock":       50,
			"categoryId":  "1",
			"imageUrls":   []string{"https://example.com/product1.jpg", "https://example.com/product2.jpg"},
		}

		rr := makeRequest(t, server, "POST", "/products", newProduct)
		assertStatus(t, rr, http.StatusCreated)

		var product map[string]any
		err := decodeJSON(rr, &product)
		require.NoError(t, err)

		assert.NotEmpty(t, product["id"])
		assert.Equal(t, newProduct["name"], product["name"])
		assert.Equal(t, newProduct["description"], product["description"])
		assert.Equal(t, newProduct["price"], product["price"])
		assert.Equal(t, float64(newProduct["stock"].(int)), product["stock"])
		assert.Equal(t, newProduct["categoryId"], product["categoryId"])
		imageUrls := product["imageUrls"].([]any)
		assert.Len(t, imageUrls, 2)
		assert.Equal(t, "https://example.com/product1.jpg", imageUrls[0])
		assert.Equal(t, "https://example.com/product2.jpg", imageUrls[1])
		assert.NotEmpty(t, product["createdAt"])
		assert.NotEmpty(t, product["updatedAt"])
	})

	t.Run("should create product without optional imageUrls", func(t *testing.T) {
		newProduct := map[string]any{
			"name":        "Product Without Images",
			"description": "No images",
			"price":       99.99,
			"stock":       10,
			"categoryId":  "1",
		}

		rr := makeRequest(t, server, "POST", "/products", newProduct)
		assertStatus(t, rr, http.StatusCreated)

		var product map[string]any
		err := decodeJSON(rr, &product)
		require.NoError(t, err)

		imageUrls := product["imageUrls"].([]any)
		assert.Empty(t, imageUrls)
	})

	// TODO: Add validation tests when implemented
	// t.Run("should return 400 for invalid product data", func(t *testing.T) {
	// 	invalidProduct := map[string]any{
	// 		"name":       "", // Empty name
	// 		"price":      -10.00, // Negative price
	// 		"stock":      -5, // Negative stock
	// 		"categoryId": "1",
	// 	}

	// 	rr := makeRequest(t, server, "POST", "/products", invalidProduct)
	// 	assertStatus(t, rr, http.StatusBadRequest)
	// 	assertErrorResponse(t, rr, "VALIDATION_ERROR")
	// })
}

func TestProductsService_Update(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should update product fields", func(t *testing.T) {
		// Create a product to update
		productID := createTestProduct(t, server, "Original Product", 100.00, 20)

		update := map[string]any{
			"name":  "Updated Product",
			"price": 150.00,
			"stock": 30,
		}

		rr := makeRequest(t, server, "PATCH", "/products/"+productID, update)
		assertStatus(t, rr, http.StatusOK)

		var product map[string]any
		err := decodeJSON(rr, &product)
		require.NoError(t, err)

		assert.Equal(t, productID, product["id"])
		assert.Equal(t, "Updated Product", product["name"])
		assert.Equal(t, float64(150.00), product["price"])
		assert.Equal(t, float64(30), product["stock"])
	})

	t.Run("should update only provided fields", func(t *testing.T) {
		productID := createTestProduct(t, server, "Partial Update Product", 200.00, 15)

		// Get original product
		getRR := makeRequest(t, server, "GET", "/products/"+productID, nil)
		var original map[string]any
		decodeJSON(getRR, &original)

		// Update only name
		update := map[string]any{
			"name": "Name Only Updated",
		}

		rr := makeRequest(t, server, "PATCH", "/products/"+productID, update)
		assertStatus(t, rr, http.StatusOK)

		var product map[string]any
		err := decodeJSON(rr, &product)
		require.NoError(t, err)

		assert.Equal(t, "Name Only Updated", product["name"])
		assert.Equal(t, original["price"], product["price"]) // Price unchanged
		assert.Equal(t, original["stock"], product["stock"]) // Stock unchanged
	})

	t.Run("should return 404 when updating non-existent product", func(t *testing.T) {
		update := map[string]any{
			"name": "Ghost Product",
		}

		rr := makeRequest(t, server, "PATCH", "/products/999", update)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})
}

func TestProductsService_Delete(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should delete a product", func(t *testing.T) {
		// Create a product to delete
		productID := createTestProduct(t, server, "Delete Me", 50.00, 5)

		// Delete the product
		rr := makeRequest(t, server, "DELETE", "/products/"+productID, nil)
		assertStatus(t, rr, http.StatusNoContent)

		// Verify product is deleted
		getRR := makeRequest(t, server, "GET", "/products/"+productID, nil)
		assertStatus(t, getRR, http.StatusNotFound)
	})

	t.Run("should return 404 when deleting non-existent product", func(t *testing.T) {
		rr := makeRequest(t, server, "DELETE", "/products/999", nil)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})
}

func TestProductsService_Integration(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should handle complex search scenarios", func(t *testing.T) {
		// Create test products with different attributes
		cat1 := createTestCategory(t, server, "Test Electronics", nil)
		cat2 := createTestCategory(t, server, "Test Clothing", nil)

		createTestProductWithCategory(t, server, "Budget Phone", 299.99, 20, cat1)
		createTestProductWithCategory(t, server, "Premium Phone", 1299.99, 5, cat1)
		createTestProductWithCategory(t, server, "Basic T-Shirt", 19.99, 100, cat2)
		createTestProductWithCategory(t, server, "Designer T-Shirt", 99.99, 10, cat2)

		// Search for electronics in specific price range
		rr := makeRequest(t, server, "GET", "/products?categoryId="+cat1+"&minPrice=200&maxPrice=500&sortBy=price&order=asc", nil)
		assertStatus(t, rr, http.StatusOK)

		var response map[string]any
		err := decodeJSON(rr, &response)
		require.NoError(t, err)

		items := response["items"].([]any)
		assert.Len(t, items, 1)
		product := items[0].(map[string]any)
		assert.Equal(t, "Budget Phone", product["name"])

		// Search for clothing sorted by price descending
		rr2 := makeRequest(t, server, "GET", "/products?categoryId="+cat2+"&sortBy=price&order=desc", nil)
		assertStatus(t, rr2, http.StatusOK)

		var response2 map[string]any
		err = decodeJSON(rr2, &response2)
		require.NoError(t, err)

		items2 := response2["items"].([]any)
		assert.GreaterOrEqual(t, len(items2), 2)
		firstProduct := items2[0].(map[string]any)
		assert.Equal(t, "Designer T-Shirt", firstProduct["name"])
	})
}

// Helper function to create product with specific category
func createTestProductWithCategory(t *testing.T, server *TestServer, name string, price float64, stock int, categoryID string) string {
	t.Helper()

	product := map[string]any{
		"name":        name,
		"description": "Test product",
		"price":       price,
		"stock":       stock,
		"categoryId":  categoryID,
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
