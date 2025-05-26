import "../models/common.tsp";

using TypeSpec.Http;

namespace ECSite;

/**
 * User address
 */
model Address {
  street: string;
  city: string;
  state: string;
  postalCode: string;
  country: string;
}

/**
 * User model
 */
model User {
  id: uuid;
  email: string;
  name: string;
  address?: Address;
  ...Timestamps;
}

/**
 * User creation request
 */
model CreateUserRequest {
  email: string;
  name: string;
  address?: Address;
}

/**
 * User update request
 */
model UpdateUserRequest {
  email?: string;
  name?: string;
  address?: Address;
}