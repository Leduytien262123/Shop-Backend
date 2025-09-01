package repo

import (
	"backend/app"
	"backend/internal/model"
	"errors"

	"gorm.io/gorm"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo() *CategoryRepo {
	return &CategoryRepo{
		db: app.GetDB(),
	}
}

// Create tạo mới một danh mục
func (r *CategoryRepo) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

// GetByID lấy danh mục theo ID
func (r *CategoryRepo) GetByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// GetBySlug lấy danh mục theo slug
func (r *CategoryRepo) GetBySlug(slug string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// GetAll lấy tất cả danh mục đang hoạt động (is_active = true)
func (r *CategoryRepo) GetAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("is_active = ?", true).Find(&categories).Error
	return categories, err
}

// GetAllWithProducts lấy tất cả danh mục kèm theo danh sách sản phẩm
func (r *CategoryRepo) GetAllWithProducts() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Preload("Products").Where("is_active = ?", true).Find(&categories).Error
	return categories, err
}

// Update cập nhật danh mục
func (r *CategoryRepo) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

// Delete xóa mềm một danh mục
func (r *CategoryRepo) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

// CheckSlugExists kiểm tra slug đã tồn tại hay chưa (loại trừ danh mục có ID = excludeID)
func (r *CategoryRepo) CheckSlugExists(slug string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.Category{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
