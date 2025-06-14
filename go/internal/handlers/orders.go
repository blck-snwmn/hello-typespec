package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/blck-snwmn/hello-typespec/go/generated"
)

// OrdersServiceList implements GET /orders
func (s *Server) OrdersServiceList(w http.ResponseWriter, r *http.Request, params generated.OrdersServiceListParams) {
	var allOrders []generated.Order

	if params.UserId != nil && *params.UserId != "" {
		allOrders = s.store.GetOrdersByUserId(*params.UserId)
	} else {
		allOrders = s.store.GetOrders()
	}

	// Apply filters
	var filteredOrders []generated.Order
	for _, order := range allOrders {
		// Status filter
		if params.Status != nil && order.Status != *params.Status {
			continue
		}

		// Date filters
		if params.StartDate != nil && order.CreatedAt.Before(*params.StartDate) {
			continue
		}
		if params.EndDate != nil && order.CreatedAt.After(*params.EndDate) {
			continue
		}

		filteredOrders = append(filteredOrders, order)
	}

	// Sort by createdAt descending
	sort.Slice(filteredOrders, func(i, j int) bool {
		return filteredOrders[i].CreatedAt.After(filteredOrders[j].CreatedAt)
	})

	// Apply pagination
	limit := int32(20)
	if params.Limit != nil {
		limit = *params.Limit
	}

	offset := int32(0)
	if params.Offset != nil {
		offset = *params.Offset
	}

	total := int32(len(filteredOrders))
	start := int(offset)
	end := int(offset + limit)

	if start > len(filteredOrders) {
		start = len(filteredOrders)
	}
	if end > len(filteredOrders) {
		end = len(filteredOrders)
	}

	paginatedOrders := filteredOrders[start:end]

	// Create response
	response := struct {
		Items  []generated.Order `json:"items"`
		Total  int32             `json:"total"`
		Limit  int32             `json:"limit"`
		Offset int32             `json:"offset"`
	}{
		Items:  paginatedOrders,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// OrdersServiceListByUser implements GET /orders/users/{userId}
func (s *Server) OrdersServiceListByUser(w http.ResponseWriter, r *http.Request, userId generated.Uuid, params generated.OrdersServiceListByUserParams) {
	allOrders := s.store.GetOrdersByUserId(userId)

	// Sort by createdAt descending
	sort.Slice(allOrders, func(i, j int) bool {
		return allOrders[i].CreatedAt.After(allOrders[j].CreatedAt)
	})

	// Apply pagination
	limit := int32(20)
	if params.Limit != nil {
		limit = *params.Limit
	}

	offset := int32(0)
	if params.Offset != nil {
		offset = *params.Offset
	}

	total := int32(len(allOrders))
	start := int(offset)
	end := int(offset + limit)

	if start > len(allOrders) {
		start = len(allOrders)
	}
	if end > len(allOrders) {
		end = len(allOrders)
	}

	paginatedOrders := allOrders[start:end]

	// Create response
	response := struct {
		Items  []generated.Order `json:"items"`
		Total  int32             `json:"total"`
		Limit  int32             `json:"limit"`
		Offset int32             `json:"offset"`
	}{
		Items:  paginatedOrders,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// OrdersServiceGet implements GET /orders/{orderId}
func (s *Server) OrdersServiceGet(w http.ResponseWriter, r *http.Request, orderId generated.Uuid) {
	order, ok := s.store.GetOrder(orderId)
	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Order not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// OrdersServiceCreate implements POST /orders/users/{userId}
func (s *Server) OrdersServiceCreate(w http.ResponseWriter, r *http.Request, userId generated.Uuid) {
	var req generated.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Validate user exists
	_, ok := s.store.GetUser(userId)
	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	// Validate items from request
	if len(req.Items) == 0 {
		errorResponse(w, http.StatusBadRequest, "EMPTY_CART", "No items in order")
		return
	}

	// Validate stock and calculate total
	var totalAmount float32
	orderItems := []generated.OrderItem{}

	for _, item := range req.Items {
		product, ok := s.store.GetProduct(item.ProductId)
		if !ok {
			errorResponse(w, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("Product %s not found", item.ProductId))
			return
		}
		if product.Stock < item.Quantity {
			errorResponse(w, http.StatusBadRequest, "INSUFFICIENT_STOCK", fmt.Sprintf("Insufficient stock for product %s", product.Name))
			return
		}

		itemPrice := product.Price
		totalAmount += itemPrice * float32(item.Quantity)
		orderItems = append(orderItems, generated.OrderItem{
			ProductId:   item.ProductId,
			Quantity:    item.Quantity,
			Price:       itemPrice,
			ProductName: product.Name,
		})

		// Update product stock
		product.Stock -= item.Quantity
		product.UpdatedAt = time.Now()
		s.store.UpdateProduct(product.Id, *product)
	}

	// Create order
	now := time.Now()
	newOrder := generated.Order{
		Id:              fmt.Sprintf("%d", now.UnixNano()),
		UserId:          userId,
		Items:           orderItems,
		TotalAmount:     totalAmount,
		Status:          generated.Pending,
		ShippingAddress: req.ShippingAddress,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	created := s.store.CreateOrder(newOrder)

	// Clear cart
	clearedCart := s.store.GetCartByUserId(userId)
	clearedCart.Items = []generated.CartItem{}
	clearedCart.UpdatedAt = now
	s.store.UpdateCart(userId, clearedCart)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// OrdersServiceUpdateStatus implements PATCH /orders/{orderId}/status
func (s *Server) OrdersServiceUpdateStatus(w http.ResponseWriter, r *http.Request, orderId generated.Uuid) {
	var req generated.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	order, ok := s.store.GetOrder(orderId)
	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Order not found")
		return
	}

	// Validate status transition
	validTransitions := map[generated.OrderStatus][]generated.OrderStatus{
		generated.Pending:    {generated.Processing, generated.Cancelled},
		generated.Processing: {generated.Shipped, generated.Cancelled},
		generated.Shipped:    {generated.Delivered},
		generated.Delivered:  {},
		generated.Cancelled:  {},
	}

	validNextStatuses := validTransitions[order.Status]
	isValidTransition := false
	for _, validStatus := range validNextStatuses {
		if validStatus == req.Status {
			isValidTransition = true
			break
		}
	}

	if !isValidTransition {
		errorResponse(w, http.StatusBadRequest, "INVALID_STATUS_TRANSITION",
			fmt.Sprintf("Cannot transition from %s to %s", order.Status, req.Status))
		return
	}

	updatedOrder := *order
	updatedOrder.Status = req.Status
	updatedOrder.UpdatedAt = time.Now()

	updated := s.store.UpdateOrder(orderId, updatedOrder)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// OrdersServiceCancel implements POST /orders/{orderId}/cancel
func (s *Server) OrdersServiceCancel(w http.ResponseWriter, r *http.Request, orderId generated.Uuid) {
	order, ok := s.store.GetOrder(orderId)
	if !ok {
		errorResponse(w, http.StatusNotFound, "NOT_FOUND", "Order not found")
		return
	}

	// Check if order can be cancelled (only pending and processing can be cancelled)
	if order.Status != generated.Pending && order.Status != generated.Processing {
		errorResponse(w, http.StatusBadRequest, "INVALID_STATUS",
			fmt.Sprintf("Cannot cancel order with status %s", order.Status))
		return
	}

	// Restore inventory
	for _, item := range order.Items {
		product, ok := s.store.GetProduct(item.ProductId)
		if ok {
			product.Stock += item.Quantity
			product.UpdatedAt = time.Now()
			s.store.UpdateProduct(product.Id, *product)
		}
	}

	updatedOrder := *order
	updatedOrder.Status = generated.Cancelled
	updatedOrder.UpdatedAt = time.Now()

	updated := s.store.UpdateOrder(orderId, updatedOrder)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}
