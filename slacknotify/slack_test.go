package slacknotify

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"
)

func TestNotifySlackChannel(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const urls = "https://hooks.slack.com/services/TKBM/BKFLY1L/tL3RAwn9EYWMaMX"

	testCase := []struct {
		Responder httpmock.Responder
		Message   string
		Error     error
	}{
		{httpmock.NewStringResponder(200, "{}"), "Hello", nil},
		{httpmock.NewStringResponder(400, "{}"), "Hello", errors.New("StatusCode: 400")},
	}

	for _, tt := range testCase {
		httpmock.RegisterResponder("POST", urls, tt.Responder)
		err := NotifySlackChannel(tt.Message, urls)
		assert.Equal(t, tt.Error, err)
	}
}
