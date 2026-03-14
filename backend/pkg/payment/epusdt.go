package payment

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// EpusdtProvider EPUSDT (USDT-TRC20) 加密货币支付
type EpusdtProvider struct {
	config EpusdtConfig
}

// EpusdtConfig EPUSDT 配置
type EpusdtConfig struct {
	APIUrl    string `json:"api_url"`    // EPUSDT API 地址
	Token     string `json:"token"`      // 认证 Token
	NotifyURL string `json:"notify_url"` // 回调地址
	ReturnURL string `json:"return_url"` // 跳转地址
}

func NewEpusdtProvider(configJSON string) (*EpusdtProvider, error) {
	var cfg EpusdtConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("解析 EPUSDT 配置失败: %w", err)
	}
	return &EpusdtProvider{config: cfg}, nil
}

// CreatePayment 创建 USDT 支付
func (p *EpusdtProvider) CreatePayment(input CreateInput) (*CreateResult, error) {
	notifyURL := input.NotifyURL
	if notifyURL == "" {
		notifyURL = p.config.NotifyURL
	}
	returnURL := input.ReturnURL
	if returnURL == "" {
		returnURL = p.config.ReturnURL
	}

	params := map[string]string{
		"order_id":   input.OrderNo,
		"amount":     fmt.Sprintf("%.2f", input.Amount),
		"notify_url": notifyURL,
		"redirect_url": returnURL,
	}

	// 签名
	sign := p.signParams(params)
	params["signature"] = sign

	// 调用 EPUSDT 创建订单 API
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	resp, err := http.PostForm(p.config.APIUrl+"/api/v1/order/create-transaction", values)
	if err != nil {
		return nil, fmt.Errorf("请求 EPUSDT 失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		Data       struct {
			TradeID        string  `json:"trade_id"`
			OrderID        string  `json:"order_id"`
			Amount         float64 `json:"amount"`
			ActualAmount   string  `json:"actual_amount"`
			Token          string  `json:"token"`
			ExpirationTime int64   `json:"expiration_time"`
			PaymentURL     string  `json:"payment_url"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 EPUSDT 响应失败: %w", err)
	}

	if result.StatusCode != 200 {
		return nil, fmt.Errorf("EPUSDT 创建订单失败: %s", result.Message)
	}

	return &CreateResult{
		PayURL:  result.Data.PaymentURL,
		QRCode:  result.Data.PaymentURL,
		TradeNo: result.Data.TradeID,
		RawResponse: map[string]interface{}{
			"trade_id":      result.Data.TradeID,
			"actual_amount": result.Data.ActualAmount,
			"token":         result.Data.Token,
		},
	}, nil
}

// VerifyNotify 验证 EPUSDT 回调
func (p *EpusdtProvider) VerifyNotify(body []byte, headers map[string]string) (*NotifyResult, error) {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("解析回调数据失败: %w", err)
	}

	status := values.Get("status")
	if status != "2" { // 2 = 支付成功
		return nil, fmt.Errorf("支付未成功: status=%s", status)
	}

	return &NotifyResult{
		TradeNo: values.Get("trade_id"),
		OrderNo: values.Get("order_id"),
		Status:  "success",
	}, nil
}

// signParams EPUSDT MD5 签名
func (p *EpusdtProvider) signParams(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "signature" || params[k] == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	for _, k := range keys {
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(params[k])
		buf.WriteByte('&')
	}
	buf.WriteString(p.config.Token)

	hash := md5.Sum([]byte(buf.String()))
	return hex.EncodeToString(hash[:])
}

var _ Provider = (*EpusdtProvider)(nil)

// 为 MockPay 添加 Provider 接口包装
type MockPayProviderWrapper struct {
	secret string
}

func NewMockPayProviderWrapper(configJSON string) (*MockPayProviderWrapper, error) {
	var cfg struct {
		SecretKey string `json:"secret_key"`
	}
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, err
	}
	return &MockPayProviderWrapper{secret: cfg.SecretKey}, nil
}

func (p *MockPayProviderWrapper) CreatePayment(input CreateInput) (*CreateResult, error) {
	tradeNo := fmt.Sprintf("MOCK%d%d", time.Now().Unix(), input.PaymentID)
	payURL := fmt.Sprintf("/mock-pay?trade_no=%s&order_no=%s&amount=%.2f&payment_id=%d",
		tradeNo, input.OrderNo, input.Amount, input.PaymentID)

	return &CreateResult{
		PayURL:  payURL,
		QRCode:  payURL,
		TradeNo: tradeNo,
		RawResponse: map[string]interface{}{
			"trade_no":   tradeNo,
			"order_no":   input.OrderNo,
			"payment_id": input.PaymentID,
			"amount":     input.Amount,
		},
	}, nil
}

func (p *MockPayProviderWrapper) VerifyNotify(body []byte, headers map[string]string) (*NotifyResult, error) {
	var data MockPayNotifyData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	if !VerifyMockPaySign(p.secret, &data) {
		return nil, fmt.Errorf("签名验证失败")
	}
	return &NotifyResult{
		TradeNo: data.TradeNo,
		OrderNo: data.OrderNo,
		Amount:  data.Amount,
		Status:  data.Status,
	}, nil
}

var _ Provider = (*MockPayProviderWrapper)(nil)
