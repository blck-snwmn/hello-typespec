import "@typespec/http";

using TypeSpec.Http;

namespace ECSite;

/**
 * Standard error codes used throughout the API
 */
enum ErrorCode {
  // Client errors (4xx)
  BAD_REQUEST: "BAD_REQUEST",
  UNAUTHORIZED: "UNAUTHORIZED",
  FORBIDDEN: "FORBIDDEN",
  NOT_FOUND: "NOT_FOUND",
  CONFLICT: "CONFLICT",
  VALIDATION_ERROR: "VALIDATION_ERROR",
  INSUFFICIENT_STOCK: "INSUFFICIENT_STOCK",
  INVALID_STATE_TRANSITION: "INVALID_STATE_TRANSITION",
  
  // Server errors (5xx)
  INTERNAL_ERROR: "INTERNAL_ERROR",
  SERVICE_UNAVAILABLE: "SERVICE_UNAVAILABLE",
}

/**
 * Common error response
 */
model ErrorResponse {
  error: {
    code: ErrorCode;
    message: string;
    details?: unknown;
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