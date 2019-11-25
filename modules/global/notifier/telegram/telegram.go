package telegram

import (
	"context"
	"io"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/copybird/copybird/core"
)

const GROUP_NAME = "global"
const TYPE_NAME = "notifier"
const MODULE_NAME = "telegram"

type GlobalNotifierTelegram struct {
	core.Module
	config *Config
	reader io.Reader
	writer io.Writer
}

func (m *GlobalNotifierTelegram) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *GlobalNotifierTelegram) GetType() core.ModuleType {
	return TYPE_NAME
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

func (m *GlobalNotifierTelegram) Run(ctx context.Context) error {
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
