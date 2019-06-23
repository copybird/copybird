package scp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	var scp SCP
	name := GetName()
	require.Equal(t, "scp", name)
}

func TestGetConfig(t *testing.T) {
	var scp SCP
	conf := GetConfig()
	require.Equal(t, Config{}, conf)
}

func TestInitPipe(t *testing.T) {
	var scp SCP
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, InitPipe(bufOutput, bufInput))
}

func TestInitModule(t *testing.T) {
	var scp SCP
	err := InitModule(Config{})
	require.NoError(t, err, "should not be any error here")
}
