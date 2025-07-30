import "../models/common.tsp";
import "../models/product.tsp";

using TypeSpec.Http;

namespace ECSite;

/**
 * Cart item
 */
model CartItem {
  @doc("ID of the product in the cart")
  productId: uuid;

  @doc("Quantity of the product")
  quantity: int32;

  @doc("Product details (populated when fetching cart)")
  product?: Product;
}

/**
 * Shopping cart
 */
model Cart {
  @doc("Unique identifier for the cart")
  id: uuid;

  @doc("ID of the user who owns this cart")
  userId: uuid;

  @doc("List of items in the cart")
  items: CartItem[];

  ...Timestamps;
}

/**
 * Add item to cart request
 */
model AddCartItemRequest {
  @doc("ID of the product to add")
  productId: uuid;

  @doc("Quantity to add")
  quantity: int32;
}

/**
 * Update cart item request
 */
model UpdateCartItemRequest {
  @doc("New quantity for the cart item")
  quantity: int32;
}

/**
 * Cart summary with calculated totals
 */
model CartSummary extends Cart {
  @doc("Total price of all items in the cart")
  totalAmount: float32;

  @doc("Total number of items in the cart")
  totalItems: int32;
}