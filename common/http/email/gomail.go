package email

import (
	"log"

	"gopkg.in/gomail.v2"
)

type Email interface {
	SendActivation(toUserName string, toEmail string, urlActivation string) error
}

type email struct {
	m *gomail.Dialer
}

func New(m *gomail.Dialer) Email {
	return &email{m: m}
}

const CONFIG_SENDER_NAME = "PT. Makmur Subur Jaya"
const CONFIG_SENDER_EMAIL = "wokdev@gmail.com"
const CONFIG_WEB_NAME = "WOKDEV"
const CONFIG_WEB_LOGO = "www.github.com"
const CONFIG_SOCIAL_GITHUB_URL = "www.github.com"
const CONFIG_SOCIAL_FACEBOOK_URL = "www.facebook.com"
const CONFIG_SOCIAL_EMAIL_URL = "www.gmail.com"
const CONFIG_SOCIAL_TWITTER_URL = "www.twitter.com"
const CONFIG_SOCIAL_LINKEDIN_URL = "www.linkedin.com"
const CONFIG_SOCIAL_YOUTUBE_URL = "www.youtube.com"
const CONFIG_SOCIAL_INSTAGRAM_URL = "www.instagram.com"
const CONFIG_TERM_OF_USE_URL = "www.github.com"
const CONFIG_PRIVACY_POLICY_URL = "www.github.com"

func (m *email) SendActivation(toUserName string, toEmail string, urlActivation string) error {
	activationUserHTML := GenerateTemplateActivationAccount(urlActivation, toUserName, toEmail)
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME+" <"+CONFIG_SENDER_EMAIL+">")
	mailer.SetHeader("To", toEmail)
	mailer.SetAddressHeader("Cc", "dewa.ketut.satriawan@gmail.com", "Admin WokDev")
	mailer.SetHeader("Subject", "Please Activated your account in WokDev")
	mailer.SetBody("text/html", activationUserHTML)
	// mailer.Attach("./sample.png")
	err := m.m.DialAndSend(mailer)
	if err != nil {
		log.Printf("Err GOMAIL: %v", err)
		return err
	}
	return nil
}
