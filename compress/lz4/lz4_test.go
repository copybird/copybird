package lz4

import (
	"bytes"
	"testing"

	"gotest.tools/assert"
)

func TestCompressLZ4(t *testing.T) {

	confLevel := 1
	comp := CompressLZ4{}

	rb := new(bytes.Buffer)
	rb.WriteString("hello, world.")

	wb := new(bytes.Buffer)
	lenBC := rb.Len()

	assert.Assert(t, comp.GetConfig() != nil)
	assert.NilError(t, comp.InitPipe(rb, wb, nil))
	assert.NilError(t, comp.InitModule(&Config{level: confLevel}))

	err := comp.Run()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	lenAC := wb.Len() // Length after compress

	if lenBC == lenAC {
		t.Fatalf("Bad compress!")
	}

	assert.NilError(t, comp.Close())
}
