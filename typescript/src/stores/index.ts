import type { components } from '../types/api'

type Product = components['schemas']['Product']
type Category = components['schemas']['Category']
type User = components['schemas']['User']
type Cart = components['schemas']['Cart']
type Order = components['schemas']['Order']
type CartItem = components['schemas']['CartItem']

class DataStore {
  private products: Map<string, Product> = new Map()
  private categories: Map<string, Category> = new Map()
  private users: Map<string, User> = new Map()
  private carts: Map<string, Cart> = new Map()
  private orders: Map<string, Order> = new Map()

  constructor() {
    this.initializeMockData()
  }

  private initializeMockData() {
    // Categories
    this.categories.set('1', {
      id: '1',
      name: 'Electronics',
      parentId: null,
    })
    this.categories.set('2', {
      id: '2',
      name: 'Laptops',
      parentId: '1',
    })
    this.categories.set('3', {
      id: '3',
      name: 'Smartphones',
      parentId: '1',
    })
    this.categories.set('4', {
      id: '4',
      name: 'Clothing',
      parentId: null,
    })

    // Products
    this.products.set('1', {
      id: '1',
      name: 'MacBook Pro 16"',
      description: 'Apple MacBook Pro with M3 chip',
      price: 2499.99,
      stock: 10,
      categoryId: '2',
      imageUrls: ['https://example.com/macbook.jpg'],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    })
    this.products.set('2', {
      id: '2',
      name: 'iPhone 15 Pro',
      description: 'Latest iPhone with titanium design',
      price: 999.99,
      stock: 25,
      categoryId: '3',
      imageUrls: ['https://example.com/iphone.jpg'],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    })
    this.products.set('3', {
      id: '3',
      name: 'T-Shirt',
      description: 'Comfortable cotton t-shirt',
      price: 29.99,
      stock: 100,
      categoryId: '4',
      imageUrls: ['https://example.com/tshirt.jpg'],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    })

    // Users
    this.users.set('1', {
      id: '1',
      email: 'user1@example.com',
      name: 'Test User 1',
      address: '123 Test St, Test City, TC 12345',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    })
    this.users.set('2', {
      id: '2',
      email: 'user2@example.com',
      name: 'Test User 2',
      address: '456 Demo Ave, Demo City, DC 67890',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    })

    // Initialize empty carts for users
    this.carts.set('1', {
      userId: '1',
      items: [],
      updatedAt: new Date().toISOString(),
    })
    this.carts.set('2', {
      userId: '2',
      items: [],
      updatedAt: new Date().toISOString(),
    })
  }

  // Products
  getProducts() { return Array.from(this.products.values()) }
  getProduct(id: string) { return this.products.get(id) }
  createProduct(product: Product) { 
    this.products.set(product.id, product)
    return product
  }
  updateProduct(id: string, product: Product) {
    this.products.set(id, product)
    return product
  }
  deleteProduct(id: string) {
    const product = this.products.get(id)
    this.products.delete(id)
    return product
  }

  // Categories
  getCategories() { return Array.from(this.categories.values()) }
  getCategory(id: string) { return this.categories.get(id) }
  createCategory(category: Category) {
    this.categories.set(category.id, category)
    return category
  }
  updateCategory(id: string, category: Category) {
    this.categories.set(id, category)
    return category
  }
  deleteCategory(id: string) {
    const category = this.categories.get(id)
    this.categories.delete(id)
    return category
  }

  // Users
  getUsers() { return Array.from(this.users.values()) }
  getUser(id: string) { return this.users.get(id) }
  createUser(user: User) {
    this.users.set(user.id, user)
    return user
  }
  updateUser(id: string, user: User) {
    this.users.set(id, user)
    return user
  }
  deleteUser(id: string) {
    const user = this.users.get(id)
    this.users.delete(id)
    return user
  }

  // Carts
  getCartByUserId(userId: string) {
    return this.carts.get(userId) || {
      userId,
      items: [],
      updatedAt: new Date().toISOString(),
    }
  }
  updateCart(userId: string, cart: Cart) {
    this.carts.set(userId, cart)
    return cart
  }

  // Orders
  getOrders() { return Array.from(this.orders.values()) }
  getOrder(id: string) { return this.orders.get(id) }
  getOrdersByUserId(userId: string) {
    return Array.from(this.orders.values()).filter(order => order.userId === userId)
  }
  createOrder(order: Order) {
    this.orders.set(order.id, order)
    return order
  }
  updateOrder(id: string, order: Order) {
    this.orders.set(id, order)
    return order
  }
}

export const store = new DataStore()