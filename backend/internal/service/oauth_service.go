package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/x-store/backend/internal/config"
	"github.com/x-store/backend/internal/model"
	"github.com/x-store/backend/internal/repo"
	"github.com/x-store/backend/pkg/crypto"
)

type OAuthService struct {
	userRepo         *repo.UserRepo
	oauthProviderRepo *repo.OAuthProviderRepo
}

func NewOAuthService() *OAuthService {
	return &OAuthService{
		userRepo:         repo.NewUserRepo(),
		oauthProviderRepo: repo.NewOAuthProviderRepo(),
	}
}

// OAuthUserInfo 第三方平台返回的用户信息
type OAuthUserInfo struct {
	Provider string // github | google
	ID       string // 平台用户 ID
	Email    string
	Name     string
	Avatar   string
}

// ======================== GitHub ========================

// GetGitHubAuthURL 获取 GitHub 授权跳转地址
func (s *OAuthService) GetGitHubAuthURL(state string) string {
	provider, err := s.oauthProviderRepo.GetByProvider("github")
	if err != nil || !provider.Enabled {
		return ""
	}
	params := url.Values{
		"client_id":    {provider.ClientID},
		"redirect_uri": {provider.RedirectURL},
		"scope":        {"read:user user:email"},
		"state":        {state},
	}
	return "https://github.com/login/oauth/authorize?" + params.Encode()
}

// HandleGitHubCallback 处理 GitHub OAuth 回调
func (s *OAuthService) HandleGitHubCallback(code string) (*AuthResp, error) {
	provider, err := s.oauthProviderRepo.GetByProvider("github")
	if err != nil || !provider.Enabled {
		return nil, fmt.Errorf("GitHub 登录未启用")
	}

	// 1. 用 code 换 access_token
	tokenURL := "https://github.com/login/oauth/access_token"
	data := url.Values{
		"client_id":     {provider.ClientID},
		"client_secret": {provider.ClientSecret},
		"code":          {code},
		"redirect_uri":  {provider.RedirectURL},
	}

	req, _ := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GitHub token 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("解析 GitHub token 失败: %w", err)
	}
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("GitHub 授权失败: %s", tokenResp.ErrorDesc)
	}

	// 2. 用 access_token 获取用户信息
	userInfo, err := s.fetchGitHubUser(tokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	// 3. 查找或创建用户
	return s.findOrCreateOAuthUser(userInfo)
}

// fetchGitHubUser 获取 GitHub 用户信息
func (s *OAuthService) fetchGitHubUser(accessToken string) (*OAuthUserInfo, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取 GitHub 用户信息失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ghUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.Unmarshal(body, &ghUser); err != nil {
		return nil, fmt.Errorf("解析 GitHub 用户信息失败: %w", err)
	}

	// 如果主邮箱为空，尝试获取邮箱列表
	email := ghUser.Email
	if email == "" {
		email = s.fetchGitHubEmail(accessToken)
	}
	if email == "" {
		email = fmt.Sprintf("github_%d@x-store.local", ghUser.ID)
	}

	name := ghUser.Name
	if name == "" {
		name = ghUser.Login
	}

	return &OAuthUserInfo{
		Provider: "github",
		ID:       fmt.Sprintf("%d", ghUser.ID),
		Email:    email,
		Name:     name,
		Avatar:   ghUser.AvatarURL,
	}, nil
}

// fetchGitHubEmail 获取 GitHub 用户邮箱
func (s *OAuthService) fetchGitHubEmail(accessToken string) string {
	req, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.Unmarshal(body, &emails); err != nil {
		return ""
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email
		}
	}
	if len(emails) > 0 {
		return emails[0].Email
	}
	return ""
}

// ======================== Google ========================

// GetGoogleAuthURL 获取 Google 授权跳转地址
func (s *OAuthService) GetGoogleAuthURL(state string) string {
	provider, err := s.oauthProviderRepo.GetByProvider("google")
	if err != nil || !provider.Enabled {
		return ""
	}
	params := url.Values{
		"client_id":     {provider.ClientID},
		"redirect_uri":  {provider.RedirectURL},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"state":         {state},
		"access_type":   {"offline"},
	}
	return "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()
}

// HandleGoogleCallback 处理 Google OAuth 回调
func (s *OAuthService) HandleGoogleCallback(code string) (*AuthResp, error) {
	provider, err := s.oauthProviderRepo.GetByProvider("google")
	if err != nil || !provider.Enabled {
		return nil, fmt.Errorf("Google 登录未启用")
	}

	// 1. 用 code 换 access_token
	data := url.Values{
		"client_id":     {provider.ClientID},
		"client_secret": {provider.ClientSecret},
		"code":          {code},
		"redirect_uri":  {provider.RedirectURL},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return nil, fmt.Errorf("Google token 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("解析 Google token 失败: %w", err)
	}
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("Google 授权失败: %s", tokenResp.ErrorDesc)
	}

	// 2. 用 access_token 获取用户信息
	userInfo, err := s.fetchGoogleUser(tokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	// 3. 查找或创建用户
	return s.findOrCreateOAuthUser(userInfo)
}

// fetchGoogleUser 获取 Google 用户信息
func (s *OAuthService) fetchGoogleUser(accessToken string) (*OAuthUserInfo, error) {
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取 Google 用户信息失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var gUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.Unmarshal(body, &gUser); err != nil {
		return nil, fmt.Errorf("解析 Google 用户信息失败: %w", err)
	}

	return &OAuthUserInfo{
		Provider: "google",
		ID:       gUser.ID,
		Email:    gUser.Email,
		Name:     gUser.Name,
		Avatar:   gUser.Picture,
	}, nil
}

// ======================== 通用 ========================

// findOrCreateOAuthUser 查找或创建 OAuth 用户并返回 JWT
func (s *OAuthService) findOrCreateOAuthUser(info *OAuthUserInfo) (*AuthResp, error) {
	// 1. 先按 OAuth 信息查找
	user, err := s.userRepo.GetByOAuth(info.Provider, info.ID)
	if err == nil && user != nil {
		// 已存在，更新头像和昵称
		user.Avatar = info.Avatar
		user.Nickname = info.Name
		_ = s.userRepo.Update(user)
		return s.generateAuthResp(user)
	}

	// 2. 再按邮箱查找（已有账号绑定 OAuth）
	user, err = s.userRepo.GetByEmail(info.Email)
	if err == nil && user != nil {
		// 邮箱已存在，绑定 OAuth
		user.OAuthProvider = info.Provider
		user.OAuthID = info.ID
		user.Avatar = info.Avatar
		if user.Nickname == "" || user.Nickname == user.Username {
			user.Nickname = info.Name
		}
		_ = s.userRepo.Update(user)
		return s.generateAuthResp(user)
	}

	// 3. 全新用户，自动创建
	username := fmt.Sprintf("%s_%s", info.Provider, info.ID)
	user = &model.User{
		Username:      username,
		Email:         info.Email,
		Password:      "", // OAuth 用户无密码
		Nickname:      info.Name,
		Avatar:        info.Avatar,
		Role:          "buyer",
		Status:        1,
		OAuthProvider: info.Provider,
		OAuthID:       info.ID,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return s.generateAuthResp(user)
}

// generateAuthResp 生成认证响应
func (s *OAuthService) generateAuthResp(user *model.User) (*AuthResp, error) {
	token, err := crypto.GenerateToken(
		config.Global.JWT.Secret,
		config.Global.JWT.ExpireHours,
		user.ID, user.Username, user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("生成 Token 失败")
	}

	user.Password = ""
	return &AuthResp{Token: token, User: user}, nil
}

// GetEnabledProviders 返回已启用的 OAuth 提供商列表
func (s *OAuthService) GetEnabledProviders() []map[string]string {
	var result []map[string]string
	providers, err := s.oauthProviderRepo.ListEnabled()
	if err != nil {
		return result
	}
	
	for _, p := range providers {
		var authURL string
		if p.Provider == "github" {
			authURL = s.GetGitHubAuthURL("")
		} else if p.Provider == "google" {
			authURL = s.GetGoogleAuthURL("")
		}
		
		result = append(result, map[string]string{
			"name":     p.Provider,
			"label":    p.Name,
			"auth_url": authURL,
		})
	}
	return result
}
