package slack

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jarcoal/httpmock"
)

func TestGetName(t *testing.T) {
	var conf Config
	name := conf.GetName()
	require.Equal(t, MODULE_NAME, name)
}

func TestGetConfig(t *testing.T) {
	config := GetConfig()
	require.Equal(t, Config{}, config)
}
func TestClose(t *testing.T) {
	var conf Config
	assert.Equal(t, nil, conf.Close())
}

func TestInitPipe(t *testing.T) {
	var local Local
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, local.InitPipe(bufOutput, bufInput))
}

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
		conf := Config{Hook: tt.Hook, MessageSuccess: tt.Message, Success: tt.Success}
		err := conf.NotifySlackChannel()
		assert.Equal(t, tt.Error, err)
	}
}
