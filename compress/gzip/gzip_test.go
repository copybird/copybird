package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var compressor Compress
var cfg Config


func TestCompress_InitCompress_Default_Compress(t *testing.T) {
	cfg.level = -1
	err := compressor.InitModule(cfg)
	assert.Equal(t, err, nil)
	assert.Equal(t, compressor.level, -1)
}

func TestCompress_InitCompress_Compress_Level_Out_Of_range(t *testing.T) {
	cfg.level = 10
	err := compressor.InitModule(cfg)
	assert.NotEqual(t, err, nil)
}

func TestCompress_Run_Success_Compress(t *testing.T) {
	cfg.level = -1

	rb := bytes.NewBufferString("hello, world.")
	wb := new(bytes.Buffer)

	_ = compressor.InitModule(cfg)
	_ = compressor.InitPipe(wb, rb)
	err := compressor.Run()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	data := make([]byte, 4096)
	var buff2  = new(bytes.Buffer)
	gr, err := gzip.NewReader(wb)
	defer gr.Close()

	data, err = ioutil.ReadAll(gr)
	if err != nil {
		print(err)
	}
	buff2.Write(data)
	out := bytes.Trim(buff2.Bytes(), "\x00")
	assert.Equal(t, string(out), "hello, world.")
}
