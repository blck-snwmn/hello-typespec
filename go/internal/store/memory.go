package store

import (
	"sort"
	"sync"
	"time"

	"github.com/blck-snwmn/hello-typespec/go/generated"
)

// MemoryStore implements the Store interface with in-memory storage
type MemoryStore struct {
	mu         sync.RWMutex
	products   map[string]generated.Product
	categories map[string]generated.Category
	users      map[string]generated.User
	carts      map[string]generated.Cart
	orders     map[string]generated.Order
}

// NewMemoryStore creates a new in-memory store with mock data
func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		products:   make(map[string]generated.Product),
		categories: make(map[string]generated.Category),
		users:      make(map[string]generated.User),
		carts:      make(map[string]generated.Cart),
		orders:     make(map[string]generated.Order),
	}
	store.initializeMockData()
	return store
}

func (s *MemoryStore) initializeMockData() {
	now := time.Now()

	// Categories
	s.categories["1"] = generated.Category{
		Id:        "1",
		Name:      "Electronics",
		ParentId:  nil,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.categories["2"] = generated.Category{
		Id:        "2",
		Name:      "Laptops",
		ParentId:  stringPtr("1"),
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.categories["3"] = generated.Category{
		Id:        "3",
		Name:      "Smartphones",
		ParentId:  stringPtr("1"),
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.categories["4"] = generated.Category{
		Id:        "4",
		Name:      "Clothing",
		ParentId:  nil,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Products
	s.products["1"] = generated.Product{
		Id:          "1",
		Name:        "MacBook Pro 16\"",
		Description: "Apple MacBook Pro with M3 chip",
		Price:       2499.99,
		Stock:       10,
		CategoryId:  "2",
		ImageUrls:   []string{"https://example.com/macbook.jpg"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.products["2"] = generated.Product{
		Id:          "2",
		Name:        "iPhone 15 Pro",
		Description: "Latest iPhone with titanium design",
		Price:       999.99,
		Stock:       25,
		CategoryId:  "3",
		ImageUrls:   []string{"https://example.com/iphone.jpg"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.products["3"] = generated.Product{
		Id:          "3",
		Name:        "T-Shirt",
		Description: "Comfortable cotton t-shirt",
		Price:       29.99,
		Stock:       100,
		CategoryId:  "4",
		ImageUrls:   []string{"https://example.com/tshirt.jpg"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Users
	s.users["1"] = generated.User{
		Id:    "1",
		Email: "user1@example.com",
		Name:  "Test User 1",
		Address: &generated.Address{
			Street:     "123 Test St",
			City:       "Test City",
			State:      "TC",
			PostalCode: "12345",
			Country:    "USA",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.users["2"] = generated.User{
		Id:    "2",
		Email: "user2@example.com",
		Name:  "Test User 2",
		Address: &generated.Address{
			Street:     "456 Demo Ave",
			City:       "Demo City",
			State:      "DC",
			PostalCode: "67890",
			Country:    "USA",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Initialize empty carts for users
	s.carts["1"] = generated.Cart{
		Id:        "cart-1",
		UserId:    "1",
		Items:     []generated.CartItem{},
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.carts["2"] = generated.Cart{
		Id:        "cart-2",
		UserId:    "2",
		Items:     []generated.CartItem{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Products
func (s *MemoryStore) GetProducts() []generated.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]generated.Product, 0, len(s.products))
	for _, product := range s.products {
		products = append(products, product)
	}
	return products
}

func (s *MemoryStore) GetProduct(id string) (*generated.Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product, ok := s.products[id]
	if !ok {
		return nil, false
	}
	return &product, true
}

func (s *MemoryStore) CreateProduct(product generated.Product) generated.Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.products[product.Id] = product
	return product
}

func (s *MemoryStore) UpdateProduct(id string, product generated.Product) generated.Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.products[id] = product
	return product
}

func (s *MemoryStore) DeleteProduct(id string) (*generated.Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, ok := s.products[id]
	if !ok {
		return nil, false
	}
	delete(s.products, id)
	return &product, true
}

// Categories
func (s *MemoryStore) GetCategories() []generated.Category {
	s.mu.RLock()
	defer s.mu.RUnlock()

	categories := make([]generated.Category, 0, len(s.categories))
	for _, category := range s.categories {
		categories = append(categories, category)
	}
	return categories
}

func (s *MemoryStore) GetCategory(id string) (*generated.Category, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	category, ok := s.categories[id]
	if !ok {
		return nil, false
	}
	return &category, true
}

func (s *MemoryStore) CreateCategory(category generated.Category) generated.Category {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.categories[category.Id] = category
	return category
}

func (s *MemoryStore) UpdateCategory(id string, category generated.Category) generated.Category {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.categories[id] = category
	return category
}

func (s *MemoryStore) DeleteCategory(id string) (*generated.Category, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	category, ok := s.categories[id]
	if !ok {
		return nil, false
	}
	delete(s.categories, id)
	return &category, true
}

// Users
func (s *MemoryStore) GetUsers() []generated.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]generated.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	// Sort by ID for consistent ordering
	sort.Slice(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})

	return users
}

func (s *MemoryStore) GetUser(id string) (*generated.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, false
	}
	return &user, true
}

func (s *MemoryStore) CreateUser(user generated.User) generated.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[user.Id] = user
	return user
}

func (s *MemoryStore) UpdateUser(id string, user generated.User) generated.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[id] = user
	return user
}

func (s *MemoryStore) DeleteUser(id string) (*generated.User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]
	if !ok {
		return nil, false
	}
	delete(s.users, id)
	return &user, true
}

// Carts
func (s *MemoryStore) GetCartByUserId(userId string) generated.Cart {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cart, ok := s.carts[userId]
	if !ok {
		// Return a new empty cart if not exists
		now := time.Now()
		return generated.Cart{
			Id:        "cart-" + userId,
			UserId:    userId,
			Items:     []generated.CartItem{},
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return cart
}

func (s *MemoryStore) UpdateCart(userId string, cart generated.Cart) generated.Cart {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.carts[userId] = cart
	return cart
}

// Orders
func (s *MemoryStore) GetOrders() []generated.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]generated.Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	return orders
}

func (s *MemoryStore) GetOrder(id string) (*generated.Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[id]
	if !ok {
		return nil, false
	}
	return &order, true
}

func (s *MemoryStore) GetOrdersByUserId(userId string) []generated.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]generated.Order, 0)
	for _, order := range s.orders {
		if order.UserId == userId {
			orders = append(orders, order)
		}
	}
	return orders
}

func (s *MemoryStore) CreateOrder(order generated.Order) generated.Order {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.Id] = order
	return order
}

func (s *MemoryStore) UpdateOrder(id string, order generated.Order) generated.Order {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[id] = order
	return order
}

// Helper function to create a pointer to a string
func stringPtr(s string) *string {
	return &s
}
