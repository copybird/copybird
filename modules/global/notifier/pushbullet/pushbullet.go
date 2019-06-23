package pushbullet

import (
	"io"

	"github.com/xconstruct/go-pushbullet"
)

type Local struct {
	config *Config
	reader io.Reader
	writer io.Writer
}

type SlackMessage struct {
	Text string `json:"text"`
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
	if err := l.config.NotifyPushbulletChannel(); err != nil {
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
