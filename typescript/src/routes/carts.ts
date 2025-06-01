import { Hono } from 'hono'
import type { components, operations } from '../types/api'
import { store } from '../stores'

type Cart = components['schemas']['Cart']
type CartAddItemRequest = operations['CartsService_addItem']['requestBody']['content']['application/json']
type CartUpdateItemRequest = operations['CartsService_updateItem']['requestBody']['content']['application/json']

const carts = new Hono()

// GET /carts/users/{userId}
carts.get('/users/:userId', (c) => {
  const userId = c.req.param('userId')
  const cart = store.getCartByUserId(userId)
  return c.json(cart)
})

// POST /carts/users/{userId}/items
carts.post('/users/:userId/items', async (c) => {
  const userId = c.req.param('userId')
  const body = await c.req.json<CartAddItemRequest>()
  
  const cart = store.getCartByUserId(userId)
  const product = store.getProduct(body.productId)

  if (!product) {
    return c.json({ error: { code: 'NOT_FOUND', message: 'Product not found' } }, 404)
  }

  if (product.stock < body.quantity) {
    return c.json({ error: { code: 'INSUFFICIENT_STOCK', message: 'Insufficient stock' } }, 400)
  }

  // Check if item already exists in cart
  const existingItemIndex = cart.items.findIndex(item => item.productId === body.productId)
  
  if (existingItemIndex >= 0) {
    // Update quantity
    cart.items[existingItemIndex].quantity += body.quantity
  } else {
    // Add new item
    cart.items.push({
      productId: body.productId,
      quantity: body.quantity,
      price: product.price,
    })
  }

  cart.updatedAt = new Date().toISOString()
  const updated = store.updateCart(userId, cart)
  return c.json(updated)
})

// PATCH /carts/users/{userId}/items/{productId}
carts.patch('/users/:userId/items/:productId', async (c) => {
  const userId = c.req.param('userId')
  const productId = c.req.param('productId')
  const body = await c.req.json<CartUpdateItemRequest>()
  
  const cart = store.getCartByUserId(userId)
  const product = store.getProduct(productId)

  if (!product) {
    return c.json({ error: { code: 'NOT_FOUND', message: 'Product not found' } }, 404)
  }

  const itemIndex = cart.items.findIndex(item => item.productId === productId)
  if (itemIndex < 0) {
    return c.json({ error: { code: 'NOT_FOUND', message: 'Item not found in cart' } }, 404)
  }

  if (product.stock < body.quantity) {
    return c.json({ error: { code: 'INSUFFICIENT_STOCK', message: 'Insufficient stock' } }, 400)
  }

  cart.items[itemIndex].quantity = body.quantity
  cart.updatedAt = new Date().toISOString()
  
  const updated = store.updateCart(userId, cart)
  return c.json(updated)
})

// DELETE /carts/users/{userId}/items/{productId}
carts.delete('/users/:userId/items/:productId', (c) => {
  const userId = c.req.param('userId')
  const productId = c.req.param('productId')
  
  const cart = store.getCartByUserId(userId)
  const itemIndex = cart.items.findIndex(item => item.productId === productId)
  
  if (itemIndex < 0) {
    return c.json({ error: { code: 'NOT_FOUND', message: 'Item not found in cart' } }, 404)
  }

  cart.items.splice(itemIndex, 1)
  cart.updatedAt = new Date().toISOString()
  
  store.updateCart(userId, cart)
  return c.body(null, 204)
})

// DELETE /carts/users/{userId}/items
carts.delete('/users/:userId/items', (c) => {
  const userId = c.req.param('userId')
  
  const cart = store.getCartByUserId(userId)
  cart.items = []
  cart.updatedAt = new Date().toISOString()
  
  store.updateCart(userId, cart)
  return c.body(null, 204)
})

export default carts