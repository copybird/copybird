package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	return nil
}
func (l *Local) Run() error {
	if err := l.config.NotifySlackChannel(); err != nil {
		return err
	}

	return nil
}
func (l *Local) GetConfig() interface{} {
	return Config{}
}

func (l *Local) Close() error {
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
