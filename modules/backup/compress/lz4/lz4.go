package lz4

import (
	"errors"
	"fmt"
	compress2 "github.com/copybird/copybird/modules/backup/compress"
	"io"

	"github.com/pierrec/lz4"
)

var (
	errCompLevel       = errors.New("compression level must be between -1 and 9")
	errNotCompressible = errors.New("is not compressible")
)

const MODULE_NAME = "lz4"

// CompressLZ4 represents ...
type CompressLZ4 struct {
	compress2.Output
	reader io.Reader
	writer io.Writer
	level  int
}

func (c *CompressLZ4) GetName() string {
	return MODULE_NAME
}

func (c *CompressLZ4) GetConfig() interface{} {
	return &Config{Level: 2}
}

func (c *CompressLZ4) InitPipe(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

func (c *CompressLZ4) InitModule(_cfg interface{}) error {
	cfg := _cfg.(*Config)
	if cfg.Level < -1 || cfg.Level > 9 {
		return errCompLevel
	}
	c.level = cfg.Level
	return nil
}

func (c *CompressLZ4) Run() error {
	lw := lz4.NewWriter(c.writer)
	lw.Header = lz4.Header{CompressionLevel: c.level}
	defer lw.Close()

	_, err := io.Copy(lw, c.reader)
	if err != nil {
		return fmt.Errorf("copy error: %s", err)
	}
	return nil
}

// Close closes compressor
func (c *CompressLZ4) Close() error {
	return nil
}
