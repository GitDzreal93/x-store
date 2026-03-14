package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/x-store/backend/internal/config"
)

type EmailService struct {
	config *config.EmailConfig
}

func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{config: cfg}
}

type OrderEmailData struct {
	OrderNo   string
	Email     string
	ProductTitle string
	Amount    float64
	CardKeys  []string
	CreatedAt string
}

// SendOrderNotification 发送订单通知邮件
func (s *EmailService) SendOrderNotification(data OrderEmailData) error {
	if !s.config.Enabled {
		return nil // 邮件未启用，跳过
	}

	subject := fmt.Sprintf("订单 %s 支付成功 - X-Store", data.OrderNo)
	body, err := s.renderOrderTemplate(data)
	if err != nil {
		return fmt.Errorf("render template failed: %w", err)
	}

	return s.send(data.Email, subject, body)
}

func (s *EmailService) renderOrderTemplate(data OrderEmailData) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .order-info { background: white; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .card-key { background: #e8f5e9; padding: 15px; margin: 10px 0; border-left: 4px solid #4caf50; font-family: monospace; font-size: 16px; }
        .footer { text-align: center; color: #999; margin-top: 30px; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 支付成功</h1>
            <p>感谢您的购买！</p>
        </div>
        <div class="content">
            <div class="order-info">
                <h2>订单信息</h2>
                <p><strong>订单号：</strong>{{.OrderNo}}</p>
                <p><strong>商品：</strong>{{.ProductTitle}}</p>
                <p><strong>金额：</strong>¥{{.Amount}}</p>
                <p><strong>下单时间：</strong>{{.CreatedAt}}</p>
            </div>
            
            {{if .CardKeys}}
            <div class="order-info">
                <h2>🔑 您的卡密</h2>
                <p>请妥善保管以下卡密信息：</p>
                {{range .CardKeys}}
                <div class="card-key">{{.}}</div>
                {{end}}
            </div>
            {{end}}
            
            <p style="margin-top: 20px; color: #666;">
                如有任何问题，请联系客服。<br>
                此邮件为系统自动发送，请勿直接回复。
            </p>
        </div>
        <div class="footer">
            <p>© 2026 X-Store. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("order").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *EmailService) send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	return smtp.SendMail(addr, auth, s.config.From, []string{to}, msg)
}
