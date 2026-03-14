---
sidebar_position: 1
title: 后端概览
---

# 后端概览

X-Store 后端基于 **Go + Gin + GORM** 构建，提供 RESTful API 服务。

## 入口文件

```go title="backend/cmd/main.go"
func main() {
    // 1. 加载配置
    config.LoadConfig()

    // 2. 初始化数据库（PostgreSQL + AutoMigrate）
    config.InitDatabase()

    // 3. 初始化 Redis
    config.InitRedis()

    // 4. 注册路由 & 启动 HTTP 服务
    r := router.SetupRouter()
    r.Run(fmt.Sprintf(":%d", config.Global.Server.Port))
}
```

## 启动流程

```
LoadConfig()         读取 config.yaml 配置
      ↓
InitDatabase()       连接 PostgreSQL，AutoMigrate 建表
      ↓
InitRedis()          连接 Redis
      ↓
SetupRouter()        注册所有路由 + 中间件
      ↓
Run(:8082)           启动 HTTP 服务
```

## 代码组织原则

1. **Handler 层**：只负责参数解析和响应输出，不写业务逻辑
2. **Service 层**：所有业务逻辑在这里处理，包括事务管理
3. **Repo 层**：只负责数据库 CRUD，不包含业务判断
4. **Model 层**：纯数据结构定义，对应数据库表

## 统一响应格式

所有接口返回统一的 JSON 结构：

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

| code | 含义 |
|------|------|
| `0` | 成功 |
| `400` | 参数错误 |
| `401` | 未认证 |
| `403` | 无权限 |
| `404` | 资源不存在 |
| `500` | 服务器内部错误 |

## 编译 & 运行

```bash
# 开发模式
cd backend && go run cmd/main.go

# 生产编译
go build -o x-store-backend ./cmd/main.go
./x-store-backend
```
