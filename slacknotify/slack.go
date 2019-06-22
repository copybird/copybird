package slacknotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type SlackMessage struct {
	Text string `json:"text"`
}

func GetCongif() Config {
	conf := Config{}
	conf.MessageSuccess = os.Getenv("SLACK_NOTIFY_MESSAGE_SUCCESS")
	conf.MessageFail = os.Getenv("SLACK_NOTIFY_MESSAGE_FAIL")
	conf.Hook = os.Getenv("SLACK_HOOK")
	return conf
}

const (
	SlackHookSite       = "https://hooks.slack.com/services"
	HeaderContentType   = "Content-Type"
	MIMEApplicationJSON = "application/json"
)

func (c Config) NotifySlackChannel(Success bool) error {

	urls := fmt.Sprintf("%s/%s", SlackHookSite, c.Hook)

	var slackMessage []byte
	var err error

	client := &http.Client{}
	if Success {
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
