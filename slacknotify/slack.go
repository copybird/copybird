package slacknotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Sirupsen/logrus"
)

type SlackMessage struct {
	Text string `json:"text"`
}

const (
	HeaderContentType   = "Content-Type"
	MIMEApplicationJSON = "application/json"
)

func NotifySlackChannel(message, urls string) error {
	client := &http.Client{}

	slackMessage, err := json.Marshal(SlackMessage{Text: "@channel " + message})
	if err != nil {
		logrus.Error("Function 'notifySlackChannel', json.Marshal, error: ", err)
	}

	req, err := http.NewRequest(http.MethodPost, urls, bytes.NewBuffer(slackMessage))

	if err != nil {
		logrus.Error("Function 'notifySlackChannel', err: ", err)
		return err
	}

	req.Header.Set(HeaderContentType, MIMEApplicationJSON)

	resp, err := client.Do(req)

	if err != nil {
		logrus.Error("Function 'notifySlackChannel', Client error: ", err)
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Error("Function 'notifySlackChannel', Invalid request, status code: ", resp.StatusCode)
		return errors.New("Invalid request")
	}

	return nil
}
