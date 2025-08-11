package mail

import (
	"github.com/hanifbg/landing_backend/config"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	Dialer *gomail.Dialer
}

func Init(config *config.AppConfig) (*Mailer, error) {
	dialer := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)
	mailer := &Mailer{
		Dialer: dialer,
	}
	// //this is for check smtp connection
	// close, err := mailer.Dialer.Dial()
	// if err != nil {
	// 	return nil, err
	// }
	// //close smtp test connection
	// defer close.Close()

	return mailer, nil
}
