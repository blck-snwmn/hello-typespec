import "@typespec/http";

using TypeSpec.Http;

namespace ECSite;

/**
 * Common error response
 */
model ErrorResponse {
  error: {
    code: string;
    message: string;
  };
}

/**
 * Pagination parameters
 */
model PaginationParams {
  @query
  @doc("Maximum number of items to return")
  limit?: int32 = 20;

  @query
  @doc("Number of items to skip")
  offset?: int32 = 0;
}

/**
 * Paginated response wrapper
 */
model PaginatedResponse<T> {
  items: T[];
  total: int32;
  limit: int32;
  offset: int32;
}

/**
 * Common timestamp fields
 */
model Timestamps {
  createdAt: utcDateTime;
  updatedAt: utcDateTime;
}

/**
 * UUID type alias
 */
scalar uuid extends string;