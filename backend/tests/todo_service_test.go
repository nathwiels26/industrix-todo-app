package tests

import (
	"testing"
	"time"

	"github.com/industrix-todo-app/backend/internal/models"
	"github.com/industrix-todo-app/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTodoRepository is a mock implementation of TodoRepository
type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) Create(todo *models.Todo) error {
	args := m.Called(todo)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	todo.ID = 1
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	return nil
}

func (m *MockTodoRepository) GetAll(filter models.TodoFilter) ([]models.Todo, int64, error) {
	args := m.Called(filter)
	return args.Get(0).([]models.Todo), args.Get(1).(int64), args.Error(2)
}

func (m *MockTodoRepository) GetByID(id uint) (*models.Todo, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Todo), args.Error(1)
}

func (m *MockTodoRepository) Update(todo *models.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestTodoService_Create(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := services.NewTodoService(mockRepo)

	t.Run("successful creation", func(t *testing.T) {
		req := models.CreateTodoRequest{
			Title:       "Test Todo",
			Description: "Test Description",
			Priority:    models.PriorityHigh,
		}

		mockRepo.On("Create", mock.AnythingOfType("*models.Todo")).Return(nil).Once()
		mockRepo.On("GetByID", uint(1)).Return(&models.Todo{
			ID:          1,
			Title:       req.Title,
			Description: req.Description,
			Priority:    req.Priority,
		}, nil).Once()

		todo, err := service.Create(req)

		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.Equal(t, req.Title, todo.Title)
		assert.Equal(t, req.Description, todo.Description)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty title error", func(t *testing.T) {
		req := models.CreateTodoRequest{
			Title: "",
		}

		todo, err := service.Create(req)

		assert.Error(t, err)
		assert.Nil(t, todo)
		assert.Equal(t, services.ErrTodoTitleRequired, err)
	})

	t.Run("invalid priority error", func(t *testing.T) {
		req := models.CreateTodoRequest{
			Title:    "Test",
			Priority: "invalid",
		}

		todo, err := service.Create(req)

		assert.Error(t, err)
		assert.Nil(t, todo)
		assert.Equal(t, services.ErrInvalidPriority, err)
	})
}

func TestTodoService_GetAll(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := services.NewTodoService(mockRepo)

	t.Run("successful get all with pagination", func(t *testing.T) {
		filter := models.TodoFilter{
			Page:  1,
			Limit: 10,
		}

		expectedTodos := []models.Todo{
			{ID: 1, Title: "Todo 1"},
			{ID: 2, Title: "Todo 2"},
		}

		mockRepo.On("GetAll", mock.AnythingOfType("models.TodoFilter")).Return(expectedTodos, int64(2), nil).Once()

		response, err := service.GetAll(filter)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, response.Pagination.CurrentPage)
		assert.Equal(t, int64(2), response.Pagination.Total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("default pagination values", func(t *testing.T) {
		filter := models.TodoFilter{
			Page:  0,
			Limit: 0,
		}

		mockRepo.On("GetAll", mock.AnythingOfType("models.TodoFilter")).Return([]models.Todo{}, int64(0), nil).Once()

		response, err := service.GetAll(filter)

		assert.NoError(t, err)
		assert.Equal(t, 1, response.Pagination.CurrentPage)
		assert.Equal(t, 10, response.Pagination.PerPage)
		mockRepo.AssertExpectations(t)
	})
}

func TestTodoService_GetByID(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := services.NewTodoService(mockRepo)

	t.Run("successful get by id", func(t *testing.T) {
		expectedTodo := &models.Todo{
			ID:    1,
			Title: "Test Todo",
		}

		mockRepo.On("GetByID", uint(1)).Return(expectedTodo, nil).Once()

		todo, err := service.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.Equal(t, expectedTodo.Title, todo.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found error", func(t *testing.T) {
		mockRepo.On("GetByID", uint(999)).Return(nil, services.ErrTodoNotFound).Once()

		todo, err := service.GetByID(999)

		assert.Error(t, err)
		assert.Nil(t, todo)
		mockRepo.AssertExpectations(t)
	})
}

func TestTodoService_Update(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := services.NewTodoService(mockRepo)

	t.Run("successful update", func(t *testing.T) {
		existingTodo := &models.Todo{
			ID:    1,
			Title: "Old Title",
		}

		req := models.UpdateTodoRequest{
			Title: "New Title",
		}

		mockRepo.On("GetByID", uint(1)).Return(existingTodo, nil).Times(2)
		mockRepo.On("Update", mock.AnythingOfType("*models.Todo")).Return(nil).Once()

		todo, err := service.Update(1, req)

		assert.NoError(t, err)
		assert.NotNil(t, todo)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found error", func(t *testing.T) {
		mockRepo.On("GetByID", uint(999)).Return(nil, services.ErrTodoNotFound).Once()

		todo, err := service.Update(999, models.UpdateTodoRequest{})

		assert.Error(t, err)
		assert.Nil(t, todo)
		mockRepo.AssertExpectations(t)
	})
}

func TestTodoService_Delete(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := services.NewTodoService(mockRepo)

	t.Run("successful delete", func(t *testing.T) {
		existingTodo := &models.Todo{ID: 1}

		mockRepo.On("GetByID", uint(1)).Return(existingTodo, nil).Once()
		mockRepo.On("Delete", uint(1)).Return(nil).Once()

		err := service.Delete(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found error", func(t *testing.T) {
		mockRepo.On("GetByID", uint(999)).Return(nil, services.ErrTodoNotFound).Once()

		err := service.Delete(999)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestTodoService_ToggleComplete(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := services.NewTodoService(mockRepo)

	t.Run("toggle from incomplete to complete", func(t *testing.T) {
		existingTodo := &models.Todo{
			ID:        1,
			Completed: false,
		}

		mockRepo.On("GetByID", uint(1)).Return(existingTodo, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*models.Todo")).Return(nil).Once()

		todo, err := service.ToggleComplete(1)

		assert.NoError(t, err)
		assert.True(t, todo.Completed)
		mockRepo.AssertExpectations(t)
	})

	t.Run("toggle from complete to incomplete", func(t *testing.T) {
		existingTodo := &models.Todo{
			ID:        2,
			Completed: true,
		}

		mockRepo.On("GetByID", uint(2)).Return(existingTodo, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*models.Todo")).Return(nil).Once()

		todo, err := service.ToggleComplete(2)

		assert.NoError(t, err)
		assert.False(t, todo.Completed)
		mockRepo.AssertExpectations(t)
	})
}
