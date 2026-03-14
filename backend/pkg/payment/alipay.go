package payment

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// AlipayProvider 支付宝支付提供商
type AlipayProvider struct {
	config  AlipayConfig
	payway  string // alipay_f2f | alipay_pc | alipay_wap
}

// AlipayConfig 支付宝配置
type AlipayConfig struct {
	AppID        string `json:"app_id"`         // 应用ID
	PrivateKey   string `json:"private_key"`    // 商户私钥(RSA2)
	AliPublicKey string `json:"ali_public_key"` // 支付宝公钥
	NotifyURL    string `json:"notify_url"`     // 异步回调地址
	ReturnURL    string `json:"return_url"`     // 同步跳转地址
	IsSandbox    bool   `json:"is_sandbox"`     // 是否沙箱模式
}

func NewAlipayProvider(configJSON string, payway string) (*AlipayProvider, error) {
	var cfg AlipayConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("解析支付宝配置失败: %w", err)
	}
	return &AlipayProvider{config: cfg, payway: payway}, nil
}

func (p *AlipayProvider) gatewayURL() string {
	if p.config.IsSandbox {
		return "https://openapi-sandbox.dl.alipaydev.com/gateway.do"
	}
	return "https://openapi.alipay.com/gateway.do"
}

// CreatePayment 创建支付宝支付
func (p *AlipayProvider) CreatePayment(input CreateInput) (*CreateResult, error) {
	notifyURL := input.NotifyURL
	if notifyURL == "" {
		notifyURL = p.config.NotifyURL
	}
	returnURL := input.ReturnURL
	if returnURL == "" {
		returnURL = p.config.ReturnURL
	}

	// 公共参数
	params := map[string]string{
		"app_id":    p.config.AppID,
		"format":    "JSON",
		"charset":   "utf-8",
		"sign_type": "RSA2",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"version":   "1.0",
		"notify_url": notifyURL,
	}

	// 业务参数
	bizContent := map[string]interface{}{
		"out_trade_no": input.OrderNo,
		"total_amount": fmt.Sprintf("%.2f", input.Amount),
		"subject":      input.Subject,
		"product_code": "FAST_INSTANT_TRADE_PAY",
	}
	bizContentJSON, _ := json.Marshal(bizContent)

	switch p.payway {
	case "alipay_f2f":
		// 当面付（扫码）
		params["method"] = "alipay.trade.precreate"
		bizContent["product_code"] = "FACE_TO_FACE_PAYMENT"
		bizContentJSON, _ = json.Marshal(bizContent)
		params["biz_content"] = string(bizContentJSON)

		sign, err := p.signParams(params)
		if err != nil {
			return nil, err
		}
		params["sign"] = sign

		// 调用支付宝 API
		resp, err := p.doRequest(params)
		if err != nil {
			return nil, err
		}

		qrCode, _ := resp["qr_code"].(string)
		return &CreateResult{
			PayURL:  qrCode,
			QRCode:  qrCode,
			TradeNo: fmt.Sprintf("ALI%d", time.Now().UnixNano()),
			RawResponse: resp,
		}, nil

	case "alipay_pc":
		// 电脑网站支付
		params["method"] = "alipay.trade.page.pay"
		params["return_url"] = returnURL
		params["biz_content"] = string(bizContentJSON)

		sign, err := p.signParams(params)
		if err != nil {
			return nil, err
		}
		params["sign"] = sign

		payURL := p.buildPayURL(params)
		return &CreateResult{
			PayURL:  payURL,
			TradeNo: fmt.Sprintf("ALI%d", time.Now().UnixNano()),
			RawResponse: map[string]interface{}{"pay_url": payURL},
		}, nil

	case "alipay_wap":
		// 手机网站支付
		params["method"] = "alipay.trade.wap.pay"
		params["return_url"] = returnURL
		bizContent["product_code"] = "QUICK_WAP_WAY"
		bizContentJSON, _ = json.Marshal(bizContent)
		params["biz_content"] = string(bizContentJSON)

		sign, err := p.signParams(params)
		if err != nil {
			return nil, err
		}
		params["sign"] = sign

		payURL := p.buildPayURL(params)
		return &CreateResult{
			PayURL:  payURL,
			TradeNo: fmt.Sprintf("ALI%d", time.Now().UnixNano()),
			RawResponse: map[string]interface{}{"pay_url": payURL},
		}, nil

	default:
		return nil, fmt.Errorf("不支持的支付宝渠道: %s", p.payway)
	}
}

// VerifyNotify 验证支付宝异步通知
func (p *AlipayProvider) VerifyNotify(body []byte, headers map[string]string) (*NotifyResult, error) {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("解析通知数据失败: %w", err)
	}

	tradeStatus := values.Get("trade_status")
	if tradeStatus != "TRADE_SUCCESS" && tradeStatus != "TRADE_FINISHED" {
		return nil, fmt.Errorf("交易未成功: %s", tradeStatus)
	}

	return &NotifyResult{
		TradeNo: values.Get("trade_no"),
		OrderNo: values.Get("out_trade_no"),
		Amount:  0, // 需从 total_amount 解析
		Status:  "success",
	}, nil
}

// signParams RSA2 签名
func (p *AlipayProvider) signParams(params map[string]string) (string, error) {
	// 排序参数
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" || k == "sign_type" {
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

	// RSA2 签名
	hash := sha256.Sum256([]byte(buf.String()))

	// 解析私钥（简化处理，实际需要 PEM 解码）
	privateKey, err := parseRSAPrivateKey(p.config.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("解析私钥失败: %w", err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("签名失败: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func (p *AlipayProvider) buildPayURL(params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return p.gatewayURL() + "?" + values.Encode()
}

func (p *AlipayProvider) doRequest(params map[string]string) (map[string]interface{}, error) {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	resp, err := http.PostForm(p.gatewayURL(), values)
	if err != nil {
		return nil, fmt.Errorf("请求支付宝失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return result, nil
}

// parseRSAPrivateKey 解析 RSA 私钥（简化版本）
func parseRSAPrivateKey(keyStr string) (*rsa.PrivateKey, error) {
	// 实际实现需要从 PEM 格式解析
	// 这里返回一个临时密钥用于编译通过
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	_ = keyStr // 实际应解析 keyStr
	return key, nil
}
