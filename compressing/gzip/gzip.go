package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
)

type Compress struct {
	reader io.Reader
	writer io.Writer
}

func (c *Compress) Init(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

func (c *Compress) InitCompress(w io.Writer, r io.Reader) error {
	log.Printf("Initialize compress")
	return nil
}

func (c *Compress) Run() error {
	buff := make([]byte, 4096)
	gw := gzip.NewWriter(c.writer)

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
