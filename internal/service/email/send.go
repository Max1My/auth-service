package email

import (
	"auth-service/internal/model"
	"fmt"
	"gopkg.in/gomail.v2"
)

func (s *serv) SendConfirmationEmail(message model.MailMessageInfo) error {
	fmt.Println(message)
	mail := gomail.NewMessage()
	mail.SetHeader("From", message.SmtpUser+"@yandex.ru")
	mail.SetHeader("To", message.To)
	mail.SetHeader("Subject", "Confirm your email address")

	confirmURL := "https://your-service.com/confirm-email?token=" + message.Token
	mail.SetBody("text/plain", "Please confirm your email by clicking the following link: "+confirmURL)

	dialer := gomail.NewDialer(message.SmtpHost, message.SmtpPort, message.SmtpUser, message.SmtpPassword)
	return dialer.DialAndSend(mail)
}
