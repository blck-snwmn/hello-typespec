import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { logger } from 'hono/logger'
import products from './routes/products'
import categories from './routes/categories'
import users from './routes/users'
import carts from './routes/carts'
import orders from './routes/orders'
import { globalErrorHandler } from './types/errors'

const app = new Hono()

// Global error handler
app.onError(globalErrorHandler)

app.use('*', logger())
app.use('*', cors())

app.get('/health', (c) => {
  return c.json({ status: 'ok' })
})

// Mount routes
app.route('/products', products)
app.route('/categories', categories)
app.route('/users', users)
app.route('/carts', carts)
app.route('/orders', orders)

export default app