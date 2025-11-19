package models

import (
	"time"
)

type Priority string

const (
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

type Todo struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Completed   bool       `gorm:"default:false" json:"completed"`
	Priority    Priority   `gorm:"size:20;default:'medium'" json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CategoryID  *uint      `json:"category_id,omitempty"`
	Category    *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTodoRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=255"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	CategoryID  *uint      `json:"category_id"`
}

type UpdateTodoRequest struct {
	Title       string     `json:"title" binding:"omitempty,min=1,max=255"`
	Description string     `json:"description"`
	Completed   *bool      `json:"completed"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	CategoryID  *uint      `json:"category_id"`
}

type TodoFilter struct {
	Search     string
	CategoryID *uint
	Completed  *bool
	Priority   Priority
	Page       int
	Limit      int
	SortBy     string
	SortOrder  string
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
}
