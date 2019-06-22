package gzip

import (
	"bufio"
	"compress/gzip"
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

func (c *Compress) Run() error  {
	gw := gzip.NewWriter(c.writer)

	for {
		buff := bufio.NewReader(c.reader)

		_, err := buff.WriteTo(gw)
		if err != nil {
			return err
		}
	}
}

func (c *Compress) Close() error  {
	return nil
}