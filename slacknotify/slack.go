package slacknotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
