package payment

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

// WechatPayProvider 微信支付提供商
type WechatPayProvider struct {
	config AlipayConfig2 // 使用独立命名避免冲突
	payway string        // wechat_native | wechat_h5 | wechat_jsapi
}

// WechatPayConfig 微信支付配置
type WechatPayConfig struct {
	AppID     string `json:"app_id"`     // 微信 AppID
	MchID     string `json:"mch_id"`     // 商户号
	APIKey    string `json:"api_key"`    // API密钥（V2）
	APIKeyV3  string `json:"api_key_v3"` // API密钥（V3）
	SerialNo  string `json:"serial_no"`  // 证书序列号
	NotifyURL string `json:"notify_url"` // 回调地址
}

func NewWechatPayProvider(configJSON string, payway string) (*WechatPayProvider, error) {
	var cfg AlipayConfig2
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("解析微信支付配置失败: %w", err)
	}
	return &WechatPayProvider{config: cfg, payway: payway}, nil
}

// AlipayConfig2 微信支付配置（复用结构体名称）
type AlipayConfig2 = WechatPayConfig

// wechatUnifiedOrderReq 微信统一下单请求 (V2 XML)
type wechatUnifiedOrderReq struct {
	XMLName        xml.Name `xml:"xml"`
	AppID          string   `xml:"appid"`
	MchID          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	Body           string   `xml:"body"`
	OutTradeNo     string   `xml:"out_trade_no"`
	TotalFee       int      `xml:"total_fee"` // 单位: 分
	SpbillCreateIP string   `xml:"spbill_create_ip"`
	NotifyURL      string   `xml:"notify_url"`
	TradeType      string   `xml:"trade_type"` // NATIVE | MWEB | JSAPI
	SceneInfo      string   `xml:"scene_info,omitempty"`
}

// wechatUnifiedOrderResp 微信统一下单响应
type wechatUnifiedOrderResp struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
	ResultCode string   `xml:"result_code"`
	ErrCode    string   `xml:"err_code"`
	ErrCodeDes string   `xml:"err_code_des"`
	PrepayID   string   `xml:"prepay_id"`
	CodeURL    string   `xml:"code_url"`  // NATIVE 扫码链接
	MwebURL    string   `xml:"mweb_url"`  // H5 支付链接
}

// CreatePayment 创建微信支付
func (p *WechatPayProvider) CreatePayment(input CreateInput) (*CreateResult, error) {
	notifyURL := input.NotifyURL
	if notifyURL == "" {
		notifyURL = p.config.NotifyURL
	}

	// 金额转为分
	totalFee := int(input.Amount * 100)

	var tradeType string
	switch p.payway {
	case "wechat_native":
		tradeType = "NATIVE"
	case "wechat_h5":
		tradeType = "MWEB"
	case "wechat_jsapi":
		tradeType = "JSAPI"
	default:
		return nil, fmt.Errorf("不支持的微信支付渠道: %s", p.payway)
	}

	nonceStr := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%d%d", time.Now().UnixNano(), input.PaymentID))))

	reqData := map[string]string{
		"appid":            p.config.AppID,
		"mch_id":           p.config.MchID,
		"nonce_str":        nonceStr,
		"body":             input.Subject,
		"out_trade_no":     input.OrderNo,
		"total_fee":        fmt.Sprintf("%d", totalFee),
		"spbill_create_ip": "127.0.0.1",
		"notify_url":       notifyURL,
		"trade_type":       tradeType,
	}

	// H5 支付需要 scene_info
	if tradeType == "MWEB" {
		reqData["scene_info"] = `{"h5_info":{"type":"Wap","wap_url":"https://x-store.com","wap_name":"X-Store"}}`
	}

	// 生成签名
	sign := p.signMD5(reqData)
	reqData["sign"] = sign

	// 构建 XML
	req := wechatUnifiedOrderReq{
		AppID:          reqData["appid"],
		MchID:          reqData["mch_id"],
		NonceStr:       reqData["nonce_str"],
		Sign:           sign,
		Body:           reqData["body"],
		OutTradeNo:     reqData["out_trade_no"],
		TotalFee:       totalFee,
		SpbillCreateIP: reqData["spbill_create_ip"],
		NotifyURL:      reqData["notify_url"],
		TradeType:      tradeType,
	}
	if tradeType == "MWEB" {
		req.SceneInfo = reqData["scene_info"]
	}

	xmlData, _ := xml.Marshal(req)

	// 发送请求
	resp, err := http.Post(
		"https://api.mch.weixin.qq.com/pay/unifiedorder",
		"application/xml",
		strings.NewReader(string(xmlData)),
	)
	if err != nil {
		return nil, fmt.Errorf("请求微信支付失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var wxResp wechatUnifiedOrderResp
	if err := xml.Unmarshal(respBody, &wxResp); err != nil {
		return nil, fmt.Errorf("解析微信响应失败: %w", err)
	}

	if wxResp.ReturnCode != "SUCCESS" || wxResp.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("微信支付失败: %s %s", wxResp.ErrCode, wxResp.ErrCodeDes)
	}

	result := &CreateResult{
		TradeNo: fmt.Sprintf("WX%d", time.Now().UnixNano()),
		RawResponse: map[string]interface{}{
			"prepay_id": wxResp.PrepayID,
		},
	}

	switch tradeType {
	case "NATIVE":
		result.PayURL = wxResp.CodeURL
		result.QRCode = wxResp.CodeURL
	case "MWEB":
		result.PayURL = wxResp.MwebURL
	case "JSAPI":
		result.PayURL = wxResp.PrepayID
	}

	return result, nil
}

// VerifyNotify 验证微信支付回调
func (p *WechatPayProvider) VerifyNotify(body []byte, headers map[string]string) (*NotifyResult, error) {
	var notify struct {
		XMLName    xml.Name `xml:"xml"`
		ReturnCode string   `xml:"return_code"`
		ResultCode string   `xml:"result_code"`
		OutTradeNo string   `xml:"out_trade_no"`
		TransactionID string `xml:"transaction_id"`
		TotalFee   int      `xml:"total_fee"`
	}

	if err := xml.Unmarshal(body, &notify); err != nil {
		return nil, fmt.Errorf("解析回调数据失败: %w", err)
	}

	if notify.ReturnCode != "SUCCESS" || notify.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("支付未成功")
	}

	return &NotifyResult{
		TradeNo: notify.TransactionID,
		OrderNo: notify.OutTradeNo,
		Amount:  float64(notify.TotalFee) / 100,
		Status:  "success",
	}, nil
}

// signMD5 微信支付 MD5 签名
func (p *WechatPayProvider) signMD5(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" || params[k] == "" {
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
	buf.WriteString("key=")
	buf.WriteString(p.config.APIKey)

	hash := md5.Sum([]byte(buf.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}
