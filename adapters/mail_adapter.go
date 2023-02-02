package adapters

import "net/smtp"

type MailAdapter struct {
}
type EmailSender interface {
	Send(to []string, body []byte) error
}

type EmailConfig struct {
	Password   string
	ServerHost string
	ServerPort string
	Username   string
	SenderAddr string
}

type emailSender struct {
	conf EmailConfig
	send func(string, smtp.Auth, string, []string, []byte) error
}

func (e *emailSender) Send(to []string, body []byte) error {
	addr := e.conf.ServerHost + ":" + e.conf.ServerPort
	auth := smtp.PlainAuth("", e.conf.Username, e.conf.Password, e.conf.ServerHost)
	return e.send(addr, auth, e.conf.SenderAddr, to, body)
}
func NewEmailSender(conf EmailConfig) EmailSender {
	return &emailSender{conf, smtp.SendMail}
}

/*type Message struct {
	To      []string
	Body    string
	Subject string
	From    string
}

func (m Message) Send() error {
	auth := smtp.PlainAuth("", m.From, config.Config.Mail.Password, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, m.From, m.To, []byte("Subject:"+m.Subject+"\r\n"+m.Body))
	if err != nil {
		return err
	}
	return nil
}

func (s MailAdapter) SendEmail(from string, to []string, body string, subject string) error {

	// 전송할 메일 메세지생성
	message := Message{
		From:    from,
		To:      to,
		Body:    body,
		Subject: subject,
	}
	// 메일전송
	return message.Send()
}
*/
