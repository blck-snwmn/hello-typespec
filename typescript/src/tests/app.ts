import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { logger } from 'hono/logger'
import products from '../routes/products'
import categories from '../routes/categories'
import users from '../routes/users'
import carts from '../routes/carts'
import orders from '../routes/orders'
import auth from '../routes/auth'
import { mockAuthMiddleware } from './auth-helper'
import { globalErrorHandler } from '../types/errors'

// Test version of app with mock auth
export function createTestApp() {
  const app = new Hono()

  // Global error handler
  app.onError(globalErrorHandler)

  app.use('*', logger())
  app.use('*', cors())

  app.get('/health', (c) => {
    return c.json({ status: 'ok' })
  })

  // Mount auth routes (no middleware needed)
  app.route('/auth', auth)

  // Mount public routes
  app.route('/products', products)
  app.route('/categories', categories)

  // Mount protected routes with MOCK authentication
  app.use('/users/*', mockAuthMiddleware)
  app.route('/users', users)

  app.use('/carts/*', mockAuthMiddleware)
  app.route('/carts', carts)

  app.use('/orders/*', mockAuthMiddleware)
  app.route('/orders', orders)

  return app
}