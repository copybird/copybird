package lz4

import (
	"errors"
	"fmt"
	"github.com/copybird/copybird/core"
	"io"

	// "os"

	"github.com/pierrec/lz4"
)

var (
	errCompLevel       = errors.New("compression level must be between -1 and 9")
	errNotCompressible = errors.New("is not compressible")
)

const MODULE_NAME = "lz4"

// RestoreDecompressLz4 represents ...
type RestoreDecompressLz4 struct {
	core.Module
	reader io.Reader
	writer io.Writer
}

func (m *RestoreDecompressLz4) GetName() string {
	return MODULE_NAME
}

func (m *RestoreDecompressLz4) GetConfig() interface{} {
	return &Config{}
}

func (m *RestoreDecompressLz4) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *RestoreDecompressLz4) InitModule(_cfg interface{}) error {
	return nil
}

func (m *RestoreDecompressLz4) Run() error {
	// make a buffer to keep chunks that are read
	lr := lz4.NewReader(m.reader)

	_, err := io.Copy(m.writer, lr)
	if err != nil {
		return fmt.Errorf("copy error: %s", err)
	}

	return nil
}

// Close closes compressor
func (m *RestoreDecompressLz4) Close() error {
	return nil
}
