package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/x-store/backend/internal/service"
	"github.com/x-store/backend/pkg/response"
)

type StatsHandler struct {
	svc *service.StatsService
}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{svc: service.NewStatsService()}
}

// GetDashboardStats 获取仪表盘统计数据
func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.svc.GetDashboardStats()
	if err != nil {
		response.ServerError(c, "获取统计数据失败")
		return
	}
	response.OK(c, stats)
}
