---
sidebar_position: 1
title: C 端商城
---

# C 端商城前端

C 端商城基于 **Next.js 14 (App Router) + TailwindCSS + shadcn/ui** 构建，面向买家用户。

## 页面结构

| 路由 | 页面 | 说明 |
|------|------|------|
| `/` | 商城首页 | 商品列表、分类筛选、搜索 |
| `/auth` | 登录/注册 | 用户名密码 + OAuth 登录 |
| `/auth/callback` | OAuth 回调 | 处理第三方登录回调 |
| `/profile` | 个人中心 | 用户信息 + 订单记录 |

## 技术要点

### App Router

使用 Next.js 14 的 App Router，所有页面为 Client Component：

```
src/app/
├── page.tsx              # 首页
├── layout.tsx            # 根布局
├── globals.css           # 全局样式
├── auth/
│   ├── page.tsx          # 登录注册
│   └── callback/
│       └── page.tsx      # OAuth 回调
└── profile/
    └── page.tsx          # 个人中心
```

### shadcn/ui 组件

项目使用 shadcn/ui 作为组件库，位于 `src/components/ui/`：

- `Button` - 按钮组件
- `Input` - 输入框
- `Card` - 卡片容器
- `Tabs` - 标签页切换

### API 调用

前端通过 `fetch` 直接调用后端 API：

```tsx
const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8082';

// 获取商品列表
const res = await fetch(`${API_BASE}/api/products`);
const data = await res.json();

// 带认证的请求
const token = localStorage.getItem('token');
const res = await fetch(`${API_BASE}/api/user/profile`, {
  headers: { Authorization: `Bearer ${token}` },
});
```

### OAuth 登录集成

登录页动态获取已启用的 OAuth 提供商：

```tsx
useEffect(() => {
  fetch(`${API_BASE}/api/oauth/providers`)
    .then(res => res.json())
    .then(data => {
      if (data.code === 0 && data.data) {
        setOauthProviders(data.data);
      }
    });
}, []);
```

OAuth 回调页处理 token 存储和跳转：

```tsx
// /auth/callback?token=xxx
const token = searchParams.get('token');
if (token) {
  localStorage.setItem('token', token);
  router.push('/profile');
}
```

### 搜索功能

首页集成了商品搜索，通过 URL 参数传递：

```tsx
const handleSearch = (keyword: string) => {
  fetch(`${API_BASE}/api/products?keyword=${keyword}`)
    .then(res => res.json())
    .then(data => setProducts(data.data));
};
```

## 开发命令

```bash
cd frontend-store

# 开发模式（热更新）
npm run dev

# 生产构建
npm run build

# 启动生产服务
npm run start
```

## 环境变量

在 `.env.local` 中配置：

```bash
NEXT_PUBLIC_API_URL=http://localhost:8082
```

## 样式定制

项目使用 TailwindCSS，全局样式在 `globals.css` 中定义。主题色通过 CSS 变量控制：

```css title="src/app/globals.css"
:root {
  --background: 0 0% 100%;
  --foreground: 222.2 84% 4.9%;
  --primary: 222.2 47.4% 11.2%;
  /* ... */
}
```
