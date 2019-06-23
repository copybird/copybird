package webcallback

import (
	"errors"
	"fmt"
	"testing"

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
		FailMsg    string
		Success    bool
		Error      error
	}{
		{httpmock.NewStringResponder(200, "{}"), "google.com", "Succes", "Fail", true, nil},
		{httpmock.NewStringResponder(400, "{}"), "google.com", "Succces", "Fail", false, errors.New("StatusCode: 400")},
	}

	for _, tc := range testCase {
		urls := fmt.Sprintf("%s", tc.TargetUrl)
		httpmock.RegisterResponder("GET", urls, tc.Responder)
		conf := GlobalNotifierWebcallback{Config: &Config{TargetUrl: tc.TargetUrl, SuccessMsg: tc.SuccessMsg, FailMsg: tc.FailMsg}}
		err := SendNotification()
		assert.Equal(t, tc.Error, err)
	}
}
