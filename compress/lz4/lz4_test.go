package lz4

import (
	"bytes"
	// "github.com/stretchr/testify/assert"
	"testing"
)

func TestCompressor_Compress_(t *testing.T) {
	var comp Compressor

	rb := new(bytes.Buffer)
	wb := new(bytes.Buffer)

	rb.WriteString("hello, world.")
	lenBC := rb.Len() // Length before compress

	_ = comp.InitLevel(-1)
	_ = comp.Init(wb, rb)

	err := comp.Compress()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	lenAC := wb.Len() // Length after compress

	if lenBC == lenAC {
		t.Fatalf("Bad compress!")
	}
}

func TestCompressor_CompressFile_(t *testing.T) {
	var comp Compressor

	rb := new(bytes.Buffer)
	wb := new(bytes.Buffer)

	_ = comp.InitLevel(-1)
	_ = comp.Init(wb, rb)

	input := "./test.txt"
	output := "./compressed.txt"
	comp.CompressFile(input, output)
}
