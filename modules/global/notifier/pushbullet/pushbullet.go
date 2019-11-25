package pushbullet

import (
	"context"
	"io"

	"github.com/copybird/copybird/core"

	"github.com/xconstruct/go-pushbullet"
)

const GROUP_NAME = "global"
const TYPE_NAME = "notifier"
const MODULE_NAME = "pushbullet"

type GlobalNotifierPushbullet struct {
	core.Module
	config *Config
	reader io.Reader
	writer io.Writer
}

type Message struct {
	Text string `json:"text"`
}

func (m *GlobalNotifierPushbullet) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *GlobalNotifierPushbullet) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *GlobalNotifierPushbullet) GetName() string {
	return MODULE_NAME
}

func (m *GlobalNotifierPushbullet) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *GlobalNotifierPushbullet) InitModule(_cfg interface{}) error {
	m.config = _cfg.(*Config)
	return nil
}

func (m *GlobalNotifierPushbullet) Run(ctx context.Context) error {
	if err := m.config.NotifyPushbulletChannel(); err != nil {
		return err
	}

	return nil
}
func (m *GlobalNotifierPushbullet) GetConfig() interface{} {
	return &Config{}
}

func (m *GlobalNotifierPushbullet) Close() error {
	return nil
}

func (c *Config) NotifyPushbulletChannel() error {
	pb := pushbullet.New(c.APIKey)
	devs, err := pb.Devices()
	if err != nil {
		return err
	}

	err = pb.PushNote(devs[0].Iden, c.MessageTitle, c.MessageBody)
	if err != nil {
		return err
	}

	//SMS test
	user, err := pb.Me()
	if err != nil {
		return err
	}

	err = pb.PushSMS(user.Iden, devs[0].Iden, c.PhoneNumber, "Sms text")
	if err != nil {
		return err
	}
	return nil
}
