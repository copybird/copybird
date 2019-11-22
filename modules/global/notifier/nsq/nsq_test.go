package nsq

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	n := &GlobalNotifierNSQ{}
	require.Equal(t, MODULE_NAME, n.GetName())
}

func TestGetConfig(t *testing.T) {
	n := &GlobalNotifierNSQ{}
	require.Equal(t, &Config{}, n.GetConfig())
}
func TestClose(t *testing.T) {
	n := &GlobalNotifierNSQ{}
	assert.Equal(t, nil, n.Close())
}

func TestInitPipe(t *testing.T) {
	n := &GlobalNotifierNSQ{}
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, n.InitPipe(bufOutput, bufInput))
}

func TestNotifySlackChannel(t *testing.T) {

	n := &GlobalNotifierNSQ{}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCase := []struct {
		Responder httpmock.Responder
		TopicName string
		Message   string
		Error     error
	}{
		{httpmock.NewStringResponder(200, "{}"), "hi", "hi", nil},
		{httpmock.NewStringResponder(400, "{}"), "ss", "ss", errors.New("StatusCode: 400")},
	}

	for _, tt := range testCase {
		urls := fmt.Sprintf("%s=%s", NSQUrlSite, tt.TopicName)
		httpmock.RegisterResponder("POST", urls, tt.Responder)
		assert.NoError(t, n.InitModule(&Config{TopicName: tt.TopicName, Message: tt.Message}))
		err := n.Run(context.TODO())
		assert.Equal(t, tt.Error, err)
	}
}
