import request from '../utils/request';
import type { Order, PageResult } from '../types';

export const getOrders = (params: {
  page?: number;
  size?: number;
  status?: number;
  keyword?: string;
}) =>
  request.get<any, { code: number; data: PageResult<Order> }>('/admin/orders', { params });

export const getOrder = (orderNo: string) =>
  request.get<any, { code: number; data: Order }>(`/orders/${orderNo}`);

export const getCardKeyCount = (productId: number) =>
  request.get<any, { code: number; data: { product_id: number; available: number } }>(`/admin/cardkeys/count/${productId}`);

export const importCardKeys = (data: { product_id: number; content: string }) =>
  request.post<any, { code: number; data: { imported: number; message: string } }>('/admin/cardkeys/import', data);
