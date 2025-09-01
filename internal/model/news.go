package model

import (
	"time"

	"gorm.io/gorm"
)

type News struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string         `json:"title" gorm:"not null;size:200;index"`
	Slug        string         `json:"slug" gorm:"unique;not null;size:200;index"`
	Summary     string         `json:"summary" gorm:"size:500"`
	Content     string         `json:"content" gorm:"type:longtext;not null"`
	ImageURL    string         `json:"image_url" gorm:"size:500"`
	AuthorID    uint           `json:"author_id" gorm:"not null;index"`
	IsPublished bool           `json:"is_published" gorm:"default:false;index"`
	PublishedAt *time.Time     `json:"published_at"`
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	Tags        string         `json:"tags" gorm:"size:500"`
	MetaTitle   string         `json:"meta_title" gorm:"size:200"`
	MetaDesc    string         `json:"meta_description" gorm:"size:300"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Quan hệ
	Author *User `json:"author,omitempty" gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName chỉ định tên bảng cho model News
func (News) TableName() string {
	return "news"
}

type NewsInput struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Slug        string `json:"slug" binding:"required,min=1,max=200"`
	Summary     string `json:"summary" binding:"max=500"`
	Content     string `json:"content" binding:"required,min=1"`
	ImageURL    string `json:"image_url" binding:"max=500"`
	IsPublished bool   `json:"is_published"`
	Tags        string `json:"tags" binding:"max=500"`
	MetaTitle   string `json:"meta_title" binding:"max=200"`
	MetaDesc    string `json:"meta_description" binding:"max=300"`
}

type NewsResponse struct {
	ID          uint          `json:"id"`
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	Summary     string        `json:"summary"`
	Content     string        `json:"content"`
	ImageURL    string        `json:"image_url"`
	AuthorID    uint          `json:"author_id"`
	Author      *UserResponse `json:"author,omitempty"`
	IsPublished bool          `json:"is_published"`
	PublishedAt *time.Time    `json:"published_at"`
	ViewCount   int           `json:"view_count"`
	Tags        string        `json:"tags"`
	MetaTitle   string        `json:"meta_title"`
	MetaDesc    string        `json:"meta_description"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// ToResponse chuyển News thành NewsResponse
func (n *News) ToResponse() NewsResponse {
	response := NewsResponse{
		ID:          n.ID,
		Title:       n.Title,
		Slug:        n.Slug,
		Summary:     n.Summary,
		Content:     n.Content,
		ImageURL:    n.ImageURL,
		AuthorID:    n.AuthorID,
		IsPublished: n.IsPublished,
		PublishedAt: n.PublishedAt,
		ViewCount:   n.ViewCount,
		Tags:        n.Tags,
		MetaTitle:   n.MetaTitle,
		MetaDesc:    n.MetaDesc,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
	}

	// Bao gồm thông tin tác giả nếu đã được nạp
	if n.Author != nil {
		authorResponse := n.Author.ToResponse()
		response.Author = &authorResponse
	}

	return response
}