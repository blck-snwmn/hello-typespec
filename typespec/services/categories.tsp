import "@typespec/rest";
import "@typespec/openapi3";
import "../models/common.tsp";
import "../models/category.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.OpenAPI;

namespace ECSite;

@route("/categories")
@tag("Categories")
interface CategoriesService {
  /**
   * List all categories
   */
  @get
  list(): Category[] | ErrorResponse;

  /**
   * Get category tree (with nested children)
   */
  @get
  @route("/tree")
  tree(): CategoryTree[] | ErrorResponse;

  /**
   * Get a category by ID
   */
  @get
  @route("/{categoryId}")
  get(@path categoryId: uuid): Category | ErrorResponse;

  /**
   * Create a new category (Admin only)
   */
  @post
  @useAuth(TypeSpec.Http.BearerAuth)
  create(@body category: CreateCategoryRequest): Category | ErrorResponse;

  /**
   * Update a category (Admin only)
   */
  @patch
  @route("/{categoryId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  update(
    @path categoryId: uuid,
    @body category: UpdateCategoryRequest
  ): Category | ErrorResponse;

  /**
   * Delete a category (Admin only)
   */
  @delete
  @route("/{categoryId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  delete(@path categoryId: uuid): void | ErrorResponse;
}