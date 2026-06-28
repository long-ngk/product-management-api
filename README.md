# Product Management Server

A RESTful API for managing products built with Go, Gin framework, and PostgreSQL.

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.21+ |
| HTTP Framework | Gin |
| Database | PostgreSQL 15+ |
| DB Driver | pgx v5 (via database/sql) |
| Testing | Go testing + testify + rapid (property-based testing) |

## Architecture

The project follows a clean layered architecture with clear separation of concerns:

```
Client HTTP → Gin Router → Handler → Service → Repository → PostgreSQL
```

### Layer Responsibilities

| Layer | Responsibility | Restrictions |
|-------|---------------|-------------|
| **Handler** | Parse HTTP requests, route params, JSON body; call Service; return JSON response | No business logic, no direct DB access |
| **Service** | Business logic, input validation, error classification | No HTTP imports (gin, net/http) |
| **Repository** | SQL execution (INSERT, SELECT, UPDATE, DELETE) | No business logic, no validation |
| **Infrastructure** | DB connection pool and configuration | Provides connection to Repository |

### Project Structure

```
product-management-server/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point, dependency injection
├── internal/
│   ├── handler/
│   │   ├── product_handler.go       # HTTP handlers
│   │   └── product_handler_test.go  # Handler unit + property tests
│   ├── service/
│   │   ├── product_service.go       # Service interface
│   │   ├── product_service_impl.go  # Business logic & validation
│   │   └── product_service_test.go  # Service property-based tests
│   ├── repository/
│   │   ├── product_repository.go    # Repository interface
│   │   └── postgres_product_repository.go  # PostgreSQL implementation
│   ├── model/
│   │   ├── product.go               # Domain model & DTOs
│   │   ├── errors.go                # Custom error types
│   │   └── product_test.go          # JSON serialization property test
│   └── infrastructure/
│       └── database.go              # DB connection pool setup
├── migrations/
│   └── 001_create_products.sql      # Database schema
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/products` | Create a new product |
| GET | `/products` | List all products (optional `?keyword=` filter) |
| GET | `/products/:id` | Get a product by ID |
| PUT | `/products/:id` | Update a product |
| DELETE | `/products/:id` | Delete a product |

### Request/Response Examples

**Create Product**
```bash
POST /products
Content-Type: application/json

{
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "quantity": 10
}
```

Response (201):
```json
{
  "id": 1,
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "quantity": 10,
  "created_at": "2026-06-28T10:00:00Z",
  "updated_at": "2026-06-28T10:00:00Z"
}
```

**Search Products**
```bash
GET /products?keyword=laptop
```

**Error Response**
```json
{
  "message": "name is required"
}
```

### Validation Rules

| Field | Rule | Error Message |
|-------|------|---------------|
| name | Required (non-empty) | "name is required" |
| name | Min 3 characters | "name must be at least 3 characters" |
| price | Required (non-nil) | "price is required" |
| price | Must be > 0 | "price must be greater than 0" |
| quantity | Required (non-nil) | "quantity is required" |
| quantity | Must be >= 0 | "quantity must be greater than or equal to 0" |

### Error Codes

| Status | Meaning |
|--------|---------|
| 400 | Validation error or malformed JSON |
| 404 | Product not found or invalid ID |
| 500 | Internal server error |

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker (optional, for running PostgreSQL)

### 1. Start PostgreSQL

Using Docker:
```bash
docker run --name product-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=product_management \
  -p 5432:5432 \
  -d postgres:15
```

### 2. Run Database Migration

```bash
psql -U postgres -d product_management -f migrations/001_create_products.sql
```

### 3. Configure Environment

Copy the example env file and edit with your credentials:

```bash
cp .env.example .env
```

Edit `.env`:
```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/product_management?sslmode=disable
PORT=:8080
```

The server automatically loads `.env` on startup (using `godotenv`). You can also set environment variables directly if you prefer.

### 4. Run the Server

```bash
go run ./cmd/server/
```

The server will start on `http://localhost:8080`.

## Testing with Bruno

[Bruno](https://www.usebruno.com/) is an offline API client for testing server endpoints. The collection is pre-configured in the `bruno/` directory.

### Installation

Download and install Bruno from [https://www.usebruno.com/downloads](https://www.usebruno.com/downloads).

### Open the Collection

1. Open Bruno
2. Click **Open Collection**
3. Navigate to the `bruno/` folder in this project

### Configure the Environment

The collection uses the **local** environment with two variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `baseUrl` | `http://localhost:8080` | Base URL of the server |
| `productId` | `1` | Product ID used by GET/PUT/DELETE requests |

Select the environment from the dropdown in the top-right corner of Bruno.

> **Tip:** The **Create Product** request (seq 1) automatically saves the newly created product's `id` into the `productId` variable, so subsequent requests use the correct ID.

### Test Cases

Requests are numbered (`seq`) to follow a logical flow:

| Seq | Name | Method | Endpoint | Description |
|-----|------|--------|----------|-------------|
| 1 | Create Product | POST | `/products` | Create a product and save its `productId` to env |
| 2 | Create Product (no description) | POST | `/products` | Create a product without a description |
| 3 | Create Product - Validation Errors | POST | `/products` | Verify validation errors (empty name, negative price) |
| 4 | Create Product - Invalid JSON | POST | `/products` | Verify error on malformed JSON body |
| 5 | Get All Products | GET | `/products` | Retrieve all products |
| 6 | Get Products by Keyword | GET | `/products?keyword=laptop` | Search products by keyword |
| 7 | Get Product by ID | GET | `/products/:id` | Retrieve a product by ID |
| 8 | Get Product - Not Found | GET | `/products/999999` | Verify 404 when ID does not exist |
| 9 | Get Product - Invalid ID | GET | `/products/abc` | Verify 404 on non-numeric ID |
| 10 | Update Product | PUT | `/products/:id` | Update product details |
| 11 | Update Product - Not Found | PUT | `/products/999999` | Verify 404 when updating a non-existent product |
| 12 | Delete Product | DELETE | `/products/:id` | Delete a product |
| 13 | Delete Product - Not Found | DELETE | `/products/999999` | Verify 404 when deleting a non-existent product |

### Run the Entire Collection

To run all requests in order:

1. Right-click the **products** collection in the sidebar
2. Select **Run**
3. Bruno will execute each request and display pass/fail results per test

---

## Unit & Property-based Tests

Run all tests:
```bash
go test ./... -count=1
```

Run tests with verbose output:
```bash
go test ./... -v -count=1
```

Run only property-based tests:
```bash
go test ./internal/service/ -run TestProperty -v -count=1
```

### Test Coverage

The project uses a dual testing approach:

- **Property-based tests** (using `pgregory.net/rapid`) — verify universal correctness properties across randomized inputs (100 iterations each)
- **Unit tests** (using `testify`) — verify specific examples, edge cases, and integration points

#### Correctness Properties Tested

| # | Property | Package |
|---|----------|---------|
| 1 | Validation rejects all invalid inputs | service |
| 2 | Valid input produces a complete product | service |
| 3 | Create then get by ID round-trip | service |
| 4 | Keyword search returns only matching products sorted by ID | service |
| 5 | Update reflects new values | service |
| 6 | Delete makes product unretrievable | service |
| 7 | JSON serialization preserves format invariants | model |
| 8 | Error responses have consistent structure | handler |

## Database Schema

```sql
CREATE TABLE products (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    price       NUMERIC(12,2) NOT NULL CHECK (price > 0),
    quantity    INT NOT NULL CHECK (quantity >= 0),
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

The `updated_at` column is automatically updated via a database trigger on every row update.

## Design Decisions

- **Pointer types for DTOs**: `*float64` and `*int` in request DTOs distinguish between "field missing" (nil) and "field is zero" (0)
- **Custom JSON marshaling**: Ensures price always has 2 decimal places and timestamps use RFC 3339 with "Z" suffix
- **Interface-based DI**: Handler depends on Service interface, Service depends on Repository interface — enabling independent testing with mocks
- **First-error-wins validation**: Validation checks run in a defined order and return the first violation encountered
- **Parameterized SQL queries**: All database operations use `$1`, `$2` placeholders to prevent SQL injection
