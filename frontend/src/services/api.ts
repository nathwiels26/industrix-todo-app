import axios from 'axios';
import {
  Todo,
  Category,
  CreateTodoRequest,
  UpdateTodoRequest,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  PaginatedResponse,
  TodoFilter,
} from '../types';

const API_BASE_URL = '/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Todo API
export const todoApi = {
  getAll: async (filter: TodoFilter = {}): Promise<PaginatedResponse<Todo>> => {
    const params = new URLSearchParams();

    if (filter.search) params.append('search', filter.search);
    if (filter.category_id) params.append('category_id', filter.category_id.toString());
    if (filter.completed !== undefined) params.append('completed', filter.completed.toString());
    if (filter.priority) params.append('priority', filter.priority);
    if (filter.page) params.append('page', filter.page.toString());
    if (filter.limit) params.append('limit', filter.limit.toString());
    if (filter.sort_by) params.append('sort_by', filter.sort_by);
    if (filter.sort_order) params.append('sort_order', filter.sort_order);

    const response = await api.get<PaginatedResponse<Todo>>(`/todos?${params.toString()}`);
    return response.data;
  },

  getById: async (id: number): Promise<Todo> => {
    const response = await api.get<Todo>(`/todos/${id}`);
    return response.data;
  },

  create: async (data: CreateTodoRequest): Promise<Todo> => {
    const response = await api.post<Todo>('/todos', data);
    return response.data;
  },

  update: async (id: number, data: UpdateTodoRequest): Promise<Todo> => {
    const response = await api.put<Todo>(`/todos/${id}`, data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/todos/${id}`);
  },

  toggleComplete: async (id: number): Promise<Todo> => {
    const response = await api.patch<Todo>(`/todos/${id}/complete`);
    return response.data;
  },
};

// Category API
export const categoryApi = {
  getAll: async (): Promise<Category[]> => {
    const response = await api.get<Category[]>('/categories');
    return response.data;
  },

  getById: async (id: number): Promise<Category> => {
    const response = await api.get<Category>(`/categories/${id}`);
    return response.data;
  },

  create: async (data: CreateCategoryRequest): Promise<Category> => {
    const response = await api.post<Category>('/categories', data);
    return response.data;
  },

  update: async (id: number, data: UpdateCategoryRequest): Promise<Category> => {
    const response = await api.put<Category>(`/categories/${id}`, data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/categories/${id}`);
  },
};

export default api;
