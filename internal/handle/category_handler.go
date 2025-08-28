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

type CategoryHandler struct {
	categoryRepo *repo.CategoryRepo
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categoryRepo: repo.NewCategoryRepo(),
	}
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input model.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Normalize slug
	input.Slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(input.Slug), " ", "-"))

	// Check if slug already exists
	exists, err := h.categoryRepo.CheckSlugExists(input.Slug, 0)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}
	if exists {
		helpers.ErrorResponse(c, http.StatusConflict, "Slug already exists", errors.New("category with this slug already exists"))
		return
	}

	category := model.Category{
		Name:        input.Name,
		Description: input.Description,
		Slug:        input.Slug,
		IsActive:    true,
	}

	if err := h.categoryRepo.Create(&category); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category", err)
		return
	}

	c.JSON(http.StatusCreated, helpers.Response{
		Success: true,
		Message: "Category created successfully",
		Data:    category.ToResponse(),
	})
}

// GetCategories retrieves all categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	withProducts := c.Query("with_products") == "true"

	var categories []model.Category
	var err error

	if withProducts {
		categories, err = h.categoryRepo.GetAllWithProducts()
	} else {
		categories, err = h.categoryRepo.GetAll()
	}

	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve categories", err)
		return
	}

	var response []model.CategoryResponse
	for _, category := range categories {
		response = append(response, category.ToResponse())
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Categories retrieved successfully",
		Data:    response,
	})
}

// GetCategoryByID retrieves a category by ID
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", errors.New("category ID must be a valid number"))
		return
	}

	category, err := h.categoryRepo.GetByID(uint(id))
	if err != nil {
		if err.Error() == "category not found" {
			helpers.ErrorResponse(c, http.StatusNotFound, "Category not found", err)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Category retrieved successfully",
		Data:    category.ToResponse(),
	})
}

// GetCategoryBySlug retrieves a category by slug
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid slug", errors.New("slug cannot be empty"))
		return
	}

	category, err := h.categoryRepo.GetBySlug(slug)
	if err != nil {
		if err.Error() == "category not found" {
			helpers.ErrorResponse(c, http.StatusNotFound, "Category not found", err)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Category retrieved successfully",
		Data:    category.ToResponse(),
	})
}

// UpdateCategory updates a category
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", errors.New("category ID must be a valid number"))
		return
	}

	var input model.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Get existing category
	category, err := h.categoryRepo.GetByID(uint(id))
	if err != nil {
		if err.Error() == "category not found" {
			helpers.ErrorResponse(c, http.StatusNotFound, "Category not found", err)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	// Normalize slug
	input.Slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(input.Slug), " ", "-"))

	// Check if slug already exists (excluding current category)
	exists, err := h.categoryRepo.CheckSlugExists(input.Slug, uint(id))
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}
	if exists {
		helpers.ErrorResponse(c, http.StatusConflict, "Slug already exists", errors.New("another category with this slug already exists"))
		return
	}

	// Update category
	category.Name = input.Name
	category.Description = input.Description
	category.Slug = input.Slug

	if err := h.categoryRepo.Update(category); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to update category", err)
		return
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Category updated successfully",
		Data:    category.ToResponse(),
	})
}

// DeleteCategory deletes a category
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", errors.New("category ID must be a valid number"))
		return
	}

	// Check if category exists
	_, err = h.categoryRepo.GetByID(uint(id))
	if err != nil {
		if err.Error() == "category not found" {
			helpers.ErrorResponse(c, http.StatusNotFound, "Category not found", err)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	if err := h.categoryRepo.Delete(uint(id)); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete category", err)
		return
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Category deleted successfully",
		Data:    nil,
	})
}
