package webcallback

import (
	"fmt"
	"testing"
	"errors"

	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"
)

func TestSendNotification(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCase := []struct {
		Responder  httpmock.Responder
		TargetUrl  string
		SuccessMsg string
		FailMsg	   string
		Success    bool
		Error     error
	}{
		{httpmock.NewStringResponder(200, "{}"), "google.com", "Succes", "Fail", true, nil},
		{httpmock.NewStringResponder(400, "{}"), "google.com", "Succces", "Fail", false, errors.New("StatusCode: 400")},
	}

	for _, tt := range testCase {
		urls := fmt.Sprintf("%s", tt.TargetUrl)
		httpmock.RegisterResponder("GET", urls, tt.Responder)
		conf := Callback{Config: &Config{TargetUrl: tt.TargetUrl, SuccessMsg: tt.SuccessMsg, FailMsg: tt.FailMsg}}
		err := conf.SendNotification()
		assert.Equal(t, tt.Error, err)
	}
}
