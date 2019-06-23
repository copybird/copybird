package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/copybird/copybird/core"
	"io"
	"net/http"
)

type GlobalNotifierSlack struct {
	core.Module
	config *Config
	reader io.Reader
	writer io.Writer
}

type SlackMessage struct {
	Text string `json:"text"`
}

func (m *GlobalNotifierSlack) GetName() string {
	return MODULE_NAME
}

func (m *GlobalNotifierSlack) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *GlobalNotifierSlack) InitModule(_cfg interface{}) error {
	m.config = _cfg.(*Config)
	return nil
}

func (m *GlobalNotifierSlack) Run() error {
	if err := m.config.NotifySlackChannel(); err != nil {
		return err
	}

	return nil
}
func (m *GlobalNotifierSlack) GetConfig() interface{} {
	return &Config{}
}

func (m *GlobalNotifierSlack) Close() error {
	return nil
}

const (
	MODULE_NAME         = "slack"
	SlackHookSite       = "https://hooks.slack.com/services"
	HeaderContentType   = "Content-Type"
	MIMEApplicationJSON = "application/json"
)

func (c *Config) NotifySlackChannel() error {

	urls := fmt.Sprintf("%s/%s", SlackHookSite, c.Hook)

	var slackMessage []byte
	var err error

	client := &http.Client{}
	if c.Success {
		slackMessage, err = json.Marshal(SlackMessage{Text: "<!channel> " + c.MessageSuccess})
	} else {
		slackMessage, err = json.Marshal(SlackMessage{Text: "<!channel> " + c.MessageFail})
	}

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, urls, bytes.NewBuffer(slackMessage))

	if err != nil {
		return err
	}

	req.Header.Set(HeaderContentType, MIMEApplicationJSON)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		statusCode := fmt.Sprintf("%v", resp.StatusCode)
		return errors.New("StatusCode: " + statusCode)
	}

	return nil
}
