/**
 * 商品相关 API
 */

import { apiClient, ApiResponse } from "../api-client";

// 商品类型定义
export interface Product {
  id: number;
  category_id: number;
  title: string;
  cover: string;
  price: number;
  original_price?: number;
  stock: number;
  sales: number;
  delivery_type: "auto" | "manual";
  tags: string[];
  sort: number;
  status: number;
  is_new: boolean;
  is_hot: boolean;
  created_at: string;
  updated_at: string;
}

export interface ProductDetail extends Product {
  description: string;
  notice: string;
}

export interface Category {
  id: number;
  name: string;
  icon: string;
  sort: number;
  status: number;
  created_at: string;
  updated_at: string;
}

export interface ProductsListParams {
  page?: number;
  size?: number;
  category_id?: number;
  keyword?: string;
  sort?: "price_asc" | "price_desc" | "sales_desc" | "newest";
}

export interface ProductsListResponse {
  list: Product[];
  total: number;
  page: number;
  size: number;
}

/**
 * 获取商品列表
 */
export async function getProducts(params?: ProductsListParams): Promise<ProductsListResponse> {
  const response = await apiClient.get<ProductsListResponse>("/api/products", params);
  return response.data;
}

/**
 * 获取商品详情
 */
export async function getProductDetail(id: number): Promise<ProductDetail> {
  const response = await apiClient.get<ProductDetail>(`/api/products/${id}`);
  return response.data;
}

/**
 * 获取分类列表
 */
export async function getCategories(): Promise<Category[]> {
  const response = await apiClient.get<Category[]>("/api/categories");
  return response.data;
}

/**
 * 获取单个分类
 */
export async function getCategory(id: number): Promise<Category> {
  const response = await apiClient.get<Category>(`/api/categories/${id}`);
  return response.data;
}

/**
 * 将后端商品数据转换为前端组件使用的格式
 */
export function transformProduct(product: Product): {
  id: string;
  title: string;
  tags: string[];
  price: number;
  sales: number;
  likes: number;
  updatedAt: string;
  isNew: boolean;
  isHot: boolean;
} {
  // 解析 tags JSON 字符串
  let tags: string[] = [];
  if (typeof product.tags === "string") {
    try {
      tags = JSON.parse(product.tags);
    } catch {
      tags = ["自动发货"];
    }
  } else if (Array.isArray(product.tags)) {
    tags = product.tags;
  }

  // 格式化日期
  const updatedAt = new Date(product.updated_at).toLocaleDateString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).replace(/\//g, ".");

  return {
    id: String(product.id),
    title: product.title,
    tags,
    price: product.price,
    sales: product.sales,
    likes: Math.floor(product.sales * 0.02), // 模拟点赞数
    updatedAt,
    isNew: product.is_new,
    isHot: product.is_hot,
  };
}
