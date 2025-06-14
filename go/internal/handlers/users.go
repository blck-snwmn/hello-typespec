package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/blck-snwmn/hello-typespec/go/generated"
)

// UsersServiceList implements GET /users
func (s *Server) UsersServiceList(w http.ResponseWriter, r *http.Request, params generated.UsersServiceListParams) {
	// Get all users
	allUsers := s.store.GetUsers()

	// Apply pagination
	limit := int32(20) // Default from TypeSpec definition
	if params.Limit != nil {
		limit = *params.Limit
	}

	offset := int32(0)
	if params.Offset != nil {
		offset = *params.Offset
	}

	total := int32(len(allUsers))
	start := int(offset)
	end := int(offset + limit)

	if start > len(allUsers) {
		start = len(allUsers)
	}
	if end > len(allUsers) {
		end = len(allUsers)
	}

	paginatedUsers := allUsers[start:end]

	// Create response
	response := struct {
		Items  []generated.User `json:"items"`
		Total  int32            `json:"total"`
		Limit  int32            `json:"limit"`
		Offset int32            `json:"offset"`
	}{
		Items:  paginatedUsers,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UsersServiceGet implements GET /users/{userId}
func (s *Server) UsersServiceGet(w http.ResponseWriter, r *http.Request, userId generated.Uuid) {
	user, ok := s.store.GetUser(userId)
	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UsersServiceCreate implements POST /users
func (s *Server) UsersServiceCreate(w http.ResponseWriter, r *http.Request) {
	var req generated.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Create new user
	now := time.Now()
	newUser := generated.User{
		Id:        fmt.Sprintf("%d", now.UnixNano()),
		Email:     req.Email,
		Name:      req.Name,
		Address:   req.Address,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created := s.store.CreateUser(newUser)

	// Initialize empty cart for new user
	s.store.UpdateCart(created.Id, generated.Cart{
		Id:        fmt.Sprintf("cart-%s", created.Id),
		UserId:    created.Id,
		Items:     []generated.CartItem{},
		CreatedAt: now,
		UpdatedAt: now,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// UsersServiceUpdate implements PATCH /users/{userId}
func (s *Server) UsersServiceUpdate(w http.ResponseWriter, r *http.Request, userId generated.Uuid) {
	existing, ok := s.store.GetUser(userId)
	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	var req generated.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Update fields if provided
	updatedUser := *existing
	if req.Email != nil {
		updatedUser.Email = *req.Email
	}
	if req.Name != nil {
		updatedUser.Name = *req.Name
	}
	if req.Address != nil {
		updatedUser.Address = req.Address
	}
	updatedUser.UpdatedAt = time.Now()

	updated := s.store.UpdateUser(userId, updatedUser)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// UsersServiceDelete implements DELETE /users/{userId}
func (s *Server) UsersServiceDelete(w http.ResponseWriter, r *http.Request, userId generated.Uuid) {
	_, ok := s.store.DeleteUser(userId)
	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}