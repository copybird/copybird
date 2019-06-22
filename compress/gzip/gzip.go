package gzip

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"

	"github.com/copybird/copybird/compress"
)

const MODULE_NAME = "GZIP"

type Compress struct {
	compress.Output
	reader io.Reader
	writer io.Writer
	level  int
}

func (c *Compress) GetName() string {
	return MODULE_NAME
}

func (c *Compress) GetConfig() interface{} {
	return nil
}

func (c *Compress) InitPipe(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

func (c *Compress) InitModule(_cfg interface{}) error {
	cfg := _cfg.(Config)
	level := cfg.level
	if level < -1 || level > 9 {
		return errors.New("compression level must be between -1 and 9")
	}
	c.level = level
	return nil
}

func (c *Compress) Run() error {
	buff := make([]byte, 4096)
	gw, err := gzip.NewWriterLevel(c.writer, c.level)
	if err != nil {
		return fmt.Errorf("cant start writer with error: %s", err)
	}
	defer gw.Close()

	for {
		_, err := c.reader.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read error: %s", err)
		}
		_, err = gw.Write(buff)
		if err != nil {
			return fmt.Errorf("write error: %s", err)
		}
	}
	return nil
}

func (c *Compress) Close() error {
	return nil
}
