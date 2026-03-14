---
sidebar_position: 6
title: 支付系统
---

# 支付系统

X-Store 支持 10 种支付方式，通过统一的 `Provider` 接口抽象实现。

## 支持的支付方式

| 提供商 | 渠道类型 | 说明 |
|--------|---------|------|
| **支付宝** | 当面付 / PC网页 / 手机网页 | 国内主流支付 |
| **微信支付** | Native / H5 / JSAPI | 国内主流支付 |
| **Stripe** | 国际信用卡 | 海外用户 |
| **PayPal** | PayPal 余额/信用卡 | 海外用户 |
| **EPUSDT** | USDT-TRC20 | 加密货币 |
| **易支付** | 支付宝/微信 | 个人免签 |
| **Payjs** | 微信个人 | 个人免签 |
| **虎皮椒** | 支付宝/微信 | 个人免签 |
| **V免签** | 支付宝/微信 | 个人免签 |
| **MockPay** | 模拟支付 | 开发测试用 |

## 统一 Provider 接口

```go title="backend/pkg/payment/provider.go"
type Provider interface {
    // CreatePayment 创建支付订单
    CreatePayment(input CreateInput) (*PaymentResult, error)

    // VerifyNotify 验证支付回调通知
    VerifyNotify(body []byte, headers map[string]string) (*NotificationResult, error)
}

type CreateInput struct {
    OrderNo   string
    PaymentID uint
    Amount    float64
    Subject   string
    NotifyURL string
    ReturnURL string
}

type PaymentResult struct {
    PayURL      string
    QRCode      string
    TradeNo     string
    RawResponse map[string]interface{}
}

type NotificationResult struct {
    OrderNo string
    TradeNo string
    Amount  float64
    Status  string // "success" | "failed"
}
```

## 支付流程

```
1. 用户下单 → 创建 Order（待支付）
2. 用户选择支付渠道 → 创建 Payment 记录
3. 调用 Provider.CreatePayment() → 获取支付链接/二维码
4. 用户完成支付 → 支付平台发送回调
5. Provider.VerifyNotify() → 验签
6. 更新订单状态 → 自动发货（分配卡密）
```

## 添加新支付方式

1. 在 `backend/pkg/payment/` 下创建新文件（如 `newpay.go`）
2. 实现 `Provider` 接口的两个方法
3. 在 `payment_service.go` 的 `createProvider()` 中注册
4. 在数据库 `payment_channels` 表中添加配置

示例：

```go title="backend/pkg/payment/newpay.go"
type NewPayProvider struct {
    apiKey    string
    apiSecret string
}

func NewNewPayProvider(configJSON string) (*NewPayProvider, error) {
    var cfg struct {
        APIKey    string `json:"api_key"`
        APISecret string `json:"api_secret"`
    }
    json.Unmarshal([]byte(configJSON), &cfg)
    return &NewPayProvider{apiKey: cfg.APIKey, apiSecret: cfg.APISecret}, nil
}

func (p *NewPayProvider) CreatePayment(input CreateInput) (*PaymentResult, error) {
    // 调用第三方 API 创建支付
    return &PaymentResult{PayURL: "https://..."}, nil
}

func (p *NewPayProvider) VerifyNotify(body []byte, headers map[string]string) (*NotificationResult, error) {
    // 验证签名，解析回调数据
    return &NotificationResult{OrderNo: "...", Status: "success"}, nil
}
```

## 支付渠道配置

每个渠道的配置存储在 `payment_channels` 表的 `config_json` 字段中，JSON 格式：

```json
{
  "app_id": "2021xxxx",
  "private_key": "MIIEv...",
  "public_key": "MIIBIj...",
  "notify_url": "http://your-domain/api/webhooks/payment/1",
  "return_url": "http://your-domain/payment/success"
}
```

不同支付方式的配置字段不同，具体参见各自的实现文件。
