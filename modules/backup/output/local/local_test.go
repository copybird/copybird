package local

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	l := &BackupOutputLocal{}
	require.Equal(t, "local", l.GetName())
}

func TestGetConfig(t *testing.T) {
	l := &BackupOutputLocal{}
	conf := l.GetConfig()
	require.Equal(t, &Config{
		DefaultMask: os.O_APPEND | os.O_CREATE | os.O_WRONLY,
		File:        "output",
	}, conf)
}

func TestInitPipe(t *testing.T) {
	l := &BackupOutputLocal{}
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, l.InitPipe(bufOutput, bufInput))
}

func TestRun(t *testing.T) {
	l := &BackupOutputLocal{}
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, l.InitPipe(bufOutput, bufInput))
	conf := &Config{
		DefaultMask: os.O_APPEND | os.O_CREATE | os.O_WRONLY,
		File:        "test.txt",
	}
	err := l.InitModule(conf)
	require.NoError(t, err)
	err = l.Run(context.TODO())
	require.NoError(t, err)
	os.Remove("test.txt")
}

func TestInitModule(t *testing.T) {
	l := &BackupOutputLocal{}
	err := l.InitModule(&Config{File: "test.sql"})
	require.NoError(t, err, "should not be any error here")
}
