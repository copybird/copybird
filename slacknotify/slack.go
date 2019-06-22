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

type Config struct {
	Hook    string
	Message string
}

func GetCongifForSlack() Config {
	conf := Config{}
	conf.Message = os.Getenv("SLACK_NOTIFY_MESSAGE")
	conf.Message = os.Getenv("SLACK_HOOK")
	return conf
}

const (
	SlackHookSite       = "https://hooks.slack.com/services"
	HeaderContentType   = "Content-Type"
	MIMEApplicationJSON = "application/json"
)

func (c Config) NotifySlackChannel() error {

	urls := fmt.Sprintf("%s/%s", SlackHookSite, c.Hook)

	client := &http.Client{}

	slackMessage, err := json.Marshal(SlackMessage{Text: "<!channel> " + c.Message})
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
