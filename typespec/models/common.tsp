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
  @doc("Error information")
  error: {
    @doc("Error code identifying the type of error")
    code: ErrorCode;

    @doc("Human-readable error message")
    message: string;

    @doc("Additional error details")
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
  @doc("Array of items in the current page")
  items: T[];

  @doc("Total number of items")
  total: int32;

  @doc("Maximum number of items per page")
  limit: int32;

  @doc("Number of items skipped")
  offset: int32;
}

/**
 * Common timestamp fields
 */
model Timestamps {
  @doc("Timestamp when the resource was created")
  createdAt: utcDateTime;

  @doc("Timestamp when the resource was last updated")
  updatedAt: utcDateTime;
}

/**
 * UUID type alias
 */
scalar uuid extends string;

/**
 * Login request
 */
model LoginRequest {
  @doc("User's email address")
  email: string;

  @doc("User's password")
  password: string;
}

/**
 * Login response with access token
 */
model LoginResponse {
  @doc("JWT access token")
  accessToken: string;

  @doc("Token type (always Bearer)")
  tokenType: "Bearer";

  @doc("Token expiration time in seconds")
  expiresIn: int32;

  @doc("Authenticated user information")
  user: {
    @doc("User ID")
    id: uuid;

    @doc("User's email address")
    email: string;

    @doc("User's full name")
    name: string;
  };
}

/**
 * Authenticated user context
 */
model AuthUser {
  @doc("User ID")
  id: uuid;

  @doc("User's email address")
  email: string;

  @doc("User's full name")
  name: string;
}