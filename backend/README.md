# X-Store Backend API

极致优雅、安全可靠的数字商品交易平台后端服务

## 技术栈

- **Go 1.24+** - 高性能后端语言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **MySQL 8.0** - 主数据库
- **Redis 7** - 缓存 + 库存 + 限流
- **JWT** - 身份认证
- **HMAC-SHA256** - 防重放签名

## 快速开始

### 1. 启动依赖服务

```bash
# 在项目根目录启动 MySQL 和 Redis
cd /Volumes/dz/code/x-store
docker-compose up -d

# 等待服务启动完成（约 10 秒）
docker-compose ps
```

### 2. 启动后端

```bash
cd backend
go run cmd/main.go
```

启动成功后会看到：
```
[Config] Loaded successfully (mode=debug, port=8080)
[Database] Connected to 127.0.0.1:3306/x_store
[Migration] Database tables migrated successfully
[Redis] Connected to 127.0.0.1:6379
[Stock] Synced 0 products to Redis cache
[Server] Starting on :8080
```

### 3. 测试 API

```bash
# 健康检查
curl http://localhost:8080/api/health

# 获取分类列表
curl http://localhost:8080/api/categories

# 获取商品列表
curl http://localhost:8080/api/products
```

## API 文档

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/health` | 健康检查 |
| GET | `/api/categories` | 分类列表 |
| GET | `/api/categories/:id` | 分类详情 |
| GET | `/api/products` | 商品列表（支持分页/筛选/排序） |
| GET | `/api/products/:id` | 商品详情 |
| GET | `/api/orders/:order_no` | 订单详情 |

### 管理员接口（需要 JWT + Admin 权限）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/admin/categories` | 创建分类 |
| PUT | `/api/admin/categories/:id` | 更新分类 |
| DELETE | `/api/admin/categories/:id` | 删除分类 |
| POST | `/api/admin/products` | 创建商品 |
| PUT | `/api/admin/products/:id` | 更新商品 |
| DELETE | `/api/admin/products/:id` | 删除商品 |
| POST | `/api/admin/cardkeys/import` | 批量导入卡密 |
| GET | `/api/admin/cardkeys/count/:product_id` | 统计可用卡密 |

### 签名保护接口（需要防重放签名 + 限流）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/orders` | 创建订单 |
| POST | `/api/orders/:order_no/pay` | 支付订单 |
| POST | `/api/orders/:order_no/cancel` | 取消订单 |

## 完整测试流程

### 1. 创建分类（需要管理员 Token）

```bash
# TODO: 先实现登录接口获取 JWT Token
# 暂时跳过，直接在数据库插入测试数据
```

### 2. 创建商品

```bash
# 插入测试商品
mysql -h 127.0.0.1 -u root -p -e "
USE x_store;
INSERT INTO categories (name, icon, sort, status, created_at, updated_at) 
VALUES ('AI 大模型', '🤖', 100, 1, NOW(), NOW());

INSERT INTO products (category_id, title, price, stock, delivery_type, tags, status, is_new, created_at, updated_at)
VALUES (1, 'ChatGPT Plus 独享账号', 120.00, 100, 'auto', '[\"自动发货\",\"独享资源\"]', 1, 1, NOW(), NOW());
"
```

### 3. 导入卡密

```bash
# TODO: 实现管理员登录后测试
```

### 4. 测试下单流程

```bash
# 1. 查询商品
curl http://localhost:8080/api/products

# 2. 创建订单（需要签名）
# TODO: 前端实现签名逻辑后测试

# 3. 支付订单
# TODO: 测试支付流程

# 4. 查询订单
curl http://localhost:8080/api/orders/XS20260314001
```

## 核心特性

### 🔒 防超卖机制

- **Redis Lua 脚本** - 原子性库存扣减
- **数据库行锁** - 卡密预占锁定
- **双重保障** - Redis + MySQL 双重验证

### 🛡️ 安全防护

- **JWT 认证** - 用户身份验证
- **防重放签名** - HMAC-SHA256 + 时间戳 + Nonce
- **Redis 限流** - IP 级别请求限流
- **CORS 跨域** - 安全的跨域配置

### ⚡ 高性能

- **Redis 缓存** - 库存热数据缓存
- **索引优化** - 数据库查询优化
- **连接池** - MySQL/Redis 连接复用

## 数据库表结构

- `users` - 用户表
- `categories` - 分类表
- `products` - 商品主表
- `product_details` - 商品详情表
- `orders` - 订单表
- `payments` - 支付记录表
- `card_keys` - 卡密表

## 配置说明

编辑 `config.yaml` 修改配置：

```yaml
server:
  port: 8080
  mode: debug # debug | release

database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: root
  dbname: x_store

redis:
  host: 127.0.0.1
  port: 6379
```

## 开发进度

- [x] Phase 2: 后端核心基础建设
  - [x] 配置管理
  - [x] 数据库模型
  - [x] 中间件（JWT/签名/限流/CORS）
  - [x] 分类模块 CRUD
  - [x] 商品模块 CRUD
- [x] Phase 3: 交易链路与风控
  - [x] Redis Lua 库存扣减
  - [x] 订单创建与状态流转
  - [x] 卡密管理与发货
  - [ ] 支付网关对接（待实现）
  - [ ] Asynq 延时队列（待实现）
- [ ] Phase 4: 用户认证与授权
  - [ ] 注册/登录接口
  - [ ] 邮箱验证
  - [ ] 密码加密
- [ ] Phase 5: 完善与优化
  - [ ] 订单超时自动取消
  - [ ] 邮件发送服务
  - [ ] 数据统计看板

## 下一步

1. 实现用户注册/登录接口
2. 完善管理员权限体系
3. 对接真实支付网关
4. 实现 Asynq 延时队列处理订单超时
5. 添加邮件发送服务
