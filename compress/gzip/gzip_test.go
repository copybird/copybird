package gzip

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var compressor Compress
var cfg Config

func TestCompress_InitCompress_Default_Compress(t *testing.T) {
	cfg.Level = -1
	err := compressor.InitModule(cfg)
	assert.Equal(t, err, nil)
	assert.Equal(t, compressor.level, -1)
}

func TestCompress_InitCompress_Compress_Level_Out_Of_range(t *testing.T) {
	cfg.Level = 10
	err := compressor.InitModule(cfg)
	assert.NotEqual(t, err, nil)
}

func TestCompress_Run_Success_Compress(t *testing.T) {
	cfg.Level = -1

	rb := bytes.NewReader([]byte("hello, world."))
	wb := new(bytes.Buffer)

	_ = compressor.InitModule(cfg)
	_ = compressor.InitPipe(wb, rb)
	err := compressor.Run()
	assert.Equal(t, err, nil)

	data := make([]byte, 4096)
	var buff2 = new(bytes.Buffer)
	gr, err := gzip.NewReader(wb)
	defer gr.Close()

	data, err = ioutil.ReadAll(gr)
	assert.Equal(t, err, nil)
	buff2.Write(data)
	out := bytes.Trim(buff2.Bytes(), "\x00")
	assert.Equal(t, string(out), "hello, world.")
}

func TestCompress_Unzip_Success_Unzip(t *testing.T) {
	cfg.Level = -1

	rd := bytes.NewReader([]byte{
		0x1f, 0x8b, 0x08, 0x08, 0xc8, 0x58, 0x13, 0x4a,
		0x00, 0x03, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e,
		0x74, 0x78, 0x74, 0x00, 0xcb, 0x48, 0xcd, 0xc9,
		0xc9, 0x57, 0x28, 0xcf, 0x2f, 0xca, 0x49, 0xe1,
		0x02, 0x00, 0x2d, 0x3b, 0x08, 0xaf, 0x0c, 0x00,
		0x00, 0x00,
	})
	var wr bytes.Buffer

	_ = compressor.InitModule(cfg)
	_ = compressor.InitPipe(&wr, rd)
	err := compressor.Unzip()
	assert.Equal(t, err, nil)
	assert.Equal(t, wr.String(), "hello world")
}
