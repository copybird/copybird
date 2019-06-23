package nsq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/copybird/copybird/core"
)

const (
	MODULE_NAME         = "nsq"
	GROUP_NAME          = "global"
	TYPE_NAME           = "notifier"
	HeaderContentType   = "Content-Type"
	MIMEApplicationJSON = "application/json"
	NSQUrlSite          = "http://127.0.0.1:4151/pub?topic"
)

const ()

func (m *GlobalNotifierNSQ) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *GlobalNotifierNSQ) GetType() core.ModuleType {
	return TYPE_NAME
}

type GlobalNotifierNSQ struct {
	core.Module
	config *Config
	reader io.Reader
	writer io.Writer
}

func (m *GlobalNotifierNSQ) GetName() string {
	return MODULE_NAME
}

func (m *GlobalNotifierNSQ) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *GlobalNotifierNSQ) InitModule(_cfg interface{}) error {
	m.config = _cfg.(*Config)
	return nil
}

func (m *GlobalNotifierNSQ) Run() error {
	if err := m.config.NotifyNSQ(); err != nil {
		return err
	}

	return nil
}

func (m *GlobalNotifierNSQ) GetConfig() interface{} {
	return &Config{}
}

func (m *GlobalNotifierNSQ) Close() error {
	return nil
}

type NSQMessage struct {
	Message string `json:"message"`
}

func (c *Config) NotifyNSQ() error {
	urls := fmt.Sprintf("%s=%s", NSQUrlSite, c.TopicName)

	message, err := json.Marshal(NSQMessage{Message: c.Message})

	if err != nil {
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, urls, bytes.NewBuffer(message))

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
