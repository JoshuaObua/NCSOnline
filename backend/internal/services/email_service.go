package services

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
)

type EmailService struct{ cfg *config.Config }

func NewEmailService(cfg *config.Config) *EmailService { return &EmailService{cfg: cfg} }

func (s *EmailService) Send(n *repository.Notification) error {
	if s.cfg.SMTPHost == "" {
		return fmt.Errorf("SMTP_HOST is not configured")
	}
	var payload map[string]string
	if err := json.Unmarshal(n.Payload, &payload); err != nil {
		return err
	}
	subject := "NCS notification"
	body := "You have a new notification from the National Council of Sports."
	if n.TemplateCode == "ORGANISATION_INVITATION" {
		subject = "Activate your NCS organisation profile"
		link := s.cfg.PublicAppURL + payload["activation_path"]
		body = fmt.Sprintf("Your %s profile has been approved. Sign in with this email address and use the secure link below within 72 hours:\r\n\r\n%s\r\n\r\nNCS will never email you a password.", payload["organisation_name"], link)
	}
	fromAddress := s.cfg.SMTPFrom
	if i := strings.LastIndex(fromAddress, "<"); i >= 0 && strings.HasSuffix(fromAddress, ">") {
		fromAddress = fromAddress[i+1 : len(fromAddress)-1]
	}
	message := []byte("From: " + s.cfg.SMTPFrom + "\r\nTo: " + n.RecipientAddress + "\r\nSubject: " + subject + "\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n" + body)
	var auth smtp.Auth
	if s.cfg.SMTPUser != "" {
		auth = smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	}
	return smtp.SendMail(s.cfg.SMTPHost+":"+s.cfg.SMTPPort, auth, fromAddress, []string{n.RecipientAddress}, message)
}
