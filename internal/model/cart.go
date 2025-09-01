package model

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User      *User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CartItems []CartItem `json:"cart_items,omitempty" gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type CartItem struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CartID    uint           `json:"cart_id" gorm:"not null;index"`
	ProductID uint           `json:"product_id" gorm:"not null;index"`
	Quantity  int            `json:"quantity" gorm:"not null;default:1"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Cart    *Cart    `json:"cart,omitempty" gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName specifies the table name for Cart model
func (Cart) TableName() string {
	return "carts"
}

// TableName specifies the table name for CartItem model
func (CartItem) TableName() string {
	return "cart_items"
}

type CartItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type CartResponse struct {
	ID         uint                `json:"id"`
	UserID     uint                `json:"user_id"`
	CartItems  []CartItemResponse  `json:"cart_items"`
	TotalItems int                 `json:"total_items"`
	TotalPrice float64             `json:"total_price"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

type CartItemResponse struct {
	ID        uint             `json:"id"`
	CartID    uint             `json:"cart_id"`
	ProductID uint             `json:"product_id"`
	Product   *ProductResponse `json:"product,omitempty"`
	Quantity  int              `json:"quantity"`
	SubTotal  float64          `json:"sub_total"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// ToResponse converts Cart to CartResponse
func (c *Cart) ToResponse() CartResponse {
	response := CartResponse{
		ID:        c.ID,
		UserID:    c.UserID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	totalItems := 0
	totalPrice := 0.0

	if len(c.CartItems) > 0 {
		for _, item := range c.CartItems {
			itemResponse := CartItemResponse{
				ID:        item.ID,
				CartID:    item.CartID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}

			if item.Product != nil {
				productResponse := item.Product.ToResponse()
				itemResponse.Product = &productResponse
				itemResponse.SubTotal = float64(item.Quantity) * item.Product.Price
				totalPrice += itemResponse.SubTotal
			}

			totalItems += item.Quantity
			response.CartItems = append(response.CartItems, itemResponse)
		}
	}

	response.TotalItems = totalItems
	response.TotalPrice = totalPrice

	return response
}