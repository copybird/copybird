package s3

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	var s BackupOutputS3
	name := GetName()
	require.Equal(t, "s3", name)
}

func TestGetConfig(t *testing.T) {
	var s BackupOutputS3
	conf := GetConfig()
	require.Equal(t, Config{}, conf)
}

func TestInitPipe(t *testing.T) {
	var s BackupOutputS3
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, InitPipe(bufOutput, bufInput))
}

func TestInitModule(t *testing.T) {
	var s BackupOutputS3
	err := InitModule(Config{Region: "us-east-1"})
	require.NoError(t, err, "should not be any error here")
}
