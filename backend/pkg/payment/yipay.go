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

// YipayProvider 易支付（彩虹易支付）提供商
type YipayProvider struct {
	config YipayConfig
	payway string // yipay_alipay | yipay_wechat
}

// YipayConfig 易支付配置
type YipayConfig struct {
	APIUrl    string `json:"api_url"`    // 易支付 API 地址 (如 https://pay.example.com)
	PID       string `json:"pid"`        // 商户 ID
	Key       string `json:"key"`        // 商户密钥
	NotifyURL string `json:"notify_url"` // 异步回调
	ReturnURL string `json:"return_url"` // 同步跳转
}

func NewYipayProvider(configJSON string, payway string) (*YipayProvider, error) {
	var cfg YipayConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("解析易支付配置失败: %w", err)
	}
	return &YipayProvider{config: cfg, payway: payway}, nil
}

// CreatePayment 创建易支付订单
func (p *YipayProvider) CreatePayment(input CreateInput) (*CreateResult, error) {
	notifyURL := input.NotifyURL
	if notifyURL == "" {
		notifyURL = p.config.NotifyURL
	}
	returnURL := input.ReturnURL
	if returnURL == "" {
		returnURL = p.config.ReturnURL
	}

	// 支付类型映射
	payType := "alipay"
	if p.payway == "yipay_wechat" {
		payType = "wxpay"
	}

	params := map[string]string{
		"pid":        p.config.PID,
		"type":       payType,
		"out_trade_no": input.OrderNo,
		"notify_url": notifyURL,
		"return_url": returnURL,
		"name":       input.Subject,
		"money":      fmt.Sprintf("%.2f", input.Amount),
	}

	// MD5 签名
	sign := p.signParams(params)
	params["sign"] = sign
	params["sign_type"] = "MD5"

	// 构建支付跳转 URL
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	submitURL := p.config.APIUrl + "/submit.php"
	payURL := submitURL + "?" + values.Encode()

	return &CreateResult{
		PayURL:  payURL,
		TradeNo: fmt.Sprintf("YP%d", time.Now().UnixNano()),
		RawResponse: map[string]interface{}{
			"pay_url":  payURL,
			"pay_type": payType,
		},
	}, nil
}

// VerifyNotify 验证易支付回调
func (p *YipayProvider) VerifyNotify(body []byte, headers map[string]string) (*NotifyResult, error) {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("解析回调数据失败: %w", err)
	}

	tradeStatus := values.Get("trade_status")
	if tradeStatus != "TRADE_SUCCESS" {
		return nil, fmt.Errorf("交易未成功: %s", tradeStatus)
	}

	// 验证签名
	params := map[string]string{}
	for k, v := range values {
		if k != "sign" && k != "sign_type" && len(v) > 0 {
			params[k] = v[0]
		}
	}
	expectedSign := p.signParams(params)
	if expectedSign != values.Get("sign") {
		return nil, fmt.Errorf("签名验证失败")
	}

	return &NotifyResult{
		TradeNo: values.Get("trade_no"),
		OrderNo: values.Get("out_trade_no"),
		Status:  "success",
	}, nil
}

// signParams 易支付 MD5 签名
func (p *YipayProvider) signParams(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" || k == "sign_type" || params[k] == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(params[k])
	}
	buf.WriteString(p.config.Key)

	hash := md5.Sum([]byte(buf.String()))
	return hex.EncodeToString(hash[:])
}

var _ Provider = (*YipayProvider)(nil)

// ================================================================
// PayjsProvider Payjs 微信个人支付
// ================================================================

type PayjsProvider struct {
	config PayjsConfig
}

type PayjsConfig struct {
	MchID     string `json:"mchid"`      // 商户号
	Key       string `json:"key"`        // 通信密钥
	NotifyURL string `json:"notify_url"` // 回调地址
}

func NewPayjsProvider(configJSON string) (*PayjsProvider, error) {
	var cfg PayjsConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("解析 Payjs 配置失败: %w", err)
	}
	return &PayjsProvider{config: cfg}, nil
}

func (p *PayjsProvider) CreatePayment(input CreateInput) (*CreateResult, error) {
	notifyURL := input.NotifyURL
	if notifyURL == "" {
		notifyURL = p.config.NotifyURL
	}

	params := map[string]string{
		"mchid":      p.config.MchID,
		"total_fee":  fmt.Sprintf("%d", int(input.Amount*100)), // 分
		"out_trade_no": input.OrderNo,
		"body":       input.Subject,
		"notify_url": notifyURL,
	}

	sign := p.signParams(params)
	params["sign"] = sign

	jsonData, _ := json.Marshal(params)
	resp, err := http.Post("https://payjs.cn/api/native", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("请求 Payjs 失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		ReturnCode int    `json:"return_code"`
		ReturnMsg  string `json:"return_msg"`
		PayjsOrderID string `json:"payjs_order_id"`
		QRCode     string `json:"qrcode"`
		CodeURL    string `json:"code_url"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if result.ReturnCode != 1 {
		return nil, fmt.Errorf("Payjs 创建失败: %s", result.ReturnMsg)
	}

	return &CreateResult{
		PayURL:  result.QRCode,
		QRCode:  result.CodeURL,
		TradeNo: result.PayjsOrderID,
		RawResponse: map[string]interface{}{
			"payjs_order_id": result.PayjsOrderID,
		},
	}, nil
}

func (p *PayjsProvider) VerifyNotify(body []byte, headers map[string]string) (*NotifyResult, error) {
	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	returnCode := data["return_code"]
	if returnCode != "1" {
		return nil, fmt.Errorf("支付未成功")
	}

	return &NotifyResult{
		TradeNo: data["payjs_order_id"],
		OrderNo: data["out_trade_no"],
		Status:  "success",
	}, nil
}

func (p *PayjsProvider) signParams(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" || params[k] == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(params[k])
	}
	buf.WriteString("&key=")
	buf.WriteString(p.config.Key)

	hash := md5.Sum([]byte(buf.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

var _ Provider = (*PayjsProvider)(nil)
