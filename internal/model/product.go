package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"not null;size:200;index"`
	Description string         `json:"description" gorm:"type:text"`
	Price       float64        `json:"price" gorm:"not null;type:decimal(10,2)"`
	SKU         string         `json:"sku" gorm:"unique;not null;size:50;index"`
	Stock       int            `json:"stock" gorm:"not null;default:0"`
	CategoryID  *uint          `json:"category_id" gorm:"index"`
	IsActive    bool           `json:"is_active" gorm:"default:true;index"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName specifies the table name for Product model
func (Product) TableName() string {
	return "products"
}

type ProductInput struct {
	Name        string  `json:"name" binding:"required,min=1,max=200"`
	Description string  `json:"description" binding:"max=1000"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	SKU         string  `json:"sku" binding:"required,min=1,max=50"`
	Stock       int     `json:"stock" binding:"gte=0"`
	CategoryID  *uint   `json:"category_id"`
}

type ProductResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	SKU         string           `json:"sku"`
	Stock       int              `json:"stock"`
	CategoryID  *uint            `json:"category_id"`
	Category    *CategoryResponse `json:"category,omitempty"`
	IsActive    bool             `json:"is_active"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// ToResponse converts Product to ProductResponse
func (p *Product) ToResponse() ProductResponse {
	response := ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		SKU:         p.SKU,
		Stock:       p.Stock,
		CategoryID:  p.CategoryID,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	// Include category information if loaded
	if p.Category != nil {
		categoryResponse := p.Category.ToResponse()
		response.Category = &categoryResponse
	}

	return response
}
