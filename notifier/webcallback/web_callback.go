package webcallback

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

const MODULE_NAME = "webcallback"

type Callback struct{
	config     *Config
}

func (c *Callback) GetName() string {
	return MODULE_NAME
}

func (c *Callback) GetConfig() interface{} {
	return &Config{}
}


func (c *Callback) sendNotification () error{

	//Set request body params
	data := url.Values{}
	data.Set("success", "true")

	req, err := http.NewRequest("GET", c.config.targetUrl, strings.NewReader(data.Encode()))
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

	defer resp.Body.Close()

	return nil
}
