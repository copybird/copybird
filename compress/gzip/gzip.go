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
	return &Config{}
}

func (c *Compress) InitPipe(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

func (c *Compress) InitModule(_cfg interface{}) error {
	cfg := _cfg.(Config)
	level := cfg.Level
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
		return fmt.Errorf("cant start gzip writer with error: %s", err)
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

func (c *Compress) Unzip() error {

	gr, err := gzip.NewReader(c.reader)
	if err != nil {
		return fmt.Errorf("cant start gzip reader with error: %s", err)
	}
	defer gr.Close()

	for {
		buff := make([]byte, 4096)

		_, err := gr.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read error: %s", err)
		}
		_, err = c.writer.Write(buff)
		if err != nil {
			return fmt.Errorf("write error: %s", err)
		}
	}
	return nil
}
