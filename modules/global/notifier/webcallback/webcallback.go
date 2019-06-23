package webcallback

import (
	"errors"
	"fmt"
	"github.com/copybird/copybird/core"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const GROUP_NAME = "global"
const TYPE_NAME = "notifier"
const MODULE_NAME = "webcallback"

type GlobalNotifierWebcallback struct {
	core.Module
	config *Config
}

func (m *GlobalNotifierWebcallback) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *GlobalNotifierWebcallback) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *GlobalNotifierWebcallback) GetName() string {
	return MODULE_NAME
}

func (m *GlobalNotifierWebcallback) GetConfig() interface{} {
	return &Config{}
}

func (m *GlobalNotifierWebcallback) InitModule(_cfg interface{}) error {
	m.config = _cfg.(*Config)
	return nil
}

func (m *GlobalNotifierWebcallback) Run() error {
	if err := m.SendNotification(); err != nil {
		return err
	}

	return nil
}

func (m *GlobalNotifierWebcallback) Close() error {
	return nil
}

func (m *GlobalNotifierWebcallback) SendNotification() error {

	//Set request body params
	data := url.Values{}
	data.Set("success", "true")

	req, err := http.NewRequest("GET", m.config.TargetUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		statusCode := fmt.Sprintf("%v", resp.StatusCode)
		return errors.New("StatusCode: " + statusCode)
	}

	defer resp.Body.Close()

	return nil
}
