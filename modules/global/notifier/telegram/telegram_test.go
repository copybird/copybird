package telegram

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	var local GlobalNotifierTelegram
	name := GetName()
	require.Equal(t, MODULE_NAME, name)
}

func TestGetConfig(t *testing.T) {
	var local GlobalNotifierTelegram
	config := GetConfig()
	require.Equal(t, Config{}, config)
}
func TestClose(t *testing.T) {
	var local GlobalNotifierTelegram
	assert.Equal(t, nil, Close())
}

func TestInitPipe(t *testing.T) {
	var local GlobalNotifierTelegram
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, InitPipe(bufOutput, bufInput))
}
