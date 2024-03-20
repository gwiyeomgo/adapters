package adapters

import (
	"github.com/gwiyeomgo/adapters/config"
	"net/smtp"
)

type EmailSender interface {
	Send(to []string, body []byte) error
}
type emailSender struct {
	send func(string, smtp.Auth, string, []string, []byte) error
}

type MailAdapter struct {
}

func (MailAdapter) NewEmailSender() EmailSender {
	return &emailSender{send: smtp.SendMail}
}

func (e *emailSender) Send(to []string, body []byte) error {
	addr := config.Config.Mail.ServerHost + ":" + config.Config.Mail.ServerPort
	auth := smtp.PlainAuth("", config.Config.Mail.Username, config.Config.Mail.Password, config.Config.Mail.ServerHost)
	return e.send(addr, auth, config.Config.Mail.SenderAddr, to, body)
}
