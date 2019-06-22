package local

import (
	"testing"
	"bytes"
	"github.com/stretchr/testify/require"
)

//InitOutput initializes S3 with session
func TestGetName(t *testing.T) {
	var loc Local
	name := loc.GetName()
	require.Equal(t, "s3", name)
}

func TestGetConfig(t *testing.T) {
	var loc Local
	conf := loc.GetConfig()
	require.Equal(t, nil, conf)
}


func TestInitPipe(t *testing.T){
	var loc Local
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, loc.InitPipe(bufOutput, bufInput))
}

func TestInitModule(t *testing.T) {
	var loc Local
	err := loc.InitModule(Config{FileName: "test."})
	require.NoError(t, err, "should not be any error here")
}
