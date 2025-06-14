# Go Implementation

This is the Go implementation of the TypeSpec-defined e-commerce API.

## Prerequisites

- Go 1.24 or higher
- oapi-codegen (for code generation)

## Setup

```bash
# Install dependencies
go mod download

# Install oapi-codegen
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

## Development

```bash
# Generate server code from OpenAPI spec
oapi-codegen -config oapi-codegen.yaml ../openapi/openapi.yaml

# Build the server
go build -o server ./cmd/server

# Run the server
./server

# Or run directly
go run ./cmd/server
```

The server will start on port 8080 by default. You can change this by setting the `PORT` environment variable.

## Project Structure

```
go/
├── cmd/
│   └── server/          # Main application entry point
├── generated/           # Generated code from OpenAPI spec
├── internal/           
│   ├── handlers/        # HTTP handlers implementation
│   └── store/          # In-memory data store
├── oapi-codegen.yaml   # Code generation configuration
├── go.mod              # Go module file
└── README.md           # This file
```

## API Endpoints

The server implements all endpoints defined in the TypeSpec specification:

- **Products**: CRUD operations, search, filtering, sorting
- **Categories**: CRUD operations, hierarchical structure
- **Users**: CRUD operations
- **Carts**: Cart management with stock validation
- **Orders**: Order creation and status management

## Testing

You can test the API using curl:

```bash
# Get all products
curl http://localhost:8080/products

# Create a new product
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name": "New Product", "description": "A new product", "price": 1000, "stock": 10, "categoryId": "1", "imageUrls": ["https://example.com/image.jpg"]}'
```