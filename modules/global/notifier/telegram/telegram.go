package telegram

import (
	"io"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Local struct {
	config *Config
	reader io.Reader
	writer io.Writer
}

func (l *Local) GetName() string {
	return MODULE_NAME
}

func (local *Local) InitPipe(w io.Writer, r io.Reader) error {
	local.reader = r
	local.writer = w
	return nil
}

func (l *Local) InitModule(_cfg interface{}) error {
	l.config = _cfg.(*Config)
	return nil
}

func (l *Local) Run() error {
	if err := l.config.NotifyTelegramChannel(); err != nil {
		return err
	}

	return nil
}
func (l *Local) GetConfig() interface{} {
	return &Config{}
}

func (l *Local) Close() error {
	return nil
}

const MODULE_NAME = "telegram"

func (conf *Config) NotifyTelegramChannel() error {

	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		return err
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	msg := tgbotapi.NewMessage(conf.ChannelID, conf.Message)
	bot.Send(msg)

	return nil
}
