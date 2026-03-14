package main

import (
	"fmt"
	"log"

	"github.com/x-store/backend/internal/config"
	"github.com/x-store/backend/internal/model"
	"github.com/x-store/backend/internal/router"
)

func main() {
	// 1. 加载配置
	config.Load("config.yaml")

	// 2. 初始化 MySQL
	config.InitDatabase()

	// 3. 自动迁移表结构
	if err := config.DB.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.ProductDetail{},
		&model.Order{},
		&model.Payment{},
		&model.CardKey{},
		&model.PaymentChannel{},
	); err != nil {
		log.Fatalf("[Migration] Failed: %v", err)
	}
	log.Println("[Migration] Database tables migrated successfully")

	// 4. 初始化 Redis
	config.InitRedis()

	// 5. 同步商品库存到 Redis 缓存
	config.InitStockCache()

	// 6. 启动 HTTP 服务
	r := router.Setup(config.Global.Server.Mode)
	addr := fmt.Sprintf(":%d", config.Global.Server.Port)
	log.Printf("[Server] Starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[Server] Failed to start: %v", err)
	}
}
