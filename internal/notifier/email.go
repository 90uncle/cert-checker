package notifier

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"
)

type EmailNotifier struct {
	Enabled    bool     `yaml:"enabled"`
	SMTPServer string   `yaml:"smtp_server"`
	SMTPPort   int      `yaml:"smtp_port"`
	SMTPUser   string   `yaml:"smtp_username"`
	SMTPPass   string   `yaml:"smtp_password"`
	From       string   `yaml:"from"`
	To         []string `yaml:"to"`
}

func (n *EmailNotifier) Send(msg AlertMessage) error {
	if !n.Enabled || len(n.To) == 0 {
		return nil
	}

	var data = MailData{
		Title:        "证书过期告警",
		AlertMessage: "SSL证书过期告警",
		Content:      msg.String(),
		CompanyName:  "KoodPower",
		Now:          time.Time{},
	}

	htmlContent, _ := RenderTemplate(data)
	mailContent, _ := creteHtmlMail(EmailContent{
		FromName:    data.CompanyName,
		FromAddress: n.From,
		To:          n.To,
		Subject:     data.Title,
		HTMLBody:    htmlContent,
	})

	return n.sendWithSSL(mailContent.String())
}
func (n *EmailNotifier) Name() string {
	return "Email"
}

func (n *EmailNotifier) IsEnabled() bool {
	return n.Enabled
}

// sendWithSSL 使用 SSL/TLS 加密发送邮件
func (n *EmailNotifier) sendWithSSL(content string) error {
	// 设置 TLS 配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         n.SMTPServer,
	}
	// 连接服务器 (使用 SSL/TLS)
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", n.SMTPServer, n.SMTPPort), tlsConfig)
	if err != nil {
		return fmt.Errorf("SSL/TLS connection failed: %v", err)
	}
	defer conn.Close()
	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, n.SMTPServer)
	if err != nil {
		return fmt.Errorf("SMTP client creation error: %v", err)
	}
	defer client.Close()
	// 认证
	auth := smtp.PlainAuth("", n.SMTPUser, n.SMTPPass, n.SMTPServer)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication error: %v", err)
	}

	// 设置发件人
	if err := client.Mail(n.From); err != nil {
		return fmt.Errorf("MAIL command error: %v", err)
	}
	// 设置收件人
	for _, recipient := range n.To {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("RCPT command error for %s: %v", recipient, err)
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command error: %v", err)
	}
	defer w.Close()
	if _, err = w.Write([]byte(content)); err != nil {
		return fmt.Errorf("mail content writing error: %v", err)
	}

	return client.Quit()
}

type EmailContent struct {
	FromName    string
	FromAddress string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	HTMLBody    string
}

func creteHtmlMail(content EmailContent) (*bytes.Buffer, error) {
	// 创建MIME多部分消息
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	// 1. 设置基础邮件头
	headers := map[string]string{
		"Message-ID":   fmt.Sprintf("<%d@%s>", time.Now().UnixNano(), "koodpower.com"),
		"Date":         time.Now().Format(time.RFC1123Z),
		"From":         fmt.Sprintf("%s <%s>", content.FromName, content.FromAddress),
		"To":           strings.Join(content.To, ", "),
		"Cc":           strings.Join(content.Cc, ", "),
		"Subject":      content.Subject,
		"MIME-Version": "1.0",
		"Content-Type": fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary()),
	}
	for k, v := range headers {
		fmt.Fprintf(body, "%s: %s\r\n", k, v)
	}
	fmt.Fprint(body, "\r\n")
	// 2. 添加正文部分 (替代文本+HTML)
	part, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"text/html; charset=utf-8"},
	})
	if err != nil {
		return body, err
	}

	part.Write([]byte(content.HTMLBody + "\r\n"))

	return body, nil
}

type MailData struct {
	Title        string
	AlertMessage string
	Content      string
	CompanyName  string
	Now          time.Time
}

func RenderTemplate(data MailData) (string, error) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f8f8; padding: 15px; text-align: center; }
        .content { padding: 20px; background-color: #ffffff; }
        .footer { margin-top: 20px; padding: 10px; text-align: center; font-size: 12px; color: #777; }
        .alert { background-color: #fff3cd; padding: 15px; margin-bottom: 15px; border-left: 5px solid #ffc107; }
        .button { background-color: #007bff; color: white; padding: 10px 15px; text-decoration: none; border-radius: 4px; display: inline-block; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>系统通知</h2>
        </div>
        <div class="content">
            <h3>{{.Title}}</h3>
            
            {{if .AlertMessage}}
            <div class="alert">
                {{.AlertMessage}}
            </div>
            {{end}}
            
            <p>{{.Content}}</p>

            <p>如果此邮件您未预期收到，请忽略或联系我们。</p>
        </div>
        <div class="footer">
            <p>© {{.Now.Year}} {{.CompanyName}}. 保留所有权利。</p>
        </div>
    </div>
</body>
</html>`
	// 添加当前时间到数据
	data.Now = time.Now()
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
