import request from '../utils/request';
import type { Category } from '../types';

export const getCategories = (all = true) =>
  request.get<any, { code: number; data: Category[] }>('/categories', { params: all ? { all: 'true' } : {} });

export const getCategory = (id: number) =>
  request.get<any, { code: number; data: Category }>(`/categories/${id}`);

export const createCategory = (data: Partial<Category>) =>
  request.post<any, { code: number; data: Category }>('/admin/categories', data);

export const updateCategory = (id: number, data: Partial<Category>) =>
  request.put<any, { code: number; data: Category }>(`/admin/categories/${id}`, data);

export const deleteCategory = (id: number) =>
  request.delete(`/admin/categories/${id}`);
