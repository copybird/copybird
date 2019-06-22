package gcp

import (
	"testing"
	"bytes"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	var gcp GCP
	name := gcp.GetName()
	require.Equal(t, "gcp", name)
}

func TestGetConfig(t *testing.T) {
	var gcp GCP
	conf := gcp.GetConfig()
	require.Equal(t, Config{}, conf)
}

func TestInitPipe(t *testing.T){
	var gcp GCP
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, gcp.InitPipe(bufOutput, bufInput))
}
//InitModule initializes S3 with session
func TestInitModule(t *testing.T) {

	var gcp GCP

	err := gcp.InitModule(Config{CredentialsFilePath: ""})
	require.Error(t, err, "Should fail to find credentials")
 
	err = gcp.InitModule(Config{CredentialsFilePath: "creds.json"})
	require.Error(t, err, "credentials file is missing")
}
