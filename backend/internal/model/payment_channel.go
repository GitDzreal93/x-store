package model

// PaymentChannel 支付渠道配置表
type PaymentChannel struct {
	BaseModel
	Name            string  `json:"name" gorm:"type:varchar(64);not null"`                  // 渠道名称
	ProviderType    string  `json:"provider_type" gorm:"type:varchar(32);not null;index"`   // 提供商类型: mock | stripe | alipay | wechat
	ChannelType     string  `json:"channel_type" gorm:"type:varchar(32);not null"`          // 渠道类型: mockpay | stripe | alipay | wechat
	InteractionMode string  `json:"interaction_mode" gorm:"type:varchar(32)"`                // 交互模式: redirect | qrcode | h5
	ConfigJSON      string  `json:"config_json" gorm:"type:text"`                            // JSON 配置（API密钥等）
	FeeRate         float64 `json:"fee_rate" gorm:"type:decimal(5,2);default:0"`             // 手续费率 %
	FixedFee        float64 `json:"fixed_fee" gorm:"type:decimal(10,2);default:0"`           // 固定手续费
	IsActive        bool    `json:"is_active" gorm:"default:true;not null;index"`            // 是否启用
	Sort            int     `json:"sort" gorm:"default:0;not null"`                          // 排序权重
}

func (PaymentChannel) TableName() string { return "payment_channels" }

// 支付提供商类型常量
const (
	ProviderTypeMock     = "mock"     // 模拟支付（测试用）
	ProviderTypeStripe   = "stripe"   // Stripe（国际信用卡）
	ProviderTypeAlipay   = "alipay"   // 支付宝
	ProviderTypeWechat   = "wechat"   // 微信支付
	ProviderTypePayPal   = "paypal"   // PayPal（国际）
	ProviderTypeEpusdt   = "epusdt"   // EPUSDT（USDT 加密货币）
	ProviderTypeYipay    = "yipay"    // 易支付（彩虹易支付）
	ProviderTypePayjs    = "payjs"    // Payjs（微信个人支付）
	ProviderTypeXunhupay = "xunhupay" // 虎皮椒支付
	ProviderTypeVmqpay   = "vmqpay"   // V免签支付
)

// 渠道类型常量
const (
	ChannelTypeMockPay     = "mockpay"      // 模拟支付
	ChannelTypeStripe      = "stripe"       // Stripe 信用卡
	ChannelTypeAlipayF2F   = "alipay_f2f"   // 支付宝当面付（扫码）
	ChannelTypeAlipayPC    = "alipay_pc"    // 支付宝电脑网站支付
	ChannelTypeAlipayWap   = "alipay_wap"   // 支付宝手机网站支付
	ChannelTypeWechatNative = "wechat_native" // 微信 Native 扫码支付
	ChannelTypeWechatH5    = "wechat_h5"    // 微信 H5 支付
	ChannelTypeWechatJSAPI = "wechat_jsapi" // 微信 JSAPI 支付
	ChannelTypePayPal      = "paypal"       // PayPal
	ChannelTypeEpusdt      = "epusdt"       // USDT(TRC20)
	ChannelTypeYipayAlipay = "yipay_alipay" // 易支付-支付宝
	ChannelTypeYipayWechat = "yipay_wechat" // 易支付-微信
	ChannelTypePayjs       = "payjs"        // Payjs 微信
	ChannelTypeXunhupayWechat = "xunhupay_wechat" // 虎皮椒-微信
	ChannelTypeXunhupayAlipay = "xunhupay_alipay" // 虎皮椒-支付宝
	ChannelTypeVmqpayWechat   = "vmqpay_wechat"   // V免签-微信
	ChannelTypeVmqpayAlipay   = "vmqpay_alipay"   // V免签-支付宝
)

// 交互模式常量
const (
	InteractionModeRedirect = "redirect" // 跳转支付
	InteractionModeQRCode   = "qrcode"   // 扫码支付
	InteractionModeH5       = "h5"       // H5支付
	InteractionModeJSAPI    = "jsapi"    // JSAPI 支付（公众号内）
)
