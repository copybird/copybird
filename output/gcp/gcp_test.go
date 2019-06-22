package gcp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
