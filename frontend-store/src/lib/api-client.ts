/**
 * API 客户端基础类
 * 封装 HTTP 请求、错误处理、请求签名等
 */

interface ApiRequestConfig {
  method: "GET" | "POST" | "PUT" | "DELETE";
  path: string;
  params?: Record<string, any>;
  body?: any;
  headers?: Record<string, string>;
  requireAuth?: boolean;
  sign?: boolean; // 是否需要签名（防重放攻击）
}

interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8082";

class ApiClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
    if (typeof window !== "undefined") {
      this.token = localStorage.getItem("auth_token");
    }
  }

  setToken(token: string) {
    this.token = token;
    if (typeof window !== "undefined") {
      localStorage.setItem("auth_token", token);
    }
  }

  clearToken() {
    this.token = null;
    if (typeof window !== "undefined") {
      localStorage.removeItem("auth_token");
    }
  }

  getToken() {
    return this.token;
  }

  /**
   * 生成请求签名（防重放攻击）
   * 签名算法: MD5(timestamp + path + body + secret)
   */
  private generateSignature(path: string, body: any, timestamp: number): string {
    const SECRET = process.env.NEXT_PUBLIC_API_SECRET || "x-store-signature-secret-key-change-in-production";

    // 构建签名内容
    const bodyStr = body ? JSON.stringify(body) : "";
    const signContent = `${timestamp}${path}${bodyStr}${SECRET}`;

    // 简单的字符串 hash (生产环境应使用 crypto-js 或 Web Crypto API)
    let hash = 0;
    for (let i = 0; i < signContent.length; i++) {
      const char = signContent.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash; // Convert to 32bit integer
    }

    return Math.abs(hash).toString(16);
  }

  private async request<T>({
    method,
    path,
    params,
    body,
    headers = {},
    requireAuth = false,
    sign = false,
  }: ApiRequestConfig): Promise<ApiResponse<T>> {
    // 构建 URL
    let url = `${this.baseURL}${path}`;
    if (params && Object.keys(params).length > 0) {
      const searchParams = new URLSearchParams();
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          searchParams.append(key, String(value));
        }
      });
      url += `?${searchParams.toString()}`;
    }

    // 构建请求头
    const requestHeaders: HeadersInit = {
      "Content-Type": "application/json",
      ...headers,
    };

    // 添加认证 Token
    if (requireAuth && this.token) {
      requestHeaders["Authorization"] = `Bearer ${this.token}`;
    }

    // 添加签名（防重放攻击）
    if (sign) {
      const timestamp = Math.floor(Date.now() / 1000);
      const signature = this.generateSignature(path, body, timestamp);
      requestHeaders["X-Timestamp"] = timestamp.toString();
      requestHeaders["X-Signature"] = signature;
    }

    // 构建请求配置
    const config: RequestInit = {
      method,
      headers: requestHeaders,
    };

    if (body && method !== "GET") {
      config.body = JSON.stringify(body);
    }

    try {
      const response = await fetch(url, config);
      const data = await response.json();

      // 处理业务错误
      if (data.code !== 0) {
        throw new ApiError(data.message || "请求失败", data.code);
      }

      return data;
    } catch (error) {
      if (error instanceof ApiError) {
        throw error;
      }

      // 网络错误或其他异常
      throw new ApiError(
        error instanceof Error ? error.message : "网络请求失败",
        -1
      );
    }
  }

  async get<T>(path: string, params?: Record<string, any>, options?: Partial<ApiRequestConfig>): Promise<ApiResponse<T>> {
    return this.request<T>({ method: "GET", path, params, ...options });
  }

  async post<T>(path: string, body?: any, options?: Partial<ApiRequestConfig>): Promise<ApiResponse<T>> {
    return this.request<T>({ method: "POST", path, body, ...options });
  }

  async put<T>(path: string, body?: any, options?: Partial<ApiRequestConfig>): Promise<ApiResponse<T>> {
    return this.request<T>({ method: "PUT", path, body, ...options });
  }

  async delete<T>(path: string, options?: Partial<ApiRequestConfig>): Promise<ApiResponse<T>> {
    return this.request<T>({ method: "DELETE", path, ...options });
  }
}

class ApiError extends Error {
  constructor(
    message: string,
    public code: number
  ) {
    super(message);
    this.name = "ApiError";
  }
}

// 创建单例
export const apiClient = new ApiClient(API_BASE_URL);

export type { ApiResponse, ApiRequestConfig };
export { ApiError };
