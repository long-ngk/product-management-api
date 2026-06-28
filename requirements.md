# Assignment: Build a Simple Product Management API

## 1. Mục tiêu
Cần xây dựng một REST API đơn giản bằng Go để quản lý danh sách sản phẩm.
Assignment này giúp thực hành các kiến thức đã học:
- Go project structure
- Variables, structs, methods
- Functions
- Error handling
- HTTP JSON API
- Database connection
- CRUD operations
- Modular code organization

## 2. Business Context
Ứng dụng mô phỏng một hệ thống nhỏ dùng để quản lý sản phẩm trong kho.
Mỗi sản phẩm có các thông tin cơ bản:
- ID
- Name
- Description
- Price
- Quantity
- CreatedAt
- UpdatedAt

Người dùng có thể:
- Tạo sản phẩm mới
- Xem danh sách sản phẩm
- Xem chi tiết một sản phẩm
- Cập nhật thông tin sản phẩm
- Xóa sản phẩm

## 3. Technical Requirements
### Backend
Học viên cần xây dựng backend API bằng:
- Go
- Gin framework hoặc net/http
- PostgreSQL hoặc MySQL
- JSON request/response

**Khuyến nghị cho training:** Gin + PostgreSQL

## 4. Database Requirement
Tạo database table `products`.

### Table: products

#### Nếu dùng PostgreSQL:
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(12, 2) NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

#### Nếu dùng MySQL:
```sql
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(12, 2) NOT NULL,
    quantity INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 5. API Requirements

### 5.1. Create Product
* **Method & URL:** `POST /products`
* **Request body:**
```json
{
  "name": "Mechanical Keyboard",
  "description": "Wireless mechanical keyboard",
  "price": 120.50,
  "quantity": 10
}
```
* **Expected response:**
```json
{
  "id": 1,
  "name": "Mechanical Keyboard",
  "description": "Wireless mechanical keyboard",
  "price": 120.50,
  "quantity": 10,
  "created_at": "2026-06-21T10:00:00Z",
  "updated_at": "2026-06-21T10:00:00Z"
}
```

### 5.2. Get Product List
* **Method & URL:** `GET /products`
* **Optional query parameters:** `GET /products?keyword=keyboard`
* **Expected response:**
```json
[
  {
    "id": 1,
    "name": "Mechanical Keyboard",
    "description": "Wireless mechanical keyboard",
    "price": 120.50,
    "quantity": 10,
    "created_at": "2026-06-21T10:00:00Z",
    "updated_at": "2026-06-21T10:00:00Z"
  }
]
```

### 5.3. Get Product Detail
* **Method & URL:** `GET /products/:id`
* **Example:** `GET /products/1`
* **Expected response:**
```json
{
  "id": 1,
  "name": "Mechanical Keyboard",
  "description": "Wireless mechanical keyboard",
  "price": 120.50,
  "quantity": 10,
  "created_at": "2026-06-21T10:00:00Z",
  "updated_at": "2026-06-21T10:00:00Z"
}
```
* **If product not found:**
```json
{
  "message": "product not found"
}
```

### 5.4. Update Product
* **Method & URL:** `PUT /products/:id`
* **Request body:**
```json
{
  "name": "Updated Keyboard",
  "description": "Updated product description",
  "price": 135.00,
  "quantity": 15
}
```
* **Expected response:**
```json
{
  "id": 1,
  "name": "Updated Keyboard",
  "description": "Updated product description",
  "price": 135.00,
  "quantity": 15,
  "created_at": "2026-06-21T10:00:00Z",
  "updated_at": "2026-06-21T10:15:00Z"
}
```

### 5.5. Delete Product
* **Method & URL:** `DELETE /products/:id`
* **Expected response:**
```json
{
  "message": "product deleted successfully"
}
```

## 6. Validation Requirements
Học viên cần validate dữ liệu đầu vào.

### Create / Update Product
* **name:**
  - Required
  - Minimum 3 characters
* **price:**
  - Required
  - Must be greater than 0
* **quantity:**
  - Required
  - Must be greater than or equal to 0

**Ví dụ response khi request sai:**
```json
{
  "message": "name is required"
}
```
Hoặc:
```json
{
  "message": "price must be greater than 0"
}
```

## 7. Recommended Project Structure
```text
go-crud-assignment/
├── go.mod
├── go.sum
├── main.go
├── README.md
├── config/
│   └── config.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── routes/
│   │   └── routes.go
│   ├── handlers/
│   │   └── product_handler.go
│   ├── services/
│   │   └── product_service.go
│   ├── repositories/
│   │   └── product_repository.go
│   ├── models/
│   │   └── product_model.go
│   └── infrastructure/
│       └── database.go
└── migrations/
    └── 001_create_products_table.sql
```

## 8. Layer Responsibilities

### Handler Layer
- Receive HTTP request
- Parse path parameters
- Decode JSON body
- Call service layer
- Return JSON response

### Service Layer
- Handle business logic
- Validate input
- Decide what error should be returned
- Call repository layer

### Repository Layer
- Connect to database
- Execute SQL queries
- Insert, update, delete, select product data

### Infrastructure Layer
- Initialize database connection
- Manage database configuration

## 9. Suggested API Flow
```text
Client 
  ↓ 
Gin Router 
  ↓
Product Handler 
  ↓ 
Product Service 
  ↓ 
Product Repository 
  ↓ 
Database
```

## 10. Required Deliverables
Học viên cần nộp:
1. Source code trên GitHub/GitLab
2. README hướng dẫn chạy project
3. SQL script tạo database/table
4. Curl commands hoặc Postman collection để test API
5. Screenshot kết quả chạy API hoặc database data

## 11. README cần có gì?
File `README.md` nên bao gồm:
- Project description
- Tech stack
- Project structure
- How to run database
- How to run application
- API endpoints
- Example curl commands
- Common errors

## 12. Example curl Commands

### Create Product
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mechanical Keyboard",
    "description": "Wireless mechanical keyboard",
    "price": 120.50,
    "quantity": 10
  }'
```

### Get Products
```bash
curl http://localhost:8080/products
```

### Get Product Detail
```bash
curl http://localhost:8080/products/1
```

### Update Product
```bash
curl -X PUT http://localhost:8080/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Keyboard",
    "description": "Updated description",
    "price": 135.00,
    "quantity": 15
  }'
```

### Delete Product
```bash
curl -X DELETE http://localhost:8080/products/1
```

## 13. Evaluation Criteria

| Criteria | Description |
| :--- | :--- |
| **Project Structure** | Code được chia layer rõ ràng |
| **CRUD Functionality** | Đủ create, read, update, delete |
| **Database Integration** | Kết nối và thao tác được với database |
| **JSON API** | Request/response đúng JSON format |
| **Error Handling** | Có xử lý lỗi cơ bản |
| **Validation** | Có validate input |
| **Code Readability** | Code dễ đọc, đặt tên rõ ràng |
| **README** | Có hướng dẫn chạy và test API |