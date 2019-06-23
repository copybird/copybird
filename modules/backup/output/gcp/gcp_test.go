package gcp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	gcp := &BackupOutputGcp{}
	name := gcp.GetName()
	require.Equal(t, "gcp", name)
}

func TestGetConfig(t *testing.T) {
	gcp := &BackupOutputGcp{}
	conf := gcp.GetConfig()
	require.Equal(t, &Config{}, conf)
}

func TestInitPipe(t *testing.T) {
	gcp := &BackupOutputGcp{}
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, gcp.InitPipe(bufOutput, bufInput))
}
func TestInitModule(t *testing.T) {
	gcp := &BackupOutputGcp{}

	err := gcp.InitModule(&Config{CredentialsFilePath: ""})
	require.Error(t, err, "Should fail to find credentials")

	err = gcp.InitModule(&Config{CredentialsFilePath: "creds.json"})
	require.Error(t, err, "credentials file is missing")
}
