package tests

import (
	"testing"
	"time"

	"github.com/industrix-todo-app/backend/internal/models"
	"github.com/industrix-todo-app/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoryRepository is a mock implementation of CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category *models.Category) error {
	args := m.Called(category)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	category.ID = 1
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	return nil
}

func (m *MockCategoryRepository) GetAll() ([]models.Category, error) {
	args := m.Called()
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByID(id uint) (*models.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCategoryService_Create(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := services.NewCategoryService(mockRepo)

	t.Run("successful creation", func(t *testing.T) {
		req := models.CreateCategoryRequest{
			Name:  "Work",
			Color: "#3B82F6",
		}

		mockRepo.On("Create", mock.AnythingOfType("*models.Category")).Return(nil).Once()

		category, err := service.Create(req)

		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, req.Name, category.Name)
		assert.Equal(t, req.Color, category.Color)
		mockRepo.AssertExpectations(t)
	})

	t.Run("default color when not provided", func(t *testing.T) {
		req := models.CreateCategoryRequest{
			Name: "Personal",
		}

		mockRepo.On("Create", mock.AnythingOfType("*models.Category")).Return(nil).Once()

		category, err := service.Create(req)

		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, "#3B82F6", category.Color)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty name error", func(t *testing.T) {
		req := models.CreateCategoryRequest{
			Name: "",
		}

		category, err := service.Create(req)

		assert.Error(t, err)
		assert.Nil(t, category)
		assert.Equal(t, services.ErrCategoryNameRequired, err)
	})
}

func TestCategoryService_GetAll(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := services.NewCategoryService(mockRepo)

	t.Run("successful get all", func(t *testing.T) {
		expectedCategories := []models.Category{
			{ID: 1, Name: "Work"},
			{ID: 2, Name: "Personal"},
		}

		mockRepo.On("GetAll").Return(expectedCategories, nil).Once()

		categories, err := service.GetAll()

		assert.NoError(t, err)
		assert.Len(t, categories, 2)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryService_GetByID(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := services.NewCategoryService(mockRepo)

	t.Run("successful get by id", func(t *testing.T) {
		expectedCategory := &models.Category{
			ID:   1,
			Name: "Work",
		}

		mockRepo.On("GetByID", uint(1)).Return(expectedCategory, nil).Once()

		category, err := service.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, expectedCategory.Name, category.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found error", func(t *testing.T) {
		mockRepo.On("GetByID", uint(999)).Return(nil, services.ErrCategoryNotFound).Once()

		category, err := service.GetByID(999)

		assert.Error(t, err)
		assert.Nil(t, category)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryService_Update(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := services.NewCategoryService(mockRepo)

	t.Run("successful update", func(t *testing.T) {
		existingCategory := &models.Category{
			ID:    1,
			Name:  "Old Name",
			Color: "#000000",
		}

		req := models.UpdateCategoryRequest{
			Name:  "New Name",
			Color: "#FFFFFF",
		}

		mockRepo.On("GetByID", uint(1)).Return(existingCategory, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*models.Category")).Return(nil).Once()

		category, err := service.Update(1, req)

		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, req.Name, category.Name)
		assert.Equal(t, req.Color, category.Color)
		mockRepo.AssertExpectations(t)
	})

	t.Run("partial update - name only", func(t *testing.T) {
		existingCategory := &models.Category{
			ID:    2,
			Name:  "Old Name",
			Color: "#000000",
		}

		req := models.UpdateCategoryRequest{
			Name: "New Name",
		}

		mockRepo.On("GetByID", uint(2)).Return(existingCategory, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*models.Category")).Return(nil).Once()

		category, err := service.Update(2, req)

		assert.NoError(t, err)
		assert.Equal(t, "New Name", category.Name)
		assert.Equal(t, "#000000", category.Color)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found error", func(t *testing.T) {
		mockRepo.On("GetByID", uint(999)).Return(nil, services.ErrCategoryNotFound).Once()

		category, err := service.Update(999, models.UpdateCategoryRequest{})

		assert.Error(t, err)
		assert.Nil(t, category)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryService_Delete(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := services.NewCategoryService(mockRepo)

	t.Run("successful delete", func(t *testing.T) {
		existingCategory := &models.Category{ID: 1}

		mockRepo.On("GetByID", uint(1)).Return(existingCategory, nil).Once()
		mockRepo.On("Delete", uint(1)).Return(nil).Once()

		err := service.Delete(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found error", func(t *testing.T) {
		mockRepo.On("GetByID", uint(999)).Return(nil, services.ErrCategoryNotFound).Once()

		err := service.Delete(999)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
