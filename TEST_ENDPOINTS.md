# TEST API ENDPOINTS (Authentication DISABLED)

## Base URL: http://localhost:8080

## ✅ CATEGORY TESTS

### 1. CREATE Category (POST)

**URL:** http://localhost:8080/api/admin/categories/
**Method:** POST
**Headers:** Content-Type: application/json
**Body:**

```json
{
  "name": "Electronics",
  "description": "Electronic devices and gadgets",
  "slug": "electronics"
}
```

### 2. CREATE Another Category

**URL:** http://localhost:8080/api/admin/categories/
**Method:** POST
**Headers:** Content-Type: application/json
**Body:**

```json
{
  "name": "Books",
  "description": "Books and educational materials",
  "slug": "books"
}
```

### 3. GET All Categories

**URL:** http://localhost:8080/api/categories/
**Method:** GET

### 4. GET Category by ID

**URL:** http://localhost:8080/api/categories/1
**Method:** GET

### 5. UPDATE Category

**URL:** http://localhost:8080/api/admin/categories/1
**Method:** PUT
**Headers:** Content-Type: application/json
**Body:**

```json
{
  "name": "Electronics & Gadgets",
  "description": "Updated: Electronic devices, gadgets and smart technology",
  "slug": "electronics-gadgets"
}
```

---

## ✅ PRODUCT TESTS

### 1. CREATE Product with Category

**URL:** http://localhost:8080/api/admin/products/
**Method:** POST
**Headers:** Content-Type: application/json
**Body:**

```json
{
  "name": "iPhone 15 Pro",
  "description": "Apple iPhone 15 Pro with A17 Pro chip",
  "price": 999.99,
  "sku": "IPHONE15PRO",
  "stock": 50,
  "category_id": 1
}
```

### 2. CREATE Product without Category

**URL:** http://localhost:8080/api/admin/products/
**Method:** POST
**Headers:** Content-Type: application/json
**Body:**

```json
{
  "name": "Generic Wireless Charger",
  "description": "Universal wireless charging pad",
  "price": 29.99,
  "sku": "WC001",
  "stock": 100
}
```

### 3. GET All Products

**URL:** http://localhost:8080/api/products/
**Method:** GET

### 4. GET Products with Pagination

**URL:** http://localhost:8080/api/products/?page=1&limit=5
**Method:** GET

### 5. GET Product by ID

**URL:** http://localhost:8080/api/products/1
**Method:** GET

### 6. UPDATE Product

**URL:** http://localhost:8080/api/admin/products/1
**Method:** PUT
**Headers:** Content-Type: application/json
**Body:**

```json
{
  "name": "iPhone 15 Pro - Updated",
  "description": "Apple iPhone 15 Pro with A17 Pro chip - Special offer",
  "price": 899.99,
  "sku": "IPHONE15PRO",
  "stock": 75,
  "category_id": 1
}
```

### 7. UPDATE Product Stock

**URL:** http://localhost:8080/api/admin/products/1/stock
**Method:** PATCH
**Headers:** Content-Type: application/json
**Body:**

```json
{
  "stock": 25
}
```

---

## 🔄 RECOMMENDED TEST FLOW:

1. **Create Categories first:**

   - Electronics
   - Books
   - Clothing

2. **Create Products:**

   - Some with category_id
   - Some without category_id (null)

3. **Test Relationships:**

   - GET categories with products: /api/categories/?with_products=true
   - GET products by category: /api/products/?category_id=1

4. **Test Updates:**

   - Update category info
   - Update product info
   - Update stock

5. **Test Validations:**
   - Try duplicate SKU (should fail)
   - Try duplicate slug (should fail)
   - Try negative price (should fail)

---

## 📝 NOTES:

- Authentication is TEMPORARILY DISABLED for testing
- All admin endpoints are accessible without token
- Remember to re-enable authentication later by uncommenting middleware lines
- Server is running on port 8080
- Database auto-creates tables on startup
