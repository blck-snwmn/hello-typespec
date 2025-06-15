import type { Context } from 'hono';
import type { StatusCode } from 'hono/utils/http-status';

// Content-returning status codes that Hono's json() method accepts
type ContentfulStatusCode = Exclude<StatusCode, 101 | 204 | 205 | 304>;

/**
 * Standard error codes matching TypeSpec definition
 */
export const ErrorCode = {
  // Client errors (4xx)
  BAD_REQUEST: 'BAD_REQUEST',
  UNAUTHORIZED: 'UNAUTHORIZED',
  FORBIDDEN: 'FORBIDDEN',
  NOT_FOUND: 'NOT_FOUND',
  CONFLICT: 'CONFLICT',
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  INSUFFICIENT_STOCK: 'INSUFFICIENT_STOCK',
  INVALID_STATE_TRANSITION: 'INVALID_STATE_TRANSITION',
  
  // Server errors (5xx)
  INTERNAL_ERROR: 'INTERNAL_ERROR',
  SERVICE_UNAVAILABLE: 'SERVICE_UNAVAILABLE',
} as const;

export type ErrorCode = typeof ErrorCode[keyof typeof ErrorCode];

/**
 * Error response type matching TypeSpec definition
 */
export interface ErrorResponse {
  error: {
    code: ErrorCode;
    message: string;
    details?: unknown;
  };
}

/**
 * Custom error class for API errors
 */
export class ApiError extends Error {
  constructor(
    public code: ErrorCode,
    public message: string,
    public statusCode: ContentfulStatusCode,
    public details?: unknown
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * Error code to HTTP status code mapping
 */
export const errorCodeToStatus: Record<ErrorCode, ContentfulStatusCode> = {
  [ErrorCode.BAD_REQUEST]: 400,
  [ErrorCode.UNAUTHORIZED]: 401,
  [ErrorCode.FORBIDDEN]: 403,
  [ErrorCode.NOT_FOUND]: 404,
  [ErrorCode.CONFLICT]: 409,
  [ErrorCode.VALIDATION_ERROR]: 400,
  [ErrorCode.INSUFFICIENT_STOCK]: 400,
  [ErrorCode.INVALID_STATE_TRANSITION]: 400,
  [ErrorCode.INTERNAL_ERROR]: 500,
  [ErrorCode.SERVICE_UNAVAILABLE]: 503,
};

/**
 * Create an error response
 */
export function createErrorResponse(
  code: ErrorCode,
  message: string,
  details?: unknown
): ErrorResponse {
  return {
    error: {
      code,
      message,
      details,
    },
  };
}

/**
 * Send an error response
 */
export function sendError(
  c: Context,
  statusCode: ContentfulStatusCode,
  code: ErrorCode,
  message: string,
  details?: unknown
) {
  return c.json<ErrorResponse>(createErrorResponse(code, message, details), statusCode);
}

/**
 * Global error handler for Hono
 */
export function globalErrorHandler(err: Error, c: Context) {
  console.error('Global error handler:', err);

  if (err instanceof ApiError) {
    return c.json<ErrorResponse>(
      createErrorResponse(err.code, err.message, err.details),
      err.statusCode
    );
  }

  // Handle validation errors from Zod (will be implemented later)
  if (err.name === 'ZodError') {
    return c.json<ErrorResponse>(
      createErrorResponse(
        ErrorCode.VALIDATION_ERROR,
        'Validation failed',
        err
      ),
      400
    );
  }

  // Default to internal server error
  return c.json<ErrorResponse>(
    createErrorResponse(
      ErrorCode.INTERNAL_ERROR,
      'An unexpected error occurred',
      process.env.NODE_ENV === 'development' ? err.message : undefined
    ),
    500
  );
}