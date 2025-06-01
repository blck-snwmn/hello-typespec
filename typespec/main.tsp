import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

// Import all services
import "./services/products.tsp";
import "./services/categories.tsp";
import "./services/users.tsp";
import "./services/carts.tsp";
import "./services/orders.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

/**
 * EC Site API
 * 
 * Sample e-commerce API built with TypeSpec
 */
@service(#{
  title: "EC Site API",
})
@server("http://localhost:3000", "Development server")
namespace ECSite;