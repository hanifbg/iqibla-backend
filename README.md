# iQibla E-commerce API

A Go web application using the Echo framework for handling e-commerce backend operations.

## Project Structure

```
├── cmd/
│   └── main.go           # Application entry point
├── config/              # Configuration files
├── internal/            # Private application code
│   ├── handler/         # HTTP handlers
│   ├── model/           # Data models
│   ├── repository/      # Data access layer
│   └── service/         # Business logic
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
└── README.md           # Project documentation
```

## Prerequisites

- Go 1.21 or later
- Git

## Getting Started

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd landing_backend
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run cmd/main.go
   ```

   The server will start on port 8080.

## Testing the API

You can test the hello endpoint using curl:

```bash
curl http://localhost:8080/hello
```

Or open in your web browser:
```
http://localhost:8080/hello
```

You should see the message: "Hello, iQibla E-commerce API!"

## Development

- The application uses the Echo framework for routing and HTTP handling
- Follows standard Go project layout conventions
- Implements graceful shutdown