package payment

import (
	"encoding/json"
	"fmt"

	"github.com/x-store/backend/internal/model"
)

// PaymentResult 支付创建结果
type PaymentResult struct {
	PaymentID   string
	RedirectURL string
	QRCodeURL   string
	Extra       map[string]interface{}
}

// NotificationResult 支付通知结果
type NotificationResult struct {
	TradeNo string
	Status  string
	Amount  float64
}

type StripeProvider struct {
	config StripeConfig
}

type StripeConfig struct {
	SecretKey      string `json:"secret_key"`
	PublishableKey string `json:"publishable_key"`
	WebhookSecret  string `json:"webhook_secret"`
}

func NewStripeProvider(configJSON string) (*StripeProvider, error) {
	var cfg StripeConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, err
	}
	return &StripeProvider{config: cfg}, nil
}

// CreatePayment 创建 Stripe 支付
func (p *StripeProvider) CreatePayment(order *model.Order, channel *model.PaymentChannel) (*PaymentResult, error) {
	// 实际生产环境需要调用 Stripe API
	// 这里返回模拟数据供前端跳转
	return &PaymentResult{
		PaymentID:   fmt.Sprintf("stripe_%d", order.ID),
		RedirectURL: fmt.Sprintf("https://checkout.stripe.com/pay/cs_test_%d", order.ID),
		QRCodeURL:   "",
		Extra: map[string]interface{}{
			"session_id": fmt.Sprintf("cs_test_%d", order.ID),
		},
	}, nil
}

// VerifyNotification 验证 Stripe Webhook 签名
func (p *StripeProvider) VerifyNotification(rawBody []byte, signature string) (*NotificationResult, error) {
	// 实际生产环境需要验证 Stripe Webhook 签名
	// stripe.ConstructEvent(rawBody, signature, p.config.WebhookSecret)
	
	// 模拟解析
	var event map[string]interface{}
	if err := json.Unmarshal(rawBody, &event); err != nil {
		return nil, err
	}

	// 检查事件类型
	eventType, ok := event["type"].(string)
	if !ok || eventType != "payment_intent.succeeded" {
		return nil, fmt.Errorf("unsupported event type")
	}

	return &NotificationResult{
		TradeNo: fmt.Sprintf("stripe_%v", event["id"]),
		Status:  "success",
		Amount:  0, // 从 event 中解析
	}, nil
}
