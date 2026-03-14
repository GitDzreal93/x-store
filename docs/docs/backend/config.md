---
sidebar_position: 2
title: 配置管理
---

# 配置管理

后端使用 [Viper](https://github.com/spf13/viper) 库加载 YAML 配置文件。

## 配置结构体

```go title="backend/internal/config/config.go"
type Config struct {
    Server    ServerConfig    `mapstructure:"server"`
    Database  DatabaseConfig  `mapstructure:"database"`
    Redis     RedisConfig     `mapstructure:"redis"`
    JWT       JWTConfig       `mapstructure:"jwt"`
    Signature SignatureConfig `mapstructure:"signature"`
    Email     EmailConfig     `mapstructure:"email"`
}
```

每个子配置都是独立的结构体，通过 `mapstructure` tag 映射 YAML 字段。

## 加载流程

```go
func LoadConfig() {
    viper.SetConfigFile("config.yaml")
    viper.SetConfigType("yaml")

    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("Failed to read config: %v", err)
    }

    if err := viper.Unmarshal(&Global); err != nil {
        log.Fatalf("Failed to unmarshal config: %v", err)
    }
}
```

`config.Global` 是全局配置实例，在其他模块中通过 `config.Global.XXX` 访问。

## 数据库初始化

```go title="backend/internal/config/database.go"
func InitDatabase() {
    DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{...})

    // AutoMigrate 自动建表
    DB.AutoMigrate(
        &model.User{},
        &model.Category{},
        &model.Product{},
        &model.Order{},
        &model.CardKey{},
        &model.Payment{},
        &model.PaymentChannel{},
        &model.OAuthProvider{},
    )
}
```

:::tip
GORM 的 AutoMigrate 只会创建缺失的表和字段，不会删除已有数据。生产环境建议使用数据库迁移工具。
:::

## Redis 初始化

Redis 用于：
- **防重放攻击**：存储 nonce，防止请求重放
- **接口限流**：基于令牌桶算法的请求频率限制

```go
func InitRedis() {
    RDB = redis.NewClient(&redis.Options{
        Addr:     cfg.Addr(),  // host:port
        Password: cfg.Password,
        DB:       cfg.DB,
    })
}
```
