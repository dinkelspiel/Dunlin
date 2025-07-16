package services

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	config, _ := LoadConfig()

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf(`"%s" <%s>`, config.GmailDisplayName, config.GmailEmail))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 465, config.GmailEmail, config.GmailPassword)
	d.SSL = true
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
