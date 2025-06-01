import { Hono } from 'hono'
import type { components, operations } from '../types/api'
import { store } from '../stores'

type Product = components['schemas']['Product']
type ProductListResponse = operations['ProductsService_list']['responses']['200']['content']['application/json']
type ProductCreateRequest = operations['ProductsService_create']['requestBody']['content']['application/json']
type ProductUpdateRequest = operations['ProductsService_update']['requestBody']['content']['application/json']

const products = new Hono()

// GET /products
products.get('/', (c) => {
  const limit = parseInt(c.req.query('limit') || '10')
  const offset = parseInt(c.req.query('offset') || '0')
  const search = c.req.query('search')
  const categoryId = c.req.query('categoryId')
  const minPrice = c.req.query('minPrice') ? parseFloat(c.req.query('minPrice')!) : undefined
  const maxPrice = c.req.query('maxPrice') ? parseFloat(c.req.query('maxPrice')!) : undefined

  let allProducts = store.getProducts()

  // Apply filters
  if (search) {
    allProducts = allProducts.filter(p => 
      p.name.toLowerCase().includes(search.toLowerCase()) ||
      p.description.toLowerCase().includes(search.toLowerCase())
    )
  }
  if (categoryId) {
    allProducts = allProducts.filter(p => p.categoryId === categoryId)
  }
  if (minPrice !== undefined) {
    allProducts = allProducts.filter(p => p.price >= minPrice)
  }
  if (maxPrice !== undefined) {
    allProducts = allProducts.filter(p => p.price <= maxPrice)
  }

  // Apply pagination
  const paginatedProducts = allProducts.slice(offset, offset + limit)

  const response: ProductListResponse = {
    products: paginatedProducts,
    total: allProducts.length,
    limit,
    offset,
  }

  return c.json(response)
})

// GET /products/{productId}
products.get('/:productId', (c) => {
  const productId = c.req.param('productId')
  const product = store.getProduct(productId)

  if (!product) {
    return c.json({ error: { code: 'NOT_FOUND', message: 'Product not found' } }, 404)
  }

  return c.json(product)
})

// POST /products
products.post('/', async (c) => {
  const body = await c.req.json<ProductCreateRequest>()
  
  const newProduct: Product = {
    id: Date.now().toString(),
    ...body,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  }

  const created = store.createProduct(newProduct)
  return c.json(created, 201)
})

// PUT /products/{productId}
products.put('/:productId', async (c) => {
  const productId = c.req.param('productId')
  const body = await c.req.json<ProductUpdateRequest>()
  
  const existing = store.getProduct(productId)
  if (!existing) {
    return c.json({ error: { code: 'NOT_FOUND', message: 'Product not found' } }, 404)
  }

  const updatedProduct: Product = {
    ...existing,
    ...body,
    id: productId,
    updatedAt: new Date().toISOString(),
  }

  const updated = store.updateProduct(productId, updatedProduct)
  return c.json(updated)
})

// DELETE /products/{productId}
products.delete('/:productId', (c) => {
  const productId = c.req.param('productId')
  const deleted = store.deleteProduct(productId)

  if (!deleted) {
    return c.json({ error: { code: 'NOT_FOUND', message: 'Product not found' } }, 404)
  }

  return c.body(null, 204)
})

export default products