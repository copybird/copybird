package s3

import (
	"testing"
	"bytes"
	"github.com/stretchr/testify/require"
)

//InitOutput initializes S3 with session
func TestGetName(t *testing.T) {
	var s S3
	name := s.GetName()
	require.Equal(t, "s3", name)
}

func TestGetConfig(t *testing.T) {
	var s S3
	conf := s.GetConfig()
	require.Equal(t, nil, conf)
}


func TestInitPipe(t *testing.T){
	var s S3
	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	require.NoError(t, s.InitPipe(bufOutput, bufInput))
}

func TestInitOutput(t *testing.T) {
	var s S3
	err := s.InitModule(map[string]string{"AWS_REGION": "us-east-1"})
	require.NoError(t, err, "should not be any error here")
}
