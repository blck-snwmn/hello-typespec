import { Hono } from 'hono'
import type { components, operations } from '../types/api'
import { store } from '../stores'
import { sendError, ErrorCode } from '../types/errors'

type Order = components['schemas']['Order']
type OrderStatus = components['schemas']['OrderStatus']
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
    items: paginatedOrders,
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
    return sendError(c, 404, ErrorCode.NOT_FOUND, 'Order not found')
  }

  return c.json(order)
})

// GET /orders/users/{userId}
orders.get('/users/:userId', (c) => {
  const userId = c.req.param('userId')
  const limit = parseInt(c.req.query('limit') || '10')
  const offset = parseInt(c.req.query('offset') || '0')
  const status = c.req.query('status')

  // Validate user exists
  const user = store.getUser(userId)
  if (!user) {
    return sendError(c, 404, ErrorCode.NOT_FOUND, 'User not found')
  }

  // Get all orders for the user
  let userOrders = store.getOrders().filter(order => order.userId === userId)

  // Apply status filter if provided
  if (status) {
    userOrders = userOrders.filter(order => order.status === status)
  }

  // Apply pagination
  const paginatedOrders = userOrders.slice(offset, offset + limit)

  const response: OrderListResponse = {
    items: paginatedOrders,
    total: userOrders.length,
    limit,
    offset,
  }

  return c.json(response)
})

// POST /orders/users/{userId}
orders.post('/users/:userId', async (c) => {
  const userId = c.req.param('userId')
  const body = await c.req.json<operations['OrdersService_create']['requestBody']['content']['application/json']>()
  
  // Validate user exists
  const user = store.getUser(userId)
  if (!user) {
    return sendError(c, 404, ErrorCode.NOT_FOUND, 'User not found')
  }

  // Get user's cart
  const cart = store.getCartByUserId(userId)
  if (cart.items.length === 0) {
    return sendError(c, 400, ErrorCode.BAD_REQUEST, 'Cart is empty')
  }

  // Validate stock and calculate total
  let totalAmount = 0
  const orderItems = []

  for (const cartItem of cart.items) {
    const product = store.getProduct(cartItem.productId)
    if (!product) {
      return sendError(c, 404, ErrorCode.NOT_FOUND, `Product ${cartItem.productId} not found`)
    }
    if (product.stock < cartItem.quantity) {
      return sendError(c, 400, ErrorCode.INSUFFICIENT_STOCK, `Insufficient stock for product ${product.name}`)
    }

    const itemPrice = product.price
    totalAmount += itemPrice * cartItem.quantity
    orderItems.push({
      productId: cartItem.productId,
      quantity: cartItem.quantity,
      price: itemPrice,
      productName: product.name,
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
    userId: userId,
    items: orderItems,
    totalAmount,
    status: 'pending',
    shippingAddress: body.shippingAddress,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  }

  const created = store.createOrder(newOrder)

  // Clear cart
  const clearedCart = store.getCartByUserId(userId)
  clearedCart.items = []
  clearedCart.updatedAt = new Date().toISOString()
  store.updateCart(userId, clearedCart)

  return c.json(created, 201)
})

// PATCH /orders/status/{orderId}
orders.patch('/status/:orderId', async (c) => {
  const orderId = c.req.param('orderId')
  const body = await c.req.json<OrderUpdateStatusRequest>()
  
  const order = store.getOrder(orderId)
  if (!order) {
    return sendError(c, 404, ErrorCode.NOT_FOUND, 'Order not found')
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
    return sendError(c, 400, ErrorCode.INVALID_STATE_TRANSITION, `Cannot transition from ${order.status} to ${body.status}`)
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