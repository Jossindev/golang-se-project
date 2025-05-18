package services

import (
	"awesomeProject/utils"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService struct {
	ApiKey    string
	FromName  string
	FromEmail string
}

func NewEmailService() *EmailService {
	return &EmailService{
		ApiKey:    utils.GetEnv("SENDGRID_API_KEY", ""),
		FromName:  utils.GetEnv("FROM_NAME", "Weather API"),
		FromEmail: utils.GetEnv("FROM_EMAIL", ""),
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
	from := mail.NewEmail(s.FromName, s.FromEmail)
	toEmail := mail.NewEmail("", to)

	plainTextContent := body
	htmlContent := body

	message := mail.NewSingleEmail(from, subject, toEmail, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.ApiKey)

	response, err := client.Send(message)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		fmt.Printf("Email sent successfully! Status code: %d\n", response.StatusCode)
	} else {
		fmt.Printf("Email API request failed. Status: %d, Body: %s\n", response.StatusCode, response.Body)
		return fmt.Errorf("email sending failed with status code %d", response.StatusCode)
	}

	return err
}
