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
//InitOutput initializes S3 with session
func TestInitOutput(t *testing.T) {

	var gcp GCP

	err := gcp.InitModule(map[string]string{"AWS_REGION": "us-east-1"})
	require.Error(t, err, "Should fail to find credentials")
 
	err = gcp.InitModule(map[string]string{"TOKEN_SOURCE": "some token source"})
	require.Error(t, err, "token source is not supported")

	err = gcp.InitModule(map[string]string{"CREDENTIALS_FILE": "creds.json"})
	require.Error(t, err, "credentials file is missing")
}
