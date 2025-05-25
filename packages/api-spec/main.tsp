import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.OpenAPI;

@service({
  title: "EC Site API",
  version: "1.0.0",
})
@server("http://localhost:3000", "Development server")
namespace ECSite;

// Common types will be imported from models/common.tsp
// Service definitions will be imported from services/*.tsp