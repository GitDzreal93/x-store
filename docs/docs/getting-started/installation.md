---
sidebar_position: 2
title: 安装部署
---

# 安装部署

## 克隆项目

```bash
git clone https://github.com/x-store/x-store.git
cd x-store
```

## 项目目录

```
x-store/
├── backend/           # Go 后端
├── frontend-store/    # C 端前端 (Next.js)
├── admin-panel/       # 管理后台 (React + Antd)
├── docs/              # 文档站 (Docusaurus)
└── start.sh           # 一键启动脚本
```

## 安装依赖

### 后端

```bash
cd backend
go mod download
```

### C 端前端

```bash
cd frontend-store
npm install
```

### 管理后台

```bash
cd admin-panel
npm install
```

### 文档站（可选）

```bash
cd docs
npm install
```

## 初始化数据库

后端首次启动会自动建表，但你可以手动导入示例数据：

```bash
psql -h localhost -p 5432 -U admin -d x_store -f backend/init_postgres.sql
```

该脚本会创建：
- 管理员账号（`admin` / `admin123`）
- 示例分类和商品
- 示例卡密数据
- 支付渠道配置
- OAuth 提供商配置

## 一键启动

项目根目录提供了一键启动脚本：

```bash
chmod +x start.sh
./start.sh
```

启动后可访问：

| 服务 | 地址 | 说明 |
|------|------|------|
| 后端 API | http://localhost:8082 | Go + Gin |
| C 端前端 | http://localhost:3000 | Next.js |
| 管理后台 | http://localhost:5173 | Vite + React |
| 文档站 | http://localhost:3001 | Docusaurus |

## 下一步

👉 前往 [配置说明](/docs/getting-started/configuration) 了解各项配置。
