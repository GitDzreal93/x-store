package handler

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/x-store/backend/internal/model"
	"github.com/x-store/backend/internal/repo"
	"github.com/x-store/backend/pkg/response"
)

type CardKeyHandler struct {
	repo *repo.CardKeyRepo
}

func NewCardKeyHandler() *CardKeyHandler {
	return &CardKeyHandler{repo: repo.NewCardKeyRepo()}
}

// BatchImport 批量导入卡密 [POST /api/admin/cardkeys/import]
func (h *CardKeyHandler) BatchImport(c *gin.Context) {
	var req struct {
		ProductID uint   `json:"product_id" binding:"required"`
		Content   string `json:"content" binding:"required"` // 多行文本，每行一个卡密
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 解析卡密（按行分割，去重）
	lines := strings.Split(req.Content, "\n")
	keysMap := make(map[string]bool)
	var keys []model.CardKey

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || keysMap[line] {
			continue
		}
		keysMap[line] = true
		keys = append(keys, model.CardKey{
			ProductID: req.ProductID,
			Content:   line,
			Status:    model.CardKeyStatusAvailable,
		})
	}

	if len(keys) == 0 {
		response.BadRequest(c, "没有有效的卡密")
		return
	}

	if err := h.repo.BatchCreate(keys); err != nil {
		response.ServerError(c, "导入失败: "+err.Error())
		return
	}

	response.OK(c, gin.H{
		"imported": len(keys),
		"message":  "导入成功",
	})
}

// CountAvailable 统计可用卡密数量 [GET /api/admin/cardkeys/count/:product_id]
func (h *CardKeyHandler) CountAvailable(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的商品ID")
		return
	}

	count, err := h.repo.CountAvailable(uint(productID))
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.OK(c, gin.H{
		"product_id": productID,
		"available":  count,
	})
}
