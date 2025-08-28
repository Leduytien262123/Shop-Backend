package handle

import (
	"backend/internal/helpers"
	"backend/internal/model"
	"backend/internal/repo"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productRepo  *repo.ProductRepo
	categoryRepo *repo.CategoryRepo
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productRepo:  repo.NewProductRepo(),
		categoryRepo: repo.NewCategoryRepo(),
	}
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input model.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Normalize SKU
	input.SKU = strings.ToUpper(strings.TrimSpace(input.SKU))

	// Check if SKU already exists
	exists, err := h.productRepo.CheckSKUExists(input.SKU, 0)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}
	if exists {
		helpers.ErrorResponse(c, http.StatusConflict, "SKU already exists", errors.New("product with this SKU already exists"))
		return
	}

	// Check if category exists (if provided)
	if input.CategoryID != nil {
		_, err := h.categoryRepo.GetByID(*input.CategoryID)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid category", errors.New("category not found"))
			return
		}
	}

	product := model.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		SKU:         input.SKU,
		Stock:       input.Stock,
		CategoryID:  input.CategoryID,
		IsActive:    true,
	}

	if err := h.productRepo.Create(&product); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to create product", err)
		return
	}

	// Load product with category for response
	createdProduct, err := h.productRepo.GetByID(product.ID)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to load created product", err)
		return
	}

	c.JSON(http.StatusCreated, helpers.Response{
		Success: true,
		Message: "Product created successfully",
		Data:    createdProduct.ToResponse(),
	})
}

// GetProducts retrieves all products with pagination
func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Parse pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	categoryIDStr := c.Query("category_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	var products []model.Product
	var total int64

	if categoryIDStr != "" {
		// Get products by category
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", errors.New("category ID must be a valid number"))
			return
		}

		products, total, err = h.productRepo.GetByCategoryID(uint(categoryID), page, limit)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve products", err)
			return
		}
	} else {
		// Get all products
		products, total, err = h.productRepo.GetAll(page, limit)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve products", err)
			return
		}
	}

	var response []model.ProductResponse
	for _, product := range products {
		response = append(response, product.ToResponse())
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Products retrieved successfully",
		Data: map[string]interface{}{
			"products":     response,
			"total":        total,
			"page":         page,
			"limit":        limit,
			"total_pages":  totalPages,
			"has_next":     page < int(totalPages),
			"has_prev":     page > 1,
		},
	})
}

// GetProductByID retrieves a product by ID
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", errors.New("product ID must be a valid number"))
		return
	}

	product, err := h.productRepo.GetByID(uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			helpers.ErrorResponse(c, http.StatusNotFound, "Product not found", err)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Product retrieved successfully",
		Data:    product.ToResponse(),
	})
}

// GetProductBySKU retrieves a product by SKU
func (h *ProductHandler) GetProductBySKU(c *gin.Context) {
	sku := c.Param("sku")
	if sku == "" {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid SKU", errors.New("SKU cannot be empty"))
		return
	}

	product, err := h.productRepo.GetBySKU(sku)
	if err != nil {
		if err.Error() == "product not found" {
			helpers.ErrorResponse(c, http.StatusNotFound, "Product not found", err)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Product retrieved successfully",
		Data:    product.ToResponse(),
	})
}

// UpdateProduct updates a product
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", errors.New("product ID must be a valid number"))
		return
	}

	var input model.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Get existing product
	product, err := h.productRepo.GetByID(uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			helpers.ErrorResponse(c, http.StatusNotFound, "Product not found", err)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	// Normalize SKU
	input.SKU = strings.ToUpper(strings.TrimSpace(input.SKU))

	// Check if SKU already exists (excluding current product)
	exists, err := h.productRepo.CheckSKUExists(input.SKU, uint(id))
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}
	if exists {
		helpers.ErrorResponse(c, http.StatusConflict, "SKU already exists", errors.New("another product with this SKU already exists"))
		return
	}

	// Check if category exists (if provided)
	if input.CategoryID != nil {
		_, err := h.categoryRepo.GetByID(*input.CategoryID)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid category", errors.New("category not found"))
			return
		}
	}

	// Update product
	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.SKU = input.SKU
	product.Stock = input.Stock
	product.CategoryID = input.CategoryID

	if err := h.productRepo.Update(product); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to update product", err)
		return
	}

	// Load updated product with category
	updatedProduct, err := h.productRepo.GetByID(product.ID)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to load updated product", err)
		return
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Product updated successfully",
		Data:    updatedProduct.ToResponse(),
	})
}

// DeleteProduct deletes a product
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", errors.New("product ID must be a valid number"))
		return
	}

	// Check if product exists
	_, err = h.productRepo.GetByID(uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			helpers.ErrorResponse(c, http.StatusNotFound, "Product not found", err)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	if err := h.productRepo.Delete(uint(id)); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete product", err)
		return
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Product deleted successfully",
		Data:    nil,
	})
}

// UpdateProductStock updates product stock
func (h *ProductHandler) UpdateProductStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", errors.New("product ID must be a valid number"))
		return
	}

	var input struct {
		Stock int `json:"stock" binding:"required,gte=0"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Check if product exists
	product, err := h.productRepo.GetByID(uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			helpers.ErrorResponse(c, http.StatusNotFound, "Product not found", err)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	if err := h.productRepo.UpdateStock(uint(id), input.Stock); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to update stock", err)
		return
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Product stock updated successfully",
		Data: map[string]interface{}{
			"product_id": product.ID,
			"old_stock":  product.Stock,
			"new_stock":  input.Stock,
		},
	})
}
