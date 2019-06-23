package awssqs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	g := &GlobalNotifierAWSSQS{}
	require.Equal(t, MODULE_NAME, g.GetName())
}

func TestGetConfig(t *testing.T) {
	g := &GlobalNotifierAWSSQS{}
	require.Equal(t, &Config{}, g.GetConfig())
}
func TestClose(t *testing.T) {
	g := &GlobalNotifierAWSSQS{}
	assert.Equal(t, nil, g.Close())
}

func TestInitPipe(t *testing.T) {
	g := &GlobalNotifierAWSSQS{}
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, g.InitPipe(bufOutput, bufInput))
}
