package scp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	s := &BackupOutputScp{}
	require.Equal(t, "scp", s.GetName())
}

func TestGetConfig(t *testing.T) {
	s := &BackupOutputScp{}
	require.Equal(t, &Config{}, s.GetConfig())
}

func TestInitPipe(t *testing.T) {
	s := &BackupOutputScp{}
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, s.InitPipe(bufOutput, bufInput))
}

func TestInitModule(t *testing.T) {
	s := &BackupOutputScp{}
	err := s.InitModule(&Config{
		Addr:     "127.0.0.1",
		Port:     2222,
		User:     "user",
		Password: "user",
	})
	require.NoError(t, err)
}
