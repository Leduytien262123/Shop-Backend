# API Testing Guide - Postman Collection

## Base URL

```
http://localhost:8080
```

## Authentication

Đối với các API admin, bạn cần:

1. Đăng ký/đăng nhập để lấy JWT token
2. Thêm header: `Authorization: Bearer <your_jwt_token>`
3. Đảm bảo user có role "admin"

---

## 1. CATEGORY APIs

### Public Endpoints (không cần authentication)

#### GET - Lấy tất cả danh mục

```
GET http://localhost:8080/api/categories/
```

#### GET - Lấy tất cả danh mục kèm sản phẩm

```
GET http://localhost:8080/api/categories/?with_products=true
```

#### GET - Lấy danh mục theo ID

```
GET http://localhost:8080/api/categories/1
```

#### GET - Lấy danh mục theo slug

```
GET http://localhost:8080/api/categories/slug/electronics
```

### Admin Endpoints (cần authentication + admin role)

#### POST - Tạo danh mục mới

```
POST http://localhost:8080/api/admin/categories/
Content-Type: application/json
Authorization: Bearer <your_jwt_token>

{
  "name": "Electronics",
  "description": "Electronic devices and gadgets",
  "slug": "electronics"
}
```

#### PUT - Cập nhật danh mục

```
PUT http://localhost:8080/api/admin/categories/1
Content-Type: application/json
Authorization: Bearer <your_jwt_token>

{
  "name": "Electronics Updated",
  "description": "Updated description for electronic devices",
  "slug": "electronics-updated"
}
```

#### DELETE - Xóa danh mục

```
DELETE http://localhost:8080/api/admin/categories/1
Authorization: Bearer <your_jwt_token>
```

---

## 2. PRODUCT APIs

### Public Endpoints (không cần authentication)

#### GET - Lấy tất cả sản phẩm (có phân trang)

```
GET http://localhost:8080/api/products/
```

#### GET - Lấy sản phẩm với phân trang

```
GET http://localhost:8080/api/products/?page=1&limit=10
```

#### GET - Lấy sản phẩm theo danh mục

```
GET http://localhost:8080/api/products/?category_id=1&page=1&limit=10
```

#### GET - Lấy sản phẩm theo ID

```
GET http://localhost:8080/api/products/1
```

#### GET - Lấy sản phẩm theo SKU

```
GET http://localhost:8080/api/products/sku/LAPTOP001
```

### Admin Endpoints (cần authentication + admin role)

#### POST - Tạo sản phẩm mới

```
POST http://localhost:8080/api/admin/products/
Content-Type: application/json
Authorization: Bearer <your_jwt_token>

{
  "name": "MacBook Air M2",
  "description": "Apple MacBook Air with M2 chip, 13-inch display",
  "price": 1299.99,
  "sku": "LAPTOP001",
  "stock": 50,
  "category_id": 1
}
```

#### POST - Tạo sản phẩm không có danh mục

```
POST http://localhost:8080/api/admin/products/
Content-Type: application/json
Authorization: Bearer <your_jwt_token>

{
  "name": "Generic Product",
  "description": "A product without category",
  "price": 99.99,
  "sku": "GENERIC001",
  "stock": 100
}
```

#### PUT - Cập nhật sản phẩm

```
PUT http://localhost:8080/api/admin/products/1
Content-Type: application/json
Authorization: Bearer <your_jwt_token>

{
  "name": "MacBook Air M2 Updated",
  "description": "Updated Apple MacBook Air with M2 chip",
  "price": 1199.99,
  "sku": "LAPTOP001",
  "stock": 75,
  "category_id": 1
}
```

#### PATCH - Cập nhật số lượng tồn kho

```
PATCH http://localhost:8080/api/admin/products/1/stock
Content-Type: application/json
Authorization: Bearer <your_jwt_token>

{
  "stock": 25
}
```

#### DELETE - Xóa sản phẩm

```
DELETE http://localhost:8080/api/admin/products/1
Authorization: Bearer <your_jwt_token>
```

---

## 3. AUTH APIs (để lấy token)

#### POST - Đăng ký user mới

```
POST http://localhost:8080/api/auth/register
Content-Type: application/json

{
  "username": "admin",
  "email": "admin@example.com",
  "password": "123456",
  "full_name": "Administrator"
}
```

#### POST - Đăng nhập

```
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "123456"
}
```

Response sẽ có token:

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {...}
  }
}
```

---

## 4. Test Flow Suggest

### Bước 1: Tạo Admin User

1. Đăng ký user mới qua `/api/auth/register`
2. Đăng nhập qua `/api/auth/login` để lấy token
3. Cập nhật role thành "admin" trong database hoặc qua admin API

### Bước 2: Test Category APIs

1. Tạo danh mục "Electronics"
2. Tạo danh mục "Books"
3. Lấy danh sách tất cả danh mục
4. Cập nhật một danh mục
5. Lấy danh mục theo slug

### Bước 3: Test Product APIs

1. Tạo sản phẩm trong danh mục Electronics
2. Tạo sản phẩm không có danh mục
3. Lấy danh sách sản phẩm với phân trang
4. Lấy sản phẩm theo danh mục
5. Cập nhật thông tin sản phẩm
6. Cập nhật số lượng tồn kho

### Bước 4: Test Relationships

1. Lấy danh mục kèm sản phẩm (`/api/categories/?with_products=true`)
2. Kiểm tra foreign key constraints
3. Test xóa danh mục có sản phẩm (sản phẩm sẽ có category_id = null)

---

## 5. Sample Data for Testing

### Categories

```json
[
  {
    "name": "Electronics",
    "description": "Electronic devices and gadgets",
    "slug": "electronics"
  },
  {
    "name": "Books",
    "description": "Books and educational materials",
    "slug": "books"
  },
  {
    "name": "Clothing",
    "description": "Fashion and clothing items",
    "slug": "clothing"
  }
]
```

### Products

```json
[
  {
    "name": "iPhone 15 Pro",
    "description": "Latest iPhone with A17 Pro chip",
    "price": 999.99,
    "sku": "PHONE001",
    "stock": 100,
    "category_id": 1
  },
  {
    "name": "MacBook Air M2",
    "description": "Apple MacBook Air with M2 chip",
    "price": 1299.99,
    "sku": "LAPTOP001",
    "stock": 50,
    "category_id": 1
  },
  {
    "name": "Programming Book",
    "description": "Learn Go programming language",
    "price": 49.99,
    "sku": "BOOK001",
    "stock": 200,
    "category_id": 2
  }
]
```

---

## 6. Error Handling

API sẽ trả về các lỗi phổ biến:

- `400 Bad Request`: Input không hợp lệ
- `401 Unauthorized`: Thiếu token hoặc token không hợp lệ
- `403 Forbidden`: Không có quyền admin
- `404 Not Found`: Resource không tồn tại
- `409 Conflict`: Slug hoặc SKU đã tồn tại
- `500 Internal Server Error`: Lỗi server

---

## 7. Health Check

```
GET http://localhost:8080/health
```

Response:

```json
{
  "status": "OK",
  "message": "Server is running"
}
```
