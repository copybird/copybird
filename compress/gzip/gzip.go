package gzip

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
)

type Compress struct {
	reader io.Reader
	writer io.Writer
	level int
}

func (c *Compress) Init(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

func (c *Compress) InitCompress(level int) error {
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
