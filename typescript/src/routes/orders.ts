import { Hono } from 'hono'
import type { components, operations } from '../types/api'
import { store } from '../stores'

type Order = components['schemas']['Order']
type OrderStatus = components['schemas']['OrderStatus']
type OrderCreateRequest = operations['OrdersService_create']['requestBody']['content']['application/json']
type OrderListResponse = operations['OrdersService_list']['responses']['200']['content']['application/json']
type OrderUpdateStatusRequest = operations['OrdersService_updateStatus']['requestBody']['content']['application/json']

const orders = new Hono()

// GET /orders
orders.get('/', (c) => {
  const limit = parseInt(c.req.query('limit') || '10')
  const offset = parseInt(c.req.query('offset') || '0')
  const userId = c.req.query('userId')
  const status = c.req.query('status') as OrderStatus | undefined

  let allOrders = userId ? store.getOrdersByUserId(userId) : store.getOrders()

  // Apply filters
  if (status) {
    allOrders = allOrders.filter(order => order.status === status)
  }

  // Sort by createdAt descending
  allOrders.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())

  // Apply pagination
  const paginatedOrders = allOrders.slice(offset, offset + limit)

  const response: OrderListResponse = {
    orders: paginatedOrders,
    total: allOrders.length,
    limit,
    offset,
  }

  return c.json(response)
})

// GET /orders/{orderId}
orders.get('/:orderId', (c) => {
  const orderId = c.req.param('orderId')
  const order = store.getOrder(orderId)

  if (!order) {
    return c.json({ error: { code: 'NOT_FOUND', message: 'Order not found' } }, 404)
  }

  return c.json(order)
})

// POST /orders
orders.post('/', async (c) => {
  const body = await c.req.json<OrderCreateRequest>()
  
  // Validate user exists
  const user = store.getUser(body.userId)
  if (!user) {
    return c.json({ error: { code: 'NOT_FOUND', message: 'User not found' } }, 404)
  }

  // Get user's cart
  const cart = store.getCartByUserId(body.userId)
  if (cart.items.length === 0) {
    return c.json({ error: { code: 'EMPTY_CART', message: 'Cart is empty' } }, 400)
  }

  // Validate stock and calculate total
  let totalAmount = 0
  const orderItems = []

  for (const cartItem of cart.items) {
    const product = store.getProduct(cartItem.productId)
    if (!product) {
      return c.json({ error: { code: 'NOT_FOUND', message: `Product ${cartItem.productId} not found` } }, 404)
    }
    if (product.stock < cartItem.quantity) {
      return c.json({ error: { code: 'INSUFFICIENT_STOCK', message: `Insufficient stock for product ${product.name}` } }, 400)
    }

    totalAmount += cartItem.price * cartItem.quantity
    orderItems.push({
      productId: cartItem.productId,
      quantity: cartItem.quantity,
      price: cartItem.price,
    })

    // Update product stock
    store.updateProduct(product.id, {
      ...product,
      stock: product.stock - cartItem.quantity,
      updatedAt: new Date().toISOString(),
    })
  }

  const newOrder: Order = {
    id: Date.now().toString(),
    userId: body.userId,
    items: orderItems,
    totalAmount,
    status: 'pending',
    shippingAddress: body.shippingAddress || user.address,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  }

  const created = store.createOrder(newOrder)

  // Clear cart
  store.updateCart(body.userId, {
    userId: body.userId,
    items: [],
    updatedAt: new Date().toISOString(),
  })

  return c.json(created, 201)
})

// PATCH /orders/{orderId}/status
orders.patch('/:orderId/status', async (c) => {
  const orderId = c.req.param('orderId')
  const body = await c.req.json<OrderUpdateStatusRequest>()
  
  const order = store.getOrder(orderId)
  if (!order) {
    return c.json({ error: { code: 'NOT_FOUND', message: 'Order not found' } }, 404)
  }

  // Validate status transition
  const validTransitions: Record<OrderStatus, OrderStatus[]> = {
    'pending': ['processing', 'cancelled'],
    'processing': ['shipped', 'cancelled'],
    'shipped': ['delivered', 'cancelled'],
    'delivered': [],
    'cancelled': [],
  }

  if (!validTransitions[order.status].includes(body.status)) {
    return c.json({ 
      error: { 
        code: 'INVALID_STATUS_TRANSITION', 
        message: `Cannot transition from ${order.status} to ${body.status}` 
      } 
    }, 400)
  }

  const updatedOrder: Order = {
    ...order,
    status: body.status,
    updatedAt: new Date().toISOString(),
  }

  const updated = store.updateOrder(orderId, updatedOrder)
  return c.json(updated)
})

export default orders