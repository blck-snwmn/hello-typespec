import { describe, it, expect, beforeEach } from 'vitest'
import app from '../index'
import { store } from '../stores'
import type { components } from '../types/api'

type CategoryListResponse = components['schemas']['Category'][]
type Category = components['schemas']['Category']
type CategoryWithChildren = Category & { children: CategoryWithChildren[] }
type CategoryTreeResponse = CategoryWithChildren[]
type ErrorResponse = { error: { code: string; message: string } }

describe('Categories API', () => {
  beforeEach(() => {
    // Reset store to initial state before each test
    const storeInstance = new (store.constructor as any)()
    Object.setPrototypeOf(store, Object.getPrototypeOf(storeInstance))
    Object.assign(store, storeInstance)
  })

  describe('GET /categories', () => {
    it('should return all categories', async () => {
      const res = await app.request('/categories')
      const json = await res.json() as CategoryListResponse

      expect(res.status).toBe(200)
      expect(Array.isArray(json)).toBe(true)
      expect(json).toHaveLength(4)
      expect(json[0]).toHaveProperty('id')
      expect(json[0]).toHaveProperty('name')
      // parentId can be undefined for root categories
    })
  })

  describe('GET /categories/tree', () => {
    it('should return categories in tree structure', async () => {
      const res = await app.request('/categories/tree')
      const json = await res.json() as CategoryTreeResponse

      expect(res.status).toBe(200)
      expect(Array.isArray(json)).toBe(true)
      
      // Root categories
      const rootCategories = json.filter(cat => !cat.parentId)
      expect(rootCategories).toHaveLength(2) // Electronics and Clothing
      
      // Check Electronics has children
      const electronics = json.find(cat => cat.name === 'Electronics')
      expect(electronics).toBeDefined()
      expect(electronics!.children).toHaveLength(2) // Laptops and Smartphones
      
      // Check nested structure
      const laptops = electronics!.children.find(child => child.name === 'Laptops')
      expect(laptops).toBeDefined()
      expect(laptops!.parentId).toBe(electronics!.id)
    })
  })

  describe('GET /categories/:id', () => {
    it('should return a category by id', async () => {
      const res = await app.request('/categories/1')
      const json = await res.json() as Category | ErrorResponse

      expect(res.status).toBe(200)
      if ('id' in json) {
        expect(json.id).toBe('1')
        expect(json.name).toBe('Electronics')
        expect(json.parentId).toBeUndefined()
      }
    })

    it('should return 404 for non-existent category', async () => {
      const res = await app.request('/categories/999')
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(404)
      expect(json.error).toHaveProperty('code', 'NOT_FOUND')
      expect(json.error).toHaveProperty('message')
    })
  })

  describe('POST /categories', () => {
    it('should create a new root category', async () => {
      const newCategory = {
        name: 'Books',
        parentId: undefined
      }

      const res = await app.request('/categories', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newCategory),
      })
      const json = await res.json() as Category

      expect(res.status).toBe(201)
      expect(json).toHaveProperty('id')
      expect(json.name).toBe(newCategory.name)
      expect(json.parentId).toBeUndefined()
    })

    it('should create a child category', async () => {
      const newCategory = {
        name: 'Gaming Laptops',
        parentId: '2' // Laptops category
      }

      const res = await app.request('/categories', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newCategory),
      })
      const json = await res.json() as Category

      expect(res.status).toBe(201)
      expect(json).toHaveProperty('id')
      expect(json.name).toBe(newCategory.name)
      expect(json.parentId).toBe(newCategory.parentId)
    })
  })

  describe('PUT /categories/:id', () => {
    it('should update a category', async () => {
      const updateData = {
        name: 'Consumer Electronics'
      }

      const res = await app.request('/categories/1', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updateData),
      })
      const json = await res.json() as Category | ErrorResponse

      expect(res.status).toBe(200)
      if ('id' in json) {
        expect(json.id).toBe('1')
        expect(json.name).toBe(updateData.name)
      }
    })

    it('should return 404 when updating non-existent category', async () => {
      const updateData = {
        name: 'Non-existent'
      }

      const res = await app.request('/categories/999', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updateData),
      })
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(404)
      expect(json.error).toHaveProperty('code', 'NOT_FOUND')
    })
  })

  describe('DELETE /categories/:id', () => {
    it('should delete a category', async () => {
      const res = await app.request('/categories/4', {
        method: 'DELETE',
      })

      expect(res.status).toBe(204)
      expect(await res.text()).toBe('')

      // Verify deletion
      const getRes = await app.request('/categories/4')
      expect(getRes.status).toBe(404)
    })

    it('should return 404 when deleting non-existent category', async () => {
      const res = await app.request('/categories/999', {
        method: 'DELETE',
      })

      expect(res.status).toBe(404)
    })
  })
})