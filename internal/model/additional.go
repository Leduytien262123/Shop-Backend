package model

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID uint           `json:"product_id" gorm:"not null;index"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Rating    int            `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Title     string         `json:"title" gorm:"size:200"`
	Comment   string         `json:"comment" gorm:"type:text"`
	IsActive  bool           `json:"is_active" gorm:"default:true;index"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Coupon struct {
	ID               uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Code             string         `json:"code" gorm:"unique;not null;size:50;index"`
	Name             string         `json:"name" gorm:"not null;size:200"`
	Description      string         `json:"description" gorm:"size:500"`
	Type             string         `json:"type" gorm:"not null;size:20;check:type IN ('percentage', 'fixed')"`
	Value            float64        `json:"value" gorm:"not null;type:decimal(10,2)"`
	MinOrderAmount   float64        `json:"min_order_amount" gorm:"type:decimal(10,2);default:0"`
	MaxDiscountValue float64        `json:"max_discount_value" gorm:"type:decimal(10,2)"`
	UsageLimit       int            `json:"usage_limit" gorm:"default:0"`
	UsedCount        int            `json:"used_count" gorm:"default:0"`
	IsActive         bool           `json:"is_active" gorm:"default:true;index"`
	StartDate        time.Time      `json:"start_date" gorm:"not null"`
	EndDate          time.Time      `json:"end_date" gorm:"not null"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

type Address struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       uint           `json:"user_id" gorm:"not null;index"`
	Name         string         `json:"name" gorm:"not null;size:100"`
	Phone        string         `json:"phone" gorm:"not null;size:20"`
	AddressLine1 string         `json:"address_line1" gorm:"not null;size:200"`
	AddressLine2 string         `json:"address_line2" gorm:"size:200"`
	City         string         `json:"city" gorm:"not null;size:100"`
	State        string         `json:"state" gorm:"not null;size:100"`
	PostalCode   string         `json:"postal_code" gorm:"not null;size:20"`
	Country      string         `json:"country" gorm:"not null;size:100;default:Vietnam"`
	IsDefault    bool           `json:"is_default" gorm:"default:false;index"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Brand struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"not null;size:100;index"`
	Slug        string         `json:"slug" gorm:"unique;not null;size:100;index"`
	Description string         `json:"description" gorm:"size:500"`
	LogoURL     string         `json:"logo_url" gorm:"size:500"`
	Website     string         `json:"website" gorm:"size:200"`
	IsActive    bool           `json:"is_active" gorm:"default:true;index"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Products []Product `json:"products,omitempty" gorm:"foreignKey:BrandID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type ProductImage struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID uint           `json:"product_id" gorm:"not null;index"`
	ImageURL  string         `json:"image_url" gorm:"not null;size:500"`
	AltText   string         `json:"alt_text" gorm:"size:200"`
	IsPrimary bool           `json:"is_primary" gorm:"default:false;index"`
	SortOrder int            `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName specifies the table names
func (Review) TableName() string { return "reviews" }
func (Coupon) TableName() string { return "coupons" }
func (Address) TableName() string { return "addresses" }
func (Brand) TableName() string { return "brands" }
func (ProductImage) TableName() string { return "product_images" }

// Input structs
type ReviewInput struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Title     string `json:"title" binding:"max=200"`
	Comment   string `json:"comment" binding:"max=1000"`
}

type CouponInput struct {
	Code             string    `json:"code" binding:"required,min=1,max=50"`
	Name             string    `json:"name" binding:"required,min=1,max=200"`
	Description      string    `json:"description" binding:"max=500"`
	Type             string    `json:"type" binding:"required,oneof=percentage fixed"`
	Value            float64   `json:"value" binding:"required,gt=0"`
	MinOrderAmount   float64   `json:"min_order_amount" binding:"gte=0"`
	MaxDiscountValue float64   `json:"max_discount_value" binding:"gte=0"`
	UsageLimit       int       `json:"usage_limit" binding:"gte=0"`
	StartDate        time.Time `json:"start_date" binding:"required"`
	EndDate          time.Time `json:"end_date" binding:"required"`
}

type AddressInput struct {
	Name         string `json:"name" binding:"required,min=1,max=100"`
	Phone        string `json:"phone" binding:"required,min=1,max=20"`
	AddressLine1 string `json:"address_line1" binding:"required,min=1,max=200"`
	AddressLine2 string `json:"address_line2" binding:"max=200"`
	City         string `json:"city" binding:"required,min=1,max=100"`
	State        string `json:"state" binding:"required,min=1,max=100"`
	PostalCode   string `json:"postal_code" binding:"required,min=1,max=20"`
	Country      string `json:"country" binding:"max=100"`
	IsDefault    bool   `json:"is_default"`
}

type BrandInput struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Slug        string `json:"slug" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
	LogoURL     string `json:"logo_url" binding:"max=500"`
	Website     string `json:"website" binding:"max=200"`
}

type ProductImageInput struct {
	ProductID uint   `json:"product_id" binding:"required"`
	ImageURL  string `json:"image_url" binding:"required,max=500"`
	AltText   string `json:"alt_text" binding:"max=200"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int    `json:"sort_order"`
}

// Response structs
type ReviewResponse struct {
	ID        uint          `json:"id"`
	ProductID uint          `json:"product_id"`
	UserID    uint          `json:"user_id"`
	User      *UserResponse `json:"user,omitempty"`
	Rating    int           `json:"rating"`
	Title     string        `json:"title"`
	Comment   string        `json:"comment"`
	IsActive  bool          `json:"is_active"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type CouponResponse struct {
	ID               uint      `json:"id"`
	Code             string    `json:"code"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Type             string    `json:"type"`
	Value            float64   `json:"value"`
	MinOrderAmount   float64   `json:"min_order_amount"`
	MaxDiscountValue float64   `json:"max_discount_value"`
	UsageLimit       int       `json:"usage_limit"`
	UsedCount        int       `json:"used_count"`
	IsActive         bool      `json:"is_active"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type AddressResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	AddressLine1 string    `json:"address_line1"`
	AddressLine2 string    `json:"address_line2"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	PostalCode   string    `json:"postal_code"`
	Country      string    `json:"country"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type BrandResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	LogoURL     string    `json:"logo_url"`
	Website     string    `json:"website"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductImageResponse struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	ImageURL  string    `json:"image_url"`
	AltText   string    `json:"alt_text"`
	IsPrimary bool      `json:"is_primary"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse methods
func (r *Review) ToResponse() ReviewResponse {
	response := ReviewResponse{
		ID:        r.ID,
		ProductID: r.ProductID,
		UserID:    r.UserID,
		Rating:    r.Rating,
		Title:     r.Title,
		Comment:   r.Comment,
		IsActive:  r.IsActive,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
	if r.User != nil {
		userResponse := r.User.ToResponse()
		response.User = &userResponse
	}
	return response
}

func (c *Coupon) ToResponse() CouponResponse {
	return CouponResponse{
		ID:               c.ID,
		Code:             c.Code,
		Name:             c.Name,
		Description:      c.Description,
		Type:             c.Type,
		Value:            c.Value,
		MinOrderAmount:   c.MinOrderAmount,
		MaxDiscountValue: c.MaxDiscountValue,
		UsageLimit:       c.UsageLimit,
		UsedCount:        c.UsedCount,
		IsActive:         c.IsActive,
		StartDate:        c.StartDate,
		EndDate:          c.EndDate,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

func (a *Address) ToResponse() AddressResponse {
	return AddressResponse{
		ID:           a.ID,
		UserID:       a.UserID,
		Name:         a.Name,
		Phone:        a.Phone,
		AddressLine1: a.AddressLine1,
		AddressLine2: a.AddressLine2,
		City:         a.City,
		State:        a.State,
		PostalCode:   a.PostalCode,
		Country:      a.Country,
		IsDefault:    a.IsDefault,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

func (b *Brand) ToResponse() BrandResponse {
	return BrandResponse{
		ID:          b.ID,
		Name:        b.Name,
		Slug:        b.Slug,
		Description: b.Description,
		LogoURL:     b.LogoURL,
		Website:     b.Website,
		IsActive:    b.IsActive,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

func (pi *ProductImage) ToResponse() ProductImageResponse {
	return ProductImageResponse{
		ID:        pi.ID,
		ProductID: pi.ProductID,
		ImageURL:  pi.ImageURL,
		AltText:   pi.AltText,
		IsPrimary: pi.IsPrimary,
		SortOrder: pi.SortOrder,
		CreatedAt: pi.CreatedAt,
		UpdatedAt: pi.UpdatedAt,
	}
}