# iQibla E-commerce API

A Go web application using the Echo framework for handling e-commerce backend operations with integrated payment processing via Midtrans.

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
- PostgreSQL database
- Midtrans account for payment processing

## Getting Started

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd iqibla-backend
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Configure the application:
   - Copy the example configuration file:
     ```bash
     cp config/app.config.json.example config/app.config.json
     ```
   - Update the configuration with your database and Midtrans credentials
   - Configure RajaOngkir API and caching settings:
     ```json
     "shipping": {
       "rajaongkir_api_key": "your_rajaongkir_api_key",
       "rajaongkir_base_url": "https://rajaongkir.komerce.id/api/v1",
       "rajaongkir_cache_enabled": true,
       "rajaongkir_cache_ttl_hours": 24,
       "rajaongkir_warmup_on_startup": true,
       "rajaongkir_warmup_timeout_secs": 30
     }
     ```

4. Run the server:
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
- Uses GORM as the ORM for database operations

## Features

### Product Management
- Product listing and details
- Product variants with different prices

### Shopping Cart
- Add items to cart
- Update item quantities
- Remove items from cart
- Apply discount codes

### Payment Processing
- Integrated with Midtrans payment gateway
- Support for various payment methods through Midtrans Snap
- Order creation and management
- Payment status tracking
- Payment notification handling

### Shipping
- Integration with RajaOngkir API for shipping rates
- In-memory caching of geographic data to reduce API calls
- Configurable cache TTL and warm-up behavior
- Support for multiple couriers

## API Documentation

Swagger documentation is available at `/swagger/index.html` when the server is running.

## Additional Documentation

- [RajaOngkir Caching Implementation](/docs/development/rajaongkir_caching_guide.md) - Details about the RajaOngkir API caching system