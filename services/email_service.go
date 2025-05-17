package services

import (
	"awesomeProject/utils"
	"fmt"
	"net/smtp"
)

type EmailService struct {
	SmtpHost string
	SmtpPort string
	Username string
	Password string
	FromName string
}

func NewEmailService() *EmailService {
	return &EmailService{
		SmtpHost: utils.GetEnv("SMTP_HOST", "smtp.example.com"),
		SmtpPort: utils.GetEnv("SMTP_PORT", "587"),
		Username: utils.GetEnv("SMTP_USERNAME", ""),
		Password: utils.GetEnv("SMTP_PASSWORD", ""),
		FromName: utils.GetEnv("FROM_NAME", "Weather API"),
	}
}

func (s *EmailService) SendConfirmationEmail(email, city, token string) error {
	apiURL := utils.GetEnv("API_URL", "http://localhost:8080")
	subject := "Confirm your Weather API subscription"
	confirmationLink := fmt.Sprintf("%s/api/confirm/%s", apiURL, token)

	body := fmt.Sprintf(`Hello,

Thank you for subscribing to weather updates for %s.

Please confirm your subscription by clicking on the following link:
%s

If you did not request this subscription, you can ignore this email.

Best regards,
The Weather API Team
`, city, confirmationLink)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) SendUnsubscribeEmail(email, city, token string) error {
	apiURL := utils.GetEnv("API_URL", "http://localhost:8080")
	subject := "Weather API - How to unsubscribe"
	unsubscribeLink := fmt.Sprintf("%s/api/unsubscribe/%s", apiURL, token)

	body := fmt.Sprintf(`Hello,

You are subscribed to weather updates for %s.

If you wish to unsubscribe, please click on the following link:
%s

Best regards,
The Weather API Team
`, city, unsubscribeLink)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	if utils.GetEnv("APP_ENV", "dev") != "production" {
		fmt.Printf("Email to: %s\nSubject: %s\nBody: %s\n", to, subject, body)
		return nil
	}

	auth := smtp.PlainAuth("", s.Username, s.Password, s.SmtpHost)

	msg := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"From: %s <%s>\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, s.FromName, s.Username, body)

	err := smtp.SendMail(
		s.SmtpHost+":"+s.SmtpPort,
		auth,
		s.Username,
		[]string{to},
		[]byte(msg),
	)

	return err
}
