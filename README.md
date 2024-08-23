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

## Protocol Buffers

This project uses Protocol Buffers for defining the gRPC service. The `.proto` file is located in `proto/recipe/recipe.proto`.

The Go code generated from this `.proto` file is committed to the repository for convenience and consistency. If you make changes to the `.proto` file, regenerate the Go code using:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/recipe/recipe.proto
```

Ensure you have the latest version of `protoc` and the necessary Go plugins installed before regenerating.
