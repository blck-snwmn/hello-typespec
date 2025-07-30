import "../models/common.tsp";

using TypeSpec.Http;

namespace ECSite;

/**
 * Product model
 */
model Product {
  @doc("Unique identifier for the product")
  id: uuid;

  @doc("Name of the product")
  name: string;

  @doc("Detailed description of the product")
  description: string;

  @doc("Price of the product")
  price: float32;

  @doc("Current stock quantity")
  stock: int32;

  @doc("ID of the category this product belongs to")
  categoryId: uuid;

  @doc("List of product image URLs")
  imageUrls: string[];

  ...Timestamps;
}

/**
 * Product creation request
 */
model CreateProductRequest {
  @doc("Name of the product")
  name: string;

  @doc("Detailed description of the product")
  description: string;

  @doc("Price of the product")
  price: float32;

  @doc("Initial stock quantity")
  stock: int32;

  @doc("ID of the category this product belongs to")
  categoryId: uuid;

  @doc("Optional list of product image URLs")
  imageUrls?: string[];
}

/**
 * Product update request
 */
model UpdateProductRequest {
  @doc("Updated name of the product")
  name?: string;

  @doc("Updated description of the product")
  description?: string;

  @doc("Updated price of the product")
  price?: float32;

  @doc("Updated stock quantity")
  stock?: int32;

  @doc("Updated category ID")
  categoryId?: uuid;

  @doc("Updated list of product image URLs")
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