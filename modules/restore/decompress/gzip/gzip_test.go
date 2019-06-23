package gzip

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var decompressor RestoreDecompressGzip

func TestCompress_InitCompress_Default_Compress(t *testing.T) {
	err := decompressor.InitModule(nil)
	assert.Equal(t, err, nil)
}

func TestCompress_Unzip_Success_Unzip(t *testing.T) {
	var wr bytes.Buffer
	rd := bytes.NewReader([]byte{
		0x1f, 0x8b, 0x08, 0x08, 0xc8, 0x58, 0x13, 0x4a,
		0x00, 0x03, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e,
		0x74, 0x78, 0x74, 0x00, 0xcb, 0x48, 0xcd, 0xc9,
		0xc9, 0x57, 0x28, 0xcf, 0x2f, 0xca, 0x49, 0xe1,
		0x02, 0x00, 0x2d, 0x3b, 0x08, 0xaf, 0x0c, 0x00,
		0x00, 0x00,
	})

	_ = decompressor.InitModule(nil)
	_ = decompressor.InitPipe(&wr, rd)
	err := decompressor.Run()
	assert.Equal(t, err, nil)
	assert.Equal(t, wr.String(), "hello world\n")
}

func TestCompress_Unzip_Unsuccess_Unzip(t *testing.T) {
	var wr bytes.Buffer
	rd := bytes.NewReader([]byte{})

	_ = decompressor.InitModule(nil)
	_ = decompressor.InitPipe(&wr, rd)
	err := decompressor.Run()
	assert.NotEqual(t, err, nil)
	assert.Equal(t, wr.String(), "")
}
