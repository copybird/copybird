package twillio

import (
	"context"
	"errors"

	"github.com/copybird/copybird/core"

	"github.com/sfreiberg/gotwilio"
)

const GROUP_NAME = "global"
const TYPE_NAME = "notifier"
const MODULE_NAME = "twillio"

type GlobalNotifierTwilio struct {
	core.Module
	config *Config
	client *gotwilio.Twilio
}

func (m *GlobalNotifierTwilio) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *GlobalNotifierTwilio) GetType() core.ModuleType {
	return TYPE_NAME
}

func (t *GlobalNotifierTwilio) GetName() string {
	return MODULE_NAME
}

func (t *GlobalNotifierTwilio) GetConfig() interface{} {
	return &Config{}
}

func (t *GlobalNotifierTwilio) InitModule(_conf interface{}) error {
	conf := _conf.(*Config)
	t.client = gotwilio.NewTwilioClient(conf.AccountSid, conf.AuthToken)
	return nil
}

func (t *GlobalNotifierTwilio) Run(ctx context.Context) error {

	_, exception, err := t.client.SendSMS(t.config.From, t.config.To, "Dump created successfully", "", "")
	if err != nil {
		return err
	}
	if exception != nil {
		return errors.New(exception.Message)
	}
	return nil
}

func (t *GlobalNotifierTwilio) Close() error {
	return nil
}
