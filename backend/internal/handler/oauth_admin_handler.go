package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/x-store/backend/internal/repo"
	"github.com/x-store/backend/pkg/response"
)

type OAuthAdminHandler struct {
	repo *repo.OAuthProviderRepo
}

func NewOAuthAdminHandler() *OAuthAdminHandler {
	return &OAuthAdminHandler{repo: repo.NewOAuthProviderRepo()}
}

// List 获取所有 OAuth 提供商配置 [GET /api/admin/oauth-providers]
func (h *OAuthAdminHandler) List(c *gin.Context) {
	providers, err := h.repo.List()
	if err != nil {
		response.ServerError(c, "获取配置失败")
		return
	}

	// 脱敏处理：隐藏 client_secret
	for i := range providers {
		if providers[i].ClientSecret != "" {
			providers[i].ClientSecret = "******"
		}
	}

	response.OK(c, providers)
}

// Update 更新 OAuth 提供商配置 [PUT /api/admin/oauth-providers/:id]
func (h *OAuthAdminHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	var req struct {
		Enabled      *bool  `json:"enabled"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURL  string `json:"redirect_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	provider, err := h.repo.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "配置不存在")
		return
	}

	// 更新字段
	if req.Enabled != nil {
		provider.Enabled = *req.Enabled
	}
	if req.ClientID != "" {
		provider.ClientID = req.ClientID
	}
	if req.ClientSecret != "" && req.ClientSecret != "******" {
		provider.ClientSecret = req.ClientSecret
	}
	if req.RedirectURL != "" {
		provider.RedirectURL = req.RedirectURL
	}

	if err := h.repo.Update(provider); err != nil {
		response.ServerError(c, "更新失败")
		return
	}

	// 返回时脱敏
	provider.ClientSecret = "******"
	response.OK(c, provider)
}

// Toggle 快速切换启用状态 [POST /api/admin/oauth-providers/:id/toggle]
func (h *OAuthAdminHandler) Toggle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	provider, err := h.repo.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "配置不存在")
		return
	}

	provider.Enabled = !provider.Enabled

	if err := h.repo.Update(provider); err != nil {
		response.ServerError(c, "更新失败")
		return
	}

	provider.ClientSecret = "******"
	response.OK(c, provider)
}

// GetDetail 获取单个配置详情（含完整 client_secret）[GET /api/admin/oauth-providers/:id]
func (h *OAuthAdminHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	provider, err := h.repo.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "配置不存在")
		return
	}

	response.OK(c, provider)
}
