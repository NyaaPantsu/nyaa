package email

import (
	"path/filepath"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/util/log"
	gomail "gopkg.in/gomail.v2"
)

type EmailError error

var (
	mailer = InitGomail()
)

func InitGomail() *gomail.Dialer {
	mailer := gomail.NewDialer(config.EmailHost, config.EmailPort, config.EmailUsername, config.EmailPassword)
	return mailer
}

func SendEmailFromAdmin(to string, subject string, body string, bodyHTML string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.EmailFrom)
	msg.SetHeader("To", to, config.EmailTestTo)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)
	msg.AddAlternative("text/html", bodyHTML)
	log.Debugf("to : %s", to)
	log.Debugf("subject : %s", subject)
	log.Debugf("body : %s", body)
	log.Debugf("bodyHTML : %s", bodyHTML)
	if config.SendEmail {
		log.Debug("SendEmail performed.")

		err := mailer.DialAndSend(msg); 
		return err
	}
	return nil
}

func SendTestEmail() error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.EmailFrom)
	msg.SetHeader("To", config.EmailTestTo)
	msg.SetAddressHeader("Cc", config.EmailTestTo, "NyaaPantsu")
	msg.SetHeader("Subject", "Hi(안녕하세요)?!")
	msg.SetBody("text/plain", "Hi(안녕하세요)?!")
	msg.AddAlternative("text/html", "<p><b>Nowplay(나우플레이)</b> means <i>Let's play</i>!!?</p>")
	path, err := filepath.Abs("img/megumin.png")
	if err != nil {
		panic(err)
	}
	msg.Attach(path)
	// SendEmail(msg)

	err = mailer.DialAndSend(msg)
	return err
}
