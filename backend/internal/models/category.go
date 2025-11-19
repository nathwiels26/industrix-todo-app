package models

import (
	"time"
)

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Color     string    `gorm:"size:20;default:'#3B82F6'" json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Todos     []Todo    `gorm:"foreignKey:CategoryID" json:"todos,omitempty"`
}

type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Color string `json:"color"`
}

type UpdateCategoryRequest struct {
	Name  string `json:"name" binding:"omitempty,min=1,max=100"`
	Color string `json:"color"`
}
