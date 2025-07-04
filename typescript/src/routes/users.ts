import { Hono } from 'hono'
import type { components, operations } from '../types/api'
import { store } from '../stores'
import { sendError, ErrorCode } from '../types/errors'

type User = components['schemas']['User']
type UserListResponse = operations['UsersService_list']['responses']['200']['content']['application/json']
type UserCreateRequest = operations['UsersService_create']['requestBody']['content']['application/json']
type UserUpdateRequest = operations['UsersService_update']['requestBody']['content']['application/json']

const users = new Hono()

// GET /users
users.get('/', (c) => {
  const limit = parseInt(c.req.query('limit') || '20')
  const offset = parseInt(c.req.query('offset') || '0')
  
  const allUsers = store.getUsers()
  
  // Apply pagination
  const paginatedUsers = allUsers.slice(offset, offset + limit)
  
  const response: UserListResponse = {
    items: paginatedUsers,
    total: allUsers.length,
    limit,
    offset,
  }
  
  return c.json(response)
})

// GET /users/{userId}
users.get('/:userId', (c) => {
  const userId = c.req.param('userId')
  const user = store.getUser(userId)

  if (!user) {
    return sendError(c, 404, ErrorCode.NOT_FOUND, 'User not found')
  }

  return c.json(user)
})

// POST /users
users.post('/', async (c) => {
  const body = await c.req.json<UserCreateRequest>()
  
  const newUser: User = {
    id: Date.now().toString(),
    ...body,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  }

  const created = store.createUser(newUser)
  
  // Initialize empty cart for new user
  store.updateCart(created.id, {
    id: `cart-${created.id}`,
    userId: created.id,
    items: [],
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  })

  return c.json(created, 201)
})

// PUT /users/{userId}
users.put('/:userId', async (c) => {
  const userId = c.req.param('userId')
  const body = await c.req.json<UserUpdateRequest>()
  
  const existing = store.getUser(userId)
  if (!existing) {
    return sendError(c, 404, ErrorCode.NOT_FOUND, 'User not found')
  }

  const updatedUser: User = {
    ...existing,
    ...body,
    id: userId,
    updatedAt: new Date().toISOString(),
  }

  const updated = store.updateUser(userId, updatedUser)
  return c.json(updated)
})

// DELETE /users/{userId}
users.delete('/:userId', (c) => {
  const userId = c.req.param('userId')
  const deleted = store.deleteUser(userId)

  if (!deleted) {
    return sendError(c, 404, ErrorCode.NOT_FOUND, 'User not found')
  }

  return c.body(null, 204)
})

export default users