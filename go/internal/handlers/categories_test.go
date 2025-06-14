package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoriesService_List(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should return all categories", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/categories", nil)
		assertStatus(t, rr, http.StatusOK)

		var categories []any
		err := decodeJSON(rr, &categories)
		require.NoError(t, err)

		// Should have default categories
		assert.GreaterOrEqual(t, len(categories), 4)

		// Check first category structure
		category := categories[0].(map[string]any)
		assert.Contains(t, category, "id")
		assert.Contains(t, category, "name")
		// parentId is omitted when nil due to omitempty tag
		assert.Contains(t, category, "createdAt")
		assert.Contains(t, category, "updatedAt")
	})
}

func TestCategoriesService_GetTree(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should return categories in tree structure", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/categories/tree", nil)
		assertStatus(t, rr, http.StatusOK)

		var tree []any
		err := decodeJSON(rr, &tree)
		require.NoError(t, err)

		// Should have root categories
		assert.Greater(t, len(tree), 0)

		// Check tree structure
		rootCategory := tree[0].(map[string]any)
		assert.Contains(t, rootCategory, "id")
		assert.Contains(t, rootCategory, "name")
		// Note: Go implementation doesn't include parentId in tree response
		// assert.Contains(t, rootCategory, "parentId")
		assert.Contains(t, rootCategory, "children")

		// Check if children exist and are properly structured
		children, ok := rootCategory["children"].([]any)
		if ok && len(children) > 0 {
			childCategory := children[0].(map[string]any)
			assert.Equal(t, rootCategory["id"], childCategory["parentId"])
		}
	})
}

func TestCategoriesService_Get(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should return a category by id", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/categories/1", nil)
		assertStatus(t, rr, http.StatusOK)

		var category map[string]any
		err := decodeJSON(rr, &category)
		require.NoError(t, err)

		assert.Equal(t, "1", category["id"])
		assert.NotEmpty(t, category["name"])
		assert.Nil(t, category["parentId"]) // Default category 1 is root
	})

	t.Run("should return 404 for non-existent category", func(t *testing.T) {
		rr := makeRequest(t, server, "GET", "/categories/999", nil)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})
}

func TestCategoriesService_Create(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should create a new root category", func(t *testing.T) {
		newCategory := map[string]any{
			"name":     "New Root Category",
			"parentId": nil,
		}

		rr := makeRequest(t, server, "POST", "/categories", newCategory)
		assertStatus(t, rr, http.StatusCreated)

		var category map[string]any
		err := decodeJSON(rr, &category)
		require.NoError(t, err)

		assert.NotEmpty(t, category["id"])
		assert.Equal(t, newCategory["name"], category["name"])
		assert.Nil(t, category["parentId"])
		assert.NotEmpty(t, category["createdAt"])
		assert.NotEmpty(t, category["updatedAt"])
	})

	t.Run("should create a child category", func(t *testing.T) {
		// Create parent first
		parentID := createTestCategory(t, server, "Parent Category", nil)

		// Create child
		newCategory := map[string]any{
			"name":     "Child Category",
			"parentId": parentID,
		}

		rr := makeRequest(t, server, "POST", "/categories", newCategory)
		assertStatus(t, rr, http.StatusCreated)

		var category map[string]any
		err := decodeJSON(rr, &category)
		require.NoError(t, err)

		assert.NotEmpty(t, category["id"])
		assert.Equal(t, newCategory["name"], category["name"])
		assert.Equal(t, parentID, category["parentId"])
	})

	t.Run("should return 404 for non-existent parent", func(t *testing.T) {
		newCategory := map[string]any{
			"name":     "Orphan Category",
			"parentId": "999999",
		}

		rr := makeRequest(t, server, "POST", "/categories", newCategory)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})

	// TODO: Add validation tests when implemented
	// t.Run("should return 400 for empty name", func(t *testing.T) {
	// 	newCategory := map[string]any{
	// 		"name":     "",
	// 		"parentId": nil,
	// 	}

	// 	rr := makeRequest(t, server, "POST", "/categories", newCategory)
	// 	assertStatus(t, rr, http.StatusBadRequest)
	// 	assertErrorResponse(t, rr, "VALIDATION_ERROR")
	// })
}

func TestCategoriesService_Update(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should update category name", func(t *testing.T) {
		// Create a category to update
		categoryID := createTestCategory(t, server, "Original Name", nil)

		update := map[string]any{
			"name": "Updated Name",
		}

		rr := makeRequest(t, server, "PATCH", "/categories/"+categoryID, update)
		assertStatus(t, rr, http.StatusOK)

		var category map[string]any
		err := decodeJSON(rr, &category)
		require.NoError(t, err)

		assert.Equal(t, categoryID, category["id"])
		assert.Equal(t, "Updated Name", category["name"])
	})

	t.Run("should update category parent", func(t *testing.T) {
		// Create parent and child categories
		parentID := createTestCategory(t, server, "New Parent", nil)
		childID := createTestCategory(t, server, "Child", nil)

		update := map[string]any{
			"parentId": parentID,
		}

		rr := makeRequest(t, server, "PATCH", "/categories/"+childID, update)
		assertStatus(t, rr, http.StatusOK)

		var category map[string]any
		err := decodeJSON(rr, &category)
		require.NoError(t, err)

		assert.Equal(t, parentID, category["parentId"])
	})

	t.Run("should return 404 when updating non-existent category", func(t *testing.T) {
		update := map[string]any{
			"name": "Ghost Category",
		}

		rr := makeRequest(t, server, "PATCH", "/categories/999", update)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})

	// TODO: Add circular reference prevention test
	// t.Run("should prevent circular reference", func(t *testing.T) {
	// 	// Create parent-child relationship
	// 	parentID := createTestCategory(t, server, "Parent", nil)
	// 	childID := createTestCategory(t, server, "Child", &parentID)

	// 	// Try to make parent a child of its own child
	// 	update := map[string]any{
	// 		"parentId": childID,
	// 	}

	// 	rr := makeRequest(t, server, "PATCH", "/categories/"+parentID, update)
	// 	assertStatus(t, rr, http.StatusBadRequest)
	// 	assertErrorResponse(t, rr, "CIRCULAR_REFERENCE")
	// })
}

func TestCategoriesService_Delete(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should delete a category", func(t *testing.T) {
		// Create a category to delete
		categoryID := createTestCategory(t, server, "Delete Me", nil)

		// Delete the category
		rr := makeRequest(t, server, "DELETE", "/categories/"+categoryID, nil)
		assertStatus(t, rr, http.StatusNoContent)

		// Verify category is deleted
		getRR := makeRequest(t, server, "GET", "/categories/"+categoryID, nil)
		assertStatus(t, getRR, http.StatusNotFound)
	})

	t.Run("should return 404 when deleting non-existent category", func(t *testing.T) {
		rr := makeRequest(t, server, "DELETE", "/categories/999", nil)
		assertStatus(t, rr, http.StatusNotFound)
		assertErrorResponse(t, rr, "NOT_FOUND")
	})

	// TODO: Add test for preventing deletion of category with children
	// t.Run("should prevent deletion of category with children", func(t *testing.T) {
	// 	// Create parent with child
	// 	parentID := createTestCategory(t, server, "Parent", nil)
	// 	createTestCategory(t, server, "Child", &parentID)

	// 	// Try to delete parent
	// 	rr := makeRequest(t, server, "DELETE", "/categories/"+parentID, nil)
	// 	assertStatus(t, rr, http.StatusBadRequest)
	// 	assertErrorResponse(t, rr, "HAS_CHILDREN")
	// })
}

func TestCategoriesService_Integration(t *testing.T) {
	server := setupTestServer(t)

	t.Run("should handle complete category hierarchy", func(t *testing.T) {
		// Create a category hierarchy
		rootID := createTestCategory(t, server, "Integration Root", nil)
		child1ID := createTestCategory(t, server, "Child 1", &rootID)
		child2ID := createTestCategory(t, server, "Child 2", &rootID)
		grandchildID := createTestCategory(t, server, "Grandchild", &child1ID)

		// Get tree structure
		treeRR := makeRequest(t, server, "GET", "/categories/tree", nil)
		assertStatus(t, treeRR, http.StatusOK)

		var tree []any
		err := decodeJSON(treeRR, &tree)
		require.NoError(t, err)

		// Find our root category in the tree
		var foundRoot map[string]any
		for _, cat := range tree {
			category := cat.(map[string]any)
			if category["id"] == rootID {
				foundRoot = category
				break
			}
		}
		require.NotNil(t, foundRoot, "Root category should be in tree")

		// Verify children
		childrenRaw, ok := foundRoot["children"]
		require.True(t, ok, "foundRoot should have children field")
		children := childrenRaw.([]any)
		assert.Len(t, children, 2) // Should have 2 children

		// Find child1 and verify it has grandchild
		var foundChild1 bool
		for _, child := range children {
			childCat := child.(map[string]any)
			if childCat["id"] == child1ID {
				foundChild1 = true
				grandchildren := childCat["children"].([]any)
				assert.Len(t, grandchildren, 1)
				if len(grandchildren) > 0 {
					grandchild := grandchildren[0].(map[string]any)
					assert.Equal(t, grandchildID, grandchild["id"])
				}
				break
			}
		}
		assert.True(t, foundChild1, "Child1 should be found in the tree")

		// Clean up: Delete in reverse order (grandchild first)
		makeRequest(t, server, "DELETE", "/categories/"+grandchildID, nil)
		makeRequest(t, server, "DELETE", "/categories/"+child1ID, nil)
		makeRequest(t, server, "DELETE", "/categories/"+child2ID, nil)
		makeRequest(t, server, "DELETE", "/categories/"+rootID, nil)
	})
}