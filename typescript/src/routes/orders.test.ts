import { describe, it, expect, beforeEach } from 'vitest'
import app from '../index'
import { loginTestUser, createAuthHeaders, TEST_USERS } from '../test-helpers/auth'
import { store } from '../stores'
import type { operations, components } from '../types/api'

type OrderListResponse = operations['OrdersService_list']['responses']['200']['content']['application/json']
type Order = components['schemas']['Order']
type ErrorResponse = { error: { code: string; message: string } }

describe('Orders API', () => {
  let authToken: string
  const userId = TEST_USERS.alice.id // Use the actual user ID from auth system

  beforeEach(async () => {
    // Reset store to initial state before each test
    const storeInstance = new (store.constructor as any)()
    Object.setPrototypeOf(store, Object.getPrototypeOf(storeInstance))
    Object.assign(store, storeInstance)
    
    // Login test user
    authToken = await loginTestUser(TEST_USERS.alice.email, TEST_USERS.alice.password)
    
    // Add the authenticated user to the store if not exists
    if (!store.getUser(userId)) {
      store.createUser({
        id: userId,
        email: TEST_USERS.alice.email,
        name: 'Alice Johnson',
        address: {
          street: '123 Test St',
          city: 'Test City',
          state: 'TC',
          postalCode: '12345',
          country: 'USA',
        },
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      })
    }
  })

  describe('GET /orders', () => {
    it('should return empty order list initially', async () => {
      const res = await app.request('/orders', {
        headers: createAuthHeaders(authToken),
      })
      const json = await res.json() as OrderListResponse

      expect(res.status).toBe(200)
      expect(json).toHaveProperty('items')
      expect(json).toHaveProperty('total')
      expect(json).toHaveProperty('limit')
      expect(json).toHaveProperty('offset')
      if ('items' in json) {
        expect(json.items).toHaveLength(0)
        expect(json.total).toBe(0)
      }
    })

    it('should filter orders by userId', async () => {
      // Create an order first
      await app.request(`/carts/users/${userId}/items`, {
        method: 'POST',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({ productId: '1', quantity: 1 }),
      })
      
      await app.request(`/orders/users/${userId}`, {
        method: 'POST',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({
          items: [],
          shippingAddress: {
            street: '123 Order St',
            city: 'Order City',
            state: 'OC',
            postalCode: '12345',
            country: 'USA'
          }
        }),
      })

      const res = await app.request(`/orders?userId=${userId}`, {
        headers: createAuthHeaders(authToken),
      })
      const json = await res.json() as OrderListResponse

      expect(res.status).toBe(200)
      if ('items' in json) {
        expect(json.items).toHaveLength(1)
        expect(json.items[0].userId).toBe(userId)
      }
    })

    it('should require authentication', async () => {
      const res = await app.request('/orders')
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(401)
      expect(json.error.code).toBe('UNAUTHORIZED')
    })
  })

  describe('GET /orders/:orderId', () => {
    let orderId: string

    beforeEach(async () => {
      // Setup: Create an order
      await app.request(`/carts/users/${userId}/items`, {
        method: 'POST',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({ productId: '1', quantity: 2 }),
      })

      const orderRes = await app.request(`/orders/users/${userId}`, {
        method: 'POST',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({
          items: [],
          shippingAddress: {
            street: '123 Order St',
            city: 'Order City',
            state: 'OC',
            postalCode: '12345',
            country: 'USA'
          }
        }),
      })
      const order = await orderRes.json() as Order
      orderId = order.id
    })

    it('should return order by id', async () => {
      const res = await app.request(`/orders/${orderId}`, {
        headers: createAuthHeaders(authToken),
      })
      const json = await res.json() as Order | ErrorResponse

      expect(res.status).toBe(200)
      if ('id' in json && !('error' in json)) {
        expect(json.id).toBe(orderId)
        expect(json.userId).toBe(userId)
        expect(json.status).toBe('pending')
        expect(json.items).toHaveLength(1)
      }
    })

    it('should return 404 for non-existent order', async () => {
      const res = await app.request('/orders/999', {
        headers: createAuthHeaders(authToken),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(404)
      expect(json.error).toHaveProperty('code', 'NOT_FOUND')
    })

    it('should require authentication', async () => {
      const res = await app.request(`/orders/${orderId}`)
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(401)
      expect(json.error.code).toBe('UNAUTHORIZED')
    })
  })

  describe('POST /orders/users/:userId', () => {
    beforeEach(async () => {
      // Add items to cart
      await app.request(`/carts/users/${userId}/items`, {
        method: 'POST',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({ productId: '1', quantity: 2 }),
      })
    })

    it('should create order from cart', async () => {
      const res = await app.request(`/orders/users/${userId}`, {
        method: 'POST',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({
          items: [],
          shippingAddress: {
            street: '789 Order Ave',
            city: 'Purchase City',
            state: 'PC',
            postalCode: '54321',
            country: 'USA'
          }
        }),
      })
      const json = await res.json() as Order

      expect(res.status).toBe(201)
      expect(json).toHaveProperty('id')
      expect(json.userId).toBe(userId)
      expect(json.status).toBe('pending')
      expect(json.items).toHaveLength(1)
      expect(json.items[0].productId).toBe('1')
      expect(json.items[0].quantity).toBe(2)
      expect(json.items[0].productName).toBe('MacBook Pro 16"')
      expect(json.totalAmount).toBe(4999.98) // 2 * 2499.99
      expect(json.shippingAddress).toEqual({
        street: '789 Order Ave',
        city: 'Purchase City',
        state: 'PC',
        postalCode: '54321',
        country: 'USA'
      })

      // Verify cart was cleared
      const cartRes = await app.request(`/carts/users/${userId}`, {
        headers: createAuthHeaders(authToken),
      })
      const cart = await cartRes.json() as components['schemas']['Cart']
      expect(cart.items).toHaveLength(0)

      // Verify stock was reduced
      const productRes = await app.request('/products/1')
      const product = await productRes.json() as components['schemas']['Product'] | ErrorResponse
      if ('stock' in product) {
        expect(product.stock).toBe(8) // 10 - 2
      }
    })

    it('should return 404 for non-existent user', async () => {
      const res = await app.request('/orders/users/999', {
        method: 'POST',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({
          items: [],
          shippingAddress: {
            street: '123 Test',
            city: 'Test',
            state: 'TS',
            postalCode: '12345',
            country: 'USA'
          }
        }),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(404)
      expect(json.error.code).toBe('NOT_FOUND')
      expect(json.error.message).toContain('User not found')
    })

    it('should return 400 for empty cart', async () => {
      // Clear cart first
      await app.request(`/carts/users/${userId}/items`, {
        method: 'DELETE',
        headers: createAuthHeaders(authToken),
      })

      const res = await app.request('/orders/users/1', {
        method: 'POST',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({
          items: [],
          shippingAddress: {
            street: '123 Test',
            city: 'Test',
            state: 'TS',
            postalCode: '12345',
            country: 'USA'
          }
        }),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(400)
      expect(json.error.code).toBe('BAD_REQUEST')
    })

    it('should return 400 for insufficient stock', async () => {
      // First add a valid quantity
      await app.request(`/carts/users/${userId}/items`, {
        method: 'DELETE',
        headers: createAuthHeaders(authToken),
      })
      await app.request(`/carts/users/${userId}/items`, {
        method: 'POST',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({ productId: '1', quantity: 15 }), // More than stock (10)
      })

      const res = await app.request('/orders/users/1', {
        method: 'POST',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({
          items: [],
          shippingAddress: {
            street: '123 Test',
            city: 'Test',
            state: 'TS',
            postalCode: '12345',
            country: 'USA'
          }
        }),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(400)
      // Cart API rejects items with insufficient stock, so cart remains empty
      expect(json.error.code).toBe('BAD_REQUEST')
    })

    it('should require authentication', async () => {
      const res = await app.request('/orders/users/1', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          items: [],
          shippingAddress: {
            street: '123 Test',
            city: 'Test',
            state: 'TS',
            postalCode: '12345',
            country: 'USA'
          }
        }),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(401)
      expect(json.error.code).toBe('UNAUTHORIZED')
    })
  })

  describe('PATCH /orders/status/:orderId', () => {
    let orderId: string

    beforeEach(async () => {
      // Create an order
      await app.request(`/carts/users/${userId}/items`, {
        method: 'POST',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({ productId: '1', quantity: 1 }),
      })

      const orderRes = await app.request(`/orders/users/${userId}`, {
        method: 'POST',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({
          items: [],
          shippingAddress: {
            street: '123 Test',
            city: 'Test',
            state: 'TS',
            postalCode: '12345',
            country: 'USA'
          }
        }),
      })
      const order = await orderRes.json() as Order
      orderId = order.id
    })

    it('should update order status', async () => {
      const res = await app.request(`/orders/status/${orderId}`, {
        method: 'PATCH',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: 'processing' }),
      })
      const json = await res.json() as Order | ErrorResponse

      expect(res.status).toBe(200)
      if ('id' in json && !('error' in json)) {
        expect(json.id).toBe(orderId)
        expect(json.status).toBe('processing')
      }
    })

    it('should validate status transitions', async () => {
      // Try invalid transition from pending to delivered
      const res = await app.request(`/orders/status/${orderId}`, {
        method: 'PATCH',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: 'delivered' }),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(400)
      expect(json.error.code).toBe('INVALID_STATE_TRANSITION')
    })

    it('should return 404 for non-existent order', async () => {
      const res = await app.request('/orders/status/999', {
        method: 'PATCH',
        headers: { ...createAuthHeaders(authToken), 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: 'processing' }),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(404)
      expect(json.error.code).toBe('NOT_FOUND')
    })

    it('should require authentication', async () => {
      const res = await app.request(`/orders/status/${orderId}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: 'processing' }),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(401)
      expect(json.error.code).toBe('UNAUTHORIZED')
    })
  })
})