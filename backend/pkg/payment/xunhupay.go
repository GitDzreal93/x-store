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

// XunhupayProvider 虎皮椒支付提供商
type XunhupayProvider struct {
	config XunhupayConfig
	payway string // xunhupay_wechat | xunhupay_alipay
}

// XunhupayConfig 虎皮椒支付配置
type XunhupayConfig struct {
	APPID      string `json:"appid"`       // 虎皮椒 AppID
	APPSecret  string `json:"appsecret"`   // 虎皮椒 AppSecret
	WechatAPIUrl  string `json:"wechat_api_url"`  // 微信支付 API (如 https://api.xunhupay.com/payment/do.html)
	AlipayAPIUrl  string `json:"alipay_api_url"`  // 支付宝 API (如 https://api.xunhupay.com/payment/do.html)
	NotifyURL  string `json:"notify_url"`
	ReturnURL  string `json:"return_url"`
}

func NewXunhupayProvider(configJSON string, payway string) (*XunhupayProvider, error) {
	var cfg XunhupayConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("解析虎皮椒配置失败: %w", err)
	}
	return &XunhupayProvider{config: cfg, payway: payway}, nil
}

func (p *XunhupayProvider) CreatePayment(input CreateInput) (*CreateResult, error) {
	notifyURL := input.NotifyURL
	if notifyURL == "" {
		notifyURL = p.config.NotifyURL
	}
	returnURL := input.ReturnURL
	if returnURL == "" {
		returnURL = p.config.ReturnURL
	}

	params := map[string]string{
		"version":       "1.1",
		"appid":         p.config.APPID,
		"trade_order_id": input.OrderNo,
		"total_fee":     fmt.Sprintf("%.2f", input.Amount),
		"title":         input.Subject,
		"time":          fmt.Sprintf("%d", time.Now().Unix()),
		"notify_url":    notifyURL,
		"return_url":    returnURL,
		"nonce_str":     fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))),
	}

	// 微信或支付宝
	if p.payway == "xunhupay_wechat" {
		params["type"] = "WAP"
	}

	sign := p.signParams(params)
	params["hash"] = sign

	// 选择 API 地址
	apiURL := p.config.WechatAPIUrl
	if p.payway == "xunhupay_alipay" {
		apiURL = p.config.AlipayAPIUrl
	}
	if apiURL == "" {
		apiURL = "https://api.xunhupay.com/payment/do.html"
	}

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	resp, err := http.PostForm(apiURL, values)
	if err != nil {
		return nil, fmt.Errorf("请求虎皮椒失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		OpenID  int    `json:"openid"`
		URLQRCode string `json:"url_qrcode"`
		URL     string `json:"url"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		Hash    string `json:"hash"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析虎皮椒响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("虎皮椒创建失败: %s", result.ErrMsg)
	}

	payURL := result.URL
	qrCode := result.URLQRCode

	return &CreateResult{
		PayURL:  payURL,
		QRCode:  qrCode,
		TradeNo: fmt.Sprintf("XHP%d", time.Now().UnixNano()),
		RawResponse: map[string]interface{}{
			"url":        result.URL,
			"url_qrcode": result.URLQRCode,
		},
	}, nil
}

func (p *XunhupayProvider) VerifyNotify(body []byte, headers map[string]string) (*NotifyResult, error) {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("解析回调数据失败: %w", err)
	}

	status := values.Get("status")
	if status != "OD" { // OD = 已支付
		return nil, fmt.Errorf("支付未完成: status=%s", status)
	}

	// 验证签名
	params := map[string]string{}
	for k, v := range values {
		if k != "hash" && len(v) > 0 {
			params[k] = v[0]
		}
	}
	expectedSign := p.signParams(params)
	if expectedSign != values.Get("hash") {
		return nil, fmt.Errorf("签名验证失败")
	}

	return &NotifyResult{
		TradeNo: values.Get("transaction_id"),
		OrderNo: values.Get("trade_order_id"),
		Status:  "success",
	}, nil
}

func (p *XunhupayProvider) signParams(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "hash" || params[k] == "" {
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
	buf.WriteString(p.config.APPSecret)

	hash := md5.Sum([]byte(buf.String()))
	return hex.EncodeToString(hash[:])
}

var _ Provider = (*XunhupayProvider)(nil)
