import { Hono } from 'hono'
import type { operations } from '../types/api'
import { store } from '../stores'
import { sendError, ErrorCode } from '../types/errors'

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
    return sendError(c, ErrorCode.NOT_FOUND, 'Product not found')
  }

  if (product.stock < body.quantity) {
    return sendError(c, ErrorCode.INSUFFICIENT_STOCK, 'Insufficient stock')
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
    return sendError(c, ErrorCode.NOT_FOUND, 'Product not found')
  }

  const itemIndex = cart.items.findIndex(item => item.productId === productId)
  if (itemIndex < 0) {
    return sendError(c, ErrorCode.NOT_FOUND, 'Item not found in cart')
  }

  if (product.stock < body.quantity) {
    return sendError(c, ErrorCode.INSUFFICIENT_STOCK, 'Insufficient stock')
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
    return sendError(c, ErrorCode.NOT_FOUND, 'Item not found in cart')
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