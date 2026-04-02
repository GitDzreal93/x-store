/**
 * 用户认证上下文
 * 管理全局认证状态和用户信息
 */

"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import {
  login,
  register,
  logout,
  getUserProfile,
  isAuthenticated,
  getToken,
  type User,
  type LoginParams,
  type RegisterParams,
} from "@/lib/api";

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (params: LoginParams) => Promise<void>;
  register: (params: RegisterParams) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // 初始化：检查是否已登录
  useEffect(() => {
    const checkAuth = async () => {
      if (isAuthenticated()) {
        try {
          const profile = await getUserProfile();
          setUser(profile);
        } catch (error) {
          console.error("获取用户信息失败:", error);
          // Token 可能已过期，清除本地状态
          localStorage.removeItem("auth_token");
        }
      }
      setIsLoading(false);
    };

    checkAuth();
  }, []);

  const handleLogin = async (params: LoginParams) => {
    const result = await login(params);
    setUser(result.user);
  };

  const handleRegister = async (params: RegisterParams) => {
    const result = await register(params);
    setUser(result.user);
  };

  const handleLogout = async () => {
    await logout();
    setUser(null);
  };

  const refreshUser = async () => {
    if (isAuthenticated()) {
      const profile = await getUserProfile();
      setUser(profile);
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login: handleLogin,
        register: handleRegister,
        logout: handleLogout,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
