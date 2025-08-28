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

// Create creates a new category
func (r *CategoryRepo) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

// GetByID retrieves a category by ID
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

// GetBySlug retrieves a category by slug
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

// GetAll retrieves all active categories
func (r *CategoryRepo) GetAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("is_active = ?", true).Find(&categories).Error
	return categories, err
}

// GetAllWithProducts retrieves all categories with their products
func (r *CategoryRepo) GetAllWithProducts() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Preload("Products").Where("is_active = ?", true).Find(&categories).Error
	return categories, err
}

// Update updates a category
func (r *CategoryRepo) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

// Delete soft deletes a category
func (r *CategoryRepo) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

// CheckSlugExists checks if a slug already exists (for different category)
func (r *CategoryRepo) CheckSlugExists(slug string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.Category{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
