package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/blck-snwmn/hello-typespec/go/generated"
)

// CategoryWithChildren represents a category with its child categories
type CategoryWithChildren struct {
	generated.Category
	Children []*CategoryWithChildren `json:"children"`
}

// CategoriesServiceList implements GET /categories
func (s *Server) CategoriesServiceList(w http.ResponseWriter, r *http.Request) {
	categories := s.store.GetCategories()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// CategoriesServiceTree implements GET /categories/tree
func (s *Server) CategoriesServiceTree(w http.ResponseWriter, r *http.Request) {
	allCategories := s.store.GetCategories()

	// Build tree structure
	categoryMap := make(map[string]*CategoryWithChildren)
	var rootCategories []*CategoryWithChildren

	// First pass: create CategoryWithChildren objects
	for _, cat := range allCategories {
		categoryMap[cat.Id] = &CategoryWithChildren{
			Category: cat,
			Children: []*CategoryWithChildren{},
		}
	}

	// Second pass: build tree and collect roots
	for _, catWithChildren := range categoryMap {
		if catWithChildren.ParentId == nil {
			rootCategories = append(rootCategories, catWithChildren)
		} else {
			parent, exists := categoryMap[*catWithChildren.ParentId]
			if exists {
				parent.Children = append(parent.Children, catWithChildren)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rootCategories)
}

// CategoriesServiceGet implements GET /categories/{categoryId}
func (s *Server) CategoriesServiceGet(w http.ResponseWriter, r *http.Request, categoryId generated.Uuid) {
	category, ok := s.store.GetCategory(categoryId)
	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Category not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// CategoriesServiceCreate implements POST /categories
func (s *Server) CategoriesServiceCreate(w http.ResponseWriter, r *http.Request) {
	var req generated.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Validate parent category exists if provided
	if req.ParentId != nil {
		_, exists := s.store.GetCategory(*req.ParentId)
		if !exists {
			errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Parent category not found")
			return
		}
	}

	// Create new category
	now := time.Now()
	newCategory := generated.Category{
		Id:        fmt.Sprintf("%d", now.UnixNano()),
		Name:      req.Name,
		ParentId:  req.ParentId,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created := s.store.CreateCategory(newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// CategoriesServiceUpdate implements PATCH /categories/{categoryId}
func (s *Server) CategoriesServiceUpdate(w http.ResponseWriter, r *http.Request, categoryId generated.Uuid) {
	existing, ok := s.store.GetCategory(categoryId)
	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Category not found")
		return
	}

	var req generated.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Update fields if provided
	updatedCategory := *existing
	if req.Name != nil {
		updatedCategory.Name = *req.Name
	}
	if req.ParentId != nil {
		updatedCategory.ParentId = req.ParentId
	}
	updatedCategory.UpdatedAt = time.Now()

	updated := s.store.UpdateCategory(categoryId, updatedCategory)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// CategoriesServiceDelete implements DELETE /categories/{categoryId}
func (s *Server) CategoriesServiceDelete(w http.ResponseWriter, r *http.Request, categoryId generated.Uuid) {
	_, ok := s.store.DeleteCategory(categoryId)
	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Category not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}