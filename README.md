# DrynkLab Recipe Service

This is the Recipe Service for the DrynkLab application. It provides gRPC endpoints for managing cocktail recipes.

## Features

- Create, read, update, and delete recipes
- List recipes with pagination
- Input validation for recipe creation and updates

## Prerequisites

- Go 1.16 or later
- Protocol Buffers compiler
- gRPC

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/drynklab-recipe-service.git
   ```

2. Change to the project directory:
   ```
   cd drynklab-recipe-service
   ```

3. Install dependencies:
   ```
   go mod tidy
   ```

## Usage

1. Run the server:
   ```
   go run cmd/server/main.go
   ```

2. The gRPC server will start on port 50051.

## Development

To regenerate gRPC code after modifying the protocol buffer definition:

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/recipe/recipe.proto
```

## Testing

(Add information about running tests once they are implemented)

## Contributing

(Add contribution guidelines if applicable)

## License

(Add license information)