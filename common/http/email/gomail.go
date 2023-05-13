package email

import (
	"log"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type Email interface {
	SendActivation(toUserName string, toEmail string, urlActivation string) error
}

type email struct {
	m      *gomail.Dialer
	config *viper.Viper
}

func New(m *gomail.Dialer, config *viper.Viper) Email {
	return &email{
		m:      m,
		config: config,
	}
}

func (m *email) SendActivation(toUserName string, toEmail string, urlActivation string) error {
	activationUserHTML := GenerateTemplateActivationAccount(urlActivation, toUserName, toEmail, m.config)
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", m.config.Sub("general").GetString("company_name")+" <"+m.config.Sub("general").GetString("company_email")+">")
	mailer.SetHeader("To", toEmail)
	mailer.SetAddressHeader("Cc", "dewa.ketut.satriawan@gmail.com", "Admin")
	mailer.SetHeader("Subject", "Silahkan aktivasi akun mu")
	mailer.SetBody("text/html", activationUserHTML)
	// mailer.Attach("./sample.png")
	err := m.m.DialAndSend(mailer)
	if err != nil {
		log.Printf("Err GOMAIL: %v", err)
		return err
	}
	return nil
}
