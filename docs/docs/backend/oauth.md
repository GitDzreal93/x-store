---
sidebar_position: 7
title: OAuth 登录
---

# OAuth 第三方登录

X-Store 支持 GitHub 和 Google OAuth 2.0 登录，配置存储在数据库中，可通过管理后台动态开关。

## 授权流程

```
用户点击「使用 GitHub 登录」
  → 前端跳转: GET /api/oauth/github
    → 后端重定向到 GitHub 授权页
      → 用户授权
        → GitHub 回调: GET /api/oauth/github/callback?code=xxx
          → 后端用 code 换 access_token
            → 获取 GitHub 用户信息（ID/邮箱/昵称/头像）
              → 查找或创建用户 + 生成 JWT
                → 重定向前端: /auth/callback?token=xxx
                  → 前端保存 token，跳转个人中心
```

Google 流程完全一致，只是 OAuth 端点不同。

## 用户绑定策略

`OAuthService.findOrCreateOAuthUser()` 按以下优先级匹配用户：

1. **OAuth ID 匹配** → 直接登录（之前已绑定）
2. **邮箱匹配** → 自动将 OAuth 绑定到已有账号
3. **全新用户** → 自动创建账号（用户名格式 `github_12345`，无密码）

```go
func (s *OAuthService) findOrCreateOAuthUser(info *OAuthUserInfo) (*AuthResp, error) {
    // 1. 按 OAuth 信息查找
    user, err := s.userRepo.GetByOAuth(info.Provider, info.ID)
    if err == nil && user != nil {
        return s.generateAuthResp(user)
    }

    // 2. 按邮箱查找（绑定已有账号）
    user, err = s.userRepo.GetByEmail(info.Email)
    if err == nil && user != nil {
        user.OAuthProvider = info.Provider
        user.OAuthID = info.ID
        s.userRepo.Update(user)
        return s.generateAuthResp(user)
    }

    // 3. 创建新用户
    user = &model.User{
        Username:      fmt.Sprintf("%s_%s", info.Provider, info.ID),
        Email:         info.Email,
        OAuthProvider: info.Provider,
        OAuthID:       info.ID,
        Role:          "buyer",
        Status:        1,
    }
    s.userRepo.Create(user)
    return s.generateAuthResp(user)
}
```

## 数据库配置

OAuth 配置存储在 `oauth_providers` 表中：

| 字段 | 说明 |
|------|------|
| `provider` | 提供商标识：`github` / `google` |
| `enabled` | 是否启用 |
| `client_id` | OAuth 应用 ID |
| `client_secret` | OAuth 应用密钥 |
| `redirect_url` | 授权回调地址 |

### 管理后台操作

在管理后台「第三方登录」页面可以：
- 通过 Switch 开关启用/禁用
- 编辑 Client ID 和 Client Secret
- 修改回调地址

### API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/oauth-providers` | 获取所有配置（secret 脱敏） |
| GET | `/api/admin/oauth-providers/:id` | 获取完整配置 |
| PUT | `/api/admin/oauth-providers/:id` | 更新配置 |
| POST | `/api/admin/oauth-providers/:id/toggle` | 切换启用状态 |

## 配置 GitHub OAuth

1. 前往 [github.com/settings/developers](https://github.com/settings/developers)
2. 点击 **New OAuth App**
3. 填写：
   - **Application name**: `X-Store`
   - **Homepage URL**: `http://localhost:3000`
   - **Authorization callback URL**: `http://localhost:8082/api/oauth/github/callback`
4. 创建后获取 **Client ID** 和 **Client Secret**
5. 在管理后台填入并启用

## 配置 Google OAuth

1. 前往 [console.cloud.google.com](https://console.cloud.google.com)
2. 创建项目 → APIs & Services → Credentials → **Create OAuth Client ID**
3. 选择 **Web application**
4. 添加授权重定向 URI: `http://localhost:8082/api/oauth/google/callback`
5. 获取 **Client ID** 和 **Client Secret**
6. 在管理后台填入并启用

## 前端集成

C 端登录页 (`/auth`) 通过 `/api/oauth/providers` 接口获取已启用的提供商，动态渲染登录按钮：

```tsx
useEffect(() => {
  fetch(`${API_BASE}/api/oauth/providers`)
    .then(res => res.json())
    .then(data => {
      if (data.code === 0 && data.data) {
        setOauthProviders(data.data);
      }
    });
}, []);
```

如果没有启用任何 OAuth 提供商，按钮和分隔线不会显示。
