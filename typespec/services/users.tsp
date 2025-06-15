import "@typespec/rest";
import "@typespec/openapi3";
import "../models/common.tsp";
import "../models/user.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.OpenAPI;

namespace ECSite;

@route("/users")
@tag("Users")
interface UsersService {
  /**
   * List all users (Admin only)
   */
  @get
  @useAuth(TypeSpec.Http.BearerAuth)
  list(...PaginationParams): PaginatedResponse<User> | ErrorResponse;

  /**
   * Get a user by ID
   */
  @get
  @route("/{userId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  get(@path userId: uuid): User | ErrorResponse;

  /**
   * Create a new user
   */
  @post
  create(@body user: CreateUserRequest): User | ErrorResponse;

  /**
   * Update a user
   */
  @patch
  @route("/{userId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  update(
    @path userId: uuid,
    @body user: UpdateUserRequest
  ): User | ErrorResponse;

  /**
   * Delete a user
   */
  @delete
  @route("/{userId}")
  @useAuth(TypeSpec.Http.BearerAuth)
  delete(@path userId: uuid): void | ErrorResponse;
}