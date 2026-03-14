# X-Store - 数字商品交易平台

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)]()
[![Go Version](https://img.shields.io/badge/go-1.25-blue)]()
[![Node Version](https://img.shields.io/badge/node-20+-green)]()
[![License](https://img.shields.io/badge/license-MIT-blue)]()

一个完整的数字商品（卡密、账号、激活码）交易平台，包含 C 端商城、管理后台和后端 API 服务。

## ✨ 特性

- 🛍️ **完整的电商流程** - 商品展示、购物车、订单、支付、发货
- 🔐 **安全可靠** - JWT 认证、防重放攻击、限流保护
- 📧 **邮件通知** - 订单发货自动发送邮件
- 📊 **数据统计** - GMV、订单趋势、热销商品分析
- 🎨 **现代化 UI** - 响应式设计、移动端适配
- ⚡ **高性能** - Redis 缓存、库存锁定机制
- 🔧 **易于部署** - Docker 支持、一键启动脚本

## 🚀 快速开始

### 环境要求

- Go 1.25+
- Node.js 20+
- PostgreSQL 15+
- Redis 7+

### 1. 克隆项目

```bash
git clone <your-repo-url>
cd x-store
```

### 2. 初始化数据库

```bash
# 创建数据库
docker run --rm -e PGPASSWORD='Postgres@2026' postgres:15 \
  psql -h host.docker.internal -p 5432 -U admin -d postgres \
  -c "CREATE DATABASE x_store;"

# 导入初始数据
docker run --rm -e PGPASSWORD='Postgres@2026' -i postgres:15 \
  psql -h host.docker.internal -p 5432 -U admin -d x_store \
  < backend/init_postgres.sql
```

### 3. 启动服务

```bash
# 方式一：使用启动脚本（推荐）
chmod +x start.sh

# 在不同终端分别运行
./start.sh backend    # 启动后端
./start.sh frontend   # 启动 C 端
./start.sh admin      # 启动管理后台

# 方式二：手动启动
# 后端
cd backend && go build -o x-store-backend ./cmd/main.go && ./x-store-backend

# C 端
cd frontend-store && npm install && npm run dev

# 管理后台
cd admin-panel && npm install && npm run dev
```

### 4. 访问应用

- **C 端商城**: http://localhost:3000
- **管理后台**: http://localhost:5174
- **后端 API**: http://localhost:8082

### 5. 测试账号

**管理员**:
- 用户名: `admin`
- 密码: `admin123`

**普通用户**:
- 用户名: `testuser`
- 密码: `admin123`

## 📁 项目结构

```
x-store/
├── backend/              # Go 后端服务
│   ├── cmd/             # 入口文件
│   ├── internal/        # 内部代码
│   │   ├── handler/    # HTTP 处理器
│   │   ├── service/    # 业务逻辑
│   │   ├── repo/       # 数据仓库
│   │   └── model/      # 数据模型
│   └── pkg/            # 公共包
│
├── frontend-store/      # Next.js C 端商城
│   └── src/
│       ├── app/        # 页面
│       └── components/ # 组件
│
├── admin-panel/         # React 管理后台
│   └── src/
│       ├── pages/      # 页面
│       └── api/        # API 调用
│
└── project-docs/        # 项目文档
```

## 🎯 核心功能

### C 端商城
- ✅ 商品浏览（瀑布流）
- ✅ 分类筛选
- ✅ 商品搜索
- ✅ 购物车
- ✅ 订单创建
- ✅ 支付流程
- ✅ 卡密展示（刮刮卡效果）
- ✅ 用户登录/注册
- ✅ 个人中心
- ✅ 订单历史

### 管理后台
- ✅ 数据统计（GMV、订单趋势）
- ✅ 商品管理（CRUD、上下架）
- ✅ 分类管理
- ✅ 订单管理（搜索、筛选、详情）
- ✅ 卡密管理（批量导入、库存预警）
- ✅ 支付渠道管理
- ✅ 手动发货
- ✅ 订单退款

### 后端 API
- ✅ RESTful API 设计
- ✅ JWT 认证
- ✅ 权限控制
- ✅ 防重放攻击
- ✅ 限流保护
- ✅ 邮件通知
- ✅ 支付集成（模拟支付 + Stripe）

## 📚 文档

- [API 文档](./API_DOCUMENTATION.md)
- [项目完成报告](./PROJECT_COMPLETION_REPORT.md)
- [需求文档](./project-docs/requirements.md)
- [技术规格](./project-docs/tech-specs.md)

## 🛠️ 技术栈

### 后端
- **框架**: Gin (Go)
- **数据库**: PostgreSQL
- **缓存**: Redis
- **ORM**: GORM
- **认证**: JWT

### C 端前台
- **框架**: Next.js 14
- **语言**: TypeScript
- **样式**: Tailwind CSS
- **组件**: shadcn/ui

### 管理后台
- **框架**: React 18 + Vite
- **语言**: TypeScript
- **UI 库**: Ant Design 5
- **状态管理**: Zustand
- **数据请求**: React Query

## 🔧 配置

### 后端配置 (backend/config.yaml)

```yaml
server:
  port: 8082
  mode: debug

database:
  host: localhost
  port: 5432
  user: admin
  password: "Postgres@2026"
  dbname: x_store

redis:
  host: 127.0.0.1
  port: 6379

email:
  enabled: false  # 设置为 true 启用邮件
  host: smtp.gmail.com
  port: 587
  username: your-email@gmail.com
  password: your-app-password
```

## 🐳 Docker 部署

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## 📊 数据库表

- `users` - 用户表
- `categories` - 分类表
- `products` - 商品表
- `product_details` - 商品详情表
- `orders` - 订单表
- `card_keys` - 卡密表
- `payments` - 支付记录表
- `payment_channels` - 支付渠道表

## 🔐 安全性

- JWT Token 认证
- 密码 bcrypt 加密
- 防重放攻击（签名验证）
- 限流保护
- CORS 跨域控制
- SQL 注入防护（GORM 参数化查询）

## 🎨 界面预览

### C 端商城
- 响应式设计，支持移动端
- 商品瀑布流展示
- 刮刮卡效果展示卡密

### 管理后台
- Ant Design 5 现代化界面
- 数据可视化图表
- 批量操作功能

## 📝 开发计划

- [ ] 集成真实支付（Stripe、支付宝、微信支付）
- [ ] 多语言支持（i18n）
- [ ] SEO 优化
- [ ] 单元测试
- [ ] CI/CD 自动部署
- [ ] 监控告警系统

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 License

MIT License

## 👨‍💻 作者

Cascade AI Assistant

## 🙏 致谢

感谢所有开源项目的贡献者！

---

**注意**: 这是一个演示项目，生产环境使用前请修改所有密钥和密码！
