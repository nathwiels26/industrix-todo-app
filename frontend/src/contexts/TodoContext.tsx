import React, { createContext, useContext, useState, useCallback, ReactNode } from 'react';
import { message } from 'antd';
import {
  Todo,
  Category,
  CreateTodoRequest,
  UpdateTodoRequest,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  TodoFilter,
  Pagination,
} from '../types';
import { todoApi, categoryApi } from '../services/api';

interface TodoContextType {
  // State
  todos: Todo[];
  categories: Category[];
  pagination: Pagination;
  loading: boolean;
  filter: TodoFilter;

  // Todo actions
  fetchTodos: (filter?: TodoFilter) => Promise<void>;
  createTodo: (data: CreateTodoRequest) => Promise<void>;
  updateTodo: (id: number, data: UpdateTodoRequest) => Promise<void>;
  deleteTodo: (id: number) => Promise<void>;
  toggleComplete: (id: number) => Promise<void>;

  // Category actions
  fetchCategories: () => Promise<void>;
  createCategory: (data: CreateCategoryRequest) => Promise<void>;
  updateCategory: (id: number, data: UpdateCategoryRequest) => Promise<void>;
  deleteCategory: (id: number) => Promise<void>;

  // Filter actions
  setFilter: (filter: TodoFilter) => void;
}

const TodoContext = createContext<TodoContextType | undefined>(undefined);

interface TodoProviderProps {
  children: ReactNode;
}

export const TodoProvider: React.FC<TodoProviderProps> = ({ children }) => {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [pagination, setPagination] = useState<Pagination>({
    current_page: 1,
    per_page: 10,
    total: 0,
    total_pages: 0,
  });
  const [loading, setLoading] = useState(false);
  const [filter, setFilterState] = useState<TodoFilter>({
    page: 1,
    limit: 10,
  });

  // Fetch todos with filter
  const fetchTodos = useCallback(async (newFilter?: TodoFilter) => {
    setLoading(true);
    try {
      const appliedFilter = newFilter || filter;
      const response = await todoApi.getAll(appliedFilter);
      setTodos(response.data);
      setPagination(response.pagination);
    } catch (error) {
      message.error('Failed to fetch todos');
      console.error('Error fetching todos:', error);
    } finally {
      setLoading(false);
    }
  }, [filter]);

  // Create todo
  const createTodo = useCallback(async (data: CreateTodoRequest) => {
    try {
      await todoApi.create(data);
      message.success('Todo created successfully');
      await fetchTodos();
    } catch (error) {
      message.error('Failed to create todo');
      console.error('Error creating todo:', error);
      throw error;
    }
  }, [fetchTodos]);

  // Update todo
  const updateTodo = useCallback(async (id: number, data: UpdateTodoRequest) => {
    try {
      await todoApi.update(id, data);
      message.success('Todo updated successfully');
      await fetchTodos();
    } catch (error) {
      message.error('Failed to update todo');
      console.error('Error updating todo:', error);
      throw error;
    }
  }, [fetchTodos]);

  // Delete todo
  const deleteTodo = useCallback(async (id: number) => {
    try {
      await todoApi.delete(id);
      message.success('Todo deleted successfully');
      await fetchTodos();
    } catch (error) {
      message.error('Failed to delete todo');
      console.error('Error deleting todo:', error);
      throw error;
    }
  }, [fetchTodos]);

  // Toggle complete
  const toggleComplete = useCallback(async (id: number) => {
    try {
      await todoApi.toggleComplete(id);
      await fetchTodos();
    } catch (error) {
      message.error('Failed to toggle todo status');
      console.error('Error toggling todo:', error);
      throw error;
    }
  }, [fetchTodos]);

  // Fetch categories
  const fetchCategories = useCallback(async () => {
    try {
      const data = await categoryApi.getAll();
      setCategories(data);
    } catch (error) {
      message.error('Failed to fetch categories');
      console.error('Error fetching categories:', error);
    }
  }, []);

  // Create category
  const createCategory = useCallback(async (data: CreateCategoryRequest) => {
    try {
      await categoryApi.create(data);
      message.success('Category created successfully');
      await fetchCategories();
    } catch (error) {
      message.error('Failed to create category');
      console.error('Error creating category:', error);
      throw error;
    }
  }, [fetchCategories]);

  // Update category
  const updateCategory = useCallback(async (id: number, data: UpdateCategoryRequest) => {
    try {
      await categoryApi.update(id, data);
      message.success('Category updated successfully');
      await fetchCategories();
    } catch (error) {
      message.error('Failed to update category');
      console.error('Error updating category:', error);
      throw error;
    }
  }, [fetchCategories]);

  // Delete category
  const deleteCategory = useCallback(async (id: number) => {
    try {
      await categoryApi.delete(id);
      message.success('Category deleted successfully');
      await fetchCategories();
    } catch (error) {
      message.error('Failed to delete category');
      console.error('Error deleting category:', error);
      throw error;
    }
  }, [fetchCategories]);

  // Set filter
  const setFilter = useCallback((newFilter: TodoFilter) => {
    setFilterState(prev => ({ ...prev, ...newFilter }));
  }, []);

  const value: TodoContextType = {
    todos,
    categories,
    pagination,
    loading,
    filter,
    fetchTodos,
    createTodo,
    updateTodo,
    deleteTodo,
    toggleComplete,
    fetchCategories,
    createCategory,
    updateCategory,
    deleteCategory,
    setFilter,
  };

  return <TodoContext.Provider value={value}>{children}</TodoContext.Provider>;
};

export const useTodo = (): TodoContextType => {
  const context = useContext(TodoContext);
  if (context === undefined) {
    throw new Error('useTodo must be used within a TodoProvider');
  }
  return context;
};

export default TodoContext;
