---
sidebar_position: 3
title: 目录结构
---

# 目录结构

## 项目根目录

```
x-store/
├── backend/             # Go 后端服务
├── frontend-store/      # C 端买家前端
├── admin-panel/         # B 端管理后台
├── docs/                # 项目文档站
├── start.sh             # 一键启动脚本
└── README.md
```

## 后端目录

```
backend/
├── cmd/
│   └── main.go              # 入口文件
├── config.yaml              # 配置文件
├── init_postgres.sql        # 数据库初始化脚本
├── internal/
│   ├── config/
│   │   ├── config.go        # 配置结构体 & 加载
│   │   └── database.go      # 数据库初始化 (GORM)
│   ├── handler/             # HTTP Handler 层
│   │   ├── user_handler.go
│   │   ├── product_handler.go
│   │   ├── order_handler.go
│   │   ├── category_handler.go
│   │   ├── cardkey_handler.go
│   │   ├── payment_handler.go
│   │   ├── stats_handler.go
│   │   ├── oauth_handler.go
│   │   └── oauth_admin_handler.go
│   ├── middleware/           # 中间件
│   │   ├── auth.go          # JWT 认证
│   │   ├── admin.go         # 管理员权限
│   │   ├── anti_replay.go   # 防重放攻击
│   │   ├── cors.go          # 跨域处理
│   │   └── rate_limit.go    # 接口限流
│   ├── model/               # 数据模型
│   │   ├── base.go          # BaseModel (ID + 时间戳)
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── order.go
│   │   ├── card_key.go
│   │   ├── payment.go
│   │   ├── payment_channel.go
│   │   ├── category.go
│   │   └── oauth_provider.go
│   ├── repo/                # 数据访问层
│   │   ├── user_repo.go
│   │   ├── product_repo.go
│   │   ├── order_repo.go
│   │   └── ...
│   ├── router/
│   │   └── router.go        # 路由注册
│   └── service/             # 业务逻辑层
│       ├── user_service.go
│       ├── product_service.go
│       ├── order_service.go
│       ├── payment_service.go
│       ├── oauth_service.go
│       ├── email_service.go
│       └── stats_service.go
└── pkg/                     # 公共工具包
    ├── crypto/
    │   └── jwt.go           # JWT 生成/解析
    ├── payment/             # 支付提供商实现
    │   ├── provider.go      # 统一接口定义
    │   ├── alipay.go
    │   ├── wechatpay.go
    │   ├── stripe.go
    │   ├── paypal.go
    │   └── ...
    └── response/
        └── response.go      # 统一 HTTP 响应
```

## C 端前端目录

```
frontend-store/
├── src/
│   ├── app/                 # Next.js App Router
│   │   ├── page.tsx         # 商城首页
│   │   ├── auth/
│   │   │   ├── page.tsx     # 登录/注册页
│   │   │   └── callback/
│   │   │       └── page.tsx # OAuth 回调页
│   │   ├── profile/
│   │   │   └── page.tsx     # 个人中心
│   │   ├── layout.tsx       # 根布局
│   │   └── globals.css      # 全局样式
│   ├── components/
│   │   └── ui/              # shadcn/ui 组件
│   └── lib/
│       └── utils.ts         # 工具函数
├── tailwind.config.ts
├── next.config.mjs
└── package.json
```

## 管理后台目录

```
admin-panel/
├── src/
│   ├── App.tsx              # 应用入口
│   ├── layouts/
│   │   └── AdminLayout.tsx  # 管理后台布局
│   ├── pages/               # 页面组件
│   │   ├── Dashboard.tsx    # 仪表盘
│   │   ├── Products.tsx     # 商品管理
│   │   ├── Categories.tsx   # 分类管理
│   │   ├── Orders.tsx       # 订单管理
│   │   ├── CardKeys.tsx     # 卡密管理
│   │   ├── PaymentChannels.tsx  # 支付渠道
│   │   ├── OAuthProviders.tsx   # 第三方登录
│   │   └── Login.tsx        # 登录页
│   ├── router/
│   │   └── index.tsx        # 路由配置
│   ├── stores/
│   │   └── auth.ts          # Zustand 状态管理
│   └── utils/
│       └── request.ts       # Axios 请求封装
├── vite.config.ts
└── package.json
```
