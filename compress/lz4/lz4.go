package lz4

import (
	"errors"
	"io"

	"github.com/copybird/copybird/compress"
	"github.com/pierrec/lz4"
)

var (
	errCompLevel = errors.New("compression level must be between -1 and 9")
)

const MODULE_NAME = "lz4"

// CompressLZ4 represents ...
type CompressLZ4 struct {
	compress.Output
	reader io.Reader
	writer io.Writer
	lz4    *lz4.Writer
}

func (c *CompressLZ4) GetName() string {
	return MODULE_NAME
}

func (c *CompressLZ4) GetConfig() interface{} {
	return &Config{}
}

func (c *CompressLZ4) InitPipe(w io.Writer, r io.Reader, _cfg interface{}) error {
	c.reader = r
	c.writer = w
	return nil
}

func (c *CompressLZ4) InitModule(_cfg interface{}) error {
	cfg := _cfg.(*Config)
	c.lz4 = lz4.NewWriter(c.writer)
	c.lz4.Header = lz4.Header{CompressionLevel: cfg.level}
	return nil
}

func (c *CompressLZ4) Run() error {

	// make a buffer to keep chunks that are read
	buf := make([]byte, 12)

	for {
		// read a chunk
		n, err := c.reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := c.lz4.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}

// Close closes compressor
func (c *CompressLZ4) Close() error {
	return nil
}
