import "@typespec/rest";
import "@typespec/openapi3";
import "../models/common.tsp";
import "../models/order.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.OpenAPI;

namespace ECSite;

@route("/orders")
@tag("Orders")
interface OrdersService {
  /**
   * List all orders with optional filtering
   */
  @get
  list(...OrderSearchParams): PaginatedResponse<Order> | ErrorResponse;

  /**
   * Get an order by ID
   */
  @get
  @route("/{orderId}")
  get(@path orderId: uuid): Order | ErrorResponse;

  /**
   * Create a new order from user's cart
   */
  @post
  @route("/users/{userId}")
  create(
    @path userId: uuid,
    @body order: CreateOrderRequest
  ): Order | ErrorResponse;

  /**
   * Update order status
   */
  @patch
  @route("/status/{orderId}")
  updateStatus(
    @path orderId: uuid,
    @body status: UpdateOrderStatusRequest
  ): Order | ErrorResponse;

  /**
   * Cancel an order
   */
  @post
  @route("/cancel/{orderId}")
  cancel(@path orderId: uuid): Order | ErrorResponse;

  /**
   * Get orders by user ID
   */
  @get
  @route("/users/{userId}")
  listByUser(
    @path userId: uuid,
    ...PaginationParams
  ): PaginatedResponse<Order> | ErrorResponse;
}