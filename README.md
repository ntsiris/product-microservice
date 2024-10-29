# Product Microservice

This is a RESTful JSON API CRUD microservice for managing products, developed in Go. The microservice provides a complete API for creating, reading, updating, and deleting products, along with MySQL-backed storage and structured error handling.

## Project Structure
```
  .
  ├── api # Contains the HTTP handlers and routing logic
  │
  ├── product-api.go # Main handlers for product CRUD operations
  │
  ├── server.go # APIServer setup and routes registration
  ├── bin # Compiled binaries
  │
  └── product
  ├── cmd
  │
  └── main.go # Application entry point
  ├── config # Configuration and environment variable management
  │
  ├── config.go
  │
  └── env.go
  ├── docker # Docker configuration files
  │
  ├── Dockerfile
  │
  └── docker-compose.yaml
  ├── internal # Internal packages with business logic, storage, and utility functions
  │
  ├── service # Business logic for product management
  │
  ├── storage # Storage layer with MySQL integration
  │
  ├── types # Error handling and custom types
  │
  └── utils # JSON handling and data validation utilities
  ├── migrations # Database migration scripts
  └── Makefile # Build and management commands
```

## Features

- **CRUD Operations**: Create, retrieve, update, and delete products.
- **MySQL Database Integration**: Persistent storage for product information with MySQL.
- **Automated Migrations**: Database schema management through migration scripts.
- **Validation**: Request validation using `go-playground/validator`.
- **Error Handling**: Consistent error responses with detailed messages.

## Requirements

- **Go** (1.22.7 or later)
- **MySQL** (8.0 or compatible)
- **Docker & Docker Compose** (for containerized deployment)

## Setup

### 1. Environment Variables

Configure the application using environment variables, either in a `.env` file or directly in the Docker Compose file:

```shell
# .env file example
PUBLIC_HOST=localhost
PORT=8080
DB_DRIVER=mysql
DB_USER=root
DB_PASSWORD=a-secret-password
DB_HOST=db
DB_PORT=3306
DB_NAME=productDB
MIGRATE_UP=true
MIGRATE_DOWN=false
MIGRATION_PATH=migrations/
```

### 2. Build and Run

Using Docker Compose:

```shell
make docker-run
```

This command builds the Docker image (if not already built) and runs the services.

Alternatively, you can run the application locally:

```shell
# Build
make build

# Run
make run
```

### 3. Running Migrations
Database migrations are automatically run on startup if `MIGRATE_UP=true`. To run migrations manually:

```shell
make migrate-up
make migrate-down
```

## API Endpoints

### Base URL

`/api/v1`

### Product Endpoints


| Method | Endpoint             | Description                   |
| ------ | -------------------- | ----------------------------- |
| POST   | /product/create      | Create a new product          |
| GET    | /product/{id}        | Retrieve a specific product   |
| GET    | /product             | List all products (paginated) |
| PUT    | /product/update      | Update an existing product    |
| DELETE | /product/delete/{id} | Delete a product              |

### Sample Product JSON

```json
{
    "name": "Sample Product",
    "description": "A sample description",
    "price": 29.99,
    "quantity": 100,
    "discount": 0.1
}
```

## Testing

Run tests with:
```shell
make test
```
## Project Components

### `api` Module
Handles HTTP requests, route registration, and request validation. Provides structured error handling using `types.APIError`.

### `service` Module
Implements core product logic and types.

### `storage` Module
Defines a `ProductStore` interface for data persistence, with `MySQLStore` as the implementation. It supports **CRUD** operations and migrations.

### `config` Module
Manages environment variables and application configuration, including storage (database) and server settings.

### `types` Module
Contains reusable types, like `APIError` for structured error handling and logging.

### `utils` Module
Utility functions for JSON parsing, writing, and validation.

## Deployment

### Docker
The application can be deployed using Docker Compose. The `docker-compose.yml` file configures both the application and the MySQL database.
```shell
docker-compose -f docker/docker-compose.yaml up -d
```

### Health Check
The application provides a basic health check configured in Docker Compose, polling the main endpoint to ensure the application is responsive.
