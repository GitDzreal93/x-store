---
sidebar_position: 2
title: 管理后台
---

# B 端管理后台

管理后台基于 **React 18 + Ant Design 5 + Vite + TypeScript** 构建，面向商户管理员。

## 功能模块

| 页面 | 路径 | 功能 |
|------|------|------|
| 仪表盘 | `/dashboard` | GMV 统计、订单趋势、TOP 商品 |
| 商品管理 | `/products` | 商品 CRUD、上下架、排序 |
| 分类管理 | `/categories` | 分类 CRUD、图标、排序 |
| 订单管理 | `/orders` | 订单列表、退款、手动发货 |
| 卡密管理 | `/card-keys` | 批量导入、库存查看 |
| 支付渠道 | `/payment-channels` | 10 种支付方式管理 |
| 第三方登录 | `/oauth-providers` | GitHub/Google OAuth 开关配置 |

## 技术架构

### 路由

使用 React Router v6，路由守卫检查登录状态：

```tsx title="admin-panel/src/router/index.tsx"
const RequireAuth: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const token = localStorage.getItem('token');
  if (!token) return <Navigate to="/login" replace />;
  return <>{children}</>;
};
```

### 状态管理

使用 Zustand 管理全局认证状态：

```tsx title="admin-panel/src/stores/auth.ts"
export const useAuthStore = create((set) => ({
  token: localStorage.getItem('token'),
  user: null,
  login: (token, user) => {
    localStorage.setItem('token', token);
    set({ token, user });
  },
  logout: () => {
    localStorage.removeItem('token');
    set({ token: null, user: null });
  },
}));
```

### HTTP 请求

Axios 封装，自动携带 Token 和处理错误：

```tsx title="admin-panel/src/utils/request.ts"
const request = axios.create({
  baseURL: 'http://localhost:8082/api',
  timeout: 10000,
});

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

request.interceptors.response.use(
  (res) => res.data,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(err);
  }
);
```

### 布局组件

`AdminLayout` 提供侧边栏导航 + 顶部工具栏：

```tsx title="admin-panel/src/layouts/AdminLayout.tsx"
const menuItems = [
  { key: '/dashboard', icon: <DashboardOutlined />, label: '仪表盘' },
  { key: '/products', icon: <ShoppingOutlined />, label: '商品管理' },
  { key: '/categories', icon: <AppstoreOutlined />, label: '分类管理' },
  { key: '/orders', icon: <OrderedListOutlined />, label: '订单管理' },
  { key: '/card-keys', icon: <KeyOutlined />, label: '卡密管理' },
  { key: '/payment-channels', icon: <CreditCardOutlined />, label: '支付渠道' },
  { key: '/oauth-providers', icon: <LoginOutlined />, label: '第三方登录' },
];
```

## 页面示例

### 商品管理

- 表格列表：名称、价格、库存、销量、状态
- 新建/编辑 Modal：表单校验
- 上架/下架操作
- 支持分页和筛选

### OAuth 管理

- 表格展示 GitHub/Google 配置
- Switch 开关快速启用/禁用
- 编辑 Modal 配置 Client ID/Secret
- Client Secret 脱敏显示

## 开发命令

```bash
cd admin-panel

# 开发模式
npm run dev

# 生产构建
npm run build

# 预览构建结果
npm run preview
```

## 自定义主题

Ant Design 5 支持通过 ConfigProvider 定制主题：

```tsx title="admin-panel/src/App.tsx"
<ConfigProvider locale={zhCN} theme={{
  token: {
    colorPrimary: '#1677ff',
  },
}}>
  <RouterProvider router={router} />
</ConfigProvider>
```
