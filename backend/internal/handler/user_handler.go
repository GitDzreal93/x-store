package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/x-store/backend/internal/service"
	"github.com/x-store/backend/pkg/response"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{svc: service.NewUserService()}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	resp, err := h.svc.Register(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, resp)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	resp, err := h.svc.Login(&req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.OK(c, resp)
}

// GetProfile 获取用户信息
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.svc.GetByID(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.OK(c, user)
}

// GetOrders 获取用户订单列表
func (h *UserHandler) GetOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "20")

	orders, total, err := h.svc.GetUserOrders(userID, page, size)
	if err != nil {
		response.ServerError(c, "获取订单失败")
		return
	}
	response.OKPage(c, orders, total, 1, 20)
}
