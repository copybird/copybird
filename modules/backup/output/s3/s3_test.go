package s3

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	s := &BackupOutputS3{}
	require.Equal(t, "s3", s.GetName())
}

func TestGetConfig(t *testing.T) {
	s := &BackupOutputS3{}
	require.Equal(t, &Config{}, s.GetConfig())
}

func TestInitPipe(t *testing.T) {
	s := &BackupOutputS3{}
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, s.InitPipe(bufOutput, bufInput))
}

func TestInitModule(t *testing.T) {
	s := &BackupOutputS3{}
	err := s.InitModule(&Config{Region: "us-east-1"})
	require.NoError(t, err, "should not be any error here")
}
