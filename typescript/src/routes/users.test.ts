import { describe, it, expect, beforeEach } from 'vitest'
import app from '../index'
import { store } from '../stores'
import type { operations, components } from '../types/api'

type UserListResponse = components['schemas']['User'][]
type User = components['schemas']['User']
type ErrorResponse = { error: { code: string; message: string } }

describe('Users API', () => {
  beforeEach(() => {
    // Reset store to initial state before each test
    const storeInstance = new (store.constructor as any)()
    Object.setPrototypeOf(store, Object.getPrototypeOf(storeInstance))
    Object.assign(store, storeInstance)
  })

  describe('GET /users', () => {
    it('should return all users', async () => {
      const res = await app.request('/users')
      const json = await res.json() as UserListResponse

      expect(res.status).toBe(200)
      expect(Array.isArray(json)).toBe(true)
      expect(json).toHaveLength(2)
      expect(json[0]).toHaveProperty('id')
      expect(json[0]).toHaveProperty('email')
      expect(json[0]).toHaveProperty('name')
      expect(json[0]).toHaveProperty('address')
    })
  })

  describe('GET /users/:id', () => {
    it('should return a user by id', async () => {
      const res = await app.request('/users/1')
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

    it('should return 404 for non-existent user', async () => {
      const res = await app.request('/users/999')
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(404)
      expect(json.error).toHaveProperty('code', 'NOT_FOUND')
      expect(json.error).toHaveProperty('message')
    })
  })

  describe('POST /users', () => {
    it('should create a new user with address', async () => {
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
        headers: { 'Content-Type': 'application/json' },
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
      const cartRes = await app.request(`/carts/users/${json.id}`)
      const cart = await cartRes.json()
      expect(cart).toHaveProperty('userId', json.id)
      expect(cart.items).toHaveLength(0)
    })

    it('should create a new user without address', async () => {
      const newUser = {
        email: 'minimal@example.com',
        name: 'Minimal User'
      }

      const res = await app.request('/users', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newUser),
      })
      const json = await res.json() as User

      expect(res.status).toBe(201)
      expect(json).toHaveProperty('id')
      expect(json.email).toBe(newUser.email)
      expect(json.name).toBe(newUser.name)
      expect(json.address).toBeUndefined()
    })
  })

  describe('PUT /users/:id', () => {
    it('should update user details', async () => {
      const updateData = {
        name: 'Updated User Name',
        email: 'updated@example.com'
      }

      const res = await app.request('/users/1', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
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

    it('should update user address', async () => {
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
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updateData),
      })
      const json = await res.json() as User | ErrorResponse

      expect(res.status).toBe(200)
      if ('id' in json) {
        expect(json.address).toEqual(updateData.address)
      }
    })

    it('should return 404 when updating non-existent user', async () => {
      const updateData = {
        name: 'Non-existent'
      }

      const res = await app.request('/users/999', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updateData),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(404)
      expect(json.error).toHaveProperty('code', 'NOT_FOUND')
    })
  })

  describe('DELETE /users/:id', () => {
    it('should delete a user', async () => {
      const res = await app.request('/users/2', {
        method: 'DELETE',
      })

      expect(res.status).toBe(204)
      expect(await res.text()).toBe('')

      // Verify deletion
      const getRes = await app.request('/users/2')
      expect(getRes.status).toBe(404)
    })

    it('should return 404 when deleting non-existent user', async () => {
      const res = await app.request('/users/999', {
        method: 'DELETE',
      })

      expect(res.status).toBe(404)
    })
  })
})