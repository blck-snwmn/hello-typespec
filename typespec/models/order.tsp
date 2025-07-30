import "../models/common.tsp";
import "../models/user.tsp";

using TypeSpec.Http;

namespace ECSite;

/**
 * Order status enum
 */
enum OrderStatus {
  @doc("Order has been placed but not yet processed")
  pending: "pending",

  @doc("Order is being processed")
  processing: "processing",

  @doc("Order has been shipped")
  shipped: "shipped",

  @doc("Order has been delivered")
  delivered: "delivered",

  @doc("Order has been cancelled")
  cancelled: "cancelled",
}

/**
 * Order item
 */
model OrderItem {
  @doc("ID of the ordered product")
  productId: uuid;

  @doc("Quantity ordered")
  quantity: int32;

  @doc("Price at the time of order")
  price: float32;

  @doc("Name of the product at the time of order")
  productName: string;
}

/**
 * Order model
 */
model Order {
  @doc("Unique identifier for the order")
  id: uuid;

  @doc("ID of the user who placed the order")
  userId: uuid;

  @doc("List of items in the order")
  items: OrderItem[];

  @doc("Total amount of the order")
  totalAmount: float32;

  @doc("Current status of the order")
  status: OrderStatus;

  @doc("Shipping address for the order")
  shippingAddress: Address;

  ...Timestamps;
}

/**
 * Create order request
 */
model CreateOrderRequest {
  @doc("List of items to order")
  items: OrderItem[];

  @doc("Shipping address for the order")
  shippingAddress: Address;
}

/**
 * Update order status request
 */
model UpdateOrderStatusRequest {
  @doc("New status for the order")
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