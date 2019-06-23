package lz4

import (
	"errors"
	"fmt"
	"io"

	// "os"

	"github.com/copybird/copybird/backup/compress"
	"github.com/pierrec/lz4"
)

var (
	errCompLevel       = errors.New("compression level must be between -1 and 9")
	errNotCompressible = errors.New("is not compressible")
)

const MODULE_NAME = "lz4"

// DecompressLZ4 represents ...
type DecompressLZ4 struct {
	compress.Output
	reader io.Reader
	writer io.Writer
}

func (c *DecompressLZ4) GetName() string {
	return MODULE_NAME
}

func (c *DecompressLZ4) GetConfig() interface{} {
	return &Config{}
}

func (c *DecompressLZ4) InitPipe(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

func (c *DecompressLZ4) InitModule(_cfg interface{}) error {
	return nil
}

func (c *DecompressLZ4) Run() error {
	// make a buffer to keep chunks that are read
	lr := lz4.NewReader(c.reader)

	_, err := io.Copy(c.writer, lr)
	if err != nil {
		return fmt.Errorf("copy error: %s", err)
	}

	return nil
}

// Close closes compressor
func (c *DecompressLZ4) Close() error {
	return nil
}
