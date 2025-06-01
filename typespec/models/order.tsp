import "../models/common.tsp";
import "../models/user.tsp";

using TypeSpec.Http;

namespace ECSite;

/**
 * Order status enum
 */
enum OrderStatus {
  pending: "pending",
  processing: "processing",
  shipped: "shipped",
  delivered: "delivered",
  cancelled: "cancelled",
}

/**
 * Order item
 */
model OrderItem {
  productId: uuid;
  quantity: int32;
  price: float32;
  productName: string;
}

/**
 * Order model
 */
model Order {
  id: uuid;
  userId: uuid;
  items: OrderItem[];
  totalAmount: float32;
  status: OrderStatus;
  shippingAddress: Address;
  ...Timestamps;
}

/**
 * Create order request
 */
model CreateOrderRequest {
  items: OrderItem[];
  shippingAddress: Address;
}

/**
 * Update order status request
 */
model UpdateOrderStatusRequest {
  status: OrderStatus;
}

/**
 * Order search parameters
 */
model OrderSearchParams {
  ...PaginationParams;
  
  @query
  @doc("Filter by order status")
  status?: OrderStatus;
  
  @query
  @doc("Filter by user ID")
  userId?: uuid;
  
  @query
  @doc("Start date for order creation")
  startDate?: utcDateTime;
  
  @query
  @doc("End date for order creation")
  endDate?: utcDateTime;
}