package payment

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// VmqpayProvider V免签支付提供商
type VmqpayProvider struct {
	config VmqpayConfig
	payway string // vmqpay_wechat | vmqpay_alipay
}

// VmqpayConfig V免签配置
type VmqpayConfig struct {
	APIUrl    string `json:"api_url"`    // V免签 API 地址
	Key       string `json:"key"`        // 通信密钥
	NotifyURL string `json:"notify_url"` // 回调地址
	ReturnURL string `json:"return_url"` // 跳转地址
}

func NewVmqpayProvider(configJSON string, payway string) (*VmqpayProvider, error) {
	var cfg VmqpayConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("解析V免签配置失败: %w", err)
	}
	return &VmqpayProvider{config: cfg, payway: payway}, nil
}

func (p *VmqpayProvider) CreatePayment(input CreateInput) (*CreateResult, error) {
	notifyURL := input.NotifyURL
	if notifyURL == "" {
		notifyURL = p.config.NotifyURL
	}
	returnURL := input.ReturnURL
	if returnURL == "" {
		returnURL = p.config.ReturnURL
	}

	// 支付类型: 1=微信, 2=支付宝
	payType := "1"
	if p.payway == "vmqpay_alipay" {
		payType = "2"
	}

	price := fmt.Sprintf("%.2f", input.Amount)

	// V免签签名: md5(payId + param + type + price + key)
	signStr := fmt.Sprintf("%s%s%s%s%s", input.OrderNo, "", payType, price, p.config.Key)
	hash := md5.Sum([]byte(signStr))
	sign := hex.EncodeToString(hash[:])

	params := url.Values{
		"payId":      {input.OrderNo},
		"type":       {payType},
		"price":      {price},
		"sign":       {sign},
		"notifyUrl":  {notifyURL},
		"returnUrl":  {returnURL},
		"isHtml":     {"0"},
	}

	apiURL := p.config.APIUrl + "/createOrder"
	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		return nil, fmt.Errorf("请求V免签失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			PayURL  string `json:"payUrl"`
			OrderID string `json:"orderId"`
			ReallyPrice string `json:"reallyPrice"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析V免签响应失败: %w", err)
	}

	if result.Code != 1 {
		return nil, fmt.Errorf("V免签创建失败: %s", result.Msg)
	}

	return &CreateResult{
		PayURL:  result.Data.PayURL,
		QRCode:  result.Data.PayURL,
		TradeNo: result.Data.OrderID,
		RawResponse: map[string]interface{}{
			"pay_url":      result.Data.PayURL,
			"order_id":     result.Data.OrderID,
			"really_price": result.Data.ReallyPrice,
		},
	}, nil
}

func (p *VmqpayProvider) VerifyNotify(body []byte, headers map[string]string) (*NotifyResult, error) {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("解析回调数据失败: %w", err)
	}

	// V免签回调验签: md5(payId + reallyPrice + key)
	payID := values.Get("payId")
	reallyPrice := values.Get("reallyPrice")
	signStr := fmt.Sprintf("%s%s%s", payID, reallyPrice, p.config.Key)
	hash := md5.Sum([]byte(signStr))
	expectedSign := hex.EncodeToString(hash[:])

	if expectedSign != values.Get("sign") {
		return nil, fmt.Errorf("签名验证失败")
	}

	return &NotifyResult{
		TradeNo: fmt.Sprintf("VMQ%d", time.Now().UnixNano()),
		OrderNo: payID,
		Status:  "success",
	}, nil
}

var _ Provider = (*VmqpayProvider)(nil)
