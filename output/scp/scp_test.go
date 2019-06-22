package scp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	var scp SCP
	name := scp.GetName()
	require.Equal(t, "scp", name)
}

func TestGetConfig(t *testing.T) {
	var scp SCP
	conf := scp.GetConfig()
	require.Equal(t, Config{}, conf)
}

func TestInitPipe(t *testing.T) {
	var scp SCP
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, scp.InitPipe(bufOutput, bufInput))
}

func TestInitModule(t *testing.T) {
	var scp SCP
	err := scp.InitModule(Config{})
	require.NoError(t, err, "should not be any error here")
}