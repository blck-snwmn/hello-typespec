import "../models/common.tsp";

using TypeSpec.Http;

namespace ECSite;

/**
 * User address
 */
model Address {
  @doc("Street address")
  street: string;

  @doc("City name")
  city: string;

  @doc("State or province")
  state: string;

  @doc("Postal or ZIP code")
  postalCode: string;

  @doc("Country name")
  country: string;
}

/**
 * User model
 */
model User {
  @doc("Unique identifier for the user")
  id: uuid;

  @doc("User's email address")
  email: string;

  @doc("User's full name")
  name: string;

  @doc("User's shipping address")
  address?: Address;

  ...Timestamps;
}

/**
 * User creation request
 */
model CreateUserRequest {
  @doc("User's email address")
  email: string;

  @doc("User's full name")
  name: string;

  @doc("Optional shipping address")
  address?: Address;
}

/**
 * User update request
 */
model UpdateUserRequest {
  @doc("Updated email address")
  email?: string;

  @doc("Updated user name")
  name?: string;

  @doc("Updated shipping address")
  address?: Address;
}