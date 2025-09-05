import "@typespec/http";
import "@typespec/rest";
import "../models/common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace ECSite;

/**
 * Authentication endpoints
 */
@route("/auth")
interface AuthService {
  /**
   * Login with email and password
   */
  @post
  @route("/login")
  login(@body request: LoginRequest): LoginResponse | ErrorResponse;

  /**
   * Logout (invalidate token)
   */
  @post
  @route("/logout")
  @useAuth(TypeSpec.Http.BearerAuth)
  logout(): OkResponse | ErrorResponse;

  /**
   * Get current user information
   */
  @get
  @route("/me")
  @useAuth(TypeSpec.Http.BearerAuth)
  getCurrentUser(): AuthUser | ErrorResponse;
}

/**
 * Simple OK response
 */
model OkResponse {
  message: string;
}
