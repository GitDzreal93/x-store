package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/x-store/backend/internal/service"
	"github.com/x-store/backend/pkg/response"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{svc: service.NewOrderService()}
}

// Create 创建订单 [POST /api/orders]
func (h *OrderHandler) Create(c *gin.Context) {
	var req service.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从 JWT 中获取用户 ID（如果已登录）
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		req.UserID = &uid
	}

	resp, err := h.svc.CreateOrder(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, resp)
}

// Get 查询订单详情 [GET /api/orders/:order_no]
func (h *OrderHandler) Get(c *gin.Context) {
	orderNo := c.Param("order_no")
	order, err := h.svc.GetByOrderNo(orderNo)
	if err != nil {
		response.NotFound(c, "订单不存在")
		return
	}
	response.OK(c, order)
}

// Pay 支付订单 [POST /api/orders/:order_no/pay]
func (h *OrderHandler) Pay(c *gin.Context) {
	orderNo := c.Param("order_no")
	
	var req struct {
		PayMethod string `json:"pay_method" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.svc.PayOrder(c.Request.Context(), orderNo, req.PayMethod); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "支付成功")
}

// List 管理后台订单列表 [GET /api/admin/orders]
func (h *OrderHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	keyword := c.Query("keyword")

	var status *int
	if s := c.Query("status"); s != "" {
		v, err := strconv.Atoi(s)
		if err == nil {
			status = &v
		}
	}

	list, total, err := h.svc.ListAll(page, size, status, keyword)
	if err != nil {
		response.ServerError(c, "获取订单列表失败")
		return
	}
	response.OKPage(c, list, total, page, size)
}

// Cancel 取消订单 [POST /api/orders/:order_no/cancel]
func (h *OrderHandler) Cancel(c *gin.Context) {
	orderNo := c.Param("order_no")
	if err := h.svc.CancelOrder(c.Request.Context(), orderNo); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "订单已取消")
}

// Refund 退款订单 [POST /api/admin/orders/:order_no/refund]
func (h *OrderHandler) Refund(c *gin.Context) {
	orderNo := c.Param("order_no")
	if err := h.svc.RefundOrder(c.Request.Context(), orderNo); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "退款成功")
}

// ManualDeliver 手动发货 [POST /api/admin/orders/:order_no/deliver]
func (h *OrderHandler) ManualDeliver(c *gin.Context) {
	orderNo := c.Param("order_no")
	order, err := h.svc.GetByOrderNo(orderNo)
	if err != nil {
		response.NotFound(c, "订单不存在")
		return
	}
	if err := h.svc.DeliverOrder(c.Request.Context(), order.ID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "发货成功")
}
