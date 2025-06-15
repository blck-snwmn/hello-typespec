import { describe, it, expect, beforeEach } from 'vitest'
import app from '../index'
import { store } from '../stores'
import { loginTestUser, createAuthHeaders, TEST_USERS } from '../test-helpers/auth'
import type { components } from '../types/api'

type Cart = components['schemas']['Cart']
type ErrorResponse = { error: { code: string; message: string } }

describe('Carts API', () => {
  let authToken: string

  beforeEach(async () => {
    // Reset store to initial state before each test
    const storeInstance = new (store.constructor as any)()
    Object.setPrototypeOf(store, Object.getPrototypeOf(storeInstance))
    Object.assign(store, storeInstance)
    
    // Login to get auth token
    authToken = await loginTestUser(TEST_USERS.alice.email, TEST_USERS.alice.password)
  })

  describe('GET /carts/users/:userId', () => {
    it('should return cart for existing user', async () => {
      const res = await app.request('/carts/users/1', {
        headers: createAuthHeaders(authToken)
      })
      const json = await res.json() as Cart

      expect(res.status).toBe(200)
      expect(json).toHaveProperty('userId', '1')
      expect(json).toHaveProperty('items')
      expect(json.items).toHaveLength(0) // Initially empty
      expect(json).toHaveProperty('updatedAt')
    })

    it('should create empty cart for new user', async () => {
      const res = await app.request('/carts/users/999', {
        headers: createAuthHeaders(authToken)
      })
      const json = await res.json() as Cart

      expect(res.status).toBe(200)
      expect(json).toHaveProperty('userId', '999')
      expect(json.items).toHaveLength(0)
    })
  })

  describe('POST /carts/users/:userId/items', () => {
    it('should add item to cart', async () => {
      const addItemRequest = {
        productId: '1',
        quantity: 2
      }

      const res = await app.request('/carts/users/1/items', {
        method: 'POST',
        headers: createAuthHeaders(authToken),
        body: JSON.stringify(addItemRequest),
      })
      const json = await res.json() as Cart

      expect(res.status).toBe(200)
      expect(json.items).toHaveLength(1)
      expect(json.items[0]).toEqual({
        productId: '1',
        quantity: 2
      })
    })

    it('should increase quantity when adding existing item', async () => {
      // First, add item
      await app.request('/carts/users/1/items', {
        method: 'POST',
        headers: createAuthHeaders(authToken),
        body: JSON.stringify({ productId: '1', quantity: 2 }),
      })

      // Add same item again
      const res = await app.request('/carts/users/1/items', {
        method: 'POST',
        headers: createAuthHeaders(authToken),
        body: JSON.stringify({ productId: '1', quantity: 3 }),
      })
      const json = await res.json() as Cart

      expect(res.status).toBe(200)
      expect(json.items).toHaveLength(1)
      expect(json.items[0].quantity).toBe(5) // 2 + 3
    })

    it('should return 404 for non-existent product', async () => {
      const addItemRequest = {
        productId: '999',
        quantity: 1
      }

      const res = await app.request('/carts/users/1/items', {
        method: 'POST',
        headers: createAuthHeaders(authToken),
        body: JSON.stringify(addItemRequest),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(404)
      expect(json.error).toHaveProperty('code', 'NOT_FOUND')
    })

    it('should return 400 for insufficient stock', async () => {
      const addItemRequest = {
        productId: '1',
        quantity: 100 // More than available stock (10)
      }

      const res = await app.request('/carts/users/1/items', {
        method: 'POST',
        headers: createAuthHeaders(authToken),
        body: JSON.stringify(addItemRequest),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(400)
      expect(json.error).toHaveProperty('code', 'INSUFFICIENT_STOCK')
    })
  })

  describe('PATCH /carts/users/:userId/items/:productId', () => {
    beforeEach(async () => {
      // Add item to cart first
      await app.request('/carts/users/1/items', {
        method: 'POST',
        headers: createAuthHeaders(authToken),
        body: JSON.stringify({ productId: '1', quantity: 2 }),
      })
    })

    it('should update item quantity', async () => {
      const updateRequest = {
        quantity: 5
      }

      const res = await app.request('/carts/users/1/items/1', {
        method: 'PATCH',
        headers: createAuthHeaders(authToken),
        body: JSON.stringify(updateRequest),
      })
      const json = await res.json() as Cart

      expect(res.status).toBe(200)
      expect(json.items[0].quantity).toBe(5)
    })

    it('should return 404 for item not in cart', async () => {
      const updateRequest = {
        quantity: 5
      }

      const res = await app.request('/carts/users/1/items/999', {
        method: 'PATCH',
        headers: createAuthHeaders(authToken),
        body: JSON.stringify(updateRequest),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(404)
      expect(json.error).toHaveProperty('code', 'NOT_FOUND')
      expect(json.error.message).toContain('not found')
    })

    it('should return 400 for insufficient stock', async () => {
      const updateRequest = {
        quantity: 100 // More than available stock
      }

      const res = await app.request('/carts/users/1/items/1', {
        method: 'PATCH',
        headers: createAuthHeaders(authToken),
        body: JSON.stringify(updateRequest),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(400)
      expect(json.error).toHaveProperty('code', 'INSUFFICIENT_STOCK')
    })
  })

  describe('DELETE /carts/users/:userId/items/:productId', () => {
    beforeEach(async () => {
      // Add items to cart first
      await app.request('/carts/users/1/items', {
        method: 'POST',
        headers: createAuthHeaders(authToken),
        body: JSON.stringify({ productId: '1', quantity: 2 }),
      })
      await app.request('/carts/users/1/items', {
        method: 'POST',
        headers: createAuthHeaders(authToken),
        body: JSON.stringify({ productId: '2', quantity: 1 }),
      })
    })

    it('should remove item from cart', async () => {
      const res = await app.request('/carts/users/1/items/1', {
        method: 'DELETE',
        headers: createAuthHeaders(authToken)
      })

      expect(res.status).toBe(204)

      // Verify item was removed
      const cartRes = await app.request('/carts/users/1', { headers: createAuthHeaders(authToken) })
      const cart = await cartRes.json() as Cart
      expect(cart.items).toHaveLength(1)
      expect(cart.items[0].productId).toBe('2')
    })

    it('should return 404 for item not in cart', async () => {
      const res = await app.request('/carts/users/1/items/999', {
        method: 'DELETE',
        headers: createAuthHeaders(authToken)
      })

      expect(res.status).toBe(404)
    })
  })

  describe('DELETE /carts/users/:userId/items', () => {
    beforeEach(async () => {
      // Add items to cart first
      await app.request('/carts/users/1/items', {
        method: 'POST',
        headers: createAuthHeaders(authToken),
        body: JSON.stringify({ productId: '1', quantity: 2 }),
      })
      await app.request('/carts/users/1/items', {
        method: 'POST',
        headers: createAuthHeaders(authToken),
        body: JSON.stringify({ productId: '2', quantity: 1 }),
      })
    })

    it('should clear all items from cart', async () => {
      const res = await app.request('/carts/users/1/items', {
        method: 'DELETE',
        headers: createAuthHeaders(authToken)
      })

      expect(res.status).toBe(204)

      // Verify cart is empty
      const cartRes = await app.request('/carts/users/1', { headers: createAuthHeaders(authToken) })
      const cart = await cartRes.json() as Cart
      expect(cart.items).toHaveLength(0)
    })
  })

  describe('Authentication', () => {
    it('should return 401 for unauthenticated requests', async () => {
      const res = await app.request('/carts/users/1')
      
      expect(res.status).toBe(401)
      const json = await res.json() as ErrorResponse
      expect(json.error.code).toBe('UNAUTHORIZED')
    })

    it('should return 401 for invalid token', async () => {
      const res = await app.request('/carts/users/1', {
        headers: createAuthHeaders('invalid-token')
      })
      
      expect(res.status).toBe(401)
      const json = await res.json() as ErrorResponse
      expect(json.error.code).toBe('UNAUTHORIZED')
    })
  })
})