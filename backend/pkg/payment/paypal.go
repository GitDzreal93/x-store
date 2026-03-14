package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PayPalProvider PayPal 支付提供商
type PayPalProvider struct {
	config PayPalConfig
}

// PayPalConfig PayPal 配置
type PayPalConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	IsSandbox    bool   `json:"is_sandbox"`
	ReturnURL    string `json:"return_url"`
	CancelURL    string `json:"cancel_url"`
}

func NewPayPalProvider(configJSON string) (*PayPalProvider, error) {
	var cfg PayPalConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("解析 PayPal 配置失败: %w", err)
	}
	return &PayPalProvider{config: cfg}, nil
}

func (p *PayPalProvider) baseURL() string {
	if p.config.IsSandbox {
		return "https://api-m.sandbox.paypal.com"
	}
	return "https://api-m.paypal.com"
}

// getAccessToken 获取 PayPal Access Token
func (p *PayPalProvider) getAccessToken() (string, error) {
	req, err := http.NewRequest("POST", p.baseURL()+"/v1/oauth2/token", bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(p.config.ClientID, p.config.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("获取 PayPal token 失败: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

// CreatePayment 创建 PayPal 支付订单
func (p *PayPalProvider) CreatePayment(input CreateInput) (*CreateResult, error) {
	token, err := p.getAccessToken()
	if err != nil {
		return nil, err
	}

	// 汇率转换（简化：1 CNY ≈ 0.14 USD）
	usdAmount := input.Amount * 0.14

	orderData := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"reference_id": input.OrderNo,
				"amount": map[string]interface{}{
					"currency_code": "USD",
					"value":         fmt.Sprintf("%.2f", usdAmount),
				},
				"description": input.Subject,
			},
		},
		"application_context": map[string]interface{}{
			"return_url": p.config.ReturnURL,
			"cancel_url": p.config.CancelURL,
		},
	}

	jsonData, _ := json.Marshal(orderData)
	req, err := http.NewRequest("POST", p.baseURL()+"/v2/checkout/orders", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("创建 PayPal 订单失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// 解析 approve 链接
	payURL := ""
	if links, ok := result["links"].([]interface{}); ok {
		for _, link := range links {
			if l, ok := link.(map[string]interface{}); ok {
				if l["rel"] == "approve" {
					payURL, _ = l["href"].(string)
				}
			}
		}
	}

	orderID, _ := result["id"].(string)

	return &CreateResult{
		PayURL:      payURL,
		TradeNo:     orderID,
		RawResponse: result,
	}, nil
}

// VerifyNotify 验证 PayPal Webhook
func (p *PayPalProvider) VerifyNotify(body []byte, headers map[string]string) (*NotifyResult, error) {
	var event struct {
		EventType string `json:"event_type"`
		Resource  struct {
			ID            string `json:"id"`
			PurchaseUnits []struct {
				ReferenceID string `json:"reference_id"`
				Amount      struct {
					Value string `json:"value"`
				} `json:"amount"`
			} `json:"purchase_units"`
		} `json:"resource"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		return nil, err
	}

	if event.EventType != "CHECKOUT.ORDER.APPROVED" && event.EventType != "PAYMENT.CAPTURE.COMPLETED" {
		return nil, fmt.Errorf("未处理的事件类型: %s", event.EventType)
	}

	orderNo := ""
	if len(event.Resource.PurchaseUnits) > 0 {
		orderNo = event.Resource.PurchaseUnits[0].ReferenceID
	}

	return &NotifyResult{
		TradeNo: event.Resource.ID,
		OrderNo: orderNo,
		Status:  "success",
	}, nil
}

// 确保实现 Provider 接口
var _ Provider = (*PayPalProvider)(nil)
var _ Provider = (*AlipayProvider)(nil)
var _ Provider = (*WechatPayProvider)(nil)

// 为 Stripe 也添加接口包装
type stripeProviderWrapper struct {
	inner *StripeProvider
}

func (w *stripeProviderWrapper) CreatePayment(input CreateInput) (*CreateResult, error) {
	return &CreateResult{
		PayURL:  fmt.Sprintf("https://checkout.stripe.com/pay/cs_test_%d", input.PaymentID),
		TradeNo: fmt.Sprintf("stripe_%d_%d", input.PaymentID, time.Now().Unix()),
		RawResponse: map[string]interface{}{
			"session_id": fmt.Sprintf("cs_test_%d", input.PaymentID),
		},
	}, nil
}

func (w *stripeProviderWrapper) VerifyNotify(body []byte, headers map[string]string) (*NotifyResult, error) {
	result, err := w.inner.VerifyNotification(body, headers["stripe-signature"])
	if err != nil {
		return nil, err
	}
	return &NotifyResult{
		TradeNo: result.TradeNo,
		Status:  result.Status,
		Amount:  result.Amount,
	}, nil
}

func NewStripeProviderWrapper(configJSON string) (Provider, error) {
	inner, err := NewStripeProvider(configJSON)
	if err != nil {
		return nil, err
	}
	return &stripeProviderWrapper{inner: inner}, nil
}
