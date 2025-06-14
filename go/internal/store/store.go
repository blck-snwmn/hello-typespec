package store

import (
	"github.com/blck-snwmn/hello-typespec/go/generated"
)

// Store defines the interface for data storage operations
type Store interface {
	// Products
	GetProducts() []generated.Product
	GetProduct(id string) (*generated.Product, bool)
	CreateProduct(product generated.Product) generated.Product
	UpdateProduct(id string, product generated.Product) generated.Product
	DeleteProduct(id string) (*generated.Product, bool)

	// Categories
	GetCategories() []generated.Category
	GetCategory(id string) (*generated.Category, bool)
	CreateCategory(category generated.Category) generated.Category
	UpdateCategory(id string, category generated.Category) generated.Category
	DeleteCategory(id string) (*generated.Category, bool)

	// Users
	GetUsers() []generated.User
	GetUser(id string) (*generated.User, bool)
	CreateUser(user generated.User) generated.User
	UpdateUser(id string, user generated.User) generated.User
	DeleteUser(id string) (*generated.User, bool)

	// Carts
	GetCartByUserId(userId string) generated.Cart
	UpdateCart(userId string, cart generated.Cart) generated.Cart

	// Orders
	GetOrders() []generated.Order
	GetOrder(id string) (*generated.Order, bool)
	GetOrdersByUserId(userId string) []generated.Order
	CreateOrder(order generated.Order) generated.Order
	UpdateOrder(id string, order generated.Order) generated.Order
}