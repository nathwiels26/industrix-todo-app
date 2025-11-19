package services

import (
	"errors"
	"math"

	"github.com/industrix-todo-app/backend/internal/models"
	"github.com/industrix-todo-app/backend/internal/repository"
)

var (
	ErrTodoNotFound     = errors.New("todo not found")
	ErrTodoTitleRequired = errors.New("todo title is required")
	ErrInvalidPriority  = errors.New("invalid priority value")
)

type TodoService interface {
	Create(req models.CreateTodoRequest) (*models.Todo, error)
	GetAll(filter models.TodoFilter) (*models.PaginatedResponse, error)
	GetByID(id uint) (*models.Todo, error)
	Update(id uint, req models.UpdateTodoRequest) (*models.Todo, error)
	Delete(id uint) error
	ToggleComplete(id uint) (*models.Todo, error)
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) Create(req models.CreateTodoRequest) (*models.Todo, error) {
	if req.Title == "" {
		return nil, ErrTodoTitleRequired
	}

	if req.Priority != "" && !isValidPriority(req.Priority) {
		return nil, ErrInvalidPriority
	}

	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		CategoryID:  req.CategoryID,
	}

	if todo.Priority == "" {
		todo.Priority = models.PriorityMedium
	}

	if err := s.repo.Create(todo); err != nil {
		return nil, err
	}

	// Reload to get category data
	return s.repo.GetByID(todo.ID)
}

func (s *todoService) GetAll(filter models.TodoFilter) (*models.PaginatedResponse, error) {
	// Set default pagination values
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	todos, total, err := s.repo.GetAll(filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &models.PaginatedResponse{
		Data: todos,
		Pagination: models.Pagination{
			CurrentPage: filter.Page,
			PerPage:     filter.Limit,
			Total:       total,
			TotalPages:  totalPages,
		},
	}, nil
}

func (s *todoService) GetByID(id uint) (*models.Todo, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrTodoNotFound
	}
	return todo, nil
}

func (s *todoService) Update(id uint, req models.UpdateTodoRequest) (*models.Todo, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrTodoNotFound
	}

	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.Description != "" {
		todo.Description = req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	if req.Priority != "" {
		if !isValidPriority(req.Priority) {
			return nil, ErrInvalidPriority
		}
		todo.Priority = req.Priority
	}
	if req.DueDate != nil {
		todo.DueDate = req.DueDate
	}
	if req.CategoryID != nil {
		todo.CategoryID = req.CategoryID
	}

	if err := s.repo.Update(todo); err != nil {
		return nil, err
	}

	// Reload to get updated category data
	return s.repo.GetByID(todo.ID)
}

func (s *todoService) Delete(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return ErrTodoNotFound
	}

	return s.repo.Delete(id)
}

func (s *todoService) ToggleComplete(id uint) (*models.Todo, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrTodoNotFound
	}

	todo.Completed = !todo.Completed

	if err := s.repo.Update(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func isValidPriority(p models.Priority) bool {
	return p == models.PriorityHigh || p == models.PriorityMedium || p == models.PriorityLow
}
