import "@typespec/rest";
import "@typespec/openapi3";
import "../models/common.tsp";
import "../models/cart.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.OpenAPI;

namespace ECSite;

@route("/carts")
@tag("Carts")
interface CartsService {
  /**
   * Get cart by user ID
   */
  @get
  @route("/users/{userId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  getByUser(@path userId: uuid): CartSummary | ErrorResponse;

  /**
   * Add item to cart
   */
  @post
  @route("/users/{userId}/items")
  @useAuth(TypeSpec.Http.BearerAuth)
  addItem(
    @path userId: uuid,
    @body item: AddCartItemRequest
  ): CartSummary | ErrorResponse;

  /**
   * Update cart item quantity
   */
  @patch
  @route("/users/{userId}/items/{productId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  updateItem(
    @path userId: uuid,
    @path productId: uuid,
    @body item: UpdateCartItemRequest
  ): CartSummary | ErrorResponse;

  /**
   * Remove item from cart
   */
  @delete
  @route("/users/{userId}/items/{productId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  removeItem(
    @path userId: uuid,
    @path productId: uuid
  ): CartSummary | ErrorResponse;

  /**
   * Clear all items from cart
   */
  @delete
  @route("/users/{userId}/items")
  @useAuth(TypeSpec.Http.BearerAuth)
  clear(@path userId: uuid): void | ErrorResponse;
}