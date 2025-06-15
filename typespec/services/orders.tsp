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
   * List all orders with optional filtering (Admin only)
   */
  @get
  @useAuth(TypeSpec.Http.BearerAuth)
  list(...OrderSearchParams): PaginatedResponse<Order> | ErrorResponse;

  /**
   * Get an order by ID
   */
  @get
  @route("/{orderId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  get(@path orderId: uuid): Order | ErrorResponse;

  /**
   * Create a new order from user's cart
   */
  @post
  @route("/users/{userId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  create(
    @path userId: uuid,
    @body order: CreateOrderRequest
  ): Order | ErrorResponse;

  /**
   * Update order status (Admin only)
   */
  @patch
  @route("/status/{orderId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  updateStatus(
    @path orderId: uuid,
    @body status: UpdateOrderStatusRequest
  ): Order | ErrorResponse;

  /**
   * Cancel an order
   */
  @post
  @route("/cancel/{orderId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  cancel(@path orderId: uuid): Order | ErrorResponse;

  /**
   * Get orders by user ID
   */
  @get
  @route("/users/{userId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  listByUser(
    @path userId: uuid,
    ...PaginationParams
  ): PaginatedResponse<Order> | ErrorResponse;
}