package pushbullet

import (
	"github.com/copybird/copybird/core"
	"io"

	"github.com/xconstruct/go-pushbullet"
)

type GlobalNotifierPushbuller struct {
	core.Module
	config *Config
	reader io.Reader
	writer io.Writer
}

type Message struct {
	Text string `json:"text"`
}

func (m *GlobalNotifierPushbuller) GetName() string {
	return MODULE_NAME
}

func (m *GlobalNotifierPushbuller) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *GlobalNotifierPushbuller) InitModule(_cfg interface{}) error {
	m.config = _cfg.(*Config)
	return nil
}

func (m *GlobalNotifierPushbuller) Run() error {
	if err := m.config.NotifyPushbulletChannel(); err != nil {
		return err
	}

	return nil
}
func (m *GlobalNotifierPushbuller) GetConfig() interface{} {
	return &Config{}
}

func (m *GlobalNotifierPushbuller) Close() error {
	return nil
}

const MODULE_NAME = "pushbullet"

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
