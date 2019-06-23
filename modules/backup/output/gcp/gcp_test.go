package gcp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	var gcp BackupOutputGcp
	name := GetName()
	require.Equal(t, "gcp", name)
}

func TestGetConfig(t *testing.T) {
	var gcp BackupOutputGcp
	conf := GetConfig()
	require.Equal(t, Config{}, conf)
}

func TestInitPipe(t *testing.T) {
	var gcp BackupOutputGcp
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, InitPipe(bufOutput, bufInput))
}
func TestInitModule(t *testing.T) {

	var gcp BackupOutputGcp

	err := InitModule(Config{CredentialsFilePath: ""})
	require.Error(t, err, "Should fail to find credentials")

	err = InitModule(Config{CredentialsFilePath: "creds.json"})
	require.Error(t, err, "credentials file is missing")
}
