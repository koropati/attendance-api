package email

import (
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGrid interface {
	SendActivation(toName string, toEmail string, token string) error
}

type sendGrid struct {
	client *sendgrid.Client
}

func NewSendGrid(c *sendgrid.Client) SendGrid {
	return &sendGrid{client: c}
}

func (c *sendGrid) SendActivation(toName string, toEmail string, token string) error {
	subject := "Activate Your Account"
	if toName == "" {
		toName = toEmail
	}

	plainText := "Hello " + toName + ", please activate your account! Token : " + token
	htmlText := ""

	from := mail.NewEmail(CONFIG_SENDER_NAME, CONFIG_SENDER_EMAIL)

	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, subject, to, plainText, htmlText)
	_, err := c.client.Send(message)
	if err != nil {
		log.Printf("Err GOMAIL: %v", err)
		return err
	}
	return nil
}
