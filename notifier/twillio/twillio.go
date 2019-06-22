package twillio

import (
	"errors"

	"github.com/sfreiberg/gotwilio"
)

const MODULE_NAME = "twillio"

type Twillio struct {
	Config *Config
	client *gotwilio.Twilio
}

func (t *Twillio) GetName() string {
	return MODULE_NAME
}

func (t *Twillio) GetConfig() interface{} {
	return Config{}
}

func (t *Twillio) InitModule(_conf interface{}) error {
	conf := _conf.(Config)
	t.client = gotwilio.NewTwilioClient(conf.AccountSid, conf.AuthToken)
	return nil
}

func (t *Twillio) Run() error {

	_, exception, err := t.client.SendSMS(t.Config.From, t.Config.To, "Dump created successfully", "", "")
	if err != nil {
		return err
	}
	if exception != nil {
		return errors.New(exception.Message)
	}
	return nil
}

func (t *Twillio) Close() error {
	return nil
}
