package web_callback

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

const MODULE_NAME = "web_callback"

type callback struct{
	targetUrl string
}

func (e *callback) GetName() string {
	return MODULE_NAME
}

func (e *callback) GetConfig() interface{} {
	return nil
}

func (c *callback) sendNotification (targetUrl string) error{

	//Set request body params
	data := url.Values{}
	data.Set("success", "true")

	req, err := http.NewRequest("GET", targetUrl, strings.NewReader(data.Encode()))
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
