package service

import (
	"time"

	"github.com/x-store/backend/internal/config"
	"github.com/x-store/backend/internal/model"
)

type StatsService struct{}

func NewStatsService() *StatsService {
	return &StatsService{}
}

type DashboardStats struct {
	TodayGMV      float64            `json:"today_gmv"`
	WeekGMV       float64            `json:"week_gmv"`
	MonthGMV      float64            `json:"month_gmv"`
	TotalOrders   int64              `json:"total_orders"`
	PendingOrders int64              `json:"pending_orders"`
	PaidOrders    int64              `json:"paid_orders"`
	TotalUsers    int64              `json:"total_users"`
	OrderTrend    []OrderTrendItem   `json:"order_trend"`
	TopProducts   []TopProductItem   `json:"top_products"`
}

type OrderTrendItem struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
	GMV   float64 `json:"gmv"`
}

type TopProductItem struct {
	ProductID uint    `json:"product_id"`
	Title     string  `json:"title"`
	Sales     int     `json:"sales"`
	Revenue   float64 `json:"revenue"`
}

// GetDashboardStats 获取仪表盘统计数据
func (s *StatsService) GetDashboardStats() (*DashboardStats, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -7)
	monthStart := todayStart.AddDate(0, -1, 0)

	stats := &DashboardStats{}

	// 今日 GMV
	config.DB.Model(&model.Order{}).
		Where("created_at >= ? AND status IN ?", todayStart, []int{model.OrderStatusPaid, model.OrderStatusDelivered}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.TodayGMV)

	// 本周 GMV
	config.DB.Model(&model.Order{}).
		Where("created_at >= ? AND status IN ?", weekStart, []int{model.OrderStatusPaid, model.OrderStatusDelivered}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.WeekGMV)

	// 本月 GMV
	config.DB.Model(&model.Order{}).
		Where("created_at >= ? AND status IN ?", monthStart, []int{model.OrderStatusPaid, model.OrderStatusDelivered}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.MonthGMV)

	// 订单统计
	config.DB.Model(&model.Order{}).Count(&stats.TotalOrders)
	config.DB.Model(&model.Order{}).Where("status = ?", model.OrderStatusPending).Count(&stats.PendingOrders)
	config.DB.Model(&model.Order{}).Where("status IN ?", []int{model.OrderStatusPaid, model.OrderStatusDelivered}).Count(&stats.PaidOrders)

	// 用户统计
	config.DB.Model(&model.User{}).Count(&stats.TotalUsers)

	// 近7天订单趋势
	stats.OrderTrend = s.getOrderTrend(7)

	// 热销商品 TOP 5
	stats.TopProducts = s.getTopProducts(5)

	return stats, nil
}

// getOrderTrend 获取订单趋势
func (s *StatsService) getOrderTrend(days int) []OrderTrendItem {
	var items []OrderTrendItem
	now := time.Now()

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		dayEnd := dayStart.Add(24 * time.Hour)

		var count int64
		var gmv float64

		config.DB.Model(&model.Order{}).
			Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
			Count(&count)

		config.DB.Model(&model.Order{}).
			Where("created_at >= ? AND created_at < ? AND status IN ?", dayStart, dayEnd, []int{model.OrderStatusPaid, model.OrderStatusDelivered}).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&gmv)

		items = append(items, OrderTrendItem{
			Date:  dateStr,
			Count: count,
			GMV:   gmv,
		})
	}

	return items
}

// getTopProducts 获取热销商品
func (s *StatsService) getTopProducts(limit int) []TopProductItem {
	var items []TopProductItem

	config.DB.Model(&model.Product{}).
		Select("id as product_id, title, sales, sales * price as revenue").
		Order("sales DESC").
		Limit(limit).
		Scan(&items)

	return items
}
