import "../models/common.tsp";
import "../models/product.tsp";

using TypeSpec.Http;

namespace ECSite;

/**
 * Cart item
 */
model CartItem {
  productId: uuid;
  quantity: int32;
  product?: Product;
}

/**
 * Shopping cart
 */
model Cart {
  id: uuid;
  userId: uuid;
  items: CartItem[];
  ...Timestamps;
}

/**
 * Add item to cart request
 */
model AddCartItemRequest {
  productId: uuid;
  quantity: int32;
}

/**
 * Update cart item request
 */
model UpdateCartItemRequest {
  quantity: int32;
}

/**
 * Cart summary with calculated totals
 */
model CartSummary extends Cart {
  totalAmount: float32;
  totalItems: int32;
}