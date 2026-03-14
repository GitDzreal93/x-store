package payment

// Provider 统一支付提供商接口
type Provider interface {
	// CreatePayment 创建支付
	CreatePayment(input CreateInput) (*CreateResult, error)
	// VerifyNotify 验证回调通知
	VerifyNotify(body []byte, headers map[string]string) (*NotifyResult, error)
}

// CreateInput 统一创建支付输入
type CreateInput struct {
	OrderNo   string  // 订单号
	PaymentID uint    // 支付记录ID
	Amount    float64 // 金额（人民币）
	Subject   string  // 商品名称
	NotifyURL string  // 异步回调地址
	ReturnURL string  // 同步跳转地址
}

// CreateResult 统一创建支付结果
type CreateResult struct {
	PayURL      string                 // 支付链接（跳转）
	QRCode      string                 // 二维码内容
	TradeNo     string                 // 第三方流水号
	RawResponse map[string]interface{} // 原始响应
}

// NotifyResult 统一回调通知结果
type NotifyResult struct {
	TradeNo string  // 第三方流水号
	OrderNo string  // 我方订单号
	Amount  float64 // 实际支付金额
	Status  string  // 状态: success | failed
}
