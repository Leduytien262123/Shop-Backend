package repo

import (
	"backend/app"
	"backend/internal/model"
	"errors"

	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{
		db: app.GetDB(),
	}
}

// Create creates a new product
func (r *ProductRepo) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

// GetByID retrieves a product by ID with category
func (r *ProductRepo) GetByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// GetBySKU retrieves a product by SKU
func (r *ProductRepo) GetBySKU(sku string) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").Where("sku = ?", sku).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// GetAll retrieves all active products with pagination
func (r *ProductRepo) GetAll(page, limit int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	// Count total records
	if err := r.db.Model(&model.Product{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get products with category
	err := r.db.Preload("Category").
		Where("is_active = ?", true).
		Offset(offset).
		Limit(limit).
		Find(&products).Error

	return products, total, err
}

// GetByCategoryID retrieves products by category ID
func (r *ProductRepo) GetByCategoryID(categoryID uint, page, limit int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	// Count total records
	query := r.db.Model(&model.Product{}).Where("category_id = ? AND is_active = ?", categoryID, true)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get products with category
	err := r.db.Preload("Category").
		Where("category_id = ? AND is_active = ?", categoryID, true).
		Offset(offset).
		Limit(limit).
		Find(&products).Error

	return products, total, err
}

// Update updates a product
func (r *ProductRepo) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

// Delete soft deletes a product
func (r *ProductRepo) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

// CheckSKUExists checks if a SKU already exists (for different product)
func (r *ProductRepo) CheckSKUExists(sku string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.Product{}).Where("sku = ?", sku)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// UpdateStock updates product stock
func (r *ProductRepo) UpdateStock(id uint, stock int) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).Update("stock", stock).Error
}
