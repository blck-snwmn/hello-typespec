import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

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

// Common types will be imported from models/common.tsp
// Service definitions will be imported from services/*.tsp