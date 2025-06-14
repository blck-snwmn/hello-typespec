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
   * List all users
   */
  @get
  list(...PaginationParams): PaginatedResponse<User>;

  /**
   * Get a user by ID
   */
  @get
  @route("/{userId}")
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
  update(
    @path userId: uuid,
    @body user: UpdateUserRequest
  ): User | ErrorResponse;

  /**
   * Delete a user
   */
  @delete
  @route("/{userId}")
  delete(@path userId: uuid): void | ErrorResponse;
}