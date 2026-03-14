---
sidebar_position: 2
title: 技术栈
---

# 技术栈

## 后端

| 技术 | 版本 | 用途 |
|------|------|------|
| [Go](https://go.dev/) | 1.21+ | 主开发语言 |
| [Gin](https://gin-gonic.com/) | v1.9 | HTTP Web 框架 |
| [GORM](https://gorm.io/) | v2 | ORM 框架，数据库操作 |
| [PostgreSQL](https://www.postgresql.org/) | 14+ | 关系型数据库 |
| [Redis](https://redis.io/) | 6+ | 缓存、防重放、限流 |
| [JWT](https://github.com/golang-jwt/jwt) | v5 | 用户认证 |
| [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) | - | 密码哈希 |
| [Viper](https://github.com/spf13/viper) | v1 | 配置管理 |

## C 端前端（买家商城）

| 技术 | 版本 | 用途 |
|------|------|------|
| [Next.js](https://nextjs.org/) | 14 | React 全栈框架（App Router） |
| [React](https://react.dev/) | 18 | UI 库 |
| [TailwindCSS](https://tailwindcss.com/) | 3 | 原子化 CSS 框架 |
| [shadcn/ui](https://ui.shadcn.com/) | - | 高质量组件库 |
| [TypeScript](https://www.typescriptlang.org/) | 5 | 类型安全 |

## B 端管理后台

| 技术 | 版本 | 用途 |
|------|------|------|
| [React](https://react.dev/) | 18 | UI 库 |
| [Ant Design](https://ant.design/) | 5 | 企业级组件库 |
| [Vite](https://vitejs.dev/) | 5 | 构建工具 |
| [React Router](https://reactrouter.com/) | 6 | 路由管理 |
| [Zustand](https://zustand-demo.pmnd.rs/) | 4 | 状态管理 |
| [Axios](https://axios-http.com/) | 1 | HTTP 请求 |
| [TypeScript](https://www.typescriptlang.org/) | 5 | 类型安全 |

## 文档站

| 技术 | 版本 | 用途 |
|------|------|------|
| [Docusaurus](https://docusaurus.io/) | 3 | 静态文档生成器 |
| [React](https://react.dev/) | 18 | 自定义组件 |
| [MDX](https://mdxjs.com/) | 3 | Markdown + JSX |

## 为什么选择这些技术？

### Go + Gin

- **高性能**：编译型语言，天然适合 API 服务
- **简洁**：语法简单，学习曲线平缓
- **Gin 生态**：中间件丰富，性能优秀

### Next.js

- **SSR/SSG**：搜索引擎友好
- **App Router**：最新的 React 服务端组件
- **开发体验**：热更新、文件路由

### Ant Design

- **企业级**：组件丰富，开箱即用
- **国际化**：原生中文支持
- **主题定制**：Token 系统灵活可配

### PostgreSQL

- **功能强大**：JSON 支持、全文搜索、CTE
- **可靠性**：ACID 事务，数据安全
- **生态成熟**：运维工具丰富
