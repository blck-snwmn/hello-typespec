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
  getByUser(@path userId: uuid): CartSummary | ErrorResponse;

  /**
   * Add item to cart
   */
  @post
  @route("/users/{userId}/items")
  addItem(
    @path userId: uuid,
    @body item: AddCartItemRequest
  ): CartSummary | ErrorResponse;

  /**
   * Update cart item quantity
   */
  @patch
  @route("/users/{userId}/items/{productId}")
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
  removeItem(
    @path userId: uuid,
    @path productId: uuid
  ): CartSummary | ErrorResponse;

  /**
   * Clear all items from cart
   */
  @delete
  @route("/users/{userId}/items")
  clear(@path userId: uuid): void | ErrorResponse;
}