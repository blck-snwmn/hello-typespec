import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { logger } from 'hono/logger'
import products from './routes/products'
import categories from './routes/categories'
import users from './routes/users'
import carts from './routes/carts'
import orders from './routes/orders'
import auth from './routes/auth'
import { authMiddleware } from './middleware/auth'
import { globalErrorHandler } from './types/errors'

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

// Mount protected routes (require authentication)
app.use('/users/*', authMiddleware)
app.route('/users', users)

app.use('/carts/*', authMiddleware)
app.route('/carts', carts)

app.use('/orders/*', authMiddleware)
app.route('/orders', orders)

export default app