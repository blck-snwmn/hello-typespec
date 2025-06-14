package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/blck-snwmn/hello-typespec/go/generated"
)

// CartsServiceGetByUser implements GET /carts/users/{userId}
func (s *Server) CartsServiceGetByUser(w http.ResponseWriter, r *http.Request, userId generated.Uuid) {
	cart := s.store.GetCartByUserId(userId)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// CartsServiceAddItem implements POST /carts/users/{userId}/items
func (s *Server) CartsServiceAddItem(w http.ResponseWriter, r *http.Request, userId generated.Uuid) {
	var req generated.AddCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	cart := s.store.GetCartByUserId(userId)
	product, ok := s.store.GetProduct(req.ProductId)

	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
		return
	}

	if product.Stock < req.Quantity {
		errorResponse(w, http.StatusBadRequest, "INSUFFICIENT_STOCK", "Insufficient stock")
		return
	}

	// Check if item already exists in cart
	itemFound := false
	for i := range cart.Items {
		if cart.Items[i].ProductId == req.ProductId {
			cart.Items[i].Quantity += req.Quantity
			itemFound = true
			break
		}
	}

	if !itemFound {
		// Add new item
		cart.Items = append(cart.Items, generated.CartItem{
			ProductId: req.ProductId,
			Quantity:  req.Quantity,
		})
	}

	cart.UpdatedAt = time.Now()
	updated := s.store.UpdateCart(userId, cart)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// CartsServiceUpdateItem implements PATCH /carts/users/{userId}/items/{productId}
func (s *Server) CartsServiceUpdateItem(w http.ResponseWriter, r *http.Request, userId generated.Uuid, productId generated.Uuid) {
	var req generated.UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	cart := s.store.GetCartByUserId(userId)
	product, ok := s.store.GetProduct(productId)

	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
		return
	}

	itemIndex := -1
	for i := range cart.Items {
		if cart.Items[i].ProductId == productId {
			itemIndex = i
			break
		}
	}

	if itemIndex < 0 {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Item not found in cart")
		return
	}

	if product.Stock < req.Quantity {
		errorResponse(w, http.StatusBadRequest, "INSUFFICIENT_STOCK", "Insufficient stock")
		return
	}

	cart.Items[itemIndex].Quantity = req.Quantity
	cart.UpdatedAt = time.Now()

	updated := s.store.UpdateCart(userId, cart)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// CartsServiceRemoveItem implements DELETE /carts/users/{userId}/items/{productId}
func (s *Server) CartsServiceRemoveItem(w http.ResponseWriter, r *http.Request, userId generated.Uuid, productId generated.Uuid) {
	cart := s.store.GetCartByUserId(userId)

	itemIndex := -1
	for i := range cart.Items {
		if cart.Items[i].ProductId == productId {
			itemIndex = i
			break
		}
	}

	if itemIndex < 0 {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Item not found in cart")
		return
	}

	// Remove item from cart
	cart.Items = append(cart.Items[:itemIndex], cart.Items[itemIndex+1:]...)
	cart.UpdatedAt = time.Now()

	s.store.UpdateCart(userId, cart)

	w.WriteHeader(http.StatusNoContent)
}

// CartsServiceClear implements DELETE /carts/users/{userId}/items
func (s *Server) CartsServiceClear(w http.ResponseWriter, r *http.Request, userId generated.Uuid) {
	cart := s.store.GetCartByUserId(userId)
	cart.Items = []generated.CartItem{}
	cart.UpdatedAt = time.Now()

	s.store.UpdateCart(userId, cart)

	w.WriteHeader(http.StatusNoContent)
}