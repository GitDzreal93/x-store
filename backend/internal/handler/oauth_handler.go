package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/x-store/backend/internal/service"
	"github.com/x-store/backend/pkg/response"
)

type OAuthHandler struct {
	svc *service.OAuthService
}

func NewOAuthHandler() *OAuthHandler {
	return &OAuthHandler{svc: service.NewOAuthService()}
}

// ListProviders 返回已启用的 OAuth 提供商列表 [GET /api/oauth/providers]
func (h *OAuthHandler) ListProviders(c *gin.Context) {
	providers := h.svc.GetEnabledProviders()
	response.OK(c, providers)
}

// GitHubLogin 跳转到 GitHub 授权页 [GET /api/oauth/github]
func (h *OAuthHandler) GitHubLogin(c *gin.Context) {
	state := c.Query("state")
	authURL := h.svc.GetGitHubAuthURL(state)
	if authURL == "" {
		response.BadRequest(c, "GitHub 登录未启用")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GitHubCallback 处理 GitHub 授权回调 [GET /api/oauth/github/callback]
func (h *OAuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		h.redirectWithError(c, "授权码为空")
		return
	}

	authResp, err := h.svc.HandleGitHubCallback(code)
	if err != nil {
		h.redirectWithError(c, err.Error())
		return
	}

	h.redirectWithToken(c, authResp.Token)
}

// GoogleLogin 跳转到 Google 授权页 [GET /api/oauth/google]
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	state := c.Query("state")
	authURL := h.svc.GetGoogleAuthURL(state)
	if authURL == "" {
		response.BadRequest(c, "Google 登录未启用")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleCallback 处理 Google 授权回调 [GET /api/oauth/google/callback]
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		h.redirectWithError(c, "授权码为空")
		return
	}

	authResp, err := h.svc.HandleGoogleCallback(code)
	if err != nil {
		h.redirectWithError(c, err.Error())
		return
	}

	h.redirectWithToken(c, authResp.Token)
}

// redirectWithToken 带 token 跳转回前端
func (h *OAuthHandler) redirectWithToken(c *gin.Context, token string) {
	frontendURL := "http://localhost:3000"
	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", frontendURL, token)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// redirectWithError 带错误信息跳转回前端
func (h *OAuthHandler) redirectWithError(c *gin.Context, errMsg string) {
	frontendURL := "http://localhost:3000"
	redirectURL := fmt.Sprintf("%s/auth/callback?error=%s", frontendURL, errMsg)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
