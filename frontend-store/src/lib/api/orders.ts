/**
 * 订单相关 API
 */

import { apiClient, ApiResponse } from "../api-client";

// 订单类型定义
export enum OrderStatus {
  Pending = "pending",
  Paid = "paid",
  Delivered = "delivered",
  Cancelled = "cancelled",
  Refunded = "refunded",
}

export interface CreateOrderParams {
  product_id: number;
  email: string;
  quantity?: number;
}

export interface Order {
  id: number;
  order_no: string;
  user_id?: number;
  product_id: number;
  product_title: string;
  email: string;
  quantity: number;
  total_amount: number;
  status: OrderStatus;
  delivery_type: "auto" | "manual";
  created_at: string;
  updated_at: string;
  paid_at?: string;
  delivered_at?: string;
}

export interface OrderDetail extends Order {
  card_key?: string;
  product?: {
    id: number;
    title: string;
    cover: string;
    price: number;
  };
}

export interface CreatePaymentParams {
  order_no: string;
  channel_id: number;
}

export interface PaymentResult {
  payment_id: number;
  payment_url?: string;
  qr_code?: string;
}

/**
 * 创建订单（需要签名）
 */
export async function createOrder(params: CreateOrderParams): Promise<Order> {
  const response = await apiClient.post<Order>(
    "/api/orders",
    params,
    { sign: true } // 启用防重放签名
  );
  return response.data;
}

/**
 * 取消订单
 */
export async function cancelOrder(orderNo: string): Promise<void> {
  await apiClient.post(`/api/orders/${orderNo}/cancel`, {}, { sign: true });
}

/**
 * 创建支付（需要签名）
 */
export async function createPayment(params: CreatePaymentParams): Promise<PaymentResult> {
  const response = await apiClient.post<PaymentResult>(
    "/api/payments",
    params,
    { sign: true }
  );
  return response.data;
}

/**
 * 查询支付状态
 */
export async function getPaymentStatus(paymentId: number): Promise<{ status: string }> {
  const response = await apiClient.get(`/api/payments/${paymentId}/status`);
  return response.data;
}

/**
 * 获取订单详情（需要认证）
 */
export async function getOrderDetail(orderNo: string): Promise<OrderDetail> {
  const response = await apiClient.get<OrderDetail>(
    `/api/orders/${orderNo}`,
    {},
    { requireAuth: true }
  );
  return response.data;
}

/**
 * 获取用户订单列表（需要认证）
 */
export async function getUserOrders(params?: {
  page?: number;
  size?: number;
  status?: OrderStatus;
}): Promise<{ list: Order[]; total: number }> {
  const response = await apiClient.get<{ list: Order[]; total: number }>(
    "/api/user/orders",
    params,
    { requireAuth: true }
  );
  return response.data;
}

/**
 * 轮询支付状态
 * @param paymentId 支付ID
 * @param interval 轮询间隔（毫秒）
 * @param maxAttempts 最大尝试次数
 */
export async function pollPaymentStatus(
  paymentId: number,
  interval: number = 2000,
  maxAttempts: number = 30
): Promise<{ status: string }> {
  let attempts = 0;

  return new Promise((resolve, reject) => {
    const timer = setInterval(async () => {
      attempts++;

      try {
        const result = await getPaymentStatus(paymentId);

        // 支付成功或失败，停止轮询
        if (result.status === "success" || result.status === "failed" || result.status === "cancelled") {
          clearInterval(timer);
          resolve(result);
          return;
        }

        // 超过最大尝试次数
        if (attempts >= maxAttempts) {
          clearInterval(timer);
          reject(new Error("支付查询超时"));
          return;
        }
      } catch (error) {
        clearInterval(timer);
        reject(error);
      }
    }, interval);
  });
}
