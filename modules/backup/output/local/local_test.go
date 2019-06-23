package local

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	var loc Local
	name := GetName()
	require.Equal(t, "local", name)
}

func TestGetConfig(t *testing.T) {
	var loc Local
	conf := GetConfig()
	require.Equal(t, Config{
		DefaultMask: os.O_APPEND | os.O_CREATE | os.O_WRONLY,
		FileName:    "test.txt",
	}, conf)
}

func TestInitPipe(t *testing.T) {
	var loc Local
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, InitPipe(bufOutput, bufInput))
}

func TestRun(t *testing.T) {
	var loc Local
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, InitPipe(bufOutput, bufInput))
	conf := Config{
		DefaultMask: os.O_APPEND | os.O_CREATE | os.O_WRONLY,
		FileName:    "test.txt",
	}
	err := InitModule(conf)
	require.NoError(t, err)
	err = Run()
	require.NoError(t, err)
	os.Remove("test.txt")
}

func TestInitModule(t *testing.T) {
	var loc Local
	err := InitModule(Config{FileName: "test.sql"})
	require.NoError(t, err, "should not be any error here")
}
