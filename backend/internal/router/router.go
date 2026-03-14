package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/x-store/backend/internal/handler"
	"github.com/x-store/backend/internal/middleware"
	"github.com/x-store/backend/internal/service"
)

func Setup(mode string) *gin.Engine {
	gin.SetMode(mode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS())

	api := r.Group("/api")

	// ==================== 公开接口 (C 端) ====================
	{
		// 用户注册登录
		userH := handler.NewUserHandler()
		api.POST("/users/register", userH.Register)
		api.POST("/users/login", userH.Login)

		// OAuth 第三方登录（GitHub / Google）
		oauthH := handler.NewOAuthHandler()
		api.GET("/oauth/providers", oauthH.ListProviders)
		api.GET("/oauth/github", oauthH.GitHubLogin)
		api.GET("/oauth/github/callback", oauthH.GitHubCallback)
		api.GET("/oauth/google", oauthH.GoogleLogin)
		api.GET("/oauth/google/callback", oauthH.GoogleCallback)

		catH := handler.NewCategoryHandler()
		api.GET("/categories", catH.List)
		api.GET("/categories/:id", catH.Get)

		prodH := handler.NewProductHandler()
		api.GET("/products", prodH.List)
		api.GET("/products/:id", prodH.Get)

		// 订单查询（公开，通过订单号查询）
		orderH := handler.NewOrderHandler()
		api.GET("/orders/:order_no", orderH.Get)

		// 支付渠道列表（公开）
		orderSvc := service.NewOrderService()
		paymentH := handler.NewPaymentHandler(orderSvc)
		api.GET("/payment-channels", paymentH.ListChannels)
		api.GET("/payments/:id/status", paymentH.GetStatus)

		// Webhook 回调（公开，供支付网关调用）
		api.POST("/webhooks/payment/:channel_id", paymentH.Webhook)
	}

	// ==================== 需要登录的接口 ====================
	user := api.Group("/user")
	user.Use(middleware.JWTAuth())
	{
		userH := handler.NewUserHandler()
		user.GET("/profile", userH.GetProfile)
		user.GET("/orders", userH.GetOrders)
	}

	// ==================== 管理后台接口 ====================
	adminPublic := api.Group("/admin")
	{
		authH := handler.NewAuthHandler()
		adminPublic.POST("/login", authH.Login)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(), middleware.AdminOnly())
	{
		authH := handler.NewAuthHandler()
		admin.GET("/profile", authH.GetProfile)
		
		// 统计数据
		statsH := handler.NewStatsHandler()
		admin.GET("/stats/dashboard", statsH.GetDashboardStats)
		
		catH := handler.NewCategoryHandler()
		admin.POST("/categories", catH.Create)
		admin.PUT("/categories/:id", catH.Update)
		admin.DELETE("/categories/:id", catH.Delete)

		prodH := handler.NewProductHandler()
		admin.POST("/products", prodH.Create)
		admin.PUT("/products/:id", prodH.Update)
		admin.DELETE("/products/:id", prodH.Delete)

		// 商品列表（管理端，含下架商品）
		admin.GET("/products", prodH.List)

		// 订单管理
		orderAdminH := handler.NewOrderHandler()
		admin.GET("/orders", orderAdminH.List)
		admin.POST("/orders/:order_no/refund", orderAdminH.Refund)
		admin.POST("/orders/:order_no/deliver", orderAdminH.ManualDeliver)

		// 卡密管理
		cardKeyH := handler.NewCardKeyHandler()
		admin.POST("/cardkeys/import", cardKeyH.BatchImport)
		admin.GET("/cardkeys/count/:product_id", cardKeyH.CountAvailable)

		// OAuth 提供商管理
		oauthAdminH := handler.NewOAuthAdminHandler()
		admin.GET("/oauth-providers", oauthAdminH.List)
		admin.GET("/oauth-providers/:id", oauthAdminH.GetDetail)
		admin.PUT("/oauth-providers/:id", oauthAdminH.Update)
		admin.POST("/oauth-providers/:id/toggle", oauthAdminH.Toggle)
	}

	// ==================== 需要防重放签名的写接口 ====================
	signed := api.Group("")
	signed.Use(middleware.AntiReplay(), middleware.RateLimit(30, time.Minute))
	{
		// 订单创建（核心写接口，需要签名防刷）
		orderH := handler.NewOrderHandler()
		signed.POST("/orders", orderH.Create)
		signed.POST("/orders/:order_no/pay", orderH.Pay)
		signed.POST("/orders/:order_no/cancel", orderH.Cancel)

		// 支付创建（核心写接口，需要签名防刷）
		orderSvc2 := service.NewOrderService()
		paymentH2 := handler.NewPaymentHandler(orderSvc2)
		signed.POST("/payments", paymentH2.Create)
	}

	// 健康检查
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
