package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// MockPayConfig 模拟支付配置
type MockPayConfig struct {
	Secret     string `json:"secret"`      // 签名密钥
	NotifyURL  string `json:"notify_url"`  // 回调地址
	ReturnURL  string `json:"return_url"`  // 返回地址
	AutoNotify bool   `json:"auto_notify"` // 是否自动触发回调（测试用）
}

// MockPayCreateInput 创建支付输入
type MockPayCreateInput struct {
	OrderNo   string  // 订单号
	PaymentID uint    // 支付ID
	Amount    float64 // 金额
	Subject   string  // 商品标题
}

// MockPayCreateResult 创建支付结果
type MockPayCreateResult struct {
	PayURL      string                 // 支付链接
	QRCode      string                 // 二维码内容
	TradeNo     string                 // 支付流水号
	RawResponse map[string]interface{} // 原始响应
}

// MockPayNotifyData 模拟支付回调数据
type MockPayNotifyData struct {
	TradeNo   string  `json:"trade_no"`   // 支付流水号
	OrderNo   string  `json:"order_no"`   // 订单号
	PaymentID uint    `json:"payment_id"` // 支付ID
	Amount    float64 `json:"amount"`     // 金额
	Status    string  `json:"status"`     // 状态: success | failed
	PaidAt    string  `json:"paid_at"`    // 支付时间
	Sign      string  `json:"sign"`       // 签名
}

// CreateMockPayment 创建模拟支付
func CreateMockPayment(cfg MockPayConfig, input MockPayCreateInput) (*MockPayCreateResult, error) {
	if cfg.Secret == "" {
		return nil, fmt.Errorf("secret is required")
	}
	if input.OrderNo == "" || input.PaymentID == 0 || input.Amount <= 0 {
		return nil, fmt.Errorf("invalid input parameters")
	}

	// 生成模拟支付流水号
	tradeNo := fmt.Sprintf("MOCK%d%d", time.Now().Unix(), input.PaymentID)

	// 生成支付链接（实际是一个测试页面）
	payURL := fmt.Sprintf("/mock-pay?trade_no=%s&order_no=%s&amount=%.2f&payment_id=%d",
		tradeNo, input.OrderNo, input.Amount, input.PaymentID)

	// 生成二维码内容（同样是支付链接）
	qrCode := payURL

	return &MockPayCreateResult{
		PayURL:  payURL,
		QRCode:  qrCode,
		TradeNo: tradeNo,
		RawResponse: map[string]interface{}{
			"trade_no":   tradeNo,
			"order_no":   input.OrderNo,
			"payment_id": input.PaymentID,
			"amount":     input.Amount,
			"subject":    input.Subject,
			"created_at": time.Now().Format(time.RFC3339),
		},
	}, nil
}

// GenerateMockPayNotify 生成模拟支付回调数据
func GenerateMockPayNotify(cfg MockPayConfig, tradeNo, orderNo string, paymentID uint, amount float64, status string) *MockPayNotifyData {
	paidAt := time.Now().Format(time.RFC3339)
	
	// 生成签名
	sign := generateMockPaySign(cfg.Secret, tradeNo, orderNo, amount, status, paidAt)

	return &MockPayNotifyData{
		TradeNo:   tradeNo,
		OrderNo:   orderNo,
		PaymentID: paymentID,
		Amount:    amount,
		Status:    status,
		PaidAt:    paidAt,
		Sign:      sign,
	}
}

// VerifyMockPaySign 验证模拟支付签名
func VerifyMockPaySign(secret string, data *MockPayNotifyData) bool {
	expectedSign := generateMockPaySign(secret, data.TradeNo, data.OrderNo, data.Amount, data.Status, data.PaidAt)
	return hmac.Equal([]byte(expectedSign), []byte(data.Sign))
}

// generateMockPaySign 生成签名
func generateMockPaySign(secret, tradeNo, orderNo string, amount float64, status, paidAt string) string {
	// 拼接签名字符串
	signStr := fmt.Sprintf("amount=%.2f&order_no=%s&paid_at=%s&status=%s&trade_no=%s&secret=%s",
		amount, orderNo, paidAt, status, tradeNo, secret)
	
	// HMAC-SHA256
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signStr))
	return hex.EncodeToString(h.Sum(nil))
}

// ParseMockPayStatus 解析模拟支付状态到系统状态
func ParseMockPayStatus(status string) string {
	switch status {
	case "success":
		return "success"
	case "failed":
		return "failed"
	default:
		return "pending"
	}
}
