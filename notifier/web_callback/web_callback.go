package web_callback

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

const MODULE_NAME = "web_callback"

type callback struct{
	config     *Config
}

func (c *callback) GetName() string {
	return MODULE_NAME
}

func (c *callback) GetConfig() interface{} {
	return &Config{}
}


func (c *callback) sendNotification () error{

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
