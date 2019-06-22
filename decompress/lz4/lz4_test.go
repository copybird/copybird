package lz4

import (
	// "bytes"
	// "fmt"
	// "strings"
	// "testing"

	// "github.com/pierrec/lz4"
	// "gotest.tools/assert"
)

// TODO: fix this test
// func TestDecompress(t *testing.T) {
// 	var wr bytes.Buffer
// 	var decompressor DecompressLZ4

// 	s := "hello worldkkkkkk"
// 	hw := helper(s)
// 	rd := bytes.NewReader(hw)

// 	_ = decompressor.InitPipe(&wr, rd)
// 	_ = decompressor.Run()

// 	assert.Equal(t, wr.String(), "hello world")
// }

// func helper(s string) []byte {

// 	data := []byte(strings.Repeat(s, 100))
// 	buf := make([]byte, len(data))
// 	ht := make([]int, 64<<10) // buffer for the compression table

// 	n, err := lz4.CompressBlock(data, buf, ht)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	if n >= len(data) {
// 		fmt.Printf("`%s` is not compressible", s)
// 	}
// 	buf = buf[:n] // compressed data

// 	return buf
// }
