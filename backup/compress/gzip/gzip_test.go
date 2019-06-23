package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var compressor Compress
var cfg Config

func TestCompress_InitCompress_Default_Compress(t *testing.T) {
	Level = -1
	err := InitModule(cfg)
	assert.Equal(t, err, nil)
	assert.Equal(t, level, -1)
}

func TestCompress_InitCompress_Compress_Level_Out_Of_range(t *testing.T) {
	Level = 10
	err := InitModule(cfg)
	assert.NotEqual(t, err, nil)
}

func TestCompress_Run_Success_Compress(t *testing.T) {
	Level = -1

	rb := bytes.NewReader([]byte("hello, world."))
	wb := new(bytes.Buffer)

	_ = InitModule(cfg)
	_ = InitPipe(wb, rb)
	err := Run()
	assert.Equal(t, err, nil)

	var buff2 = new(bytes.Buffer)
	gr, err := gzip.NewReader(wb)
	defer gr.Close()

	_, err = io.Copy(buff2, gr)
	assert.Equal(t, err, nil)
	assert.Equal(t, buff2.String(), "hello, world.")
}
