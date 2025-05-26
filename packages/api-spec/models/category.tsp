import "../models/common.tsp";

using TypeSpec.Http;

namespace ECSite;

/**
 * Category model
 */
model Category {
  id: uuid;
  name: string;
  parentId?: uuid;
  ...Timestamps;
}

/**
 * Category with nested children
 */
model CategoryTree extends Category {
  children: CategoryTree[];
}

/**
 * Category creation request
 */
model CreateCategoryRequest {
  name: string;
  parentId?: uuid;
}

/**
 * Category update request
 */
model UpdateCategoryRequest {
  name?: string;
  parentId?: uuid;
}