import request from '../utils/request';
import type { Product, CreateProductReq, UpdateProductReq, PageResult } from '../types';

export const getProducts = (params: {
  page?: number;
  size?: number;
  category_id?: number;
  keyword?: string;
  all?: string;
}) =>
  request.get<any, { code: number; data: PageResult<Product> }>('/admin/products', { params: { all: 'true', ...params } });

export const getProduct = (id: number) =>
  request.get<any, { code: number; data: Product }>(`/products/${id}`);

export const createProduct = (data: CreateProductReq) =>
  request.post<any, { code: number; data: Product }>('/admin/products', data);

export const updateProduct = (id: number, data: UpdateProductReq) =>
  request.put<any, { code: number; data: Product }>(`/admin/products/${id}`, data);

export const deleteProduct = (id: number) =>
  request.delete(`/admin/products/${id}`);
