package utils

import (
	"net/smtp"
	"user-service/config"
)

func SendEmail(to, subject, message string) error {
	conf, err := config.LoadConfig()
	if err != nil {
		return err
	}

	from := conf.EMAIL
	password := conf.EMAIL_PASSWORD
	smtpHost := conf.EMAIL_SMTP_HOST
	smtpPort := conf.EMAIL_SMTP_PORT

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}
