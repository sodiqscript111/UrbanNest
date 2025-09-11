package email

import (
	"context"
	"github.com/resend/resend-go/v2"
)

type EmailParams struct {
	To      string
	Subject string
	Body    string
}

type ResendClient struct {
	client *resend.Client
}

func NewResendClient(apiKey string) *ResendClient {
	return &ResendClient{client: resend.NewClient(apiKey)}
}

func (c *ResendClient) SendEmail(ctx context.Context, params EmailParams) error {
	email := &resend.SendEmailRequest{
		From:    "no-reply@urban-nest.com",
		To:      []string{params.To},
		Subject: params.Subject,
		Html:    params.Body, // Use Html for rich text; adjust if plain text is preferred
	}
	_, err := c.client.Emails.SendWithContext(ctx, email)
	return err
}
