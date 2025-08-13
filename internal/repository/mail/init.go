package mail

import (
	"fmt"

	"github.com/hanifbg/landing_backend/config"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	Dialer *gomail.Dialer
}

func Init(config *config.AppConfig) (*Mailer, error) {
	fmt.Println("SMTPHost:", config.SMTPHost)
	fmt.Println("SMTPPort:", config.SMTPPort)
	fmt.Println("SMTPUsername:", config.SMTPUsername)
	fmt.Println("SMTPFrom:", config.SMTPFrom)
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
