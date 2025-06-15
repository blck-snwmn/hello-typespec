import "@typespec/rest";
import "@typespec/openapi3";
import "../models/common.tsp";
import "../models/product.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.OpenAPI;

namespace ECSite;

@route("/products")
@tag("Products")
interface ProductsService {
  /**
   * List all products with optional filtering
   */
  @get
  list(...ProductSearchParams): PaginatedResponse<Product> | ErrorResponse;

  /**
   * Get a product by ID
   */
  @get
  @route("/{productId}")
  get(@path productId: uuid): Product | ErrorResponse;

  /**
   * Create a new product (Admin only)
   */
  @post
  @useAuth(TypeSpec.Http.BearerAuth)
  create(@body product: CreateProductRequest): Product | ErrorResponse;

  /**
   * Update a product (Admin only)
   */
  @patch
  @route("/{productId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  update(
    @path productId: uuid,
    @body product: UpdateProductRequest
  ): Product | ErrorResponse;

  /**
   * Delete a product (Admin only)
   */
  @delete
  @route("/{productId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  delete(@path productId: uuid): void | ErrorResponse;
}