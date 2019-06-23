package email

import (
	"fmt"
	"github.com/copybird/copybird/core"
	"net/smtp"
)

const GROUP_NAME = "global"
const TYPE_NAME = "notifier"
const MODULE_NAME = "email"

type GlobalNotifierEmail struct {
	core.Module
	Config *Config
}

func (e *GlobalNotifierEmail) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (e *GlobalNotifierEmail) GetType() core.ModuleType {
	return TYPE_NAME
}

func (e *GlobalNotifierEmail) GetName() string {
	return MODULE_NAME
}

func (e *GlobalNotifierEmail) GetConfig() interface{} {
	return &Config{}
}

func (e *GlobalNotifierEmail) InitModule(_cfg interface{}) error {
	e.Config = _cfg.(*Config)
	return nil
}

func (e *GlobalNotifierEmail) Run() error {
	if err := e.SendEmail(); err != nil {
		return err
	}

	return nil
}

func (e *GlobalNotifierEmail) Close() error {
	return nil
}

func (e *GlobalNotifierEmail) SendEmail() error {

	from := e.Config.MailerUser
	pass := e.Config.MailerPassword
	to := e.Config.MailTo
	body := "Dump created successfully"
	subject := "Dump"

	header := ""
	header += fmt.Sprintf("From: %s\r\n", from)
	header += fmt.Sprintf("To: %s\r\n", to)
	header += fmt.Sprint("MIME-Version: 1.0\r\n")
	header += fmt.Sprint("Content-type: text/html\r\n")
	header += fmt.Sprintf("Subject: %s\r\n", subject)
	header += "\r\n" + body + "\r\n"

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(header))

	if err != nil {
		return err
	}

	return nil
}
