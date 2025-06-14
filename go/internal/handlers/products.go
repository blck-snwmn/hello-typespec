package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/blck-snwmn/hello-typespec/go/generated"
)

// ProductsServiceList implements GET /products
func (s *Server) ProductsServiceList(w http.ResponseWriter, r *http.Request, params generated.ProductsServiceListParams) {
	// Get all products
	allProducts := s.store.GetProducts()

	// Apply filters
	var filteredProducts []generated.Product
	for _, product := range allProducts {
		// Search filter (name)
		if params.Name != nil && *params.Name != "" {
			searchStr := strings.ToLower(*params.Name)
			if !strings.Contains(strings.ToLower(product.Name), searchStr) &&
				!strings.Contains(strings.ToLower(product.Description), searchStr) {
				continue
			}
		}

		// Category filter
		if params.CategoryId != nil && *params.CategoryId != "" {
			if product.CategoryId != *params.CategoryId {
				continue
			}
		}

		// Price filters
		if params.MinPrice != nil && product.Price < *params.MinPrice {
			continue
		}
		if params.MaxPrice != nil && product.Price > *params.MaxPrice {
			continue
		}

		filteredProducts = append(filteredProducts, product)
	}

	// Apply sorting
	if params.SortBy != nil {
		switch *params.SortBy {
		case generated.ProductsServiceListParamsSortByName:
			sort.Slice(filteredProducts, func(i, j int) bool {
				if params.Order != nil && *params.Order == generated.ProductsServiceListParamsOrderDesc {
					return filteredProducts[i].Name > filteredProducts[j].Name
				}
				return filteredProducts[i].Name < filteredProducts[j].Name
			})
		case generated.ProductsServiceListParamsSortByPrice:
			sort.Slice(filteredProducts, func(i, j int) bool {
				if params.Order != nil && *params.Order == generated.ProductsServiceListParamsOrderDesc {
					return filteredProducts[i].Price > filteredProducts[j].Price
				}
				return filteredProducts[i].Price < filteredProducts[j].Price
			})
		case generated.ProductsServiceListParamsSortByCreatedAt:
			sort.Slice(filteredProducts, func(i, j int) bool {
				if params.Order != nil && *params.Order == generated.ProductsServiceListParamsOrderDesc {
					return filteredProducts[i].CreatedAt.After(filteredProducts[j].CreatedAt)
				}
				return filteredProducts[i].CreatedAt.Before(filteredProducts[j].CreatedAt)
			})
		}
	}

	// Apply pagination
	limit := int32(10)
	if params.Limit != nil {
		limit = *params.Limit
	}

	offset := int32(0)
	if params.Offset != nil {
		offset = *params.Offset
	}

	total := int32(len(filteredProducts))
	start := int(offset)
	end := int(offset + limit)

	if start > len(filteredProducts) {
		start = len(filteredProducts)
	}
	if end > len(filteredProducts) {
		end = len(filteredProducts)
	}

	paginatedProducts := filteredProducts[start:end]

	// Create response
	response := struct {
		Items  []generated.Product `json:"items"`
		Total  int32               `json:"total"`
		Limit  int32               `json:"limit"`
		Offset int32               `json:"offset"`
	}{
		Items:  paginatedProducts,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ProductsServiceGet implements GET /products/{productId}
func (s *Server) ProductsServiceGet(w http.ResponseWriter, r *http.Request, productId generated.Uuid) {
	product, ok := s.store.GetProduct(productId)
	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// ProductsServiceCreate implements POST /products
func (s *Server) ProductsServiceCreate(w http.ResponseWriter, r *http.Request) {
	var req generated.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Create new product
	now := time.Now()
	newProduct := generated.Product{
		Id:          fmt.Sprintf("%d", now.UnixNano()),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryId:  req.CategoryId,
		ImageUrls:   []string{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if req.ImageUrls != nil {
		newProduct.ImageUrls = *req.ImageUrls
	}

	created := s.store.CreateProduct(newProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// ProductsServiceUpdate implements PATCH /products/{productId}
func (s *Server) ProductsServiceUpdate(w http.ResponseWriter, r *http.Request, productId generated.Uuid) {
	existing, ok := s.store.GetProduct(productId)
	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
		return
	}

	var req generated.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Update fields if provided
	updatedProduct := *existing
	if req.Name != nil {
		updatedProduct.Name = *req.Name
	}
	if req.Description != nil {
		updatedProduct.Description = *req.Description
	}
	if req.Price != nil {
		updatedProduct.Price = *req.Price
	}
	if req.Stock != nil {
		updatedProduct.Stock = *req.Stock
	}
	if req.CategoryId != nil {
		updatedProduct.CategoryId = *req.CategoryId
	}
	if req.ImageUrls != nil {
		updatedProduct.ImageUrls = *req.ImageUrls
	}
	updatedProduct.UpdatedAt = time.Now()

	updated := s.store.UpdateProduct(productId, updatedProduct)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// ProductsServiceDelete implements DELETE /products/{productId}
func (s *Server) ProductsServiceDelete(w http.ResponseWriter, r *http.Request, productId generated.Uuid) {
	_, ok := s.store.DeleteProduct(productId)
	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper function for error responses
func errorResponse(w http.ResponseWriter, statusCode int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(generated.ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	})
}