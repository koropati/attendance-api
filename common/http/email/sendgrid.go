package email

import (
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
)

type SendGrid interface {
	SendActivation(toName string, toEmail string, token string) error
}

type sendGrid struct {
	client *sendgrid.Client
	config *viper.Viper
}

func NewSendGrid(c *sendgrid.Client, config *viper.Viper) SendGrid {
	return &sendGrid{
		client: c,
		config: config,
	}
}

func (c *sendGrid) SendActivation(toName string, toEmail string, token string) error {
	subject := "Aktivasi Akun"
	if toName == "" {
		toName = toEmail
	}

	plainText := "Hello " + toName + ", silahkan aktivasi akun mu! Token : " + token
	htmlText := ""

	from := mail.NewEmail(c.config.Sub("general").GetString("company_name"), c.config.Sub("general").GetString("company_email"))

	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, subject, to, plainText, htmlText)
	_, err := c.client.Send(message)
	if err != nil {
		log.Printf("Err GOMAIL: %v", err)
		return err
	}
	return nil
}
