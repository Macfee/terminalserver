package alert

import (
	"audit-system/pkg/config"
	"audit-system/pkg/email"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// 邮件告警处理器
type EmailAlertHandler struct {
	smtpConfig config.SMTPConfig
}

func NewEmailAlertHandler(config config.SMTPConfig) *EmailAlertHandler {
	return &EmailAlertHandler{
		smtpConfig: config,
	}
}

func (h *EmailAlertHandler) HandleAlert(alert Alert) error {
	// 根据告警级别发送邮件
	subject := fmt.Sprintf("[%s] 安全告警: %s", alert.Level, alert.Type)
	body := fmt.Sprintf(
		"会话ID: %s\n用户ID: %d\n时间: %s\n消息: %s",
		alert.SessionID,
		alert.UserID,
		alert.Time.Format("2006-01-02 15:04:05"),
		alert.Message,
	)

	return email.SendEmail(h.smtpConfig, subject, body)
}

// Webhook告警处理器
type WebhookAlertHandler struct {
	webhookURL string
}

func NewWebhookAlertHandler(webhookURL string) *WebhookAlertHandler {
	return &WebhookAlertHandler{
		webhookURL: webhookURL,
	}
}

func (h *WebhookAlertHandler) HandleAlert(alert Alert) error {
	payload, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	resp, err := http.Post(h.webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook返回非200状态码: %d", resp.StatusCode)
	}

	return nil
}
