package webcallback

import (
	"fmt"
	"github.com/copybird/copybird/core"
	"net/http"
	"net/url"
	"strings"
	"time"
	"errors"
)

const MODULE_NAME = "webcallback"

type Callback struct{
	core.Module
	Config     *Config
}

func (c *Callback) GetName() string {
	return MODULE_NAME
}

func (c *Callback) GetConfig() interface{} {
	return &Config{}
}

func (c *Callback) InitModule(_cfg interface{}) error {
	c.Config = _cfg.(*Config)
	return nil
}

func (c *Callback) Run() error {
	if err := c.SendNotification(); err != nil {
		return err
	}

	return nil
}

func (c *Callback) Close() error {
	return nil
}

func (c *Callback) SendNotification () error{

	//Set request body params
	data := url.Values{}
	data.Set("success", "true")

	req, err := http.NewRequest("GET", c.Config.TargetUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	resp , err := client.Do(req)
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
