package telegram

import (
	"github.com/copybird/copybird/core"
	"io"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type GlobalNotifierTelegram struct {
	core.Module
	config *Config
	reader io.Reader
	writer io.Writer
}

func (m *GlobalNotifierTelegram) GetName() string {
	return MODULE_NAME
}

func (m *GlobalNotifierTelegram) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *GlobalNotifierTelegram) InitModule(_cfg interface{}) error {
	m.config = _cfg.(*Config)
	return nil
}

func (m *GlobalNotifierTelegram) Run() error {
	if err := m.config.NotifyTelegramChannel(); err != nil {
		return err
	}

	return nil
}
func (m *GlobalNotifierTelegram) GetConfig() interface{} {
	return &Config{}
}

func (m *GlobalNotifierTelegram) Close() error {
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
