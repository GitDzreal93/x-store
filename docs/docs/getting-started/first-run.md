---
sidebar_position: 4
title: 首次运行
---

# 首次运行

## 启动后端

```bash
cd backend
go run cmd/main.go
```

你应该看到类似输出：

```
[Database] Connected to PostgreSQL localhost:5432/x_store
[Redis] Connected to localhost:6379
[Server] Running on :8082 (debug mode)
```

后端会自动执行 GORM AutoMigrate，创建所有数据表。

## 启动 C 端前端

```bash
cd frontend-store
npm run dev
```

访问 `http://localhost:3000` 即可看到商城首页。

## 启动管理后台

```bash
cd admin-panel
npm run dev
```

访问 `http://localhost:5173`，使用默认管理员账号登录：

| 字段 | 值 |
|------|-----|
| 用户名 | `admin` |
| 密码 | `admin123` |

## 功能验证清单

启动成功后，你可以按以下顺序验证各个功能：

### 1. 管理后台

- [x] 登录管理后台
- [x] 查看仪表盘数据统计
- [x] 管理商品分类（增删改）
- [x] 管理商品（创建、上架、下架）
- [x] 批量导入卡密
- [x] 配置支付渠道
- [x] 配置第三方登录（OAuth）
- [x] 查看订单列表

### 2. C 端商城

- [x] 浏览商品列表
- [x] 搜索商品
- [x] 查看商品详情
- [x] 注册 / 登录
- [x] GitHub / Google 第三方登录
- [x] 创建订单
- [x] 模拟支付
- [x] 查看个人中心 & 订单记录

### 3. API 测试

```bash
# 健康检查
curl http://localhost:8082/api/categories

# 用户注册
curl -X POST http://localhost:8082/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@test.com","password":"123456"}'

# 查看商品
curl http://localhost:8082/api/products
```

## 常见问题

### 数据库连接失败

```
[Database] Failed to connect: ...
```

**解决**: 检查 PostgreSQL 是否运行，`config.yaml` 中的数据库配置是否正确。

### Redis 连接失败

```
[Redis] Failed to connect: ...
```

**解决**: 检查 Redis 是否运行 (`redis-cli ping`)。

### 端口被占用

```
listen tcp :8082: bind: address already in use
```

**解决**: 修改 `config.yaml` 中的 `server.port`，或终止占用端口的进程。

## 下一步

👉 前往 [项目架构](/docs/architecture/overview) 深入了解系统设计。
