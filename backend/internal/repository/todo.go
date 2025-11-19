package repository

import (
	"github.com/industrix-todo-app/backend/internal/models"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(todo *models.Todo) error
	GetAll(filter models.TodoFilter) ([]models.Todo, int64, error)
	GetByID(id uint) (*models.Todo, error)
	Update(todo *models.Todo) error
	Delete(id uint) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) GetAll(filter models.TodoFilter) ([]models.Todo, int64, error) {
	var todos []models.Todo
	var total int64

	query := r.db.Model(&models.Todo{})

	// Apply search filter
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Apply category filter
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	// Apply completed filter
	if filter.Completed != nil {
		query = query.Where("completed = ?", *filter.Completed)
	}

	// Apply priority filter
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}

	// Count total records
	query.Count(&total)

	// Apply sorting
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := filter.SortOrder
	if sortOrder == "" {
		sortOrder = "DESC"
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Apply pagination
	if filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	// Execute query with preload
	err := query.Preload("Category").Find(&todos).Error
	return todos, total, err
}

func (r *todoRepository) GetByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.Preload("Category").First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) Update(todo *models.Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Todo{}, id).Error
}
