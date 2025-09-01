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

// CreateCategory tạo danh mục mới
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input model.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ", err)
		return
	}

	// Chuẩn hóa slug
	input.Slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(input.Slug), " ", "-"))

	// Kiểm tra xem slug đã tồn tại chưa
	exists, err := h.categoryRepo.CheckSlugExists(input.Slug, 0)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Lỗi cơ sở dữ liệu", err)
		return
	}
	if exists {
		helpers.ErrorResponse(c, http.StatusConflict, "Slug đã tồn tại", errors.New("danh mục với slug này đã tồn tại"))
		return
	}

	category := model.Category{
		Name:        input.Name,
		Description: input.Description,
		Slug:        input.Slug,
		IsActive:    true,
	}

	if err := h.categoryRepo.Create(&category); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Không thể tạo danh mục", err)
		return
	}

	c.JSON(http.StatusCreated, helpers.Response{
		Success: true,
		Message: "Tạo danh mục thành công",
		Data:    category.ToResponse(),
	})
}

// GetCategories lấy tất cả danh mục
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
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Không thể lấy danh sách danh mục", err)
		return
	}

	var response []model.CategoryResponse
	for _, category := range categories {
		response = append(response, category.ToResponse())
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Lấy danh sách danh mục thành công",
		Data:    response,
	})
}

// GetCategoryByID lấy danh mục theo ID
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "ID danh mục không hợp lệ", errors.New("ID danh mục phải là số hợp lệ"))
		return
	}

	category, err := h.categoryRepo.GetByID(uint(id))
	if err != nil {
		if err.Error() == "category not found" {
			helpers.ErrorResponse(c, http.StatusNotFound, "Không tìm thấy danh mục", err)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Lỗi cơ sở dữ liệu", err)
		return
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Lấy thông tin danh mục thành công",
		Data:    category.ToResponse(),
	})
}

// GetCategoryBySlug lấy danh mục theo slug
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Slug không hợp lệ", errors.New("slug không được để trống"))
		return
	}

	category, err := h.categoryRepo.GetBySlug(slug)
	if err != nil {
		if err.Error() == "category not found" {
			helpers.ErrorResponse(c, http.StatusNotFound, "Không tìm thấy danh mục", err)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Lỗi cơ sở dữ liệu", err)
		return
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Lấy thông tin danh mục thành công",
		Data:    category.ToResponse(),
	})
}

// UpdateCategory cập nhật danh mục
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "ID danh mục không hợp lệ", errors.New("ID danh mục phải là số hợp lệ"))
		return
	}

	var input model.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ", err)
		return
	}

	// Lấy danh mục hiện tại
	category, err := h.categoryRepo.GetByID(uint(id))
	if err != nil {
		if err.Error() == "category not found" {
			helpers.ErrorResponse(c, http.StatusNotFound, "Không tìm thấy danh mục", err)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Lỗi cơ sở dữ liệu", err)
		return
	}

	// Chuẩn hóa slug
	input.Slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(input.Slug), " ", "-"))

	// Kiểm tra xem slug đã tồn tại chưa (loại trừ danh mục hiện tại)
	exists, err := h.categoryRepo.CheckSlugExists(input.Slug, uint(id))
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Lỗi cơ sở dữ liệu", err)
		return
	}
	if exists {
		helpers.ErrorResponse(c, http.StatusConflict, "Slug đã tồn tại", errors.New("danh mục khác với slug này đã tồn tại"))
		return
	}

	// Cập nhật danh mục
	category.Name = input.Name
	category.Description = input.Description
	category.Slug = input.Slug

	if err := h.categoryRepo.Update(category); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Không thể cập nhật danh mục", err)
		return
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Cập nhật danh mục thành công",
		Data:    category.ToResponse(),
	})
}

// DeleteCategory xóa danh mục
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "ID danh mục không hợp lệ", errors.New("ID danh mục phải là số hợp lệ"))
		return
	}

	// Kiểm tra xem danh mục có tồn tại không
	_, err = h.categoryRepo.GetByID(uint(id))
	if err != nil {
		if err.Error() == "category not found" {
			helpers.ErrorResponse(c, http.StatusNotFound, "Không tìm thấy danh mục", err)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Lỗi cơ sở dữ liệu", err)
		return
	}

	if err := h.categoryRepo.Delete(uint(id)); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Không thể xóa danh mục", err)
		return
	}

	c.JSON(http.StatusOK, helpers.Response{
		Success: true,
		Message: "Xóa danh mục thành công",
		Data:    nil,
	})
}
