import "../models/common.tsp";

using TypeSpec.Http;

namespace ECSite;

/**
 * Category model
 */
model Category {
  @doc("Unique identifier for the category")
  id: uuid;

  @doc("Name of the category")
  name: string;

  @doc("ID of the parent category for hierarchical structure")
  parentId?: uuid;

  ...Timestamps;
}

/**
 * Category with nested children
 */
model CategoryTree extends Category {
  @doc("List of child categories")
  children: CategoryTree[];
}

/**
 * Category creation request
 */
model CreateCategoryRequest {
  @doc("Name of the category")
  name: string;

  @doc("Optional ID of the parent category")
  parentId?: uuid;
}

/**
 * Category update request
 */
model UpdateCategoryRequest {
  @doc("Updated name of the category")
  name?: string;

  @doc("Updated parent category ID")
  parentId?: uuid;
}