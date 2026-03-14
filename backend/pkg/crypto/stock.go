package crypto

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// DeductStockLua Redis Lua 脚本：原子性扣减库存（防超卖）
// 返回: 1=成功, 0=库存不足, -1=商品不存在
const DeductStockLua = `
local key = KEYS[1]
local quantity = tonumber(ARGV[1])

if redis.call('EXISTS', key) == 0 then
    return -1
end

local stock = tonumber(redis.call('GET', key))
if stock < quantity then
    return 0
end

redis.call('DECRBY', key, quantity)
return 1
`

// ReleaseStockLua Redis Lua 脚本：释放库存（订单取消/超时）
const ReleaseStockLua = `
local key = KEYS[1]
local quantity = tonumber(ARGV[1])

if redis.call('EXISTS', key) == 0 then
    return -1
end

redis.call('INCRBY', key, quantity)
return 1
`

// StockManager 库存管理器
type StockManager struct {
	rdb           *redis.Client
	deductScript  *redis.Script
	releaseScript *redis.Script
}

func NewStockManager(rdb *redis.Client) *StockManager {
	return &StockManager{
		rdb:           rdb,
		deductScript:  redis.NewScript(DeductStockLua),
		releaseScript: redis.NewScript(ReleaseStockLua),
	}
}

// StockKey 生成库存 Redis Key
func StockKey(productID uint) string {
	return fmt.Sprintf("stock:%d", productID)
}

// InitStock 初始化商品库存到 Redis
func (m *StockManager) InitStock(ctx context.Context, productID uint, stock int) error {
	key := StockKey(productID)
	return m.rdb.Set(ctx, key, stock, 0).Err()
}

// GetStock 获取当前库存
func (m *StockManager) GetStock(ctx context.Context, productID uint) (int, error) {
	key := StockKey(productID)
	val, err := m.rdb.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, fmt.Errorf("库存不存在")
	}
	return val, err
}

// DeductStock 扣减库存（原子性）
func (m *StockManager) DeductStock(ctx context.Context, productID uint, quantity int) error {
	key := StockKey(productID)
	result, err := m.deductScript.Run(ctx, m.rdb, []string{key}, quantity).Int()
	if err != nil {
		return fmt.Errorf("扣减库存失败: %w", err)
	}
	switch result {
	case -1:
		return fmt.Errorf("商品库存不存在")
	case 0:
		return fmt.Errorf("库存不足")
	case 1:
		return nil
	default:
		return fmt.Errorf("未知错误")
	}
}

// ReleaseStock 释放库存（订单取消/超时）
func (m *StockManager) ReleaseStock(ctx context.Context, productID uint, quantity int) error {
	key := StockKey(productID)
	result, err := m.releaseScript.Run(ctx, m.rdb, []string{key}, quantity).Int()
	if err != nil {
		return fmt.Errorf("释放库存失败: %w", err)
	}
	if result == -1 {
		return fmt.Errorf("商品库存不存在")
	}
	return nil
}

// SyncStockFromDB 从数据库同步库存到 Redis（启动时调用）
func (m *StockManager) SyncStockFromDB(ctx context.Context, products map[uint]int) error {
	for productID, stock := range products {
		if err := m.InitStock(ctx, productID, stock); err != nil {
			return err
		}
	}
	return nil
}
