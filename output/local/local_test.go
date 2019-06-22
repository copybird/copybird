package local

import (
	"testing"
	"bytes"
	"os"
	"github.com/stretchr/testify/require"
)

//InitOutput initializes S3 with session
func TestGetName(t *testing.T) {
	var loc Local
	name := loc.GetName()
	require.Equal(t, "local", name)
}

func TestGetConfig(t *testing.T) {
	var loc Local
	conf := loc.GetConfig()
	require.Equal(t, Config{
		DefaultMask:os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		FileName: "test.txt",
		}, conf)
}


func TestInitPipe(t *testing.T){
	var loc Local
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, loc.InitPipe(bufOutput, bufInput))
}

func TestRun(t *testing.T){
	var loc Local
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, loc.InitPipe(bufOutput, bufInput))
	conf := Config{
		DefaultMask: os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		FileName: "test.txt",
	}
	err := loc.InitModule(conf)
	require.NoError(t, err)
	err = loc.Run()
	require.NoError(t, err)
	os.Remove("test.txt")
}

func TestInitModule(t *testing.T) {
	var loc Local
	err := loc.InitModule(Config{FileName: "test.sql"})
	require.NoError(t, err, "should not be any error here")
}
