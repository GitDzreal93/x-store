---
sidebar_position: 1
title: Docker 部署
---

# Docker 部署

使用 Docker Compose 一键部署 X-Store 全栈服务。

## docker-compose.yml

```yaml title="docker-compose.yml"
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: admin123
      POSTGRES_DB: x_store
    ports:
      - "5432:5432"
    volumes:
      - pg_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "8082:8082"
    depends_on:
      - postgres
      - redis
    environment:
      - X_STORE_DB_HOST=postgres
      - X_STORE_REDIS_HOST=redis

  frontend:
    build:
      context: ./frontend-store
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    depends_on:
      - backend

  admin:
    build:
      context: ./admin-panel
      dockerfile: Dockerfile
    ports:
      - "5173:80"
    depends_on:
      - backend

volumes:
  pg_data:
```

## 后端 Dockerfile

```dockerfile title="backend/Dockerfile"
# 构建阶段
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o x-store-backend ./cmd/main.go

# 运行阶段
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/x-store-backend .
COPY --from=builder /app/config.yaml .
EXPOSE 8082
CMD ["./x-store-backend"]
```

## 前端 Dockerfile

```dockerfile title="frontend-store/Dockerfile"
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public
EXPOSE 3000
CMD ["node", "server.js"]
```

## 管理后台 Dockerfile

```dockerfile title="admin-panel/Dockerfile"
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
```

## 启动命令

```bash
# 构建并启动所有服务
docker compose up -d --build

# 查看日志
docker compose logs -f backend

# 初始化数据库
docker compose exec backend sh -c \
  "psql -h postgres -U admin -d x_store -f init_postgres.sql"

# 停止所有服务
docker compose down
```

## 数据持久化

PostgreSQL 数据通过 Docker Volume `pg_data` 持久化，即使容器重启数据也不会丢失。

```bash
# 查看数据卷
docker volume ls

# 备份数据库
docker compose exec postgres pg_dump -U admin x_store > backup.sql
```
