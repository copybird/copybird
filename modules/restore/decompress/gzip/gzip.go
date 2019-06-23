package gzip

import (
	"compress/gzip"
	"fmt"
	"io"

	"github.com/copybird/copybird/modules/backup/compress"
)

const MODULE_NAME = "gzip_decompress"

type Decompress struct {
	compress.Output
	reader io.Reader
	writer io.Writer
}

func (c *Decompress) GetName() string {
	return MODULE_NAME
}

func (c *Decompress) GetConfig() interface{} {
	return nil
}

func (c *Decompress) InitPipe(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

func (c *Decompress) InitModule(_cfg interface{}) error {
	return nil
}

func (c *Decompress) Run() error {
	gr, err := gzip.NewReader(c.reader)
	if err != nil {
		return fmt.Errorf("cant start gzip reader with error: %s", err)
	}
	defer gr.Close()
	_, err = io.Copy(c.writer, gr)
	if err != nil {
		return fmt.Errorf("copy error: %s", err)
	}
	return nil
}

func (c *Decompress) Close() error {
	return nil
}
