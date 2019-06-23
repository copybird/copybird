package telegram

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	var local Local
	name := GetName()
	require.Equal(t, MODULE_NAME, name)
}

func TestGetConfig(t *testing.T) {
	var local Local
	config := GetConfig()
	require.Equal(t, Config{}, config)
}
func TestClose(t *testing.T) {
	var local Local
	assert.Equal(t, nil, Close())
}

func TestInitPipe(t *testing.T) {
	var local Local
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, InitPipe(bufOutput, bufInput))
}
