package notifier

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type callback struct{
	targetUrl string
}

func (c *callback) sendNotification () (){

	//Set request body params
	data := url.Values{}
	data.Set("success", "true")

	req, err := http.NewRequest("POST", c.targetUrl, strings.NewReader(data.Encode()))
	if err != nil {
		log.Print("Error reading request. ", err)
	}

	// Set headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	resp , err := client.Do(req)
	if err != nil {
		log.Print("Error reading response. ", err)
	}
}
