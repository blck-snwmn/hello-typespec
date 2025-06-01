import { describe, it, expect, beforeEach } from 'vitest'
import app from '../index'
import { store } from '../stores'
import type { operations } from '../types/api'

type ProductListResponse = operations['ProductsService_list']['responses']['200']['content']['application/json']
type Product = operations['ProductsService_get']['responses']['200']['content']['application/json']
type ErrorResponse = { error: { code: string; message: string } }

describe('Products API', () => {
  beforeEach(() => {
    // Reset store to initial state before each test
    const storeInstance = new (store.constructor as any)()
    Object.setPrototypeOf(store, Object.getPrototypeOf(storeInstance))
    Object.assign(store, storeInstance)
  })

  describe('GET /products', () => {
    it('should return all products with default pagination', async () => {
      const res = await app.request('/products')
      const json = await res.json() as ProductListResponse

      expect(res.status).toBe(200)
      expect(json).toHaveProperty('items')
      expect(json).toHaveProperty('total')
      expect(json).toHaveProperty('limit')
      expect(json).toHaveProperty('offset')
      if ('items' in json) {
        expect(json.items).toHaveLength(3)
        expect(json.total).toBe(3)
        expect(json.limit).toBe(10)
        expect(json.offset).toBe(0)
      }
    })

    it('should handle pagination correctly', async () => {
      const res = await app.request('/products?limit=2&offset=1')
      const json = await res.json() as ProductListResponse

      expect(res.status).toBe(200)
      if ('items' in json) {
        expect(json.items).toHaveLength(2)
        expect(json.limit).toBe(2)
        expect(json.offset).toBe(1)
        expect(json.items[0].id).toBe('2')
        expect(json.items[1].id).toBe('3')
      }
    })

    it('should filter by search term', async () => {
      const res = await app.request('/products?search=iPhone')
      const json = await res.json() as ProductListResponse

      expect(res.status).toBe(200)
      if ('items' in json) {
        expect(json.items).toHaveLength(1)
        expect(json.items[0].name).toContain('iPhone')
      }
    })

    it('should filter by price range', async () => {
      const res = await app.request('/products?minPrice=100&maxPrice=1000')
      const json = await res.json() as ProductListResponse

      expect(res.status).toBe(200)
      if ('items' in json) {
        expect(json.items).toHaveLength(1)
        expect(json.items[0].price).toBeGreaterThanOrEqual(100)
        expect(json.items[0].price).toBeLessThanOrEqual(1000)
      }
    })
  })

  describe('GET /products/:id', () => {
    it('should return a product by id', async () => {
      const res = await app.request('/products/1')
      const json = await res.json() as Product

      expect(res.status).toBe(200)
      if ('id' in json) {
        expect(json.id).toBe('1')
        expect(json.name).toBe('MacBook Pro 16"')
      }
    })

    it('should return 404 for non-existent product', async () => {
      const res = await app.request('/products/999')
      const json = await res.json() as ErrorResponse

      expect(res.status).toBe(404)
      expect(json.error).toHaveProperty('code', 'NOT_FOUND')
      expect(json.error).toHaveProperty('message')
    })
  })

  describe('POST /products', () => {
    it('should create a new product', async () => {
      const newProduct = {
        name: 'Test Product',
        description: 'Test description',
        price: 99.99,
        stock: 50,
        categoryId: '1',
        imageUrls: ['https://example.com/test.jpg']
      }

      const res = await app.request('/products', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newProduct),
      })
      const json = await res.json() as Product

      expect(res.status).toBe(201)
      expect(json).toHaveProperty('id')
      if ('name' in json) {
        expect(json.name).toBe(newProduct.name)
        expect(json.price).toBe(newProduct.price)
      }
      expect(json).toHaveProperty('createdAt')
      expect(json).toHaveProperty('updatedAt')
    })
  })
})