package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/x-store/backend/internal/service"
	"github.com/x-store/backend/pkg/response"
)

type PaymentHandler struct {
	svc *service.PaymentService
}

func NewPaymentHandler(orderSvc *service.OrderService) *PaymentHandler {
	return &PaymentHandler{svc: service.NewPaymentService(orderSvc)}
}

// Create 创建支付 [POST /api/payments]
func (h *PaymentHandler) Create(c *gin.Context) {
	var req service.CreatePaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	resp, err := h.svc.CreatePayment(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, resp)
}

// GetStatus 查询支付状态 [GET /api/payments/:id/status]
func (h *PaymentHandler) GetStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的支付ID")
		return
	}

	payment, err := h.svc.GetPaymentStatus(uint(id))
	if err != nil {
		response.NotFound(c, "支付记录不存在")
		return
	}
	response.OK(c, payment)
}

// Webhook 处理支付回调 [POST /api/webhooks/payment/:channel_id]
func (h *PaymentHandler) Webhook(c *gin.Context) {
	channelID, err := strconv.ParseUint(c.Param("channel_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的渠道ID")
		return
	}

	var req service.PaymentCallbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.svc.HandleCallback(uint(channelID), req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "回调处理成功")
}

// ListChannels 获取可用支付渠道 [GET /api/payment-channels]
func (h *PaymentHandler) ListChannels(c *gin.Context) {
	channels, err := h.svc.ListChannels()
	if err != nil {
		response.ServerError(c, "获取支付渠道失败")
		return
	}
	response.OK(c, channels)
}

// MockPaySuccess 模拟支付成功（仅测试环境）[POST /api/mock-pay/success]
func (h *PaymentHandler) MockPaySuccess(c *gin.Context) {
	var req struct {
		PaymentID uint `json:"payment_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 查询支付记录
	payment, err := h.svc.GetPaymentStatus(req.PaymentID)
	if err != nil {
		response.NotFound(c, "支付记录不存在")
		return
	}

	// 模拟回调
	callbackReq := service.PaymentCallbackReq{
		TradeNo:   payment.TradeNo,
		OrderNo:   "", // 从 payment 关联查询
		PaymentID: payment.ID,
		Amount:    payment.Amount,
		Status:    "success",
		PaidAt:    payment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Sign:      "mock_sign", // 模拟签名
	}

	// 需要从 payment 获取 order_no，这里简化处理
	// 实际应该通过 payment.OrderID 查询订单
	if err := h.svc.HandleCallback(1, callbackReq); err != nil {
		response.BadRequest(c, "模拟支付失败: "+err.Error())
		return
	}

	response.OKMsg(c, "模拟支付成功")
}
