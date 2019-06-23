package twillio

import (
	"errors"
	"github.com/copybird/copybird/core"

	"github.com/sfreiberg/gotwilio"
)

const GROUP_NAME = "global"
const TYPE_NAME = "notifier"
const MODULE_NAME = "twillio"

type GlobalNotifierTwillio struct {
	core.Module
	config *Config
	client *gotwilio.Twilio
}

func (m *GlobalNotifierTwillio) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *GlobalNotifierTwillio) GetType() core.ModuleType {
	return TYPE_NAME
}

func (t *GlobalNotifierTwillio) GetName() string {
	return MODULE_NAME
}

func (t *GlobalNotifierTwillio) GetConfig() interface{} {
	return &Config{}
}

func (t *GlobalNotifierTwillio) InitModule(_conf interface{}) error {
	conf := _conf.(*Config)
	t.client = gotwilio.NewTwilioClient(conf.AccountSid, conf.AuthToken)
	return nil
}

func (t *GlobalNotifierTwillio) Run() error {

	_, exception, err := t.client.SendSMS(t.config.From, t.config.To, "Dump created successfully", "", "")
	if err != nil {
		return err
	}
	if exception != nil {
		return errors.New(exception.Message)
	}
	return nil
}

func (t *GlobalNotifierTwillio) Close() error {
	return nil
}
