package lz4

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pierrec/lz4"
	"github.com/stretchr/testify/assert"
)

func TestDecompress(t *testing.T) {
	wr := &bytes.Buffer{}
	var decompressor DecompressLZ4

	s := "hello world"
	hw, err := helper(s)
	assert.NoError(t, err)
	spew.Dump(string(hw))
	rd := bytes.NewReader(hw)

	assert.NoError(t, decompressor.InitPipe(wr, rd))
	assert.NoError(t, decompressor.Run())

	assert.Equal(t, wr.String(), "hello world")
}

func helper(s string) ([]byte, error) {
	wr := &bytes.Buffer{}
	lw := lz4.NewWriter(wr)
	lw.Header = lz4.Header{CompressionLevel: 2}
	defer lw.Close()
	r := strings.NewReader(s)
	if _, err := io.Copy(lw, r); err != nil {
		return nil, err
	}
	return wr.Bytes(), nil
}
