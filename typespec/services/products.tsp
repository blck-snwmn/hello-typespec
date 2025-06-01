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
   * Create a new product
   */
  @post
  create(@body product: CreateProductRequest): Product | ErrorResponse;

  /**
   * Update a product
   */
  @patch
  @route("/{productId}")
  update(
    @path productId: uuid,
    @body product: UpdateProductRequest
  ): Product | ErrorResponse;

  /**
   * Delete a product
   */
  @delete
  @route("/{productId}")
  delete(@path productId: uuid): void | ErrorResponse;
}