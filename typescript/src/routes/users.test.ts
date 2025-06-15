import { describe, it, expect, beforeEach } from 'vitest'
import app from '../index'
import { store } from '../stores'
import type { operations, components } from '../types/api'
import { loginTestUser, createAuthHeaders, TEST_USERS } from '../test-helpers/auth'

type User = components['schemas']['User']
type Cart = components['schemas']['Cart']
type ErrorResponse = components['schemas']['ErrorResponse']
type UserListResponse = {
  items: User[];
  total: number;
  limit: number;
  offset: number;
}

describe('Users API', () => {
  let authToken: string

  beforeEach(async () => {
    // Reset store to initial state before each test
    const storeInstance = new (store.constructor as any)()
    Object.setPrototypeOf(store, Object.getPrototypeOf(storeInstance))
    Object.assign(store, storeInstance)

    // Login to get auth token
    authToken = await loginTestUser(TEST_USERS.alice.email, TEST_USERS.alice.password)
  })

  describe('GET /users', () => {
    it('should return all users with authentication', async () => {
      const res = await app.request('/users', {
        headers: createAuthHeaders(authToken),
      })
      const json = await res.json() as UserListResponse

      expect(res.status).toBe(200)
      expect(json).toHaveProperty('items')
      expect(json).toHaveProperty('total')
      expect(json).toHaveProperty('limit')
      expect(json).toHaveProperty('offset')
      expect(Array.isArray(json.items)).toBe(true)
      expect(json.items).toHaveLength(2)
      expect(json.total).toBe(2)
      expect(json.limit).toBe(20)
      expect(json.offset).toBe(0)
      expect(json.items[0]).toHaveProperty('id')
      expect(json.items[0]).toHaveProperty('email')
      expect(json.items[0]).toHaveProperty('name')
      expect(json.items[0]).toHaveProperty('address')
    })

    it('should return 401 without authentication', async () => {
      const res = await app.request('/users')
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(401)
      expect(json.error).toHaveProperty('code', 'UNAUTHORIZED')
      expect(json.error).toHaveProperty('message')
    })

    it('should support pagination with authentication', async () => {
      const res1 = await app.request('/users?limit=1&offset=0', {
        headers: createAuthHeaders(authToken),
      })
      const json1 = await res1.json() as UserListResponse

      expect(res1.status).toBe(200)
      expect(json1.items).toHaveLength(1)
      expect(json1.total).toBe(2)
      expect(json1.limit).toBe(1)
      expect(json1.offset).toBe(0)

      const res2 = await app.request('/users?limit=1&offset=1', {
        headers: createAuthHeaders(authToken),
      })
      const json2 = await res2.json() as UserListResponse

      expect(res2.status).toBe(200)
      expect(json2.items).toHaveLength(1)
      expect(json2.total).toBe(2)
      expect(json2.limit).toBe(1)
      expect(json2.offset).toBe(1)

      // Ensure different users
      expect(json1.items[0].id).not.toBe(json2.items[0].id)
    })
  })

  describe('GET /users/:id', () => {
    it('should return a user by id with authentication', async () => {
      const res = await app.request('/users/1', {
        headers: createAuthHeaders(authToken),
      })
      const json = await res.json() as User | ErrorResponse

      expect(res.status).toBe(200)
      if ('id' in json) {
        expect(json.id).toBe('1')
        expect(json.email).toBe('user1@example.com')
        expect(json.name).toBe('Test User 1')
        expect(json.address).toEqual({
          street: '123 Test St',
          city: 'Test City',
          state: 'TC',
          postalCode: '12345',
          country: 'USA',
        })
      }
    })

    it('should return 401 without authentication', async () => {
      const res = await app.request('/users/1')
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(401)
      expect(json.error).toHaveProperty('code', 'UNAUTHORIZED')
    })

    it('should return 404 for non-existent user with authentication', async () => {
      const res = await app.request('/users/999', {
        headers: createAuthHeaders(authToken),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(404)
      expect(json.error).toHaveProperty('code', 'NOT_FOUND')
      expect(json.error).toHaveProperty('message')
    })
  })

  describe('POST /users', () => {
    it('should create a new user with address and authentication', async () => {
      const newUser = {
        email: 'newuser@example.com',
        name: 'New User',
        address: {
          street: '789 New St',
          city: 'New City',
          state: 'NC',
          postalCode: '11111',
          country: 'USA'
        }
      }

      const res = await app.request('/users', {
        method: 'POST',
        headers: {
          ...createAuthHeaders(authToken),
        },
        body: JSON.stringify(newUser),
      })
      const json = await res.json() as User

      expect(res.status).toBe(201)
      expect(json).toHaveProperty('id')
      expect(json.email).toBe(newUser.email)
      expect(json.name).toBe(newUser.name)
      expect(json.address).toEqual(newUser.address)
      expect(json).toHaveProperty('createdAt')
      expect(json).toHaveProperty('updatedAt')

      // Verify cart was created for new user
      const cartRes = await app.request(`/carts/users/${json.id}`, {
        headers: createAuthHeaders(authToken),
      })
      const cart = await cartRes.json() as Cart
      expect(cart).toHaveProperty('userId', json.id)
      expect(cart.items).toHaveLength(0)
    })

    it('should create a new user without address with authentication', async () => {
      const newUser = {
        email: 'minimal@example.com',
        name: 'Minimal User'
      }

      const res = await app.request('/users', {
        method: 'POST',
        headers: {
          ...createAuthHeaders(authToken),
        },
        body: JSON.stringify(newUser),
      })
      const json = await res.json() as User

      expect(res.status).toBe(201)
      expect(json).toHaveProperty('id')
      expect(json.email).toBe(newUser.email)
      expect(json.name).toBe(newUser.name)
      expect(json.address).toBeUndefined()
    })

    it('should return 401 without authentication', async () => {
      const newUser = {
        email: 'unauth@example.com',
        name: 'Unauthorized User'
      }

      const res = await app.request('/users', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newUser),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(401)
      expect(json.error).toHaveProperty('code', 'UNAUTHORIZED')
    })
  })

  describe('PUT /users/:id', () => {
    it('should update user details with authentication', async () => {
      const updateData = {
        name: 'Updated User Name',
        email: 'updated@example.com'
      }

      const res = await app.request('/users/1', {
        method: 'PUT',
        headers: {
          ...createAuthHeaders(authToken),
        },
        body: JSON.stringify(updateData),
      })
      const json = await res.json() as User | ErrorResponse

      expect(res.status).toBe(200)
      if ('id' in json) {
        expect(json.id).toBe('1')
        expect(json.name).toBe(updateData.name)
        expect(json.email).toBe(updateData.email)
        expect(json).toHaveProperty('updatedAt')
      }
    })

    it('should update user address with authentication', async () => {
      const updateData = {
        address: {
          street: '999 Updated Ave',
          city: 'Update City',
          state: 'UC',
          postalCode: '99999',
          country: 'USA'
        }
      }

      const res = await app.request('/users/1', {
        method: 'PUT',
        headers: {
          ...createAuthHeaders(authToken),
        },
        body: JSON.stringify(updateData),
      })
      const json = await res.json() as User | ErrorResponse

      expect(res.status).toBe(200)
      if ('id' in json) {
        expect(json.address).toEqual(updateData.address)
      }
    })

    it('should return 401 without authentication', async () => {
      const updateData = {
        name: 'Unauthorized Update'
      }

      const res = await app.request('/users/1', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updateData),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(401)
      expect(json.error).toHaveProperty('code', 'UNAUTHORIZED')
    })

    it('should return 404 when updating non-existent user with authentication', async () => {
      const updateData = {
        name: 'Non-existent'
      }

      const res = await app.request('/users/999', {
        method: 'PUT',
        headers: {
          ...createAuthHeaders(authToken),
        },
        body: JSON.stringify(updateData),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(404)
      expect(json.error).toHaveProperty('code', 'NOT_FOUND')
    })
  })

  describe('DELETE /users/:id', () => {
    it('should delete a user with authentication', async () => {
      const res = await app.request('/users/2', {
        method: 'DELETE',
        headers: createAuthHeaders(authToken),
      })

      expect(res.status).toBe(204)
      expect(await res.text()).toBe('')

      // Verify deletion
      const getRes = await app.request('/users/2', {
        headers: createAuthHeaders(authToken),
      })
      expect(getRes.status).toBe(404)
    })

    it('should return 401 without authentication', async () => {
      const res = await app.request('/users/2', {
        method: 'DELETE',
      })

      expect(res.status).toBe(401)
      const json = await res.json() as ErrorResponse
      expect(json.error).toHaveProperty('code', 'UNAUTHORIZED')
    })

    it('should return 404 when deleting non-existent user with authentication', async () => {
      const res = await app.request('/users/999', {
        method: 'DELETE',
        headers: createAuthHeaders(authToken),
      })

      expect(res.status).toBe(404)
    })
  })
})