/**
 * 用户认证相关 API
 */

import { apiClient, ApiResponse } from "../api-client";

// 用户类型定义
export interface LoginParams {
  username: string;
  password: string;
}

export interface RegisterParams {
  username: string;
  email: string;
  password: string;
  nickname?: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface User {
  id: number;
  username: string;
  email: string;
  nickname?: string;
  role: "user" | "admin";
  status: number;
  created_at: string;
  updated_at: string;
}

export interface UserProfile extends User {
  phone?: string;
  avatar?: string;
  bio?: string;
}

/**
 * 管理员登录
 */
export async function adminLogin(params: LoginParams): Promise<LoginResponse> {
  const response = await apiClient.post<LoginResponse>("/api/admin/login", params);
  // 保存 token
  if (response.data.token) {
    apiClient.setToken(response.data.token);
  }
  return response.data;
}

/**
 * 用户注册
 */
export async function register(params: RegisterParams): Promise<LoginResponse> {
  const response = await apiClient.post<LoginResponse>("/api/auth/register", params, {
    sign: true,
  });
  // 保存 token
  if (response.data.token) {
    apiClient.setToken(response.data.token);
  }
  return response.data;
}

/**
 * 用户登录
 */
export async function login(params: LoginParams): Promise<LoginResponse> {
  const response = await apiClient.post<LoginResponse>("/api/auth/login", params, {
    sign: true,
  });
  // 保存 token
  if (response.data.token) {
    apiClient.setToken(response.data.token);
  }
  return response.data;
}

/**
 * 退出登录
 */
export async function logout(): Promise<void> {
  try {
    await apiClient.post("/api/auth/logout", {}, { requireAuth: true });
  } finally {
    // 无论后端是否成功，都清除本地 token
    apiClient.clearToken();
  }
}

/**
 * 获取用户信息（需要认证）
 */
export async function getUserProfile(): Promise<UserProfile> {
  const response = await apiClient.get<UserProfile>(
    "/api/user/profile",
    {},
    { requireAuth: true }
  );
  return response.data;
}

/**
 * 获取管理员信息（需要认证）
 */
export async function getAdminProfile(): Promise<UserProfile> {
  const response = await apiClient.get<UserProfile>(
    "/api/admin/profile",
    {},
    { requireAuth: true }
  );
  return response.data;
}

/**
 * 更新用户信息（需要认证）
 */
export async function updateUserProfile(params: {
  nickname?: string;
  phone?: string;
  bio?: string;
}): Promise<UserProfile> {
  const response = await apiClient.put<UserProfile>(
    "/api/user/profile",
    params,
    { requireAuth: true }
  );
  return response.data;
}

/**
 * 修改密码（需要认证）
 */
export async function changePassword(params: {
  old_password: string;
  new_password: string;
}): Promise<void> {
  await apiClient.post(
    "/api/user/change-password",
    params,
    { requireAuth: true, sign: true }
  );
}

/**
 * 检查登录状态
 */
export function isAuthenticated(): boolean {
  return !!apiClient.getToken();
}

/**
 * 获取当前 Token
 */
export function getToken(): string | null {
  return apiClient.getToken();
}
