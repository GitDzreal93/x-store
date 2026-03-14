package config

import (
	"context"
	"log"

	"github.com/x-store/backend/pkg/crypto"
)

// InitStockCache 启动时同步商品库存到 Redis
func InitStockCache() {
	ctx := context.Background()
	stockMgr := crypto.NewStockManager(RDB)

	// 从数据库查询所有商品库存
	var products []struct {
		ID    uint
		Stock int
	}
	if err := DB.Table("products").Select("id, stock").Where("status = 1").Find(&products).Error; err != nil {
		log.Printf("[Stock] Failed to load products: %v", err)
		return
	}

	// 同步到 Redis
	stockMap := make(map[uint]int)
	for _, p := range products {
		stockMap[p.ID] = p.Stock
	}

	if err := stockMgr.SyncStockFromDB(ctx, stockMap); err != nil {
		log.Printf("[Stock] Failed to sync to Redis: %v", err)
		return
	}

	log.Printf("[Stock] Synced %d products to Redis cache", len(products))
}
