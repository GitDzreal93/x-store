---
sidebar_position: 3
title: 配置说明
---

# 配置说明

后端配置文件位于 `backend/config.yaml`，采用 YAML 格式。

## 完整配置示例

```yaml title="backend/config.yaml"
server:
  port: 8082
  mode: debug          # debug | release

database:
  host: localhost
  port: 5432
  user: admin
  password: admin123
  dbname: x_store
  sslmode: disable
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-jwt-secret-change-in-production
  expire_hours: 72

signature:
  secret: x-store-signature-secret-key-change-in-production
  expire_seconds: 60

email:
  enabled: false
  host: smtp.gmail.com
  port: 587
  username: your-email@gmail.com
  password: your-app-password
  from: "X-Store <noreply@x-store.com>"
```

## 配置项详解

### Server

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `port` | int | 8082 | HTTP 监听端口 |
| `mode` | string | debug | 运行模式，`debug` 开启详细日志 |

### Database (PostgreSQL)

| 字段 | 类型 | 说明 |
|------|------|------|
| `host` | string | 数据库地址 |
| `port` | int | 数据库端口，默认 5432 |
| `user` | string | 数据库用户名 |
| `password` | string | 数据库密码 |
| `dbname` | string | 数据库名称 |
| `sslmode` | string | SSL 模式，本地开发用 `disable` |
| `max_idle_conns` | int | 最大空闲连接数 |
| `max_open_conns` | int | 最大打开连接数 |

### Redis

| 字段 | 类型 | 说明 |
|------|------|------|
| `host` | string | Redis 地址 |
| `port` | int | Redis 端口，默认 6379 |
| `password` | string | 密码，无密码留空 |
| `db` | int | 数据库编号 |

### JWT

| 字段 | 类型 | 说明 |
|------|------|------|
| `secret` | string | JWT 签名密钥，**生产环境必须修改** |
| `expire_hours` | int | Token 有效时长（小时） |

### Signature（防重放签名）

| 字段 | 类型 | 说明 |
|------|------|------|
| `secret` | string | 签名密钥 |
| `expire_seconds` | int | 签名有效期（秒） |

### Email

| 字段 | 类型 | 说明 |
|------|------|------|
| `enabled` | bool | 是否启用邮件通知 |
| `host` | string | SMTP 服务器地址 |
| `port` | int | SMTP 端口 |
| `username` | string | SMTP 用户名 |
| `password` | string | SMTP 密码或应用专用密码 |
| `from` | string | 发件人显示名称 |

### OAuth（第三方登录）

:::info
OAuth 配置已迁移到数据库管理。请在管理后台「第三方登录」页面配置 GitHub / Google OAuth 的 Client ID 和 Client Secret。
:::

## 环境变量

你也可以通过环境变量覆盖配置（优先级高于 config.yaml）：

```bash
export X_STORE_DB_HOST=localhost
export X_STORE_DB_PASSWORD=your_password
export X_STORE_JWT_SECRET=your_jwt_secret
```

## 下一步

👉 前往 [首次运行](/docs/getting-started/first-run) 启动项目。
