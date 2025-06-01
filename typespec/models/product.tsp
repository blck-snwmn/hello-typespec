import "../models/common.tsp";

using TypeSpec.Http;

namespace ECSite;

/**
 * Product model
 */
model Product {
  id: uuid;
  name: string;
  description: string;
  price: float32;
  stock: int32;
  categoryId: uuid;
  imageUrls: string[];
  ...Timestamps;
}

/**
 * Product creation request
 */
model CreateProductRequest {
  name: string;
  description: string;
  price: float32;
  stock: int32;
  categoryId: uuid;
  imageUrls?: string[];
}

/**
 * Product update request
 */
model UpdateProductRequest {
  name?: string;
  description?: string;
  price?: float32;
  stock?: int32;
  categoryId?: uuid;
  imageUrls?: string[];
}

/**
 * Product search parameters
 */
model ProductSearchParams {
  ...PaginationParams;
  
  @query
  @doc("Search by product name")
  name?: string;
  
  @query
  @doc("Filter by category ID")
  categoryId?: uuid;
  
  @query
  @doc("Minimum price")
  minPrice?: float32;
  
  @query
  @doc("Maximum price")
  maxPrice?: float32;
  
  @query
  @doc("Sort field")
  sortBy?: "name" | "price" | "createdAt" = "createdAt";
  
  @query
  @doc("Sort order")
  order?: "asc" | "desc" = "desc";
}