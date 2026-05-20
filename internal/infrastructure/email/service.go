package email

import (
	"fmt"

	"github.com/quoteyouros/backend/internal/config"
)

// EmailService interface for sending emails
type EmailService interface {
	SendEmail(to, subject, body string) error
}

// NoOpEmailService is a placeholder implementation
type NoOpEmailService struct{}

func NewNoOpEmailService() *NoOpEmailService {
	return &NoOpEmailService{}
}

func (s *NoOpEmailService) SendEmail(to, subject, body string) error {
	fmt.Printf("📧 [EMAIL] To: %s, Subject: %s\n", to, subject)
	return nil
}

// ResendEmailService sends emails using Resend API
type ResendEmailService struct {
	apiKey    string
	fromEmail string
}

func NewResendEmailService(cfg *config.EmailConfig) *ResendEmailService {
	return &ResendEmailService{
		apiKey:    cfg.ResendAPIKey,
		fromEmail: cfg.FromEmail,
	}
}

func (s *ResendEmailService) SendEmail(to, subject, body string) error {
	// TODO: Implement Resend API call
	fmt.Printf("📧 [RESEND] To: %s, Subject: %s\n", to, subject)
	return nil
}

// SMTPEmailService sends emails using SMTP
type SMTPEmailService struct {
	host      string
	port      int
	user      string
	password  string
	fromEmail string
}

func NewSMTPEmailService(cfg *config.EmailConfig) *SMTPEmailService {
	return &SMTPEmailService{
		host:      cfg.SMTPHost,
		port:      cfg.SMTPPort,
		user:      cfg.SMTPUser,
		password:  cfg.SMTPPassword,
		fromEmail: cfg.FromEmail,
	}
}

func (s *SMTPEmailService) SendEmail(to, subject, body string) error {
	// TODO: Implement SMTP email sending
	fmt.Printf("📧 [SMTP] To: %s, Subject: %s\n", to, subject)
	return nil
}

// Factory function to create appropriate email service
func NewEmailService(cfg *config.EmailConfig) EmailService {
	switch cfg.Provider {
	case "resend":
		return NewResendEmailService(cfg)
	case "smtp":
		return NewSMTPEmailService(cfg)
	default:
		return NewNoOpEmailService()
	}
}
