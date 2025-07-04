import { Hono } from 'hono'
import type { components, operations } from '../types/api'
import { store } from '../stores'
import { sendError, ErrorCode } from '../types/errors'

type Category = components['schemas']['Category']
type CategoryWithChildren = Category & { children: CategoryWithChildren[] }
type CategoryCreateRequest = operations['CategoriesService_create']['requestBody']['content']['application/json']
type CategoryUpdateRequest = operations['CategoriesService_update']['requestBody']['content']['application/json']

const categories = new Hono()

// GET /categories
categories.get('/', (c) => {
  const allCategories = store.getCategories()
  return c.json(allCategories)
})

// GET /categories/tree
categories.get('/tree', (c) => {
  const allCategories = store.getCategories()
  
  // Build tree structure
  const categoryMap = new Map<string, CategoryWithChildren>()
  const rootCategories: CategoryWithChildren[] = []

  // First pass: create CategoryWithChildren objects
  allCategories.forEach(cat => {
    categoryMap.set(cat.id, { ...cat, children: [] })
  })

  // Second pass: build tree
  categoryMap.forEach(cat => {
    if (!cat.parentId) {
      rootCategories.push(cat)
    } else {
      const parent = categoryMap.get(cat.parentId)
      if (parent) {
        parent.children.push(cat)
      }
    }
  })

  return c.json(rootCategories)
})

// GET /categories/{categoryId}
categories.get('/:categoryId', (c) => {
  const categoryId = c.req.param('categoryId')
  const category = store.getCategory(categoryId)

  if (!category) {
    return sendError(c, 404, ErrorCode.NOT_FOUND, 'Category not found')
  }

  return c.json(category)
})

// POST /categories
categories.post('/', async (c) => {
  const body = await c.req.json<CategoryCreateRequest>()
  
  const newCategory: Category = {
    id: Date.now().toString(),
    ...body,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  }

  const created = store.createCategory(newCategory)
  return c.json(created, 201)
})

// PUT /categories/{categoryId}
categories.put('/:categoryId', async (c) => {
  const categoryId = c.req.param('categoryId')
  const body = await c.req.json<CategoryUpdateRequest>()
  
  const existing = store.getCategory(categoryId)
  if (!existing) {
    return sendError(c, 404, ErrorCode.NOT_FOUND, 'Category not found')
  }

  const updatedCategory: Category = {
    ...existing,
    ...body,
    id: categoryId,
  }

  const updated = store.updateCategory(categoryId, updatedCategory)
  return c.json(updated)
})

// DELETE /categories/{categoryId}
categories.delete('/:categoryId', (c) => {
  const categoryId = c.req.param('categoryId')
  const deleted = store.deleteCategory(categoryId)

  if (!deleted) {
    return sendError(c, 404, ErrorCode.NOT_FOUND, 'Category not found')
  }

  return c.body(null, 204)
})

export default categories