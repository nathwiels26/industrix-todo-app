package services

import (
	"errors"

	"github.com/industrix-todo-app/backend/internal/models"
	"github.com/industrix-todo-app/backend/internal/repository"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryNameRequired = errors.New("category name is required")
)

type CategoryService interface {
	Create(req models.CreateCategoryRequest) (*models.Category, error)
	GetAll() ([]models.Category, error)
	GetByID(id uint) (*models.Category, error)
	Update(id uint, req models.UpdateCategoryRequest) (*models.Category, error)
	Delete(id uint) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(req models.CreateCategoryRequest) (*models.Category, error) {
	if req.Name == "" {
		return nil, ErrCategoryNameRequired
	}

	category := &models.Category{
		Name:  req.Name,
		Color: req.Color,
	}

	if category.Color == "" {
		category.Color = "#3B82F6"
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (s *categoryService) Update(id uint, req models.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Color != "" {
		category.Color = req.Color
	}

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Delete(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return ErrCategoryNotFound
	}

	return s.repo.Delete(id)
}
