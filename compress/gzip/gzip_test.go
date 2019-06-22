package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress_InitCompress_Default_Compress(t *testing.T) {
	var compressor Compress
	err := compressor.InitCompress(-1)
	assert.Equal(t, err, nil)
	assert.Equal(t, compressor.level, -1)
}

func TestCompress_InitCompress_Compress_Level_Out_Of_range(t *testing.T) {
	var compressor Compress
	err := compressor.InitCompress(10)
	assert.Equal(t, err, nil)
}

func TestCompress_Run_Success_Compress(t *testing.T) {
	var compressor Compress
	rb := new(bytes.Buffer)
	wb := new(bytes.Buffer)
	rb.WriteString("hello, world.")

	_ = compressor.InitCompress(-1)
	_ = compressor.Init(wb, rb)
	err := compressor.Run()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	or, err := gzip.NewReader(bytes.NewReader(wb.Bytes()))
	if err != nil {
		panic(err)
	}
	s, _ := ioutil.ReadAll(or)
	//assert.Equal(t, s, "hello, world.")
}
