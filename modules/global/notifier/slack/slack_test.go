package slack

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jarcoal/httpmock"
)

func TestGetName(t *testing.T) {
	n := &GlobalNotifierSlack{}
	require.Equal(t, MODULE_NAME, n.GetName())
}

func TestGetConfig(t *testing.T) {
	n := &GlobalNotifierSlack{}
	require.Equal(t, &Config{}, n.GetConfig())
}
func TestClose(t *testing.T) {
	n := &GlobalNotifierSlack{}
	assert.Equal(t, nil, n.Close())
}

func TestInitPipe(t *testing.T) {
	n := &GlobalNotifierSlack{}
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, n.InitPipe(bufOutput, bufInput))
}

func TestNotifySlackChannel(t *testing.T) {

	n := &GlobalNotifierSlack{}
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
		assert.NoError(t, n.InitModule(&Config{Hook: tt.Hook, MessageSuccess: tt.Message, Success: tt.Success}))
		err := n.Run(context.TODO())
		assert.Equal(t, tt.Error, err)
	}
}

func TestRun(t *testing.T) {

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
		n := &GlobalNotifierSlack{}
		urls := fmt.Sprintf("%s/%s", SlackHookSite, tt.Hook)
		httpmock.RegisterResponder("POST", urls, tt.Responder)
		assert.NoError(t, n.InitModule(&Config{Hook: tt.Hook, MessageSuccess: tt.Message, Success: tt.Success}))
		err := n.Run(context.TODO())
		assert.Equal(t, tt.Error, err)
	}
}
