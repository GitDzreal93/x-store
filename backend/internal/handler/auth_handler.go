package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/x-store/backend/internal/config"
	"github.com/x-store/backend/internal/model"
	"github.com/x-store/backend/pkg/crypto"
	"github.com/x-store/backend/pkg/response"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{db: config.DB}
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

// Login 管理员登录 [POST /api/admin/login]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var user model.User
	if err := h.db.Where("username = ? AND role = ?", req.Username, "admin").First(&user).Error; err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	if user.Status != 1 {
		response.Forbidden(c, "账号已禁用")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	token, err := crypto.GenerateToken(
		config.Global.JWT.Secret,
		config.Global.JWT.ExpireHours,
		user.ID,
		user.Username,
		user.Role,
	)
	if err != nil {
		response.ServerError(c, "生成令牌失败")
		return
	}

	response.OK(c, LoginResp{Token: token, User: user})
}

// GetProfile 获取当前登录用户信息 [GET /api/admin/profile]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var user model.User
	if err := h.db.First(&user, userID.(uint)).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.OK(c, user)
}
