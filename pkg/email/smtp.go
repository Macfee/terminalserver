package email

import (
	"audit-system/pkg/config"
	"fmt"
	"net/smtp"
)

func SendEmail(cfg config.SMTPConfig, subject, body string) error {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	msg := fmt.Sprintf("From: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s", cfg.From, subject, body)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	return smtp.SendMail(addr, auth, cfg.From, []string{cfg.From}, []byte(msg))
}
