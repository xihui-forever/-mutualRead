package mail

import (
	"fmt"
	"github.com/darabuchi/log"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"github.com/xihui-forever/mutualRead/config"
	"net/smtp"
)

func Send(to []string, title, text string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("ITest <%s>", viper.GetString(config.EmailSender))
	e.To = to
	e.Subject = title
	e.Text = []byte(text)
	err := e.Send(
		viper.GetString(config.EmailAddr),
		smtp.PlainAuth("", viper.GetString(config.EmailSender), viper.GetString(config.EmailPassword), viper.GetString(config.EmailHost)),
	)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return err
}
