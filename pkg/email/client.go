package email

import (
	"github.com/resend/resend-go/v2"
)

func SendEmail(apiKey, to, subject, message string) error {
	client := resend.NewClient(apiKey)
	params := &resend.SendEmailRequest{
		From:    "no-reply@yourdomain.com",
		To:      []string{to},
		Subject: subject,
		Text:    message,
	}
	_, err := client.Emails.Send(params)
	return err
}
