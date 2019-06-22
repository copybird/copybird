package slacknotify

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"
)

func TestNotifySlackChannel(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCase := []struct {
		Responder httpmock.Responder
		Hook      string
		Message   string
		Success   bool
		Error     error
	}{
		{httpmock.NewStringResponder(200, "{}"), "TKBM/BKFLY1L/tL3RAwn9EYWMaMX", "Hello", true, nil},
		{httpmock.NewStringResponder(400, "{}"), "TKBM/YWMaMX", "Hello", false, errors.New("StatusCode: 400")},
	}

	for _, tt := range testCase {
		urls := fmt.Sprintf("%s/%s", SlackHookSite, tt.Hook)
		httpmock.RegisterResponder("POST", urls, tt.Responder)
		conf := Config{Hook: tt.Hook, MessageSuccess: tt.Message}
		err := conf.NotifySlackChannel(tt.Success)
		assert.Equal(t, tt.Error, err)
	}
}
